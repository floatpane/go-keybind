package keybind

import (
	"errors"
	"testing"
)

func TestParse_SingleChar(t *testing.T) {
	tests := []struct{ in, wantBase string }{
		{"k", "k"},
		{"j", "j"},
		{"d", "d"},
		{"r", "r"},
		{"/", "/"},
		{"1", "1"},
		{"[", "["},
		{"]", "]"},
	}
	for _, tt := range tests {
		k, err := Parse(tt.in)
		if err != nil {
			t.Errorf("Parse(%q): unexpected error %v", tt.in, err)
			continue
		}
		if k.Ctrl || k.Alt || k.Shift {
			t.Errorf("Parse(%q): unexpected modifiers %+v", tt.in, k)
		}
		if k.Base != tt.wantBase {
			t.Errorf("Parse(%q).Base = %q, want %q", tt.in, k.Base, tt.wantBase)
		}
	}
}

func TestParse_Specials(t *testing.T) {
	tests := []struct{ in, wantBase string }{
		{"enter", "enter"},
		{"return", "enter"}, // alias
		{"esc", "esc"},
		{"escape", "esc"}, // alias
		{"tab", "tab"},
		{"backspace", "backspace"},
		{"bs", "backspace"}, // alias
		{"space", "space"},
		{"up", "up"},
		{"down", "down"},
		{"left", "left"},
		{"right", "right"},
		{"delete", "delete"},
		{"del", "delete"}, // alias
		{"insert", "insert"},
		{"home", "home"},
		{"end", "end"},
		{"pgup", "pgup"},
		{"pageup", "pgup"}, // alias
		{"pgdown", "pgdown"},
		{"pagedown", "pgdown"}, // alias
		{"pgdn", "pgdown"},     // alias
		{"f1", "f1"},
		{"f12", "f12"},
		{"F1", "f1"},   // case-insensitive
		{"ESC", "esc"}, // case-insensitive
	}
	for _, tt := range tests {
		k, err := Parse(tt.in)
		if err != nil {
			t.Errorf("Parse(%q): %v", tt.in, err)
			continue
		}
		if k.Base != tt.wantBase {
			t.Errorf("Parse(%q).Base = %q, want %q", tt.in, k.Base, tt.wantBase)
		}
		if k.Ctrl || k.Alt || k.Shift {
			t.Errorf("Parse(%q): unexpected modifiers", tt.in)
		}
	}
}

func TestParse_Modifiers(t *testing.T) {
	tests := []struct {
		in               string
		ctrl, alt, shift bool
		base             string
	}{
		{"ctrl+c", true, false, false, "c"},
		{"ctrl+C", true, false, false, "C"}, // uppercase preserved
		{"alt+enter", false, true, false, "enter"},
		{"meta+enter", false, true, false, "enter"}, // meta alias
		{"opt+enter", false, true, false, "enter"},  // opt alias
		{"shift+tab", false, false, true, "tab"},
		{"ctrl+shift+a", true, false, true, "a"},
		{"ctrl+alt+shift+f1", true, true, true, "f1"},
		{"CTRL+C", true, false, false, "C"},      // modifier case-insensitive
		{"Shift+Tab", false, false, true, "tab"}, // mixed case
	}
	for _, tt := range tests {
		k, err := Parse(tt.in)
		if err != nil {
			t.Errorf("Parse(%q): %v", tt.in, err)
			continue
		}
		if k.Ctrl != tt.ctrl || k.Alt != tt.alt || k.Shift != tt.shift || k.Base != tt.base {
			t.Errorf("Parse(%q) = %+v, want {Ctrl:%v Alt:%v Shift:%v Base:%q}",
				tt.in, k, tt.ctrl, tt.alt, tt.shift, tt.base)
		}
	}
}

func TestParse_Errors(t *testing.T) {
	tests := []struct {
		in      string
		wantErr error
	}{
		{"", ErrEmptyKey},
		{"   ", ErrEmptyKey},
		{"ctrl+", ErrUnknownKey}, // trailing +
		{"ctrl+unknown", ErrUnknownKey},
		{"badmod+c", ErrUnknownKey},
		{"f0", ErrUnknownKey},
		{"f13", ErrUnknownKey},
		{"delete+c", ErrUnknownKey}, // special used as modifier
	}
	for _, tt := range tests {
		_, err := Parse(tt.in)
		if err == nil {
			t.Errorf("Parse(%q): want error, got nil", tt.in)
			continue
		}
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Parse(%q): got %v, want wrapping %v", tt.in, err, tt.wantErr)
		}
	}
}

func TestKey_String(t *testing.T) {
	tests := []struct{ in, want string }{
		{"ctrl+c", "ctrl+c"},
		{"shift+tab", "shift+tab"},
		{"ctrl+shift+a", "ctrl+shift+a"},
		{"ctrl+alt+shift+f1", "ctrl+alt+shift+f1"},
		{"esc", "esc"},
		{"escape", "esc"},   // alias → canonical
		{"return", "enter"}, // alias → canonical
		{"k", "k"},
		// Modifier order is always ctrl < alt < shift regardless of input order.
		{"shift+ctrl+a", "ctrl+shift+a"},
		{"alt+ctrl+f2", "ctrl+alt+f2"},
	}
	for _, tt := range tests {
		k, err := Parse(tt.in)
		if err != nil {
			t.Errorf("Parse(%q): %v", tt.in, err)
			continue
		}
		if got := k.String(); got != tt.want {
			t.Errorf("Parse(%q).String() = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestKey_Equality(t *testing.T) {
	a, _ := Parse("esc")
	b, _ := Parse("escape")
	if a != b {
		t.Errorf("esc != escape: %v vs %v", a, b)
	}

	c, _ := Parse("ctrl+c")
	d, _ := Parse("ctrl+c")
	if c != d {
		t.Errorf("ctrl+c != ctrl+c")
	}

	e, _ := Parse("ctrl+shift+a")
	f, _ := Parse("shift+ctrl+a") // different input order
	if e != f {
		t.Errorf("ctrl+shift+a != shift+ctrl+a after parse: %v vs %v", e, f)
	}
}

func TestMustParse_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParse of invalid key should panic")
		}
	}()
	MustParse("ctrl+")
}

func TestMustParse_Valid(t *testing.T) {
	k := MustParse("ctrl+c")
	if !k.Ctrl || k.Base != "c" {
		t.Errorf("MustParse(ctrl+c) = %+v", k)
	}
}
