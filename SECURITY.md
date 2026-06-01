# Security Policy

## Supported Versions

Only the latest release of go-keybind is supported with security updates.

## Reporting a Vulnerability

If you discover a security vulnerability in go-keybind, please report it responsibly. **Do not open a public issue.**

Email us at [us@floatpane.com](mailto:us@floatpane.com) with:

- A description of the vulnerability
- Steps to reproduce the issue
- The potential impact
- Any suggested fixes (optional)

We will acknowledge your report within 48 hours and aim to provide a fix or mitigation plan within 7 days, depending on severity.

## Scope

This policy covers the go-keybind codebase and its official releases.

This library reads user-controlled key strings and config files, so of particular interest:

- **Config file path traversal** — a manipulated `cfgDir` or `filename` argument to `Load`/`Save` that reads or writes files outside the intended directory.
- **Resource exhaustion** — a maliciously large or deeply nested JSON config file that causes excessive memory use during `Load`.
- **Panic via `MustParse`** — a key string in config that is accepted by `Parse` but triggers unexpected behavior at the call site; unlikely given the simple grammar but worth reporting.

This library has no external dependencies and does not process network input directly.

## Disclosure

We ask that you give us reasonable time to address the issue before disclosing it publicly. We are committed to crediting reporters in release notes (unless you prefer to remain anonymous).
