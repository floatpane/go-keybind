<div align="center">

# go-keybind

**Key-binding config for Go TUI apps. Parse and normalize key strings, load/save JSON keybind files with merge-with-defaults semantics, validate for intra-area conflicts.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/floatpane/go-keybind)](https://golang.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/floatpane/go-keybind.svg)](https://pkg.go.dev/github.com/floatpane/go-keybind)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/floatpane/go-keybind)](https://github.com/floatpane/go-keybind/releases)
[![CI](https://github.com/floatpane/go-keybind/actions/workflows/ci.yml/badge.svg)](https://github.com/floatpane/go-keybind/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

</div>

`go-keybind` provides two focused utilities for Go TUI applications:

1. **A `Key` type** — parse and normalize keyboard shortcut strings. Catches
   typos (`"crtl+c"`), canonicalizes aliases (`"escape"` → `"esc"`, `"return"` →
   `"enter"`, `"meta"` → `"alt"`), fixes modifier order, and produces a comparable
   value you can switch on.

2. **A generic JSON config loader** — `Load[T]`, `Save[T]`, and `Validate` work
   with any app-defined struct. Fields the user didn't touch keep their compiled-in
   defaults, the file is written on first run, and `Validate` catches intra-area
   conflicts without false-positives for intentional cross-area duplicates.

Extracted from [matcha](https://github.com/floatpane/matcha)'s config package.
Zero external dependencies.

## Install

```bash
go get github.com/floatpane/go-keybind
```

Requires Go 1.26+.

## Usage

### Parse and compare key strings

```go
k, err := keybind.Parse("ctrl+shift+a")
// k.Ctrl = true, k.Shift = true, k.Base = "a"
// k.String() = "ctrl+shift+a"

// Aliases normalize to the same Key:
a, _ := keybind.Parse("escape")
b, _ := keybind.Parse("esc")
fmt.Println(a == b) // true

// Modifier order doesn't matter in input:
c, _ := keybind.Parse("shift+ctrl+a")
d, _ := keybind.Parse("ctrl+shift+a")
fmt.Println(c == d) // true
```

### Load a keybind config from disk

```go
type MyKeys struct {
    Quit   string `json:"quit"`
    Reload string `json:"reload"`
    NavUp  string `json:"nav_up"`
}

cfg, err := keybind.Load(cfgDir, "keybinds.json", MyKeys{
    Quit:   "ctrl+c",
    Reload: "r",
    NavUp:  "k",
})
// If keybinds.json didn't exist it was created with the defaults.
// Fields absent from an existing file keep their default values.
```

### Validate for conflicts

```go
conflicts := keybind.Validate(map[string]map[string]string{
    "global": {"quit": cfg.Quit, "reload": cfg.Reload},
    "nav":    {"up": cfg.NavUp, "down": cfg.NavDown},
})
// Cross-area duplicates ("d" for delete in two different views) are fine.
// Only intra-area conflicts (two actions bound to the same key) are reported.
```

### Use in a TUI event loop

```go
quitKey := keybind.MustParse(cfg.Quit) // panics on bad config value

// In your event handler:
if event.String() == quitKey.String() {
    return tea.Quit
}
```

## Supported keys

**Modifiers:** `ctrl`, `alt` (= `meta` = `opt`), `shift`

**Special keys:**

| Canonical | Aliases |
|-----------|---------|
| `enter` | `return` |
| `esc` | `escape` |
| `tab` | — |
| `backspace` | `bs` |
| `space` | — |
| `up`, `down`, `left`, `right` | — |
| `delete` | `del` |
| `insert` | — |
| `home`, `end` | — |
| `pgup` | `pageup` |
| `pgdown` | `pagedown`, `pgdn` |
| `f1` … `f12` | — |

Any single printable character is also accepted (`a`, `A`, `1`, `/`, `[`…).

## Documentation

Full API reference: [pkg.go.dev/github.com/floatpane/go-keybind](https://pkg.go.dev/github.com/floatpane/go-keybind)

Guides: see [`docs/`](docs/).

## Sister projects

| Project | Role |
|---------|------|
| [floatpane/matcha](https://github.com/floatpane/matcha) | Reference consumer — inbox, email, composer, folder keybinds. |
| [floatpane/go-icalendar](https://github.com/floatpane/go-icalendar) | Sibling extraction — iCalendar parsing and iMIP replies. |
| [floatpane/go-secretbox](https://github.com/floatpane/go-secretbox) | Sibling extraction — password-based encryption at rest. |

## Contributing

PRs welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

Report vulnerabilities privately via [SECURITY.md](SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
