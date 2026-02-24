package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/IvelOt/youtui-player/internal/search"
)

func isYouTubeURL(s string) bool {
	s = strings.TrimSpace(s)
	return strings.Contains(s, "youtube.com/watch") ||
		strings.Contains(s, "youtu.be/") ||
		strings.Contains(s, "youtube.com/shorts/") ||
		strings.Contains(s, "youtube.com/playlist") ||
		strings.Contains(s, "music.youtube.com/watch")
}

func isPlaylistURL(s string) bool {
	s = strings.TrimSpace(s)
	return strings.Contains(s, "youtube.com/playlist?list=") ||
		(isYouTubeURL(s) && strings.Contains(s, "&list="))
}

func (a *SimpleApp) onSearchDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		query := strings.TrimSpace(a.searchInput.GetText())
		if query != "" {
			if isPlaylistURL(query) {
				go a.searchPlaylistURL(query)
			} else if isYouTubeURL(query) {
				go a.searchVideoURL(query)
			} else {
				go a.doSearch(query)
			}
		}
	}
	a.AutoSaveState()
}

func (a *SimpleApp) searchVideoURL(url string) {
	a.app.QueueUpdateDraw(func() {
		a.setStatus(a.theme.Yellow, "  "+a.strings.LoadingURL)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := search.GetVideoDetails(ctx, url)
	if err != nil {
		a.app.QueueUpdateDraw(func() {
			a.setStatusf(a.theme.Red, "❌ "+a.strings.SearchError, err)
		})
		return
	}

	a.populateResults([]search.Result{*result})
}

func (a *SimpleApp) searchPlaylistURL(url string) {
	a.app.QueueUpdateDraw(func() {
		a.setStatus(a.theme.Yellow, "  "+a.strings.LoadingPlaylist)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	results, err := search.GetPlaylistVideos(ctx, url, 200)
	if err != nil {
		a.app.QueueUpdateDraw(func() {
			a.setStatusf(a.theme.Red, "❌ "+a.strings.SearchError, err)
		})
		return
	}

	a.populateResults(results)
}

func (a *SimpleApp) populateResults(results []search.Result) {
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

func (a *SimpleApp) addAllToPlaylist() {
	a.mu.Lock()
	tracks := make([]Track, len(a.tracks))
	copy(tracks, a.tracks)
	a.mu.Unlock()

	if len(tracks) == 0 {
		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Yellow, "⚠ "+a.strings.PlaylistEmpty)
		})
		return
	}

	for _, track := range tracks {
		a.addToPlaylist(track)
	}

	count := len(tracks)
	a.app.QueueUpdateDraw(func() {
		a.setStatusf(a.theme.Green, "✓ "+a.strings.PlaylistImported, count)
	})
}

func (a *SimpleApp) yankURL(focused interface{}) {
	var url string

	a.mu.Lock()
	// Try current playing track first
	if a.isPlaying && a.currentTrack >= 0 && a.currentTrack < len(a.playlistTracks) {
		url = a.playlistTracks[a.currentTrack].URL
	}
	a.mu.Unlock()

	// If nothing playing, try selected item in focused list
	if url == "" {
		switch focused {
		case a.searchResults.Flex:
			track := a.searchResults.GetCurrentTrack()
			if track != nil {
				url = track.URL
			}
		case a.playlist.Flex:
			idx := a.playlist.GetCurrentItem()
			a.mu.Lock()
			if idx >= 0 && idx < len(a.playlistTracks) {
				url = a.playlistTracks[idx].URL
			}
			a.mu.Unlock()
		}
	}

	if url == "" {
		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Yellow, "⚠ "+a.strings.NoTrackSelected)
		})
		return
	}

	if err := copyToClipboard(url); err != nil {
		a.app.QueueUpdateDraw(func() {
			a.setStatusf(a.theme.Red, "❌ "+a.strings.ClipboardError, err)
		})
		return
	}

	a.app.QueueUpdateDraw(func() {
		a.setStatusf(a.theme.Green, "  "+a.strings.URLCopied, url)
	})
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

	a.populateResults(results)
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
