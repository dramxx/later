# later — CLI cross-device clipboard via GitHub Gist

## Overview

`later` is a minimal Go CLI tool that lets you send text/links from one machine and read them on another, using a private GitHub Gist as the backend. No server, no cost, no sync client.

---

## User Flow

```
later send https://example.com
later send "read this carefully"
later inbox
later inbox --clear
later inbox --pop 1
later inbox --pop 1 3
```

---

## Project Structure

```
later/
├── main.go
├── go.mod
├── go.sum
├── config/
│   └── config.go        # load/parse ~/.config/later/config.toml
├── gist/
│   └── gist.go          # all GitHub Gist API interactions
├── cmd/
│   ├── send.go          # later send
│   ├── inbox.go         # later inbox
│   └── config.go        # later config
├── install.sh           # Linux/Mac installer
├── install.ps1          # Windows installer
└── PLAN.md
```

---

## Config File

Location:

- Linux/Mac: `~/.config/later/config.toml`
- Windows: `%APPDATA%\later\config.toml`

Format:

```toml
[gist]
token = "ghp_xxxx..."
gist_id = "1d726ef02757ca62c48defa1ab646bdc"
```

---

## Gist File Format

The gist contains a single file: `inbox.txt`

Example content:

```
LATER

[2026-03-24 19:42]  https://magazine.sebastianraschka.com/
[2026-03-24 19:45]  read the KV cache section carefully
[2026-03-24 19:46]  https://arxiv.org/abs/2501.12345
```

### Rules:

- Only lines starting with `[` are treated as entries
- All other lines (like the `LATER` header) are ignored silently
- Timestamps are in local time of the sending machine, format `[YYYY-MM-DD HH:MM]`
- One entry per line

---

## Commands

### `later send <text>`

- Reads config
- Fetches current gist content via GET
- Appends a new timestamped line: `[YYYY-MM-DD HH:MM]  <text>`
- Writes updated content back via PATCH
- Prints: `✓ saved`

### `later inbox`

- Reads config
- Fetches current gist content via GET
- Filters lines starting with `[`
- Prints them numbered:
  ```
  1  [2026-03-24 19:42]  https://example.com
  2  [2026-03-24 19:45]  read this carefully
  ```
- If empty, prints: `inbox is empty`

### `later inbox --clear`

- Reads config
- Fetches current gist content via GET
- Removes all entry lines (lines starting with `[`)
- Writes back only the non-entry lines (preserves the header)
- Prints: `✓ cleared`

### `later inbox --pop <n> [n...]`

- Reads config
- Fetches current gist content via GET
- Removes entry lines at the given 1-based indices
- Writes updated content back via PATCH
- Prints: `✓ removed 1 entry` or `✓ removed 3 entries`

### `later config`

- Resolves the config file path for the current OS
- If the config file does not exist yet, creates it with an empty template
- Opens it in the default editor:
  - Windows: `notepad.exe`
  - Linux: `gedit`
  - Mac: `open -e` (TextEdit)
- Uses `os/exec` to launch the editor process

---

## GitHub Gist API

Base URL: `https://api.github.com`

### GET gist content

```
GET /gists/{gist_id}
Authorization: Bearer {token}
Accept: application/vnd.github+json
```

Response: JSON with `files["inbox.txt"].content`

### PATCH gist content

```
PATCH /gists/{gist_id}
Authorization: Bearer {token}
Accept: application/vnd.github+json
Content-Type: application/json

{
  "files": {
    "inbox.txt": {
      "content": "<full updated content>"
    }
  }
}
```

---

## Dependencies

Standard library only — no external packages needed:

- `net/http` — HTTP client for Gist API calls
- `encoding/json` — parse/build API payloads
- `os` — file paths, environment
- `fmt` — output
- `strings` — line parsing
- `time` — timestamps
- `flag` or `os.Args` — CLI argument parsing (no cobra needed for this scope)

`go.mod` should declare `module later` and `go 1.21` (or latest stable).

---

## Error Handling

- Config file missing → print clear message: `config not found at ~/.config/later/config.toml — see README`
- Missing token or gist_id in config → print which field is missing
- HTTP error from GitHub API → print status code and response body
- Invalid `--pop` index (out of range, not a number) → print error, do not modify gist
- Network failure → print error message

All errors exit with code 1.

---

## Build & Distribution

Repository: `https://github.com/dramxx/later`

### Cross-compile

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o later .

# Windows
GOOS=windows GOARCH=amd64 go build -o later.exe .
```

Binaries are published as GitHub Releases assets, named:

- `later-linux-amd64`
- `later-windows-amd64.exe`

### GitHub Actions

A `.github/workflows/release.yml` workflow should be created that:

- Triggers on a new tag push (e.g. `v1.0.0`)
- Cross-compiles for Linux and Windows
- Uploads both binaries as Release assets

---

## Install Scripts

### `install.sh` (Linux / Mac)

Used via:

```bash
curl -sSL https://raw.githubusercontent.com/dramxx/later/main/install.sh | sh
```

Script should:

1. Detect OS and arch
2. Download the correct binary from the latest GitHub Release:
   `https://github.com/dramxx/later/releases/latest/download/later-linux-amd64`
3. Move it to `/usr/local/bin/later`
4. Make it executable (`chmod +x`)
5. Print `✓ later installed` with the version

### `install.ps1` (Windows)

Used via:

```powershell
powershell -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/dramxx/later/main/install.ps1 | iex"
```

Script should:

1. Download the binary from the latest GitHub Release:
   `https://github.com/dramxx/later/releases/latest/download/later-windows-amd64.exe`
2. Save it to `$env:LOCALAPPDATA\later\later.exe`
3. Add `$env:LOCALAPPDATA\later` to the user's PATH (permanent, via registry)
4. Print `✓ later installed` with the version

---

## README (to include in repo)

Should cover:

1. Create a private GitHub Gist with a file named `inbox.txt`
2. Generate a GitHub personal access token with only the `gist` scope
3. Create config file at the appropriate path with token and gist_id
4. Install via the one-liner for your platform
5. Usage examples

---

## Out of Scope

- Images or binary files
- Encryption of gist content
- Multiple inboxes
- Tagging or categorizing entries
- Any TUI or interactive mode
