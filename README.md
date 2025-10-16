# YouTui

Interface TUI (Terminal User Interface) para buscar e assistir vídeos do YouTube diretamente no terminal.

## Requisitos

- **Go 1.24+** (para compilar)
- **mpv** - player de vídeo
- **yt-dlp** - extrator de vídeos do YouTube

### Instalação dos requisitos

```bash
# Debian/Ubuntu
sudo apt install mpv
pip install -U yt-dlp

# Arch Linux
sudo pacman -S mpv yt-dlp

# macOS
brew install mpv yt-dlp
```

## Compilação

```bash
go build -o youtui ./cmd/youtui
```

## Uso

```bash
./youtui
```

### Interface em Grid

A interface é dividida em **4 painéis**:

```
┌─────────────────┬─────────────────┐
│  🔍 Busca       │  📋 Playlist    │
│                 │                 │
├─────────────────┼─────────────────┤
│  📺 Resultados  │  🎵 Controles   │
│                 │                 │
└─────────────────┴─────────────────┘
```

### Controles

- **Tab** - Alternar entre painéis (Busca → Playlist → Resultados → Controles)
- **Enter** - Buscar (no painel de busca) ou Reproduzir (nos resultados/playlist)
- **↑/↓** ou **j/k** - Navegar pelos resultados
- **m** - Alternar modo de reprodução (Vídeo MP4 / Áudio MP3)
- **p** - Adicionar item selecionado à playlist
- **q** ou **Ctrl+C** - Sair

### Funcionalidades

- 🎨 **Layout em Grid** - Interface dividida em 4 seções organizadas
- 📋 **Playlist** - Adicione vídeos à playlist com a tecla 'p'
- 🎵 **Modo Áudio/Vídeo** - Alterne entre reproduzir vídeo completo ou apenas áudio
- 🎯 **Navegação por Painéis** - Use Tab para focar em diferentes seções
- 📊 **Painel de Informações** - Veja estatísticas e atalhos disponíveis
- 🎨 **Visual Moderno** - Bordas coloridas indicam o painel ativo

## Notas sobre reprodução de vídeos (2025)

O YouTube começou a exigir **PO Tokens** para alguns formatos de vídeo. Este projeto usa uma estratégia que:

1. Primeiro tenta reproduzir com a melhor qualidade disponível sem PO Token
2. Se falhar, usa formato progressivo 360p (sempre disponível)

A qualidade pode variar dependendo das restrições do YouTube no momento.

## Configuração opcional

Você pode definir uma instância Invidious alternativa:

```bash
export INVIDIOUS_BASE="https://invidious.exemplo.com"
./youtui
```

Por padrão usa: `https://yewtu.be`

## Licença

MIT
