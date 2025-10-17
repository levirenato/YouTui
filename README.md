# YouTui 🎵

Player de YouTube para terminal com interface moderna e tema Catppuccin Mocha.

![Status](https://img.shields.io/badge/status-stable-green)
![Go Version](https://img.shields.io/badge/go-1.24+-blue)
![License](https://img.shields.io/badge/license-MIT-blue)

## ✨ Recursos

- Busca rápida no YouTube via `yt-dlp` (sem APIs)
- Thumbnails em alta qualidade com TrueColor
- Playlist com modos: Normal, Repeat 1, Repeat All, Shuffle
- Controles completos: Play, Pause, Stop, Next, Previous
- Barra de progresso em tempo real
- Modo áudio ou vídeo
- Interface colorida com tema Catppuccin Mocha
- Atalhos de teclado intuitivos (pressione `?` para ajuda)

## Requisitos

- **Go 1.24+**
- **mpv** - player de mídia
- **yt-dlp** - extrator de vídeos do YouTube
- **socat** - para comandos IPC do mpv

### Instalação dos requisitos

```bash
# Arch Linux
sudo pacman -S mpv yt-dlp socat
yay -S nerd-fonts-complete  # ou qualquer Nerd Font

# Debian/Ubuntu
sudo apt install mpv socat
pip install -U yt-dlp
# Baixe uma Nerd Font de: https://www.nerdfonts.com/

# macOS
brew install mpv yt-dlp socat
brew tap homebrew/cask-fonts
brew install --cask font-hack-nerd-font
```

## Instalação e Uso

```bash
# Clone o repositório
git clone https://github.com/levirenato/YouTui
cd YouTui

# Compile
go build -o youtui .

# Execute
./youtui
```

## Controles Principais

### Navegação
- `Tab` - Alternar entre painéis
- `/` - Focar na busca
- `↑/↓` - Navegar nas listas
- `?` - Mostrar ajuda

### Busca e Reprodução
- `Enter` (na busca) - Executar busca
- `Enter` (nos resultados) - Tocar imediatamente
- `a` - Adicionar à playlist
- `d` - Remover da playlist

### Player
- `c` ou `Space` - Pausar/Retomar
- `s` - Parar
- `n` - Próxima (somente em playlist)
- `b` - Anterior (somente em playlist)
- `r` - Alternar modo repetição
- `h` - Alternar shuffle
- `m` - Alternar áudio/vídeo
- `J/K` - Mover item na playlist

### Sair
- `q` - Sair da aplicação

## Licença

MIT License

---

**Desenvolvido com Go** • [tview](https://github.com/rivo/tview) • [yt-dlp](https://github.com/yt-dlp/yt-dlp) • [mpv](https://mpv.io/)
