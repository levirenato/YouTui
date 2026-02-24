# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build       # Compile the project
make run         # Build and run
make test        # Run all tests (go test ./...)
make fmt         # Format code (go fmt ./...)
make vet         # Static analysis (go vet ./...)
make check-deps  # Verify runtime deps (mpv, yt-dlp, socat)
make clean       # Remove compiled binaries
make deps        # Download and tidy Go modules
```

Run a single test: `go test ./internal/... -run TestName`

## Runtime Dependencies

The app shells out to external tools at runtime:
- `mpv` — media player backend
- `yt-dlp` — YouTube search and extraction
- `socat` — IPC communication with mpv via Unix socket

## Architecture

**Module:** `github.com/IvelOt/youtui-player`

Three internal packages:

- `internal/ui/` — Main TUI layer. The central `SimpleApp` struct (defined in `app.go`) owns all state and UI components. All subsystems are wired in `setup.go`; input handling lives in `handlers.go`.
- `internal/search/` — YouTube search via `yt-dlp` subprocess; parses JSON output into `Track` structs.
- `internal/config/` — TOML config (`~/.config/youtui-player/youtui.conf`) and JSON session state (`~/.local/state/youtui-player/state.json`).

**Key data flow:**

1. User search → `search.SearchVideos()` spawns `yt-dlp`, returns `[]Track`
2. Playback → `playTrackSimple()` spawns `mpv` with `--input-ipc-server`; progress polling uses `socat` to query the IPC socket
3. State is auto-saved to JSON after changes and restored on startup

**Theme system:** Four built-in Catppuccin variants (Latte, Frappé, Macchiato, Mocha) plus custom themes loaded from TOML files. See `THEMES.md` for the color palette format.

**i18n:** `internal/ui/i18n.go` provides PT-BR and EN string packs, auto-detected from `LC_ALL`/`LANG` environment variables or the config file.

**Playback config:** `[playback]` section in the TOML config exposes three user-facing options — `default_mode` (`"audio"`/`"video"`), `video_quality` (`"best"`, `"360"`, `"480"`, `"720"`, `"1080"`), and `video_codec` (`""`, `"vp9"`, `"av1"`). These are read at startup and cycle-able via the settings modal (Ctrl+C). The values are translated into a `--ytdl-format` selector by `buildYtdlFormat()` in `player.go` and passed to mpv at play time.

**UI framework:** `github.com/rivo/tview` + `github.com/gdamore/tcell/v2`. The custom `CustomList` widget (`custom_list.go`) extends tview with thumbnail rendering support.
