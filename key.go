package keybind

import (
	"fmt"
	"strings"
)

// Key is a parsed, normalized keyboard shortcut.
//
// The zero value is not a valid binding; [Parse] always returns a non-zero Key
// on success.
//
// Key is comparable — two Keys are equal (==) when they represent the same
// shortcut, regardless of how the original strings were written ("esc" and
// "escape" produce the same Key).
type Key struct {
	Ctrl  bool
	Alt   bool
	Shift bool
	// Base is the normalized base key name. For printable characters it is the
	// literal character (e.g. "c", "1", "/"). For special keys it is one of the
	// canonical names listed below.
	Base string
}

// Canonical base key names for special keys.
const (
	keyEnter     = "enter"
	keyEsc       = "esc"
	keyTab       = "tab"
	keyBackspace = "backspace"
	keySpace     = "space"
	keyUp        = "up"
	keyDown      = "down"
	keyLeft      = "left"
	keyRight     = "right"
	keyDelete    = "delete"
	keyInsert    = "insert"
	keyHome      = "home"
	keyEnd       = "end"
	keyPgUp      = "pgup"
	keyPgDown    = "pgdown"
)

// canonicalSpecials maps every accepted base key name — including aliases — to
// its canonical form. Single printable characters are not listed here; they are
// accepted as-is.
var canonicalSpecials = map[string]string{
	// enter
	keyEnter: keyEnter,
	"return": keyEnter,
	// escape
	keyEsc:   keyEsc,
	"escape": keyEsc,
	// tab
	keyTab: keyTab,
	// backspace
	keyBackspace: keyBackspace,
	"bs":         keyBackspace,
	// space
	keySpace: keySpace,
	// arrow keys
	keyUp:    keyUp,
	keyDown:  keyDown,
	keyLeft:  keyLeft,
	keyRight: keyRight,
	// delete / insert
	keyDelete: keyDelete,
	"del":     keyDelete,
	keyInsert: keyInsert,
	// navigation
	keyHome:    keyHome,
	keyEnd:     keyEnd,
	keyPgUp:    keyPgUp,
	"pageup":   keyPgUp,
	keyPgDown:  keyPgDown,
	"pagedown": keyPgDown,
	"pgdn":     keyPgDown,
}

func init() {
	// F1..F12
	for i := 1; i <= 12; i++ {
		name := fmt.Sprintf("f%d", i)
		canonicalSpecials[name] = name
	}
}

// Parse parses and normalizes a key string such as "ctrl+c", "shift+tab",
// "ctrl+shift+a", "f1", "esc", or "k".
//
// Modifier names (case-insensitive): ctrl, alt, meta, opt (alt aliases), shift.
// Base key names: any single printable character, or a special name from the
// table below. Modifier order in the input does not matter; [Key.String]
// always emits ctrl, alt, shift in that order.
//
// Accepted special key names (and their aliases):
//
//	enter (return)  esc (escape)  tab  backspace (bs)  space
//	up  down  left  right  delete (del)  insert
//	home  end  pgup (pageup)  pgdown (pagedown, pgdn)
//	f1 – f12
//
// Returns [ErrEmptyKey] for an empty string, and [ErrUnknownKey] for an
// unrecognized base or modifier.
func Parse(s string) (Key, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Key{}, ErrEmptyKey
	}

	parts := strings.Split(s, "+")
	// A trailing "+" means the last part is empty — the user tried to bind the
	// "+" character, which we don't support (it's ambiguous with the separator).
	if parts[len(parts)-1] == "" {
		return Key{}, fmt.Errorf("%w: %q (use shift+= for plus)", ErrUnknownKey, s)
	}

	var k Key
	// Everything except the last segment is a modifier.
	for _, part := range parts[:len(parts)-1] {
		switch strings.ToLower(part) {
		case "ctrl":
			k.Ctrl = true
		case "alt", "meta", "opt":
			k.Alt = true
		case "shift":
			k.Shift = true
		default:
			return Key{}, fmt.Errorf("%w: unknown modifier %q in %q", ErrUnknownKey, part, s)
		}
	}

	base := strings.ToLower(parts[len(parts)-1])
	if canon, ok := canonicalSpecials[base]; ok {
		k.Base = canon
	} else if len([]rune(parts[len(parts)-1])) == 1 {
		// Preserve the original case of single-character keys so that "A" (an
		// uppercase letter sent by some frameworks without a Shift modifier) and
		// "a" remain distinct.
		k.Base = parts[len(parts)-1]
	} else {
		return Key{}, fmt.Errorf("%w: %q", ErrUnknownKey, base)
	}

	return k, nil
}

// MustParse parses a key string and panics on error. Useful for package-level
// var initialization and test helpers.
func MustParse(s string) Key {
	k, err := Parse(s)
	if err != nil {
		panic("keybind.MustParse: " + err.Error())
	}
	return k
}

// String returns the canonical form of the key, suitable for comparison with
// key-event strings produced by TUI frameworks. Modifiers appear in ctrl, alt,
// shift order followed by the base: e.g. "ctrl+alt+shift+a", "ctrl+c", "esc".
func (k Key) String() string {
	var b strings.Builder
	if k.Ctrl {
		b.WriteString("ctrl+")
	}
	if k.Alt {
		b.WriteString("alt+")
	}
	if k.Shift {
		b.WriteString("shift+")
	}
	b.WriteString(k.Base)
	return b.String()
}

// Sentinel errors returned by [Parse].
var (
	// ErrEmptyKey is returned when the input string is empty or all whitespace.
	ErrEmptyKey = fmt.Errorf("keybind: empty key string")
	// ErrUnknownKey is returned when the base key or a modifier is not recognized.
	ErrUnknownKey = fmt.Errorf("keybind: unknown key")
)
