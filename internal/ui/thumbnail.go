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

type ThumbnailCache struct {
	cacheDir string
	mu       sync.RWMutex
}

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

func (tc *ThumbnailCache) hashURL(url string) string {
	h := md5.Sum([]byte(url))
	return fmt.Sprintf("%x", h)
}

func (tc *ThumbnailCache) getCachePath(url string) string {
	hash := tc.hashURL(url)
	return filepath.Join(tc.cacheDir, hash+".jpg")
}

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

func (tc *ThumbnailCache) GetThumbnailImage(url string) (image.Image, error) {
	return tc.GetThumbnailImageWithContext(context.Background(), url)
}

func (tc *ThumbnailCache) GetThumbnailImageWithContext(ctx context.Context, url string) (image.Image, error) {
	if url == "" {
		return nil, fmt.Errorf("empty URL")
	}

	cachePath := tc.getCachePath(url)

	if _, statErr := os.Stat(cachePath); statErr == nil {
		f, err := os.Open(cachePath)
		if err == nil {
			defer f.Close()
			img, _, err := image.Decode(f)
			if err == nil {
				return img, nil
			}
		}
	}

	img, err := tc.downloadImageWithContext(ctx, url)
	if err != nil {
		return nil, err
	}

	if err := tc.saveImageToCache(img, cachePath); err != nil {
	}

	return img, nil
}

func (tc *ThumbnailCache) saveImageToCache(img image.Image, cachePath string) error {
	f, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer f.Close()

	resized := resize.Thumbnail(100, 100, img, resize.Lanczos3)

	return jpeg.Encode(f, resized, &jpeg.Options{Quality: 90})
}
