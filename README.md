# YouTui 🎵

A modern YouTube player for the terminal with TUI interface.

![Go Version](https://img.shields.io/badge/go-1.24+-blue)
![License](https://img.shields.io/badge/license-MIT-blue)

## What does it do?

YouTui is a YouTube player that runs entirely in the terminal, allowing you to search, play music/videos, and manage playlists without leaving the command line. Beautiful interface with inline thumbnails, complete controls, and 4 Catppuccin themes (light + dark).

**Key features:**

- Fast YouTube search (no API keys required)
- High-quality thumbnails in terminal
- Playlist with shuffle, repeat, and navigation
- Complete controls (play, pause, next, previous)
- Real-time progress bar
- 4 Catppuccin themes (🌻 Latte, 🪴 Frappé, 🌺 Macchiato, 🌿 Mocha)
- Custom theme support
- Multilingual (PT-BR and EN)

## Dependencies

- **Go 1.24+** - Programming language
- **mpv** - Media player
- **yt-dlp** - YouTube video extractor
- **socat** - IPC communication with mpv
- **Nerd Font** (optional) - For beautiful icons

## Quick Install

```bash
# Clone the repository
git clone https://github.com/levirenato/YouTui
cd YouTui

# Install dependencies and compile (requires sudo)
make install

# Or just compile (if you already have dependencies)
make build

# Run
./youtui
```

## Main Shortcuts

| Key       | Action               |
| --------- | -------------------- |
| `/`       | Search               |
| `Enter`   | Play/Search          |
| `a`       | Add to playlist      |
| `d`       | Remove from playlist |
| `Space`   | Pause/Resume         |
| `n` / `b` | Next/Previous        |
| `h`       | Shuffle              |
| `r`       | Repeat mode          |
| `Tab`     | Switch panels        |
| `?`       | Full help            |
| `Ctrl+Q`  | Quit                 |
| `Ctrl+C`  | Settings             |

## Themes

YouTui includes 4 Catppuccin themes:

- 🌻 **Latte** - Elegant light mode
- 🪴 **Frappé** - Cool dark mode
- 🌺 **Macchiato** - Warm dark mode
- 🌿 **Mocha** - Deep dark mode (default)

**Switch theme:**
1. Press `Ctrl+C`
2. Select "Theme"
3. Choose from 4 available themes

Theme is automatically saved to `~/.config/youtui/youtui.conf`

**Custom theme:**
See [THEMES.md](THEMES.md) for instructions on how to create your own theme.

## Development

```bash
# Check dependencies
make check-deps

# Compile
make build

# Compile and run
make run

# Format code
make fmt

# Clean generated files
make clean
```

## License

MIT License

Copyright (c) 2025 LeviRenato

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
