# YouTui-player

A modern YouTube player for the terminal with TUI interface.

![Go Version](https://img.shields.io/badge/go-1.24+-blue)
![License](https://img.shields.io/badge/license-MIT-blue)

## What does it do?

YouTui-player is a YouTube player that runs entirely in the terminal, allowing you to search, play music/videos, and manage playlists without leaving the command line. Beautiful interface with inline thumbnails, complete controls, and 4 Catppuccin themes (light + dark).

**Key features:**

- Fast YouTube search (no API keys required)
- High-quality thumbnails in terminal
- Playlist with shuffle, repeat, and navigation
- Complete controls (play, pause, next, previous)
- Real-time progress bar
- Audio and video playback modes
- Terminal video mode (renders video as unicode art via `mpv --vo=tct`)
- Configurable video quality (Best, 360p, 480p, 720p, 1080p, Terminal)
- Configurable video codec (Any, VP9, AV1)
- 4 Catppuccin themes (🌻 Latte, 🪴 Frappé, 🌺 Macchiato, 🌿 Mocha)
- Custom theme support
- Multilingual (PT-BR and EN)

## Screenshots

<img width="1917" height="1045" alt="image" src="https://github.com/user-attachments/assets/94df9e10-d1d5-4065-b668-0ae003def764" />
<img width="1903" height="1036" alt="image" src="https://github.com/user-attachments/assets/e4c9957a-c14b-4c68-9bf1-7a00e3579900" />

## Dependencies

- **Go 1.24+** - Programming language
- **mpv** - Media player
- **yt-dlp** - YouTube video extractor
- **socat** - IPC communication with mpv
- **Nerd Font** (optional) - For beautiful icons

## Installation

### Arch Linux (AUR) — recommended

No Go required. The AUR package handles everything automatically.

```bash
# Using yay
yay -S youtui-player

# Using paru
paru -S youtui-player

# Manually
git clone https://aur.archlinux.org/youtui-player.git
cd youtui-player
makepkg -si
```

After install, make sure you have the runtime dependencies:

```bash
sudo pacman -S mpv yt-dlp socat
```

---

### Manual (from source)

Requires **Go 1.24+**, **mpv**, **yt-dlp** and **socat**.

```bash
# Install runtime dependencies (Arch Linux)
sudo pacman -S mpv yt-dlp socat go

# Clone and build
git clone https://github.com/IvelOt/youtui-player
cd youtui-player
make build

# Run
./youtui-player

# Or install to /usr/local/bin
sudo make install-bin
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
| `m`       | Toggle audio/video   |

## Themes

YouTui-player includes 4 Catppuccin themes:

- 🌻 **Latte** - Elegant light mode
- 🪴 **Frappé** - Cool dark mode
- 🌺 **Macchiato** - Warm dark mode
- 🌿 **Mocha** - Deep dark mode (default)

**Switch theme:**

1. Press `Ctrl+C`
2. Select "Theme"
3. Choose from 4 available themes

Theme is automatically saved to `~/.config/youtui-player/youtui.conf`

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

Copyright (c) 2025 IvelOt

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
