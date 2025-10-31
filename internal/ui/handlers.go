package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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

	case 'm':
		go a.toggleMode()
		return nil

	case '/':
		a.app.SetFocus(a.searchInput)
		a.updateCommandBar()
		return nil

	case ']':
		if focused == a.searchResults.Flex {
			go a.nextPage()
			return nil
		}

	case '[':
		if focused == a.searchResults.Flex {
			go a.prevPage()
			return nil
		}
	}

	return event
}

func (a *SimpleApp) updateCommandBar() {
	focused := a.app.GetFocus()

	a.searchInput.SetBorderColor(a.theme.Surface0)
	a.searchInput.SetBackgroundColor(a.theme.Base)
	a.searchInput.SetTitleColor(a.theme.Subtext0)
	a.searchResults.SetBorderColor(a.theme.Surface0)
	a.playlist.SetBorderColor(a.theme.Surface0)
	a.playerBox.SetBorderColor(a.theme.Surface0)

	var help string

	switch focused {
	case a.searchInput:
		a.searchInput.SetBorderColor(a.theme.Blue)
		help = a.strings.CmdSearchBar

	case a.searchResults.Flex:
		a.searchResults.SetBorderColor(a.theme.Blue)
		help = a.strings.CmdResultsBar

	case a.playlist.Flex:
		a.playlist.SetBorderColor(a.theme.Blue)
		help = a.strings.CmdPlaylistBar

	case a.playerBox:
		a.playerBox.SetBorderColor(a.theme.Blue)
		help = a.strings.CmdPlayerBar

	default:
		help = a.strings.CmdDefaultBar
	}

	a.commandBar.SetText(help)
}
