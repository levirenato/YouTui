package ui

import (
	"context"
	"fmt"
	"image"
	"time"
)

// updateSearchDetailsDebounced adiciona debounce para evitar múltiplas chamadas
// ao navegar rapidamente pelos resultados
func (a *SimpleApp) updateSearchDetailsDebounced(idx int) {
	a.detailsLoadingMutex.Lock()
	
	// Cancela timer anterior se existir
	if a.detailsDebounceTimer != nil {
		a.detailsDebounceTimer.Stop()
	}
	
	// Cria novo timer de 150ms
	a.detailsDebounceTimer = time.AfterFunc(150*time.Millisecond, func() {
		a.updateSearchDetails(idx)
	})
	
	a.detailsLoadingMutex.Unlock()
}

// updateSearchDetails atualiza o painel de detalhes com informações do item selecionado
func (a *SimpleApp) updateSearchDetails(idx int) {
	// Cancela download anterior se existir
	a.detailsLoadingMutex.Lock()
	if a.detailsCancelFunc != nil {
		a.detailsCancelFunc()
		a.detailsCancelFunc = nil
	}
	a.detailsLoadingIdx = idx
	a.detailsLoadingMutex.Unlock()

	a.mu.Lock()
	if idx < 0 || idx >= len(a.tracks) {
		a.mu.Unlock()
		// Limpa detalhes
		a.app.QueueUpdateDraw(func() {
			a.detailsText.SetText("")
			a.detailsThumb.SetImage(nil)
		})
		return
	}

	track := a.tracks[idx]
	// Faz cópia dos dados para evitar race condition
	title := track.Title
	author := track.Author
	duration := track.Duration
	thumbnailURL := track.Thumbnail
	a.mu.Unlock()

	// Valida campos para evitar panic
	if title == "" {
		title = "Sem título"
	}
	if author == "" {
		author = "Desconhecido"
	}
	if duration == "" {
		duration = "--:--"
	}

	// Mostra informações básicas IMEDIATAMENTE
	basicDetails := fmt.Sprintf(
		"[yellow::b]%s[-:-:-]\n[cyan]Canal:[-] %s\n[green]Duração:[-] %s\n\n[gray]Pressione Enter para tocar[-]",
		title,
		author,
		duration,
	)

	a.app.QueueUpdateDraw(func() {
		a.detailsText.SetText(basicDetails)
	})

	// Atualiza thumbnail em background (não bloqueia) - COM cancelamento
	if thumbnailURL != "" && a.thumbCache != nil {
		// Cria contexto cancelável para este download
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
		a.detailsLoadingMutex.Lock()
		a.detailsCancelFunc = cancel
		a.detailsLoadingMutex.Unlock()
		
		go func(url string, ctx context.Context) {
			// Canal para resultado do download
			type result struct {
				img image.Image
				err error
			}
			resultChan := make(chan result, 1)

			// Download em goroutine separada com contexto
			go func() {
				img, err := a.thumbCache.GetThumbnailImageWithContext(ctx, url)
				select {
				case resultChan <- result{img: img, err: err}:
				case <-ctx.Done():
					// Contexto cancelado, não envia resultado
				}
			}()

			// Aguarda resultado ou cancelamento
			select {
			case res := <-resultChan:
				if res.err == nil && res.img != nil {
					// Verifica se ainda é o item correto antes de atualizar
					a.detailsLoadingMutex.Lock()
					currentIdx := a.detailsLoadingIdx
					a.detailsLoadingMutex.Unlock()
					
					if currentIdx == idx {
						a.app.QueueUpdateDraw(func() {
							a.detailsThumb.SetImage(res.img)
						})
					}
				}
			case <-ctx.Done():
				// Cancelado ou timeout - não faz nada
				return
			}
			
			// Limpa a função de cancelamento
			a.detailsLoadingMutex.Lock()
			a.detailsCancelFunc = nil
			a.detailsLoadingMutex.Unlock()
		}(thumbnailURL, ctx)
	} else {
		// Limpa thumbnail se não houver URL
		a.app.QueueUpdateDraw(func() {
			a.detailsThumb.SetImage(nil)
		})
	}
}

// updateThumbnail atualiza a thumbnail do player
func (a *SimpleApp) updateThumbnail(thumbnailURL string) {
	if thumbnailURL == "" || a.thumbCache == nil {
		// Limpa o thumbnail
		a.thumbnailView.SetImage(nil)
		return
	}

	// Baixa thumbnail em goroutine para não bloquear UI
	go func() {
		img, err := a.thumbCache.GetThumbnailImage(thumbnailURL)
		if err != nil {
			// Se falhar, apenas não exibe
			return
		}

		// Atualiza na UI thread
		a.app.QueueUpdateDraw(func() {
			a.thumbnailView.SetImage(img)
		})
	}()
}
