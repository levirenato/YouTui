package ui

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/levirenato/YouTui/internal/search"
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
		return "🔀 Shuffle"
	case ModeRepeatOne:
		return "🔂 Repeat 1"
	case ModeRepeatAll:
		return "🔁 Repeat All"
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
		return "🎵 Audio"
	}
	return "🎬 Video"
}

type Track struct {
	Title  string
	Author string
	URL    string
}

type SimpleApp struct {
	app            *tview.Application
	searchInput    *tview.InputField
	searchResults  *tview.List
	playlist       *tview.List
	playerInfo     *tview.TextView
	statusBar      *tview.TextView
	commandBar     *tview.TextView
	helpModal      *tview.Modal
	
	tracks         []Track
	playlistTracks []Track
	
	mpvProcess     *exec.Cmd
	mpvSocket      string
	isPlaying      bool
	isPaused       bool
	currentTrack   int
	nowPlaying     string
	duration       float64
	position       float64
	
	playlistMode   PlaylistMode
	playMode       PlayMode
	
	progressTicker *time.Ticker
	stopProgress   chan bool
	
	skipAutoPlay   bool
	
	theme          *Theme
	
	mu             sync.Mutex
}

func NewSimpleApp() *SimpleApp {
	theme := CatppuccinMocha
	
	app := &SimpleApp{
		app:           tview.NewApplication(),
		tracks:        []Track{},
		playlistTracks: []Track{},
		playlistMode:  ModeNormal,
		playMode:      ModeAudio,
		currentTrack:  -1,
		theme:         &theme,
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
	
	app.setupSimple()
	return app
}

func (a *SimpleApp) setupSimple() {
	a.searchInput = tview.NewInputField().
		SetLabel("🔍 ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(a.theme.Surface0).
		SetFieldTextColor(a.theme.Text)
	
	a.searchInput.SetBorder(true).
		SetTitle(" Busca ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(a.theme.Blue)
	
	a.searchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			query := a.searchInput.GetText()
			if query != "" {
				go a.doSearch(query)
			}
		}
	})
	
	a.searchResults = tview.NewList().
		ShowSecondaryText(false).
		SetMainTextColor(a.theme.Text).
		SetSelectedTextColor(a.theme.Base).
		SetSelectedBackgroundColor(a.theme.Blue)
	
	a.searchResults.SetBorder(true).
		SetTitle(" Resultados [0] ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(a.theme.Surface0)
	
	a.searchResults.SetSelectedFunc(func(idx int, _ string, _ string, _ rune) {
		a.mu.Lock()
		if idx >= 0 && idx < len(a.tracks) {
			track := a.tracks[idx]
			a.mu.Unlock()
			go a.playTrackDirect(track)
		} else {
			a.mu.Unlock()
		}
	})
	
	a.playlist = tview.NewList().
		ShowSecondaryText(false).
		SetMainTextColor(a.theme.Text).
		SetSelectedTextColor(a.theme.Base).
		SetSelectedBackgroundColor(a.theme.Blue)
	
	a.playlist.SetBorder(true).
		SetTitle(" Playlist [0] ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(a.theme.Surface0)
	
	a.playlist.SetSelectedFunc(func(idx int, _ string, _ string, _ rune) {
		a.mu.Lock()
		if idx >= 0 && idx < len(a.playlistTracks) {
			track := a.playlistTracks[idx]
			a.mu.Unlock()
			go a.playTrackSimple(track, idx)
		} else {
			a.mu.Unlock()
		}
	})
	
	a.playerInfo = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(a.theme.Text)
	
	a.playerInfo.SetBorder(true).
		SetTitle(" Player ").
		SetBorderColor(a.theme.Mauve)
	
	a.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	
	a.commandBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	
	a.updatePlayerSimple()
	a.statusBar.SetText("")
	a.updateCommandBar()
	
	searchPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.searchInput, 3, 0, true).
		AddItem(a.searchResults, 0, 1, false)
	
	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(searchPanel, 0, 1, true).
		AddItem(a.playlist, 0, 1, false)
	
	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topFlex, 0, 1, true).
		AddItem(a.playerInfo, 5, 0, false).
		AddItem(a.statusBar, 1, 0, false).
		AddItem(a.commandBar, 1, 0, false)
	
	a.helpModal = tview.NewModal().
		SetText(a.getHelpText()).
		AddButtons([]string{"Fechar"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.app.SetRoot(mainLayout, true)
		})
	
	a.app.SetRoot(mainLayout, true).SetFocus(a.searchInput)
	
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		focused := a.app.GetFocus()
		
		if event.Rune() == '?' && focused != a.searchInput {
			a.app.SetRoot(a.helpModal, true)
			return nil
		}
		
		if focused == a.searchInput {
			if event.Key() == tcell.KeyTab {
				a.app.SetFocus(a.searchResults)
				a.updateCommandBar()
				return nil
			}
			return event
		}
		
		switch event.Key() {
		case tcell.KeyTab:
			switch focused {
			case a.searchResults:
				a.app.SetFocus(a.playlist)
				a.updateCommandBar()
			case a.playlist:
				a.app.SetFocus(a.searchInput)
				a.updateCommandBar()
			}
			return nil
		}
		
		switch event.Rune() {
		case 'q':
			a.cleanup()
			a.app.Stop()
			return nil
			
		case 'a':
			if focused == a.searchResults {
				idx := a.searchResults.GetCurrentItem()
				a.mu.Lock()
				if idx >= 0 && idx < len(a.tracks) {
					track := a.tracks[idx]
					a.mu.Unlock()
					go a.addToPlaylist(track)
				} else {
					a.mu.Unlock()
				}
				return nil
			}
			
		case 'd':
			if focused == a.playlist {
				idx := a.playlist.GetCurrentItem()
				go a.removeFromPlaylist(idx)
				return nil
			}
			
		case 'J':
			if focused == a.playlist {
				idx := a.playlist.GetCurrentItem()
				go a.movePlaylistItem(idx, idx+1)
				return nil
			}
			
		case 'K':
			if focused == a.playlist {
				idx := a.playlist.GetCurrentItem()
				go a.movePlaylistItem(idx, idx-1)
				return nil
			}
			
		case 'r':
			go a.cycleRepeatMode()
			return nil
			
		case 'h':
			go a.toggleShuffle()
			return nil
			
		case 'c', ' ':
			go a.togglePause()
			return nil
			
		case 's':
			go a.stopPlayback()
			return nil
			
		case 'n':
			go a.playNext()
			return nil
			
		case 'b':
			go a.playPrevious()
			return nil
			
		case 'm':
			go a.toggleMode()
			return nil
			
		case '/':
			a.app.SetFocus(a.searchInput)
			a.updateCommandBar()
			return nil
		}
		
		return event
	})
}

