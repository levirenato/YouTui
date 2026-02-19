// Package ui
package ui

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/levirenato/YouTui/internal/config"
	"github.com/rivo/tview"
)

type PlaylistMode int

const (
	ModeNormal PlaylistMode = iota
	ModeRepeatOne
	ModeRepeatAll
	ModeShuffle
)

func (m PlaylistMode) String() string {
	switch m {
	case ModeShuffle:
		return " Shuffle"
	case ModeRepeatOne:
		return "󰑘 Repeat 1"
	case ModeRepeatAll:
		return "󰑖 Repeat All"
	default:
		return "▶ Normal"
	}
}

type PlayMode int

const (
	ModeAudio PlayMode = iota
	ModeVideo
)

func (m PlayMode) String() string {
	if m == ModeAudio {
		return " Audio"
	}
	return "󰗃 Video"
}

type Track struct {
	Title       string
	Author      string
	URL         string
	Thumbnail   string
	Duration    string
	PublishedAt string
	Description string
}

type SimpleApp struct {
	app            *tview.Application
	searchInput    *tview.InputField
	searchResults  *CustomList
	playlist       *CustomList
	detailsView    *tview.Flex
	detailsThumb   *tview.Image
	detailsText    *tview.TextView
	thumbnailView  *tview.Image
	playerInfo     *tview.TextView
	playerBox      *tview.Flex
	playlistFooter *tview.TextView
	statusBar      *tview.TextView
	commandBar     *tview.TextView
	modeBadge      *tview.TextView
	helpModal      *tview.Modal
	configModal    *tview.Modal

	tracks         []Track
	playlistTracks []Track
	pagination     *Pagination

	mpvProcess   *exec.Cmd
	mpvSocket    string
	isPlaying    bool
	isPaused     bool
	currentTrack int
	nowPlaying   string
	currentThumb string
	duration     float64
	position     float64

	playlistMode PlaylistMode
	playMode     PlayMode

	stopProgress chan bool

	skipAutoPlay bool

	thumbCache *ThumbnailCache

	theme    *Theme
	language Language
	strings  Strings

	mu sync.Mutex
}

func NewSimpleApp() *SimpleApp {
	cfg, _ := config.LoadConfig()

	var theme *Theme
	if cfg.Theme.Active == "custom" && cfg.Theme.CustomPath != "" {
		customTheme, err := LoadCustomTheme(cfg.Theme.CustomPath)
		if err == nil {
			theme = customTheme
		} else {
			theme = GetThemeByID("catppuccin-mocha")
		}
	} else {
		theme = GetThemeByID(cfg.Theme.Active)
	}

	lang := LanguageEN
	thumbCache, _ := NewThumbnailCache()

	app := &SimpleApp{
		app:            tview.NewApplication(),
		tracks:         []Track{},
		playlistTracks: []Track{},
		pagination:     NewPagination(10),
		playlistMode:   ModeNormal,
		playMode:       ModeAudio,
		currentTrack:   -1,
		theme:          theme,
		language:       lang,
		strings:        GetStrings(lang),
		thumbCache:     thumbCache,
	}

	tview.Styles.PrimitiveBackgroundColor = theme.Base
	tview.Styles.ContrastBackgroundColor = theme.Surface0
	tview.Styles.MoreContrastBackgroundColor = theme.Surface1
	tview.Styles.BorderColor = theme.Surface0
	tview.Styles.TitleColor = theme.Text
	tview.Styles.GraphicsColor = theme.Blue
	tview.Styles.PrimaryTextColor = theme.Text
	tview.Styles.SecondaryTextColor = theme.Subtext1
	tview.Styles.TertiaryTextColor = theme.Subtext0
	tview.Styles.InverseTextColor = theme.Base
	tview.Styles.ContrastSecondaryTextColor = theme.Subtext0

	app.setupUI()

	go func() {
		currentVersion, _, needsUpdate := CheckYtDlpVersion()
		if needsUpdate {
			app.app.QueueUpdateDraw(func() {
				app.setStatus(app.theme.Yellow, "⚠"+app.strings.ytDlpOutdated+"("+currentVersion+")")
			})
		}
	}()

	go func() {
		if err := app.RestoreState(); err != nil {
			_ = err
		}
	}()

	return app
}

func (a *SimpleApp) Run() error {
	return a.app.Run()
}

func (a *SimpleApp) cleanup() {
	a.mu.Lock()
	if a.stopProgress != nil {
		close(a.stopProgress)
		a.stopProgress = nil
	}

	if a.mpvProcess != nil && a.mpvProcess.Process != nil {
		if KillError := a.mpvProcess.Process.Kill(); KillError == nil {
			fmt.Printf("Error: %s", KillError)
		}
		a.mpvProcess = nil
	}
	a.isPlaying = false
	a.isPaused = false
	a.nowPlaying = ""
	a.currentThumb = ""
	a.mpvSocket = ""
	a.position = 0
	a.duration = 0
	a.mu.Unlock()
}

