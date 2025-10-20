package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/levirenato/YouTui/internal/search"
)

// onSearchDone is called when Enter is pressed in search
func (a *SimpleApp) onSearchDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		query := a.searchInput.GetText()
		if query != "" {
			go a.doSearch(query)
		}
	}
}

// doSearch performs a YouTube search
func (a *SimpleApp) doSearch(query string) {
	a.app.QueueUpdateDraw(func() {
		a.statusBar.SetText("[yellow]  " + a.strings.Searching)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := search.SearchVideos(ctx, query, 30)
	if err != nil {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[red]❌ " + a.strings.SearchError, err))
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

	// Setup pagination and display first page
	a.pagination.SetTotalItems(len(tracksCopy))
	a.pagination.Reset()

	a.displayCurrentPage()
	
	a.app.QueueUpdateDraw(func() {
		a.statusBar.SetText(fmt.Sprintf("[green]✓ " + a.strings.FoundResults, 
			len(tracksCopy), 1, a.pagination.GetTotalPages()))
		a.app.SetFocus(a.searchResults.Flex)
		a.updateCommandBar()
	})
}

// displayCurrentPage displays current page items
func (a *SimpleApp) displayCurrentPage() {
	a.mu.Lock()
	start, end := a.pagination.GetPageItems()
	
	var pageItems []Track
	if start < len(a.tracks) {
		if end > len(a.tracks) {
			end = len(a.tracks)
		}
		pageItems = a.tracks[start:end]
	}
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.searchResults.Clear()
		
		for i, track := range pageItems {
			// Add item with thumbnail inline
			a.searchResults.AddItem(track, i)
			
			// Load thumbnail in background (uses cache)
			if track.Thumbnail != "" && a.thumbCache != nil {
				go func(idx int, url string) {
					img, err := a.thumbCache.GetThumbnailImage(url)
					if err == nil && img != nil {
						a.searchResults.SetThumbnail(idx, img)
					}
				}(i, track.Thumbnail)
			}
		}
		
		currentPage := a.pagination.GetCurrentPage() + 1
		totalPages := a.pagination.GetTotalPages()
		a.searchResults.SetTitle(fmt.Sprintf(" Resultados [Página %d/%d] ", currentPage, totalPages))
	})
}

// nextPage advances to the next page
func (a *SimpleApp) nextPage() {
	if a.pagination.NextPage() {
		a.displayCurrentPage()
		currentPage := a.pagination.GetCurrentPage() + 1
		totalPages := a.pagination.GetTotalPages()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[cyan]→ " + a.strings.NextPage, currentPage, totalPages))
		})
	} else {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ " + a.strings.AlreadyLastPage)
		})
	}
}

// prevPage goes to the previous page
func (a *SimpleApp) prevPage() {
	if a.pagination.PrevPage() {
		a.displayCurrentPage()
		currentPage := a.pagination.GetCurrentPage() + 1
		totalPages := a.pagination.GetTotalPages()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[cyan]← " + a.strings.PrevPage, currentPage, totalPages))
		})
	} else {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ " + a.strings.AlreadyFirstPage)
		})
	}
}
