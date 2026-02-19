.PHONY: help install install-arch install-ubuntu install-macos build run clean test install-bin uninstall deps fmt check-deps version

BINARY_NAME = youtui-player
GO          = go
PREFIX      = /usr/local
DESTDIR     =
BINDIR      = $(DESTDIR)$(PREFIX)/bin
MANDIR      = $(DESTDIR)$(PREFIX)/share/man/man1
DATADIR     = $(DESTDIR)$(PREFIX)/share

VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS     = -ldflags "-X main.Version=$(VERSION) -s -w"

UNAME_S := $(shell uname -s)

help:
	@echo "YouTui $(VERSION) - Makefile"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

install: ## Instala dependencias e compila (detecta SO automaticamente)
	@echo "Detectando sistema operacional..."
ifeq ($(UNAME_S),Linux)
	@if [ -f /etc/arch-release ]; then \
		echo "Arch Linux detectado"; \
		$(MAKE) install-arch; \
	elif [ -f /etc/debian_version ]; then \
		echo "Debian/Ubuntu detectado"; \
		$(MAKE) install-ubuntu; \
	else \
		echo "Distribuicao Linux nao suportada automaticamente"; \
		echo "Instale manualmente: mpv, yt-dlp, socat"; \
		exit 1; \
	fi
else ifeq ($(UNAME_S),Darwin)
	@echo "macOS detectado"
	@$(MAKE) install-macos
else
	@echo "Sistema operacional nao suportado: $(UNAME_S)"
	@exit 1
endif
	@echo ""
	@$(MAKE) build

install-arch: ## Instala dependencias no Arch Linux
	@echo "Instalando dependencias no Arch Linux..."
	sudo pacman -S --needed mpv yt-dlp socat go
	@echo "Dependencias instaladas!"

install-ubuntu: ## Instala dependencias no Ubuntu/Debian
	@echo "Instalando dependencias no Ubuntu/Debian..."
	sudo apt update
	sudo apt install -y mpv socat python3-pip golang
	sudo pip3 install -U yt-dlp || pip3 install --user -U yt-dlp
	@echo "Dependencias instaladas!"

install-macos: ## Instala dependencias no macOS
	@echo "Instalando dependencias no macOS..."
	@which brew > /dev/null || (echo "Homebrew nao encontrado. Instale em: https://brew.sh" && exit 1)
	brew install mpv yt-dlp socat go
	@echo "Dependencias instaladas!"

deps: ## Baixa dependencias do Go
	$(GO) mod download
	$(GO) mod tidy

build: deps ## Compila o projeto
	@echo "Compilando $(BINARY_NAME) $(VERSION)..."
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Compilado: ./$(BINARY_NAME)"

run: build ## Compila e executa
	./$(BINARY_NAME)

install-bin: build ## Instala o binario (respeita DESTDIR e PREFIX)
	@echo "Instalando $(BINARY_NAME) em $(BINDIR)..."
	install -Dm755 $(BINARY_NAME) $(BINDIR)/$(BINARY_NAME)
	@echo "Instalado!"

uninstall: ## Remove o binario instalado
	rm -f $(BINDIR)/$(BINARY_NAME)
	@echo "Desinstalado."

clean: ## Remove arquivos compilados
	rm -f $(BINARY_NAME)
	$(GO) clean

test: ## Executa testes
	$(GO) test ./...

fmt: ## Formata o codigo
	$(GO) fmt ./...

vet: ## Analisa o codigo
	$(GO) vet ./...

check-deps: ## Verifica dependencias de runtime
	@which mpv    > /dev/null && echo "OK mpv"    || echo "FALTANDO mpv"
	@which yt-dlp > /dev/null && echo "OK yt-dlp" || echo "FALTANDO yt-dlp"
	@which socat  > /dev/null && echo "OK socat"  || echo "FALTANDO socat"

version: ## Mostra a versao atual
	@echo $(VERSION)