func (a *SimpleApp) doSearch(query string) {
	a.app.QueueUpdateDraw(func() {
		a.statusBar.SetText("[yellow]🔍 Buscando...")
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
			Title:  r.Title,
			Author: r.Author,
			URL:    r.URL,
		}
	}
	tracksCopy := make([]Track, len(a.tracks))
	copy(tracksCopy, a.tracks)
	a.mu.Unlock()
	
	a.app.QueueUpdateDraw(func() {
		a.searchResults.Clear()
		for i, track := range tracksCopy {
			title := fmt.Sprintf("%d. %s - %s", i+1, track.Title, track.Author)
			a.searchResults.AddItem(title, "", 0, nil)
		}
		a.searchResults.SetTitle(fmt.Sprintf(" Resultados [%d] ", len(tracksCopy)))
		a.app.SetFocus(a.searchResults)
		a.updateCommandBar()
	})
}

func (a *SimpleApp) playTrackSimple(track Track, idx int) {
	a.mu.Lock()
	if a.stopProgress != nil {
		close(a.stopProgress)
		a.stopProgress = nil
	}
	if a.mpvProcess != nil && a.mpvProcess.Process != nil {
		a.mpvProcess.Process.Kill()
		a.mpvProcess = nil
	}
	a.mu.Unlock()
	
	socketPath := fmt.Sprintf("/tmp/mpv-socket-%d", time.Now().UnixNano())
	
	args := []string{
		"--no-terminal",
		"--really-quiet",
		"--script-opts=ytdl_hook-ytdl_path=yt-dlp",
		fmt.Sprintf("--title=%s", track.Title),
		fmt.Sprintf("--input-ipc-server=%s", socketPath),
	}
	
	a.mu.Lock()
	if a.playMode == ModeAudio {
		args = append(args, "--no-video", "--ytdl-format=bestaudio")
	}
	a.mu.Unlock()
	
	args = append(args, track.URL)
	
	cmd := exec.Command("mpv", args...)
	if err := cmd.Start(); err != nil {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[red]❌ Erro mpv: %v", err))
		})
		return
	}
	
	a.mu.Lock()
	a.mpvProcess = cmd
	a.mpvSocket = socketPath
	a.isPlaying = true
	a.isPaused = false
	a.nowPlaying = track.Title
	a.currentTrack = idx
	a.position = 0
	a.duration = 0
	a.mu.Unlock()
	
	a.app.QueueUpdateDraw(func() {
		a.updatePlayerSimple()
		a.statusBar.SetText(fmt.Sprintf("[green]▶ Tocando: %s", track.Title))
	})
	
	a.startProgressUpdater()
	
	go func() {
		cmd.Wait()
		
		time.Sleep(500 * time.Millisecond)
		
		a.mu.Lock()
		if a.skipAutoPlay {
			a.skipAutoPlay = false
			a.mu.Unlock()
			return
		}
		
		mode := a.playlistMode
		var shouldPlayNext bool
		var nextTrack Track
		var nextIdx int
		
		switch mode {
		case ModeRepeatOne:
			shouldPlayNext = true
			nextTrack = track
			nextIdx = idx
			
		case ModeRepeatAll, ModeNormal:
			if len(a.playlistTracks) > 0 {
				next := idx + 1
				if next >= len(a.playlistTracks) {
					if mode == ModeRepeatAll {
						next = 0
						shouldPlayNext = true
						nextTrack = a.playlistTracks[next]
						nextIdx = next
					} else {
						shouldPlayNext = false
					}
				} else {
					shouldPlayNext = true
					nextTrack = a.playlistTracks[next]
					nextIdx = next
				}
			}
		}
		a.mu.Unlock()
		
		if shouldPlayNext {
			go a.playTrackSimple(nextTrack, nextIdx)
		} else {
			a.mu.Lock()
			a.isPlaying = false
			a.mu.Unlock()
			
			a.app.QueueUpdateDraw(func() {
				a.updatePlayerSimple()
				a.statusBar.SetText("[yellow]Playlist finalizada")
			})
		}
	}()
}

