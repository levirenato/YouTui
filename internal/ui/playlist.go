package ui

import (
	"fmt"
)

// onPlaylistSelectedCustom é chamado quando Enter é pressionado na playlist
func (a *SimpleApp) onPlaylistSelectedCustom(idx int) {
	track := a.playlist.GetCurrentTrack()
	if track != nil {
		// Pega o índice real da playlist
		a.mu.Lock()
		realIdx := a.playlist.GetCurrentItem()
		a.mu.Unlock()
		go a.playTrackSimple(*track, realIdx)
	}
}

// addToPlaylist adiciona uma faixa à playlist
func (a *SimpleApp) addToPlaylist(track Track) {
	a.mu.Lock()
	a.playlistTracks = append(a.playlistTracks, track)
	count := len(a.playlistTracks)
	index := count - 1
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		// Adiciona item com thumbnail inline
		a.playlist.AddItem(track, index)
		
		// Carrega thumbnail em background (usa cache!)
		if track.Thumbnail != "" && a.thumbCache != nil {
			go func(idx int, url string) {
				img, err := a.thumbCache.GetThumbnailImage(url)
				if err == nil && img != nil {
					a.playlist.SetThumbnail(idx, img)
				}
			}(index, track.Thumbnail)
		}
		
		a.playlist.SetTitle(fmt.Sprintf(" Playlist [%d] ", count))
		a.statusBar.SetText(fmt.Sprintf("[green]✓ Adicionado: %s", track.Title))
	})
}

// removeFromPlaylist remove uma faixa da playlist
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
			// Adiciona item com thumbnail
			a.playlist.AddItem(t, i)
			
			// Carrega thumbnail em background (usa cache!)
			if t.Thumbnail != "" && a.thumbCache != nil {
				go func(idx int, url string) {
					img, err := a.thumbCache.GetThumbnailImage(url)
					if err == nil && img != nil {
						a.playlist.SetThumbnail(idx, img)
					}
				}(i, t.Thumbnail)
			}
		}
		a.playlist.SetTitle(fmt.Sprintf(" Playlist [%d] ", count))
		
		// Atualiza índice do item tocando após remoção
		a.mu.Lock()
		currentIdx := a.currentTrack
		a.mu.Unlock()
		a.playlist.SetPlayingIndex(currentIdx)
		
		a.statusBar.SetText("[yellow]✓ Removido da playlist")
	})
}

// movePlaylistItem move um item na playlist
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
			// Adiciona item com thumbnail
			a.playlist.AddItem(t, i)
			
			// Carrega thumbnail em background (usa cache!)
			if t.Thumbnail != "" && a.thumbCache != nil {
				go func(idx int, url string) {
					img, err := a.thumbCache.GetThumbnailImage(url)
					if err == nil && img != nil {
						a.playlist.SetThumbnail(idx, img)
					}
				}(i, t.Thumbnail)
			}
		}
		a.playlist.SetCurrentIndex(newPos)
		
		// Atualiza índice do item tocando após mover
		a.mu.Lock()
		currentIdx := a.currentTrack
		a.mu.Unlock()
		a.playlist.SetPlayingIndex(currentIdx)
		
		a.statusBar.SetText("[cyan]✓ Item movido")
	})
}

// cycleRepeatMode alterna entre os modos de repetição
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
		a.updatePlayerInfo()
		a.updatePlaylistFooter()
		a.statusBar.SetText(fmt.Sprintf("[cyan]  Modo: %s", a.playlistMode.String()))
	})
}

// toggleShuffle alterna o modo shuffle
func (a *SimpleApp) toggleShuffle() {
	if a.playlistMode == ModeShuffle {
		a.playlistMode = ModeNormal
	} else {
		a.playlistMode = ModeShuffle
	}

	a.app.QueueUpdateDraw(func() {
		a.updatePlayerInfo()
		a.updatePlaylistFooter()
		a.statusBar.SetText(fmt.Sprintf("[cyan]  Modo: %s", a.playlistMode.String()))
	})
}
