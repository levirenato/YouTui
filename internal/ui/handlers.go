package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// handleKeyPress processa as teclas pressionadas globalmente
func (a *SimpleApp) handleKeyPress(event *tcell.EventKey, focused tview.Primitive) *tcell.EventKey {
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
}

// updateCommandBar atualiza a barra de comandos com atalhos contextuais
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
