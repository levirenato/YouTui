package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/levirenato/YouTui/internal/search"
)

func (a *SimpleApp) onSearchDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		query := a.searchInput.GetText()
		if query != "" {
			go a.doSearch(query)
		}
	}
	a.AutoSaveState()
}

func (a *SimpleApp) doSearch(query string) {
	a.app.QueueUpdateDraw(func() {
		a.setStatus(a.theme.Yellow, "  "+a.strings.Searching)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := search.SearchVideos(ctx, query, 30)
	if err != nil {
		a.app.QueueUpdateDraw(func() {
			a.setStatusf(a.theme.Red, "❌ "+a.strings.SearchError, err)
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

	a.pagination.SetTotalItems(len(tracksCopy))
	a.pagination.Reset()

	a.displayCurrentPage()

	a.app.QueueUpdateDraw(func() {
		a.setStatusf(a.theme.Green, "✓ "+a.strings.FoundResults,
			len(tracksCopy), 1, a.pagination.GetTotalPages())
		a.app.SetFocus(a.searchResults.Flex)
		a.updateCommandBar()
	})

	a.AutoSaveState()
}

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
			a.searchResults.AddItem(track, i)

			if track.Thumbnail != "" && a.thumbCache != nil {
				go func(idx int, url string) {
					img, err := a.thumbCache.GetThumbnailImage(url)
					if err == nil && img != nil {
						a.app.QueueUpdateDraw(func() {
							a.searchResults.SetThumbnail(idx, img)
						})
					}
				}(i, track.Thumbnail)
			}
		}

		currentPage := a.pagination.GetCurrentPage() + 1
		totalPages := a.pagination.GetTotalPages()
		a.searchResults.SetTitle(fmt.Sprintf(" %s [%s %d/%d] ", a.strings.Results, a.strings.Page, currentPage, totalPages))
	})
}

func (a *SimpleApp) nextPage() {
	if a.pagination.NextPage() {
		a.displayCurrentPage()
		currentPage := a.pagination.GetCurrentPage() + 1
		totalPages := a.pagination.GetTotalPages()
		a.app.QueueUpdateDraw(func() {
			a.setStatusf(a.theme.Sapphire, "→ "+a.strings.NextPage, currentPage, totalPages)
		})
	} else {
		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Yellow, "⚠ "+a.strings.AlreadyLastPage)
		})
	}
}

func (a *SimpleApp) prevPage() {
	if a.pagination.PrevPage() {
		a.displayCurrentPage()
		currentPage := a.pagination.GetCurrentPage() + 1
		totalPages := a.pagination.GetTotalPages()
		a.app.QueueUpdateDraw(func() {
			a.setStatusf(a.theme.Sapphire, "← "+a.strings.PrevPage, currentPage, totalPages)
		})
	} else {
		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Yellow, "⚠ "+a.strings.AlreadyFirstPage)
		})
	}
}