func (a *SimpleApp) addToPlaylist(track Track) {
	a.mu.Lock()
	a.playlistTracks = append(a.playlistTracks, track)
	count := len(a.playlistTracks)
	a.mu.Unlock()
	
	a.app.QueueUpdateDraw(func() {
		a.playlist.AddItem(fmt.Sprintf("%d. %s", count, track.Title), "", 0, nil)
		a.playlist.SetTitle(fmt.Sprintf(" Playlist [%d] ", count))
		a.statusBar.SetText(fmt.Sprintf("[green]✓ Adicionado: %s", track.Title))
	})
}

func (a *SimpleApp) removeFromPlaylist(idx int) {
	a.mu.Lock()
	if idx < 0 || idx >= len(a.playlistTracks) {
		a.mu.Unlock()
		return
	}
	
	if idx == a.currentTrack {
		a.mu.Unlock()
		a.stopPlayback()
		a.mu.Lock()
	} else if idx < a.currentTrack {
		a.currentTrack--
	}
	
	a.playlistTracks = append(a.playlistTracks[:idx], a.playlistTracks[idx+1:]...)
	tracks := make([]Track, len(a.playlistTracks))
	copy(tracks, a.playlistTracks)
	count := len(a.playlistTracks)
	a.mu.Unlock()
	
	a.app.QueueUpdateDraw(func() {
		a.playlist.Clear()
		for i, t := range tracks {
			a.playlist.AddItem(fmt.Sprintf("%d. %s", i+1, t.Title), "", 0, nil)
		}
		a.playlist.SetTitle(fmt.Sprintf(" Playlist [%d] ", count))
		a.statusBar.SetText("[yellow]✓ Removido da playlist")
	})
}

