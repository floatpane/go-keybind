package keybind

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Load reads a JSON keybind config file from cfgDir/filename into a value of
// type T, starting from defaults so that any key absent from the file keeps
// its default value (merge semantics). If the file does not exist it is
// created with the JSON representation of defaults.
//
// T must be JSON-marshalable. A typical call looks like:
//
//	cfg, err := keybind.Load(cfgDir, "keybinds.json", MyKeys{
//	    Quit:   "ctrl+c",
//	    Reload: "r",
//	})
func Load[T any](cfgDir, filename string, defaults T) (T, error) {
	path := filepath.Join(cfgDir, filename)

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return defaults, fmt.Errorf("keybind: read %s: %w", path, err)
		}
		// File absent — create the directory and write defaults.
		if err := os.MkdirAll(cfgDir, 0o700); err != nil {
			return defaults, fmt.Errorf("keybind: mkdir %s: %w", cfgDir, err)
		}
		encoded, err := json.MarshalIndent(defaults, "", "  ")
		if err != nil {
			return defaults, fmt.Errorf("keybind: marshal defaults: %w", err)
		}
		if err := os.WriteFile(path, encoded, 0o600); err != nil {
			return defaults, fmt.Errorf("keybind: write defaults to %s: %w", path, err)
		}
		return defaults, nil
	}

	// Unmarshal into a copy of defaults so missing fields keep their values.
	cfg := defaults
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaults, fmt.Errorf("keybind: parse %s: %w", path, err)
	}
	return cfg, nil
}

// Save writes cfg as indented JSON to cfgDir/filename, creating cfgDir if
// needed. It is the write-side counterpart to [Load].
func Save[T any](cfgDir, filename string, cfg T) error {
	if err := os.MkdirAll(cfgDir, 0o700); err != nil {
		return fmt.Errorf("keybind: mkdir %s: %w", cfgDir, err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("keybind: marshal: %w", err)
	}
	path := filepath.Join(cfgDir, filename)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("keybind: write %s: %w", path, err)
	}
	return nil
}

// Validate checks for intra-area conflicts: two different actions within the
// same area mapped to the same (non-empty) key. Cross-area duplicates are
// intentional and not reported — "d" for delete in both an inbox view and an
// email view is normal.
//
// areas maps area name → (action name → key string). The shape mirrors a
// struct-per-area keybind config, and is easily built from one:
//
//	keybind.Validate(map[string]map[string]string{
//	    "global": {"quit": cfg.Global.Quit, "cancel": cfg.Global.Cancel},
//	    "inbox":  {"delete": cfg.Inbox.Delete, "archive": cfg.Inbox.Archive},
//	})
//
// Each returned string is a human-readable conflict description.
func Validate(areas map[string]map[string]string) []string {
	var conflicts []string
	for area, bindings := range areas {
		seen := make(map[string]string, len(bindings)) // key → first action
		for action, key := range bindings {
			if key == "" {
				continue
			}
			if prev, ok := seen[key]; ok {
				conflicts = append(conflicts,
					fmt.Sprintf("conflict in %s: key %q used for both %q and %q", area, key, prev, action))
			} else {
				seen[key] = action
			}
		}
	}
	return conflicts
}