func (a *SimpleApp) SaveCurrentState() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	state := &config.PlayerState{
		SearchTerm:        a.getSearchTerm(),
		SearchResults:     convertTracksToConfigTracks(a.tracks),
		Playlist:          convertTracksToConfigTracks(a.playlistTracks),
		CurrentTrackIdx:   a.currentTrack,
		PlaylistMode:      int(a.playlistMode),
		PlayMode:          int(a.playMode),
		SearchScrollIdx:   a.searchResults.GetCurrentItem(),
		PlaylistScrollIdx: a.playlist.GetCurrentItem(),
		SearchPage:        a.pagination.GetCurrentPage(),
	}

	return config.SaveState(state)
}

func (a *SimpleApp) RestoreState() error {
	state, err := config.LoadState()
	if err != nil {
		return err
	}

	if state == nil || (len(state.SearchResults) == 0 && len(state.Playlist) == 0) {
		return nil
	}

	a.mu.Lock()

	a.playlistMode = PlaylistMode(state.PlaylistMode)
	a.playMode = PlayMode(state.PlayMode)
	a.currentTrack = state.CurrentTrackIdx

	a.tracks = convertConfigTracksToTracks(state.SearchResults)
	a.playlistTracks = convertConfigTracksToTracks(state.Playlist)

	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		if state.SearchTerm != "" {
			a.searchInput.SetText(state.SearchTerm)
		}

		// Restaurar paginação antes de exibir a página correta
		a.pagination.SetTotalItems(len(a.tracks))
		a.pagination.SetCurrentPage(state.SearchPage)

		// Exibir apenas a página salva (com paginação correta)
		start, end := a.pagination.GetPageItems()
		for i, track := range a.tracks[start:end] {
			a.searchResults.AddItem(track, i)
			if track.Thumbnail != "" && a.thumbCache != nil {
				go func(idx int, url string) {
					img, err := a.thumbCache.GetThumbnailImage(url)
					if err == nil && img != nil {
						a.searchResults.SetThumbnail(idx, img)
					}
				}(start+i, track.Thumbnail)
			}
		}

		if state.SearchScrollIdx > 0 && state.SearchScrollIdx < (end-start) {
			a.searchResults.SetCurrentIndex(state.SearchScrollIdx)
		}

		currentPage := a.pagination.GetCurrentPage() + 1
		totalPages := a.pagination.GetTotalPages()
		if totalPages > 0 {
			a.searchResults.SetTitle(fmt.Sprintf(" %s [%s %d/%d] ", a.strings.Results, a.strings.Page, currentPage, totalPages))
		} else {
			a.searchResults.SetTitle(fmt.Sprintf(" %s [0] ", a.strings.Results))
		}

		a.playlist.Clear()
		for i, track := range a.playlistTracks {
			a.playlist.AddItem(track, i)

			if track.Thumbnail != "" && a.thumbCache != nil {
				go func(idx int, url string) {
					img, err := a.thumbCache.GetThumbnailImage(url)
					if err == nil && img != nil {
						a.playlist.SetThumbnail(idx, img)
					}
				}(i, track.Thumbnail)
			}
		}

		if state.PlaylistScrollIdx > 0 && state.PlaylistScrollIdx < len(a.playlistTracks) {
			a.playlist.SetCurrentIndex(state.PlaylistScrollIdx)
		}

		if state.CurrentTrackIdx >= 0 && state.CurrentTrackIdx < len(a.playlistTracks) {
			a.playlist.SetPlayingIndex(state.CurrentTrackIdx)
		}

		a.playlist.SetTitle(fmt.Sprintf(" Playlist [%d] ", len(a.playlistTracks)))

		a.updatePlayerInfo()
		a.updatePlaylistFooter()

		if len(state.Playlist) > 0 || len(state.SearchResults) > 0 {
			a.setStatus(a.theme.Green, a.strings.stateRestored)
		}
	})

	return nil
}

func convertTracksToConfigTracks(tracks []Track) []config.Track {
	result := make([]config.Track, len(tracks))
	for i, t := range tracks {
		result[i] = config.Track{
			Title:       t.Title,
			Author:      t.Author,
			URL:         t.URL,
			Thumbnail:   t.Thumbnail,
			Duration:    t.Duration,
			PublishedAt: t.PublishedAt,
			Description: t.Description,
		}
	}
	return result
}

func convertConfigTracksToTracks(tracks []config.Track) []Track {
	result := make([]Track, len(tracks))
	for i, t := range tracks {
		result[i] = Track{
			Title:       t.Title,
			Author:      t.Author,
			URL:         t.URL,
			Thumbnail:   t.Thumbnail,
			Duration:    t.Duration,
			PublishedAt: t.PublishedAt,
			Description: t.Description,
		}
	}
	return result
}

func (a *SimpleApp) getSearchTerm() string {
	return a.searchInput.GetText()
}

func (a *SimpleApp) AutoSaveState() {
	go func() {
		if err := a.SaveCurrentState(); err != nil {
			_ = err
		}
	}()
}
