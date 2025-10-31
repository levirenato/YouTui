package ui

import (
	"fmt"
)

func (a *SimpleApp) onPlaylistSelectedCustom() {
	track := a.playlist.GetCurrentTrack()
	if track != nil {
		a.mu.Lock()
		realIdx := a.playlist.GetCurrentItem()
		a.mu.Unlock()
		go a.playTrackSimple(*track, realIdx)
	}
}

func (a *SimpleApp) addToPlaylist(track Track) {
	a.mu.Lock()
	a.playlistTracks = append(a.playlistTracks, track)
	count := len(a.playlistTracks)
	index := count - 1
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.playlist.AddItem(track, index)

		if track.Thumbnail != "" && a.thumbCache != nil {
			go func(idx int, url string) {
				img, err := a.thumbCache.GetThumbnailImage(url)
				if err == nil && img != nil {
					a.playlist.SetThumbnail(idx, img)
				}
			}(index, track.Thumbnail)
		}

		a.playlist.SetTitle(fmt.Sprintf(" Playlist [%d] ", count))
		a.setStatus(a.theme.Green, "✓ "+fmt.Sprintf(a.strings.AddedToPlaylist, track.Title))
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
			a.playlist.AddItem(t, i)

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

		a.mu.Lock()
		currentIdx := a.currentTrack
		a.mu.Unlock()
		a.playlist.SetPlayingIndex(currentIdx)

		a.setStatus(a.theme.Yellow, "✓ "+a.strings.RemovedFromPlaylist)
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

	switch a.currentTrack {
	case from:
		a.currentTrack = to
	case to:
		a.currentTrack = from
	}

	tracks := make([]Track, len(a.playlistTracks))
	copy(tracks, a.playlistTracks)
	newPos := to
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.playlist.Clear()
		for i, t := range tracks {
			a.playlist.AddItem(t, i)

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

		a.mu.Lock()
		currentIdx := a.currentTrack
		a.mu.Unlock()
		a.playlist.SetPlayingIndex(currentIdx)

		a.setStatus(a.theme.Sapphire, "✓ "+a.strings.ItemMoved)
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
		a.updatePlayerInfo()
		a.updatePlaylistFooter()
		a.setStatus(a.theme.Sapphire, "  "+fmt.Sprintf(a.strings.ModeChanged, a.playlistMode.String()))
	})
}

func (a *SimpleApp) toggleShuffle() {
	if a.playlistMode == ModeShuffle {
		a.playlistMode = ModeNormal
	} else {
		a.playlistMode = ModeShuffle
	}

	a.app.QueueUpdateDraw(func() {
		a.updatePlayerInfo()
		a.updatePlaylistFooter()
		a.setStatus(a.theme.Sapphire, "  "+fmt.Sprintf(a.strings.ModeChanged, a.playlistMode.String()))
	})
}
