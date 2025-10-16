# ♫ YouTui Music Player

Player de YouTube/música para terminal (TUI) com interface moderna inspirada em Spotify.

## 🎨 Recursos

- **🎵 Player Central**: Visualização grande da música tocando
- **🔍 Busca Inteligente**: Interface não bloqueia durante busca
- **📋 Playlist Avançada**: Modos Normal, Aleatório, Repetir 1, Repetir Todas
- **🎨 UI Moderna**: Cores vibrantes, bordas animadas, componentes visuais
- **🎵 Visualizador de Áudio**: Barras animadas mostrando áudio em tempo real  
- **🎬 Background Player**: MPV roda em background sem travar o TUI
- **🛑 Controle Total**: Pause, play, stop, skip a qualquer momento
- **📝 Sistema de Logs**: Debug completo com níveis (Info/Warning/Error)
- **⚠️ Notificações**: Feedback visual colorido de todas as ações

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

**Navegação:**
- **Tab** - Alternar entre painéis (Busca → Playlist → Resultados → Visualizador → Logs*)
- **↑/↓** ou **j/k** - Navegar pelos itens (resultados ou playlist)
- **Enter** - Buscar (no painel de busca) ou Reproduzir (nos resultados/playlist)

**Reprodução:**
- **m** - Alternar modo de reprodução (Vídeo MP4 / Áudio MP3)
- **s** - Parar reprodução atual (funcionando!)
- **Space** - Iniciar reprodução da playlist completa

**Playlist:**
- **p** - Adicionar item selecionado à playlist
- **r** - Alternar modo de playlist (Normal → Aleatório → Repetir 1 → Repetir Todas)
- **d** ou **x** - Remover item da playlist (quando no painel de playlist)
- **Shift+J** - Mover item para baixo na playlist
- **Shift+K** - Mover item para cima na playlist

**Debug:**
- **l** - Alternar visualização do painel de logs

**Geral:**
- **q** ou **Ctrl+C** - Sair (mata todos os processos automaticamente)

*O painel de logs só aparece quando ativado com 'l'

### Como Usar

1. **Buscar vídeos**: Digite no campo de busca e pressione Enter
2. **Adicionar à playlist**: Navegue pelos resultados com ↑/↓ e pressione 'p'
3. **Reproduzir playlist**: Pressione Space para iniciar a reprodução automática
4. **Mudar modo de playlist**: Pressione 'r' para alternar entre modos
5. **Ver logs**: Pressione 'l' para abrir/fechar o painel de logs
6. **Parar música**: Pressione 's' a qualquer momento

### Funcionalidades

- 🎨 **Layout em Grid** - Interface dividida em 4 seções (ou 5 com logs)
- 📋 **Playlist Completa** - Adicione, remova, reordene e reproduza automaticamente
- 🔀 **Modos de Playlist** - Normal, Aleatório, Repetir Uma, Repetir Todas
- 🎵 **Modo Áudio/Vídeo** - Alterne entre reproduzir vídeo completo ou apenas áudio
- 🎯 **Navegação por Painéis** - Use Tab para focar em diferentes seções
- 🎬 **Reprodução em Background** - O TUI não trava durante a reprodução (mpv em background)
- 🎵 **Visualizador de Áudio** - Veja informações e visualização da música tocando
- 🛑 **Controle de Reprodução** - Pare a música a qualquer momento (tecla 's' funcionando!)
- 📝 **Painel de Logs** - Veja todos os eventos e erros da aplicação em tempo real
- ⚠️ **Notificações Visuais** - Indicadores coloridos de erros, warnings e info
- 🔄 **Loading Assíncrono** - A UI não trava durante buscas
- 🧹 **Cleanup Automático** - Todos os processos (mpv e cava) são finalizados ao sair
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
