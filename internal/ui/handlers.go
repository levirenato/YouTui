package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// handleKeyPress processa as teclas pressionadas globalmente
func (a *SimpleApp) handleKeyPress(event *tcell.EventKey, focused tview.Primitive) *tcell.EventKey {
	switch event.Rune() {
	case 'a':
		if focused == a.searchResults.Flex {
			track := a.searchResults.GetCurrentTrack()
			if track != nil {
				go a.addToPlaylist(*track)
			}
			return nil
		}

	case 'd':
		if focused == a.playlist.Flex {
			idx := a.playlist.GetCurrentItem()
			go a.removeFromPlaylist(idx)
			return nil
		}

	case 'J':
		if focused == a.playlist.Flex {
			idx := a.playlist.GetCurrentItem()
			go a.movePlaylistItem(idx, idx+1)
			return nil
		}

	case 'K':
		if focused == a.playlist.Flex {
			idx := a.playlist.GetCurrentItem()
			go a.movePlaylistItem(idx, idx-1)
			return nil
		}

	// Controles da PLAYLIST (r, h)
	case 'r':
		if focused == a.playlist.Flex {
			go a.cycleRepeatMode()
			return nil
		}

	case 'h':
		if focused == a.playlist.Flex {
			go a.toggleShuffle()
			return nil
		}

	// Controles do PLAYER (space, n, b, s)
	case 'c', ' ':
		if focused == a.playerBox {
			go a.togglePause()
			return nil
		}

	case 's':
		if focused == a.playerBox {
			go a.stopPlayback()
			return nil
		}

	case 'n':
		if focused == a.playerBox {
			go a.playNext()
			return nil
		}

	case 'p':
		if focused == a.playerBox {
			go a.playPrevious()
			return nil
		}

	// Controle global (m)
	case 'm':
		go a.toggleMode()
		return nil

	case '/':
		a.app.SetFocus(a.searchInput)
		a.updateCommandBar()
		return nil

	case ']':
		// Próxima página (só nos resultados)
		if focused == a.searchResults.Flex {
			go a.nextPage()
			return nil
		}

	case '[':
		// Página anterior (só nos resultados)
		if focused == a.searchResults.Flex {
			go a.prevPage()
			return nil
		}
	}

	return event
}

// updateCommandBar atualiza a barra de comandos com atalhos contextuais
func (a *SimpleApp) updateCommandBar() {
	focused := a.app.GetFocus()

	a.searchInput.SetBorderColor(a.theme.Surface0)
	a.searchResults.SetBorderColor(a.theme.Surface0)
	a.playlist.SetBorderColor(a.theme.Surface0)
	a.playerBox.SetBorderColor(a.theme.Surface0)

	var help string
	
	switch focused {
	case a.searchInput:
		a.searchInput.SetBorderColor(a.theme.Blue)
		help = "Digite para buscar | [#89b4fa]Enter[-] Buscar | [#89b4fa]Tab[-] Próximo | [#f38ba8]Ctrl+Q[-] Sair | [#f9e2af]?[-] Ajuda"

	case a.searchResults.Flex:
		a.searchResults.SetBorderColor(a.theme.Blue)
		help = "[#89b4fa]↑/↓[-] Nav | [#89b4fa]Enter[-] Play | [#a6e3a1]a[-] Add | [#cba6f7][ ][-] Pág | [#89b4fa]/[-] Buscar | [#f38ba8]Ctrl+Q[-] Sair"

	case a.playlist.Flex:
		a.playlist.SetBorderColor(a.theme.Blue)
		help = "[#89b4fa]↑/↓[-] Nav | [#89b4fa]Enter[-] Play | [#f38ba8]d[-] Del | [#cba6f7]J/K[-] Move | [#fab387]r[-] Repeat | [#94e2d5]h[-] Shuffle"

	case a.playerBox:
		a.playerBox.SetBorderColor(a.theme.Blue)
		help = "[#a6e3a1]Space[-] Pause/Play | [#89dceb]n[-] Next | [#89dceb]p[-] Prev | [#f38ba8]s[-] Stop | [#cba6f7]m[-] Modo"

	default:
		help = "[#89b4fa]Tab[-] Navegar entre painéis | [#f38ba8]Ctrl+Q[-] Sair | [#f9e2af]?[-] Ajuda"
	}

	a.commandBar.SetText(help)
}