func (a *SimpleApp) movePlaylistItem(from, to int) {
	a.mu.Lock()
	if from < 0 || from >= len(a.playlistTracks) {
		a.mu.Unlock()
		return
	}
	if to < 0 || to >= len(a.playlistTracks) {
		a.mu.Unlock()
		return
	}
	if from == to {
		a.mu.Unlock()
		return
	}
	
	a.playlistTracks[from], a.playlistTracks[to] = a.playlistTracks[to], a.playlistTracks[from]
	
	if a.currentTrack == from {
		a.currentTrack = to
	} else if a.currentTrack == to {
		a.currentTrack = from
	}
	
	tracks := make([]Track, len(a.playlistTracks))
	copy(tracks, a.playlistTracks)
	newPos := to
	a.mu.Unlock()
	
	a.app.QueueUpdateDraw(func() {
		a.playlist.Clear()
		for i, t := range tracks {
			a.playlist.AddItem(fmt.Sprintf("%d. %s", i+1, t.Title), "", 0, nil)
		}
		a.playlist.SetCurrentItem(newPos)
		a.statusBar.SetText("[cyan]✓ Item movido")
	})
}

func (a *SimpleApp) cycleRepeatMode() {
	switch a.playlistMode {
	case ModeNormal:
		a.playlistMode = ModeRepeatOne
	case ModeRepeatOne:
		a.playlistMode = ModeRepeatAll
	case ModeRepeatAll:
		a.playlistMode = ModeNormal
	}
	
	a.app.QueueUpdateDraw(func() {
		a.updatePlayerSimple()
		a.statusBar.SetText(fmt.Sprintf("[cyan]🔁 Modo: %s", a.playlistMode.String()))
	})
}

func (a *SimpleApp) toggleShuffle() {
	if a.playlistMode == ModeShuffle {
		a.playlistMode = ModeNormal
	} else {
		a.playlistMode = ModeShuffle
	}
	
	a.app.QueueUpdateDraw(func() {
		a.updatePlayerSimple()
		a.statusBar.SetText(fmt.Sprintf("[cyan]🔀 Modo: %s", a.playlistMode.String()))
	})
}

func (a *SimpleApp) togglePause() {
	a.mu.Lock()
	isPlaying := a.isPlaying
	socket := a.mpvSocket
	a.mu.Unlock()
	
	if !isPlaying || socket == "" {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[red]⚠ Estado: isPlaying=%v socket=%s", isPlaying, socket))
		})
		return
	}
	
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`echo '{ "command": ["cycle", "pause"] }' | socat - "%s" 2>&1`, socket))
	output, err := cmd.CombinedOutput()
	if err != nil {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[red]❌ Erro: %v | %s", err, string(output)))
		})
		return
	}
	
	a.mu.Lock()
	a.isPaused = !a.isPaused
	isPaused := a.isPaused
	a.mu.Unlock()
	
	a.app.QueueUpdateDraw(func() {
		a.updatePlayerSimple()
		if isPaused {
			a.statusBar.SetText("[yellow]⏸ Pausado")
		} else {
			a.statusBar.SetText("[green]▶ Tocando")
		}
	})
}

