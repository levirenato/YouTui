package ui

import (
	"context"
	"crypto/md5"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nfnt/resize"
)

// ThumbnailCache gerencia download e cache de thumbnails
type ThumbnailCache struct {
	cacheDir string
	mu       sync.RWMutex
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

// GetThumbnailImage retorna a imagem diretamente
// Usado para renderização com tview.Image (half-blocks Unicode)
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

	// Salva no cache para não baixar novamente
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
	resized := resize.Thumbnail(100, 100, img, resize.Lanczos3)

	// Salva como JPEG no cache
	return jpeg.Encode(f, resized, &jpeg.Options{Quality: 90})
}
