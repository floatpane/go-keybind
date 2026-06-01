package keybind

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testKeys struct {
	Quit   string `json:"quit"`
	Reload string `json:"reload"`
	Up     string `json:"up"`
}

var testDefaults = testKeys{
	Quit:   "ctrl+c",
	Reload: "r",
	Up:     "k",
}

func TestLoad_WritesDefaults(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir, "keybinds.json", testDefaults)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg != testDefaults {
		t.Errorf("got %+v, want defaults %+v", cfg, testDefaults)
	}
	// File should now exist.
	if _, err := os.Stat(filepath.Join(dir, "keybinds.json")); err != nil {
		t.Errorf("defaults not written: %v", err)
	}
}

func TestLoad_MergesWithDefaults(t *testing.T) {
	dir := t.TempDir()
	// Write partial config — only override Quit.
	partial := `{"quit":"ctrl+q"}`
	if err := os.WriteFile(filepath.Join(dir, "keybinds.json"), []byte(partial), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(dir, "keybinds.json", testDefaults)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Quit != "ctrl+q" {
		t.Errorf("Quit = %q, want ctrl+q", cfg.Quit)
	}
	// Fields absent from the partial file keep their default values.
	if cfg.Reload != testDefaults.Reload {
		t.Errorf("Reload = %q, want %q (default)", cfg.Reload, testDefaults.Reload)
	}
	if cfg.Up != testDefaults.Up {
		t.Errorf("Up = %q, want %q (default)", cfg.Up, testDefaults.Up)
	}
}

func TestLoad_ParsesExistingFile(t *testing.T) {
	dir := t.TempDir()
	full := testKeys{Quit: "q", Reload: "F5", Up: "up"}
	data, _ := json.Marshal(full)
	if err := os.WriteFile(filepath.Join(dir, "keybinds.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(dir, "keybinds.json", testDefaults)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg != full {
		t.Errorf("got %+v, want %+v", cfg, full)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "keybinds.json"), []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(dir, "keybinds.json", testDefaults); err == nil {
		t.Error("want error for invalid JSON")
	}
}

func TestSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	if err := Save(dir, "keybinds.json", testDefaults); err != nil {
		t.Fatalf("Save: %v", err)
	}
	cfg, err := Load(dir, "keybinds.json", testKeys{})
	if err != nil {
		t.Fatalf("Load after Save: %v", err)
	}
	if cfg != testDefaults {
		t.Errorf("round-trip: got %+v, want %+v", cfg, testDefaults)
	}
}

func TestSave_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "cfg")
	if err := Save(dir, "k.json", testDefaults); err != nil {
		t.Fatalf("Save into new dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "k.json")); err != nil {
		t.Error("file not created")
	}
}

func TestValidate_NoConflicts(t *testing.T) {
	areas := map[string]map[string]string{
		"global": {"quit": "ctrl+c", "cancel": "esc", "up": "k", "down": "j"},
		"inbox":  {"delete": "d", "archive": "a", "refresh": "r"},
	}
	if got := Validate(areas); len(got) != 0 {
		t.Errorf("want no conflicts, got %v", got)
	}
}

func TestValidate_IntraAreaConflict(t *testing.T) {
	areas := map[string]map[string]string{
		"inbox": {"delete": "d", "archive": "d"}, // same key
	}
	conflicts := Validate(areas)
	if len(conflicts) == 0 {
		t.Fatal("want conflict, got none")
	}
	if !strings.Contains(conflicts[0], "inbox") {
		t.Errorf("conflict should name area: %v", conflicts)
	}
	if !strings.Contains(conflicts[0], `"d"`) {
		t.Errorf("conflict should name key: %v", conflicts)
	}
}

func TestValidate_CrossAreaNotConflict(t *testing.T) {
	// Same key in two different areas is intentional (e.g. "d" for delete).
	areas := map[string]map[string]string{
		"inbox": {"delete": "d"},
		"email": {"delete": "d"},
	}
	if got := Validate(areas); len(got) != 0 {
		t.Errorf("cross-area duplicate should not be a conflict: %v", got)
	}
}

func TestValidate_EmptyKeySkipped(t *testing.T) {
	areas := map[string]map[string]string{
		"global": {"quit": "", "cancel": ""},
	}
	if got := Validate(areas); len(got) != 0 {
		t.Errorf("empty keys should not produce conflicts: %v", got)
	}
}

func TestValidate_MultipleConflicts(t *testing.T) {
	areas := map[string]map[string]string{
		"inbox":  {"delete": "d", "archive": "d"},
		"global": {"quit": "q", "cancel": "q"},
	}
	if got := Validate(areas); len(got) != 2 {
		t.Errorf("want 2 conflicts, got %d: %v", len(got), got)
	}
}
