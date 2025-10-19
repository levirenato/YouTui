# YouTui 🎵

Player de YouTube para terminal com interface TUI moderna.

![Go Version](https://img.shields.io/badge/go-1.24+-blue)
![License](https://img.shields.io/badge/license-MIT-blue)

## O que faz?

YouTui é um player de YouTube que roda inteiramente no terminal, permitindo buscar, tocar músicas/vídeos e gerenciar playlists sem sair da linha de comando. Interface bonita com thumbnails inline, controles completos e tema Catppuccin Mocha.

**Recursos principais:**

- Busca rápida no YouTube (sem API keys)
- Thumbnails em alta qualidade no terminal
- Playlist com shuffle, repeat e navegação
- Controles completos (play, pause, next, previous)
- Barra de progresso em tempo real
- Interface colorida e moderna

## Dependências

- **Go 1.24+** - Linguagem de programação
- **mpv** - Player de mídia
- **yt-dlp** - Extrator de vídeos do YouTube
- **socat** - Comunicação IPC com mpv
- **Nerd Font** (opcional) - Para ícones bonitos

## Instalação Rápida

```bash
# Clone o repositório
git clone https://github.com/levirenato/YouTui
cd YouTui

# Instale dependências e compile (requer sudo)
make install

# Ou apenas compile (se já tem as dependências)
make build

# Execute
./youtui
```

## Atalhos Principais

| Tecla     | Ação                 |
| --------- | -------------------- |
| `/`       | Buscar               |
| `Enter`   | Tocar/Buscar         |
| `a`       | Adicionar à playlist |
| `d`       | Remover da playlist  |
| `Space`   | Pausar/Retomar       |
| `n` / `b` | Próxima/Anterior     |
| `h`       | Shuffle              |
| `r`       | Modo repetição       |
| `Tab`     | Alternar painéis     |
| `?`       | Ajuda completa       |
| `q`       | Sair                 |

## Desenvolvimento

```bash
# Verificar dependências
make check-deps

# Compilar
make build

# Compilar e executar
make run

# Formatar código
make fmt

# Limpar arquivos gerados
make clean
```

## Licença

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
