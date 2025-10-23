package ui

import (
	"context"
	"os/exec"
	"sync"
	"time"

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
		return " Audio"
	}
	return "  Video"
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

	progressTicker *time.Ticker
	stopProgress   chan bool

	skipAutoPlay bool

	thumbCache           *ThumbnailCache
	detailsLoadingIdx    int
	detailsLoadingMutex  sync.Mutex
	detailsCancelFunc    context.CancelFunc
	detailsDebounceTimer *time.Timer

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
				app.setStatus(app.theme.Yellow, "⚠ yt-dlp desatualizado ("+currentVersion+"). Atualize: sudo yt-dlp -U")
			})
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
		a.mpvProcess.Process.Kill()
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
