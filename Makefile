.PHONY: help install install-arch install-ubuntu install-macos build run clean test

# Variáveis
BINARY_NAME=youtui
GO=go
INSTALL_DIR=/usr/local/bin

# Detecta o sistema operacional
UNAME_S := $(shell uname -s)

help: ## Mostra esta mensagem de ajuda
	@echo "YouTui - Makefile"
	@echo ""
	@echo "Uso: make [target]"
	@echo ""
	@echo "Targets disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Instala dependências e compila (detecta SO automaticamente)
	@echo "🔍 Detectando sistema operacional..."
ifeq ($(UNAME_S),Linux)
	@if [ -f /etc/arch-release ]; then \
		echo "📦 Arch Linux detectado"; \
		$(MAKE) install-arch; \
	elif [ -f /etc/debian_version ]; then \
		echo "📦 Debian/Ubuntu detectado"; \
		$(MAKE) install-ubuntu; \
	else \
		echo "❌ Distribuição Linux não suportada automaticamente"; \
		echo "   Instale manualmente: mpv, yt-dlp, socat"; \
		exit 1; \
	fi
else ifeq ($(UNAME_S),Darwin)
	@echo "📦 macOS detectado"
	@$(MAKE) install-macos
else
	@echo "❌ Sistema operacional não suportado: $(UNAME_S)"
	@exit 1
endif
	@echo ""
	@$(MAKE) build

install-arch: ## Instala dependências no Arch Linux
	@echo "📦 Instalando dependências no Arch Linux..."
	sudo pacman -S --needed mpv yt-dlp socat go
	@echo "✅ Dependências instaladas!"
	@echo "💡 Dica: Instale uma Nerd Font para ícones bonitos:"
	@echo "   yay -S ttf-nerd-fonts-symbols-mono"

install-ubuntu: ## Instala dependências no Ubuntu/Debian
	@echo "📦 Instalando dependências no Ubuntu/Debian..."
	sudo apt update
	sudo apt install -y mpv socat python3-pip golang
	@echo "📦 Instalando yt-dlp..."
	sudo pip3 install -U yt-dlp || pip3 install --user -U yt-dlp
	@echo "✅ Dependências instaladas!"
	@echo "💡 Dica: Instale uma Nerd Font para ícones bonitos:"
	@echo "   https://www.nerdfonts.com/font-downloads"

install-macos: ## Instala dependências no macOS
	@echo "📦 Instalando dependências no macOS..."
	@which brew > /dev/null || (echo "❌ Homebrew não encontrado. Instale em: https://brew.sh" && exit 1)
	brew install mpv yt-dlp socat go
	@echo "✅ Dependências instaladas!"
	@echo "💡 Dica: Instale uma Nerd Font para ícones bonitos:"
	@echo "   brew tap homebrew/cask-fonts"
	@echo "   brew install --cask font-hack-nerd-font"

deps: ## Baixa dependências do Go
	@echo "📦 Baixando dependências do Go..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✅ Dependências do Go instaladas!"

build: deps ## Compila o projeto
	@echo "🔨 Compilando $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) .
	@echo "✅ Compilado com sucesso: ./$(BINARY_NAME)"

run: build ## Compila e executa
	@echo "🚀 Executando $(BINARY_NAME)..."
	./$(BINARY_NAME)

install-bin: build ## Instala o binário em /usr/local/bin
	@echo "📥 Instalando $(BINARY_NAME) em $(INSTALL_DIR)..."
	sudo cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Instalado! Execute com: $(BINARY_NAME)"

uninstall: ## Remove o binário de /usr/local/bin
	@echo "🗑️  Removendo $(BINARY_NAME)..."
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Desinstalado!"

clean: ## Remove arquivos compilados e cache
	@echo "🧹 Limpando..."
	rm -f $(BINARY_NAME)
	rm -rf /tmp/youtui-thumbnails-*
	$(GO) clean
	@echo "✅ Limpeza concluída!"

test: ## Executa testes
	@echo "🧪 Executando testes..."
	$(GO) test ./...

fmt: ## Formata o código
	@echo "✨ Formatando código..."
	$(GO) fmt ./...
	@echo "✅ Código formatado!"

check-deps: ## Verifica se as dependências estão instaladas
	@echo "🔍 Verificando dependências..."
	@which $(GO) > /dev/null && echo "✅ Go instalado" || echo "❌ Go não encontrado"
	@which mpv > /dev/null && echo "✅ mpv instalado" || echo "❌ mpv não encontrado"
	@which yt-dlp > /dev/null && echo "✅ yt-dlp instalado" || echo "❌ yt-dlp não encontrado"
	@which socat > /dev/null && echo "✅ socat instalado" || echo "❌ socat não encontrado"
	@echo ""
	@echo "Execute 'make install' para instalar dependências faltantes"

version: ## Mostra versões das dependências
	@echo "📊 Versões:"
	@$(GO) version
	@mpv --version | head -n 1
	@yt-dlp --version
	@echo "socat: $$(socat -V 2>&1 | head -n 1)"
