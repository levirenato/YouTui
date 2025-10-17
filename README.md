# YouTui 🎵

Um player de YouTube para terminal (TUI) **moderno e bonito** com tema **Catppuccin Mocha**, construído com Go e **tview**.

> **Destaques**: Barra de progresso em tempo real • Atalhos contextuais • Controle completo de playlist • Interface colorida • Thread-safe • Sem dependência de APIs

![Status](https://img.shields.io/badge/status-stable-green)
![Go Version](https://img.shields.io/badge/go-1.24+-blue)
![License](https://img.shields.io/badge/license-MIT-blue)

## ✨ Recursos

### 🎨 Interface Visual
-  **Tema Catppuccin Mocha**: Cores harmoniosas e modernas
-  **Thumbnails de Alta Qualidade**: TrueColor + Floyd-Steinberg dithering para pixels menores e mais detalhes
-  **Painel de Detalhes**: Mostra thumbnail, título, canal e duração do vídeo selecionado instantaneamente
-  **Barra de Progresso Dinâmica**: Atualização em tempo real com tempo atual/total e largura responsiva
-  **Bordas Coloridas**: Painel ativo destacado em azul
-  **Ícones Musicais Unicode**: Símbolos musicais (♪ ♫ ♬ ♩ ▸ •) nas listas
-  **Atalhos Contextuais**: Barra inferior com atalhos específicos para cada painel
-  **Ajuda Integrada**: Pressione `?` para ver todos os atalhos

### 🎵 Funcionalidades de Player
-  **Busca Rápida**: Resultados aparecem em 2-5 segundos via yt-dlp, sem depender de APIs externas
-  **Playlist Completa**: Modos Normal, Repetir 1, Repetir Todas, Shuffle
-  **Controles Completos**: Play/Pause/Stop/Next/Previous
-  **Dois Modos de Reprodução**:
   - **Direto**: Toque músicas dos resultados instantaneamente
   - **Playlist**: Controle completo com navegação n/b
-  **Modo Áudio/Vídeo**: Alterne entre reprodução de áudio ou vídeo
-  **Reordenação**: Mova músicas na playlist com J/K

### 🔧 Técnico
-  **Thread-Safe**: Sincronização adequada com Mutex para operações concorrentes
-  **IPC com mpv**: Controle via socket Unix para pause/progresso em tempo real
-  **Sistema de Temas Desacoplado**: Preparado para múltiplos temas no futuro

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

## Instalação

```bash
go build -o youtui .
./youtui
```

## 📐 Layout

```
┌──────────────────────┬─────────────────────┐
│  🔍 Busca           │  📋 Playlist       │
│  [Input]            │  ♪ Música A        │
│                     │  ♫ Música B        │
│  📋 Resultados      │  ♬ Música C        │
│  ♪ Resultado 1      │                    │
│  ♫ Resultado 2      │                    │
│┌─────┬──────────────┐│                    │
││🖼️  │Título        ││                    │
││img │Canal: Nome   ││                    │
││    │Duração: 4:30 ││                    │
││    │Data: 12/10   ││                    │
│└─────┴──────────────┘│                    │
├─────┬────────────────┴─────────────────────┤
│ 🖼️ │  🎵 Player                          │
│img │  ▶ Now Playing - Artist             │
│    │  ████████████░░░░ 02:45/04:30       │
│    │  Áudio | Normal                      │
├─────┴──────────────────────────────────────┤
│  Status: Tocando música...                 │
├────────────────────────────────────────────┤
│  ↑/↓ Navegar | a Add | c Pause | ? Ajuda │
└────────────────────────────────────────────┘
```

## ⌨️ Controles

### 🧭 Navegação
- **Tab** - Alternar entre painéis (Busca → Resultados → Playlist)
- **/** - Focar na busca a qualquer momento
- **↑/↓** - Navegar nas listas
- **?** - Abrir ajuda com todos os atalhos

### 🔍 Busca
- Digite normalmente na caixa de busca
- **Enter** - Executar busca

### 📋 Resultados
- **↑/↓** - Navegar pelos resultados
- **Enter** - Tocar faixa diretamente (modo preview)
- **a** - Adicionar à playlist

### 📑 Playlist
- **↑/↓** - Navegar na playlist
- **Enter** ou **Space** - Tocar faixa selecionada (com controles n/b)
- **d** - Remover item da playlist
- **J** - Mover item para baixo
- **K** - Mover item para cima

### 🎮 Player (Controles Globais)
- **c** ou **Space** - Pause/Play (funciona sempre)
- **s** - Stop completo
- **n** - Próxima faixa *(só quando tocando da playlist)*
- **b** - Faixa anterior *(só quando tocando da playlist)*
- **r** - Ciclar modo repetição:
  - Normal → Repetir 1 → Repetir Todas → Normal
- **h** - Toggle Shuffle (embaralhar)
- **m** - Alternar modo áudio/vídeo

### 🚪 Geral
- **q** - Sair

> **💡 Dica**: Os atalhos **n** e **b** só funcionam quando você toca uma música **da playlist** (não dos resultados). Para usar esses controles, adicione músicas com **a** e toque da playlist com **Enter** ou **Space**.

## 🚀 Workflow Recomendado

### Para ouvir músicas avulsas (Preview):
1. Digite uma busca e pressione **Enter**
2. Navegue nos resultados com **↑/↓**
3. Pressione **Enter** para tocar imediatamente
4. Use **c** para pausar/retomar
5. Pressione **/** para nova busca

### Para criar uma playlist:
1. Faça buscas e pressione **a** em cada resultado para adicionar
2. Vá para a playlist com **Tab**
3. Reordene se necessário com **J/K**
4. Pressione **Enter** ou **Space** para iniciar
5. Use **n/b** para navegar entre faixas
6. Configure modos com **r** (repetição) e **h** (shuffle)
7. Pressione **?** para ver todos os atalhos disponíveis

## Características Técnicas

- **Framework**: tview (robusto e estável)
- **Busca**: yt-dlp NDJSON streaming (sem dependência de APIs)
- **Player**: mpv com IPC socket para controle completo
- **Concorrência**: Mutex para thread-safety
- **Auto-avanço**: Playlist contínua com modos de repetição
- **Layout Flex**: Responsivo e adaptável
- **Tema**: Catppuccin Mocha (sistema de temas desacoplado para futuras expansões)
- **Barra de Progresso**: Atualização em tempo real via IPC
- **Ajuda Integrada**: Modal com todos os atalhos (pressione `?`)

## 🎨 Sistema de Temas

YouTui usa o tema **Catppuccin Mocha** por padrão, com cores cuidadosamente selecionadas para uma experiência visual agradável e moderna:

### Paleta de Cores
- **Borda Ativa**: Azul Catppuccin (#89b4fa) - indica o painel focado em tempo real
- **Borda Inativa**: Surface0 (#313244) - painéis em segundo plano
- **Player**: Roxo Mauve (#cba6f7) - destaque especial para o player
- **Background**: Base escuro (#1e1e2e) - fundo confortável para os olhos
- **Texto**: Text (#cdd6f4) - claro e legível
- **Seleção**: Azul sobre preto - itens selecionados nas listas

### Atalhos Coloridos (Barra Inferior)
Cada tipo de ação tem sua cor específica para facilitar identificação:
- **Navegação** (↑/↓, Tab, Enter): Azul (#89b4fa)
- **Adicionar** (a, c): Verde (#a6e3a1)
- **Remover** (d, q): Vermelho (#f38ba8)
- **Mover** (J/K): Roxo (#cba6f7)
- **Repeat** (r): Laranja (#fab387)
- **Shuffle** (h): Teal (#94e2d5)
- **Next/Prev** (n/b): Sky (#89dceb)
- **Ajuda** (?): Amarelo (#f9e2af)

### Arquitetura
O sistema de temas está desacoplado em `internal/ui/theme.go`, preparado para futura implementação de:
- Múltiplos temas (Gruvbox, Nord, Dracula, etc.)
- Seleção via arquivo de configuração TOML (`~/.config/youtui/themes.toml`)
- Temas personalizados pelo usuário

Um exemplo de configuração está disponível em `themes.toml.example`.

## 🌟 Destaques de Implementação

### Barra de Progresso Dinâmica
A barra de progresso atualiza a cada 500ms consultando o mpv via IPC socket:
```
████████████████████░░░░░░░░░░░░░░░░░░░░ 02:45/04:30
```
- Mostra tempo atual e duração total
- Atualização visual em tempo real
- Funciona mesmo quando em pausa

### Atalhos Contextuais Inteligentes
A barra inferior muda automaticamente conforme o painel ativo:
- **Na Busca**: Mostra atalhos de busca e navegação
- **Nos Resultados**: Mostra como adicionar à playlist
- **Na Playlist**: Mostra controles completos (mover, remover, modos)

### Dois Modos de Reprodução
1. **Modo Direto** (dos Resultados): Preview rápido sem playlist
2. **Modo Playlist**: Controle completo com n/b, reordenação e modos

### Sistema Robusto de Concorrência
- Mutex protegendo todas as operações críticas
- Flag `skipAutoPlay` para evitar race conditions entre pulo manual e auto-play
- Gerenciamento seguro de goroutines do mpv

## ⚠️ Notas sobre reprodução de vídeos (2025)

O YouTube começou a exigir **PO Tokens** para alguns formatos de vídeo. Este projeto usa uma estratégia que:

1. Primeiro tenta reproduzir com a melhor qualidade disponível sem PO Token
2. Se falhar, usa formato progressivo 360p (sempre disponível)

A qualidade pode variar dependendo das restrições do YouTube no momento.

## 🖼️ Thumbnails

YouTui exibe **thumbnails reais em alta qualidade** dos vídeos do YouTube usando o widget `tview.Image`.

### Como Funciona
- ✅ **Dois painéis com thumbnails**:
  - **Player**: Thumbnail da música tocando (20 caracteres de largura)
  - **Detalhes**: Thumbnail do vídeo selecionado nos resultados (20 caracteres)
- ✅ **Alta Qualidade**: TrueColor (16 milhões de cores) + Floyd-Steinberg dithering
- ✅ **Download automático** das capas do YouTube (hqdefault.jpg)
- ✅ **Cache em disco** (`~/.cache/youtui/thumbnails/`)
- ✅ **Atualização em tempo real** quando você muda de música ou seleção
- ✅ **Funciona em qualquer terminal** (não requer Kitty Graphics Protocol)

### Características Técnicas
- **TrueColor**: Renderização com 16 milhões de cores para máxima fidelidade
- **Floyd-Steinberg Dithering**: Algoritmo de difusão de erro para suavizar gradientes
- **Resultado**: Pixels menores e imagem mais definida vs. 256 cores padrão
- **Download assíncrono**: Não trava a UI durante o carregamento
- **Cache inteligente**: Só baixa uma vez, reutiliza em próximas execuções

### Painel de Detalhes
Ao navegar pelos resultados de busca, o painel inferior mostra:
- **Thumbnail** do vídeo (à esquerda)
- **Título** em amarelo/negrito
- **Canal** do autor
- **Duração** do vídeo

As informações são exibidas **instantaneamente** quando você navega pelos resultados, sem necessidade de esperar carregamentos adicionais.

## 🐛 Solução de Problemas

### Pause não funciona
Certifique-se de que:
- `socat` está instalado
- A música está realmente tocando (veja o ícone ▶)
- O socket IPC do mpv foi criado corretamente

### n/b não funcionam
Esses atalhos **só funcionam quando tocando da playlist**:
1. Adicione músicas à playlist com **a**
2. Navegue até a playlist com **Tab**
3. Pressione **Enter** ou **Space** para iniciar
4. Agora **n/b** funcionarão

### Músicas pulam incorretamente
Se as músicas pularem para a última e finalizarem:
- Recompile o projeto: `go build -o youtui .`
- O bug de race condition foi corrigido na versão atual

### Ícones aparecem como quadrados
Se os ícones musicais (♪ ♫ ♬) aparecem como `□`:
- Sua fonte não suporta caracteres Unicode musicais
- Instale uma fonte que suporte Unicode completo
- Recomendado: JetBrains Mono, Fira Code, ou qualquer Nerd Font

## Configuração opcional

Você pode definir uma instância Invidious alternativa:

```bash
export INVIDIOUS_BASE="https://invidious.exemplo.com"
./youtui
```

Por padrão usa: `https://yewtu.be`

## 🗺️ Roadmap

### Futuras Implementações
- [ ] **Thumbnails assíncronos**: Download em background sem travar a UI
- [ ] **Thumbnail no player**: Exibir capa do álbum/vídeo na área do player
- [ ] Seleção de temas via arquivo TOML
- [ ] Temas adicionais (Gruvbox, Nord, Dracula, Tokyo Night)
- [ ] Histórico de músicas tocadas
- [ ] Salvar/carregar playlists
- [ ] Filtro de busca nos resultados
- [ ] Visualizador de letras (lyrics)
- [ ] Equalizer visual ASCII
- [ ] Suporte a múltiplas playlists
- [ ] Download de músicas
- [ ] Cache de resultados de busca

### Melhorias Técnicas
- [ ] Testes unitários
- [ ] CI/CD pipeline
- [ ] Binários pré-compilados para releases
- [ ] Documentação de API interna
- [ ] Profiles de performance

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se livre para:
- Reportar bugs via Issues
- Sugerir novas features
- Enviar Pull Requests
- Melhorar a documentação

### Estrutura do Projeto
```
YouTui/
├── cmd/              # Ponto de entrada da aplicação
├── internal/
│   ├── ui/          # Interface TUI (tview)
│   │   ├── simple.go   # UI principal
│   │   └── theme.go    # Sistema de temas
│   └── search/      # Integração com yt-dlp
├── go.mod
└── README.md
```

## 📝 Licença

MIT License - sinta-se livre para usar, modificar e distribuir.

## 🙏 Agradecimentos

- [tview](https://github.com/rivo/tview) - Framework TUI excepcional
- [tcell](https://github.com/gdamore/tcell) - Terminal handling robusto
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - Extrator poderoso do YouTube
- [mpv](https://mpv.io/) - Player de mídia versátil
- [Catppuccin](https://github.com/catppuccin/catppuccin) - Tema lindo e acessível

---

**Desenvolvido com ❤️ e Go**