func (a *SimpleApp) stopPlayback() {
	a.cleanup()
	a.app.QueueUpdateDraw(func() {
		a.updatePlayerSimple()
		a.statusBar.SetText("[red]⏹ Parado")
	})
}

func (a *SimpleApp) playTrackDirect(track Track) {
	a.cleanup()
	
	socketPath := fmt.Sprintf("/tmp/mpv-socket-%d", time.Now().UnixNano())
	
	args := []string{
		"--no-terminal",
		"--really-quiet",
		"--script-opts=ytdl_hook-ytdl_path=yt-dlp",
		fmt.Sprintf("--title=%s", track.Title),
		fmt.Sprintf("--input-ipc-server=%s", socketPath),
	}
	
	a.mu.Lock()
	if a.playMode == ModeAudio {
		args = append(args, "--no-video", "--ytdl-format=bestaudio")
	}
	a.mu.Unlock()
	
	args = append(args, track.URL)
	
	cmd := exec.Command("mpv", args...)
	if err := cmd.Start(); err != nil {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[red]❌ Erro mpv: %v", err))
		})
		return
	}
	
	a.mu.Lock()
	a.mpvProcess = cmd
	a.mpvSocket = socketPath
	a.isPlaying = true
	a.isPaused = false
	a.nowPlaying = track.Title
	a.currentTrack = -1
	a.position = 0
	a.duration = 0
	a.mu.Unlock()
	
	a.app.QueueUpdateDraw(func() {
		a.updatePlayerSimple()
		a.statusBar.SetText(fmt.Sprintf("[green]▶ Tocando: %s (sem playlist)", track.Title))
	})
	
	a.startProgressUpdater()
	
	go func() {
		cmd.Wait()
		
		a.mu.Lock()
		a.isPlaying = false
		a.mu.Unlock()
		
		a.app.QueueUpdateDraw(func() {
			a.updatePlayerSimple()
			a.statusBar.SetText("[yellow]Reprodução finalizada")
		})
	}()
}

func (a *SimpleApp) playNext() {
	a.mu.Lock()
	currentIsPlaying := a.isPlaying
	currentTrack := a.currentTrack
	playlistLen := len(a.playlistTracks)
	
	if playlistLen == 0 {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ Playlist vazia")
		})
		return
	}
	
	if !currentIsPlaying {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[yellow]⚠ isPlaying=%v - Inicie a playlist primeiro.", currentIsPlaying))
		})
		return
	}
	
	if currentTrack < 0 {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[yellow]⚠ currentTrack=%d - Não está tocando da playlist", currentTrack))
		})
		return
	}
	
	var next int
	
	if a.playlistMode == ModeShuffle {
		next = rand.Intn(playlistLen)
	} else {
		next = currentTrack + 1
		if next >= playlistLen {
			if a.playlistMode == ModeRepeatAll {
				next = 0
			} else {
				a.mu.Unlock()
				a.app.QueueUpdateDraw(func() {
					a.statusBar.SetText("[yellow]Já está na última música")
				})
				return
			}
		}
	}
	
	track := a.playlistTracks[next]
	a.skipAutoPlay = true
	a.mu.Unlock()
	
	a.app.QueueUpdateDraw(func() {
		a.statusBar.SetText(fmt.Sprintf("[green]▶ Pulando para: %d/%d - %s", next+1, playlistLen, track.Title))
	})
	
	go a.playTrackSimple(track, next)
}

