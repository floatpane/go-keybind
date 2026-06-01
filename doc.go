// Package keybind provides two things for Go TUI applications:
//
// 1. A [Key] type — a parsed, normalized keyboard shortcut. [Parse] validates
// a key string (catching typos like "crtl+c" or "delte"), canonicalizes
// modifier order and special-key aliases ("escape" → "esc", "return" →
// "enter"), and produces a value that can be compared with == or used in a
// switch.
//
// 2. A generic JSON keybind-config loader ([Load], [Save], [Validate]) that
// any TUI app can use with its own config struct. It reads a JSON file from a
// config directory, merges with defaults so unknown fields survive upgrades,
// writes defaults on first run, and validates intra-area key conflicts.
//
// The package is framework-agnostic and has no external dependencies.
//
// Typical usage:
//
//	// Define your app's keybind config once.
//	type MyKeys struct {
//	    Quit   string `json:"quit"`
//	    Reload string `json:"reload"`
//	}
//
//	// Load from ~/.config/myapp/keybinds.json, writing defaults if missing.
//	cfg, err := keybind.Load(cfgDir, "keybinds.json", MyKeys{
//	    Quit:   "ctrl+c",
//	    Reload: "r",
//	})
//
//	// Validate for intra-area conflicts.
//	conflicts := keybind.Validate(map[string]map[string]string{
//	    "global": {"quit": cfg.Quit, "reload": cfg.Reload},
//	})
//
//	// Parse a key string from the config and compare against incoming events.
//	k, err := keybind.Parse(cfg.Quit)
//	if err != nil { /* bad value in config */ }
//	if incomingEvent == k.String() { /* handle quit */ }
package keybind
