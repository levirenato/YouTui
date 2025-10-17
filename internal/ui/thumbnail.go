package ui

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dolmen-go/kittyimg"
	"github.com/nfnt/resize"
)

// ThumbnailCache gerencia download e cache de thumbnails
type ThumbnailCache struct {
	cacheDir string
	mu       sync.RWMutex
	cache    map[string]string // URL -> escape code
}

// NewThumbnailCache cria um novo cache de thumbnails
func NewThumbnailCache() (*ThumbnailCache, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "youtui", "thumbnails")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	return &ThumbnailCache{
		cacheDir: cacheDir,
		cache:    make(map[string]string),
	}, nil
}

// hashURL cria um hash MD5 da URL para usar como nome de arquivo
func (tc *ThumbnailCache) hashURL(url string) string {
	h := md5.Sum([]byte(url))
	return fmt.Sprintf("%x", h)
}

// getCachePath retorna o caminho do arquivo em cache
func (tc *ThumbnailCache) getCachePath(url string) string {
	hash := tc.hashURL(url)
	return filepath.Join(tc.cacheDir, hash+".jpg")
}

// downloadImage baixa a imagem da URL (versão sem contexto)
func (tc *ThumbnailCache) downloadImage(url string) (image.Image, error) {
	return tc.downloadImageWithContext(context.Background(), url)
}

// downloadImageWithContext baixa a imagem da URL com suporte a cancelamento
func (tc *ThumbnailCache) downloadImageWithContext(ctx context.Context, url string) (image.Image, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	return img, nil
}

// GetThumbnailImage retorna a imagem diretamente (sem escape codes)
// Útil para tview.Image
func (tc *ThumbnailCache) GetThumbnailImage(url string) (image.Image, error) {
	return tc.GetThumbnailImageWithContext(context.Background(), url)
}

// GetThumbnailImageWithContext retorna a imagem com suporte a cancelamento via contexto
func (tc *ThumbnailCache) GetThumbnailImageWithContext(ctx context.Context, url string) (image.Image, error) {
	if url == "" {
		return nil, fmt.Errorf("empty URL")
	}

	cachePath := tc.getCachePath(url)

	// Verifica se já existe em disco
	if _, statErr := os.Stat(cachePath); statErr == nil {
		// Carrega do cache (rápido, não precisa de contexto)
		f, err := os.Open(cachePath)
		if err == nil {
			defer f.Close()
			img, _, err := image.Decode(f)
			if err == nil {
				return img, nil
			}
		}
	}

	// Se não está em cache ou falhou ao carregar, baixa com contexto
	img, err := tc.downloadImageWithContext(ctx, url)
	if err != nil {
		return nil, err
	}

	// CRÍTICO: Salva no cache para não baixar novamente
	if err := tc.saveImageToCache(img, cachePath); err != nil {
		// Log erro mas retorna a imagem mesmo assim
		// (falha ao salvar não deve impedir uso)
	}

	return img, nil
}

// saveImageToCache salva a imagem no cache
func (tc *ThumbnailCache) saveImageToCache(img image.Image, cachePath string) error {
	f, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Redimensiona para tamanho pequeno (para uso em lista)
	resized := resize.Thumbnail(120, 90, img, resize.Lanczos3)

	// Salva como JPEG no cache (kittyimg vai usar depois)
	var buf bytes.Buffer
	if err := kittyimg.Fprint(&buf, resized); err != nil {
		return err
	}

	_, err = io.Copy(f, &buf)
	return err
}

// GetThumbnail obtém o thumbnail (baixa se necessário) e retorna o escape code do Kitty
func (tc *ThumbnailCache) GetThumbnail(url string, width, height int) (string, error) {
	if url == "" {
		return "", nil
	}

	tc.mu.RLock()
	if cached, ok := tc.cache[url]; ok {
		tc.mu.RUnlock()
		return cached, nil
	}
	tc.mu.RUnlock()

	cachePath := tc.getCachePath(url)

	// Verifica se já existe em disco
	var img image.Image
	var err error

	if _, statErr := os.Stat(cachePath); statErr == nil {
		// Carrega do cache
		f, err := os.Open(cachePath)
		if err == nil {
			defer f.Close()
			img, _, err = image.Decode(f)
		}
	}

	// Se não está em cache ou falhou ao carregar, baixa
	if img == nil {
		img, err = tc.downloadImage(url)
		if err != nil {
			return "", err
		}

		// Salva no cache
		if err := tc.saveImageToCache(img, cachePath); err != nil {
			// Não fatal, continua
			_ = err
		}
	}

	// Redimensiona para o tamanho desejado
	resized := resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)

	// Gera escape code do Kitty
	var buf bytes.Buffer
	if err := kittyimg.Fprint(&buf, resized); err != nil {
		return "", err
	}

	escapeCode := buf.String()

	// Cacheia em memória
	tc.mu.Lock()
	tc.cache[url] = escapeCode
	tc.mu.Unlock()

	return escapeCode, nil
}

// GetThumbnailIcon retorna um ícone musical se thumbnail não disponível
func (tc *ThumbnailCache) GetThumbnailIcon(url string) string {
	if url == "" {
		return " "
	}

	// Tenta obter thumbnail (assíncrono seria melhor, mas por ora síncrono)
	// Retorna ícone como fallback
	return " "
}

// IsKittyTerminal verifica se está rodando no Kitty
func IsKittyTerminal() bool {
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")

	return strings.Contains(term, "kitty") || termProgram == "kitty"
}