func (a *SimpleApp) playPrevious() {
	a.mu.Lock()
	if len(a.playlistTracks) == 0 {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ Playlist vazia")
		})
		return
	}
	
	if !a.isPlaying {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ Nada tocando. Inicie a playlist primeiro.")
		})
		return
	}
	
	if a.currentTrack < 0 {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ Não está tocando da playlist. Use Space na playlist para iniciar.")
		})
		return
	}
	
	prev := a.currentTrack - 1
	if prev < 0 {
		if a.playlistMode == ModeRepeatAll {
			prev = len(a.playlistTracks) - 1
		} else {
			a.mu.Unlock()
			a.app.QueueUpdateDraw(func() {
				a.statusBar.SetText("[yellow]Já está na primeira música")
			})
			return
		}
	}
	track := a.playlistTracks[prev]
	a.skipAutoPlay = true
	a.mu.Unlock()
	
	go a.playTrackSimple(track, prev)
}

func (a *SimpleApp) toggleMode() {
	if a.playMode == ModeAudio {
		a.playMode = ModeVideo
	} else {
		a.playMode = ModeAudio
	}
	
	a.app.QueueUpdateDraw(func() {
		a.updatePlayerSimple()
		a.statusBar.SetText(fmt.Sprintf("[cyan]🎬 Modo: %s", a.playMode.String()))
	})
}

func (a *SimpleApp) updatePlayerSimple() {
	var status string
	if a.isPlaying {
		icon := "▶"
		if a.isPaused {
			icon = "⏸"
		}
		status = fmt.Sprintf("[green::b]%s %s[-:-:-]\n", icon, a.nowPlaying)
	} else {
		status = "[gray]⏹ Nenhuma faixa tocando[-]\n"
	}
	
	var progress string
	if a.isPlaying && a.duration > 0 {
		percentage := a.position / a.duration
		if percentage > 1 {
			percentage = 1
		}
		if percentage < 0 {
			percentage = 0
		}
		
		totalBars := 40
		filledBars := int(percentage * float64(totalBars))
		emptyBars := totalBars - filledBars
		
		posMin := int(a.position / 60)
		posSec := int(a.position) % 60
		durMin := int(a.duration / 60)
		durSec := int(a.duration) % 60
		
		progress = fmt.Sprintf("[blue]%s[gray]%s[-] [cyan]%02d:%02d[-]/[white]%02d:%02d[-]",
			strings.Repeat("█", filledBars),
			strings.Repeat("░", emptyBars),
			posMin, posSec,
			durMin, durSec)
	} else {
		progress = "[gray]" + strings.Repeat("░", 40) + " --:--/--:--[-]"
	}
	
	modeInfo := fmt.Sprintf("%s | %s", a.playMode.String(), a.playlistMode.String())
	controls := "[yellow]a[-] Add [yellow]d[-] Del [yellow]r[-] Repeat [yellow]h[-] Shuffle [yellow]c[-] Pause [yellow]n/b[-] Next/Prev"
	
	a.playerInfo.SetText(fmt.Sprintf("%s%s\n[cyan]%s[-]\n%s", status, progress, modeInfo, controls))
}

func (a *SimpleApp) updateCommandBar() {
	focused := a.app.GetFocus()
	
	a.searchInput.SetBorderColor(a.theme.Surface0)
	a.searchResults.SetBorderColor(a.theme.Surface0)
	a.playlist.SetBorderColor(a.theme.Surface0)
	
	var help string
	switch focused {
	case a.searchInput:
		a.searchInput.SetBorderColor(a.theme.Blue)
		help = "Digite para buscar | [#89b4fa]Enter[-] Buscar | [#89b4fa]Tab[-] Próximo | [#f38ba8]q[-] Sair | [#f9e2af]?[-] Ajuda"
		
	case a.searchResults:
		a.searchResults.SetBorderColor(a.theme.Blue)
		help = "[#89b4fa]↑/↓[-] Navegar | [#89b4fa]Enter[-] Tocar | [#a6e3a1]a[-] Add | [#89b4fa]Tab[-] Próximo | [#89b4fa]/[-] Buscar | [#f38ba8]q[-] Sair | [#f9e2af]?[-] Ajuda"
		
	case a.playlist:
		a.playlist.SetBorderColor(a.theme.Blue)
		help = "[#89b4fa]↑/↓[-] Nav | [#89b4fa]Enter[-] Play | [#f38ba8]d[-] Del | [#cba6f7]J/K[-] Move | [#fab387]r[-] Repeat | [#94e2d5]h[-] Shuffle | [#a6e3a1]c[-] Pause | [#89dceb]n/b[-] Next/Prev | [#f9e2af]?[-] Ajuda"
		
	default:
		help = "[#89b4fa]Tab[-] Navegar | [#a6e3a1]a[-] Add | [#f38ba8]d[-] Del | [#fab387]r[-] Repeat | [#94e2d5]h[-] Shuffle | [#a6e3a1]c[-] Pause | [#89dceb]n/b[-] Next/Prev | [#f38ba8]q[-] Sair | [#f9e2af]?[-] Ajuda"
	}
	
	a.commandBar.SetText(help)
}

