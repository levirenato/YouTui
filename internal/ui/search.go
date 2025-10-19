package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/levirenato/YouTui/internal/search"
	"github.com/levirenato/YouTui/internal/ui/components"
)

// onSearchDone é chamado quando Enter é pressionado na busca
func (a *SimpleApp) onSearchDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		query := a.searchInput.GetText()
		if query != "" {
			go a.doSearch(query)
		}
	}
}

// doSearch executa a busca no YouTube
func (a *SimpleApp) doSearch(query string) {
	a.app.QueueUpdateDraw(func() {
		a.statusBar.SetText("[yellow]  Buscando...")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := search.SearchVideos(ctx, query, 30)
	if err != nil {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[red]❌ Erro: %v", err))
		})
		return
	}

	a.mu.Lock()
	a.tracks = make([]Track, len(results))
	for i, r := range results {
		a.tracks[i] = Track{
			Title:       r.Title,
			Author:      r.Author,
			URL:         r.URL,
			Thumbnail:   r.Thumbnail,
			Duration:    r.Duration,
			PublishedAt: r.PublishedAt,
			Description: r.Description,
		}
	}
	tracksCopy := make([]Track, len(a.tracks))
	copy(tracksCopy, a.tracks)
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		// CRÍTICO: Desabilita SetChangedFunc temporariamente para evitar
		// disparar 30 chamadas simultâneas ao yt-dlp ao adicionar itens
		a.searchResults.SetChangedFunc(nil)

		a.searchResults.Clear()
		for i, track := range tracksCopy {
			icon := components.GetTrackIcon(i)
			title := fmt.Sprintf("%s %s - %s", icon, track.Title, track.Author)
			a.searchResults.AddItem(title, "", 0, nil)
		}
		a.searchResults.SetTitle(fmt.Sprintf(" Resultados [%d] ", len(tracksCopy)))
		a.statusBar.SetText(fmt.Sprintf("[green]✓ Encontrados %d resultados", len(tracksCopy)))

		// Reabilita o handler DEPOIS de adicionar todos os itens
		a.searchResults.SetChangedFunc(a.onResultChanged)

		a.app.SetFocus(a.searchResults)
		a.updateCommandBar()

		// Carrega detalhes do primeiro item de forma assíncrona
		// para não bloquear a UI thread
		if len(tracksCopy) > 0 {
			go a.updateSearchDetails(0)
		}
	})
}