func (a *SimpleApp) startProgressUpdater() {
	a.mu.Lock()
	if a.stopProgress != nil {
		select {
		case <-a.stopProgress:
		default:
			close(a.stopProgress)
		}
	}
	a.stopProgress = make(chan bool)
	stopChan := a.stopProgress
	a.mu.Unlock()
	
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				a.mu.Lock()
				isPlaying := a.isPlaying
				a.mu.Unlock()
				
				if !isPlaying {
					return
				}
				a.updateProgress()
			case <-stopChan:
				return
			}
		}
	}()
}

func (a *SimpleApp) updateProgress() {
	if a.mpvSocket == "" {
		return
	}
	
	posCmd := exec.Command("sh", "-c", fmt.Sprintf(`echo '{ "command": ["get_property", "time-pos"] }' | socat - UNIX-CONNECT:%s 2>/dev/null | grep -o '"data":[0-9.]*' | cut -d: -f2`, a.mpvSocket))
	posOut, _ := posCmd.Output()
	if len(posOut) > 0 {
		if pos, err := strconv.ParseFloat(strings.TrimSpace(string(posOut)), 64); err == nil {
			a.position = pos
		}
	}
	
	durCmd := exec.Command("sh", "-c", fmt.Sprintf(`echo '{ "command": ["get_property", "duration"] }' | socat - UNIX-CONNECT:%s 2>/dev/null | grep -o '"data":[0-9.]*' | cut -d: -f2`, a.mpvSocket))
	durOut, _ := durCmd.Output()
	if len(durOut) > 0 {
		if dur, err := strconv.ParseFloat(strings.TrimSpace(string(durOut)), 64); err == nil && dur > 0 {
			a.duration = dur
		}
	}
	
	a.app.QueueUpdateDraw(func() {
		a.updatePlayerSimple()
	})
}

func (a *SimpleApp) getHelpText() string {
	return `ATALHOS DO YOUTUI

NAVEGAÇÃO:
  Tab       Alternar entre painéis
  /         Focar na busca
  ↑/↓       Navegar nas listas
  ?         Mostrar esta ajuda

BUSCA:
  Digite    Texto para buscar
  Enter     Executar busca

RESULTADOS:
  Enter     Tocar faixa diretamente (sem playlist)
  a         Adicionar à playlist

PLAYLIST:
  Enter     Tocar faixa da playlist
  Space     Tocar playlist do início
  d         Remover item
  J         Mover item para baixo
  K         Mover item para cima

PLAYER (Global):
  c/Space   Pause/Play
  s         Stop
  n         Próxima (só funciona tocando da playlist)
  b         Anterior (só funciona tocando da playlist)
  m         Alternar áudio/vídeo
  r         Ciclar repetição
  h         Toggle shuffle

IMPORTANTE:
  n/b só funcionam quando tocando da PLAYLIST!
  Para tocar da playlist: Entre na Playlist e pressione Enter ou Space

GERAL:
  q         Sair
`
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
	a.mpvSocket = ""
	a.position = 0
	a.duration = 0
	a.mu.Unlock()
}

func (a *SimpleApp) Run() error {
	return a.app.Run()
}
