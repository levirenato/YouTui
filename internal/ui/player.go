package ui

import (
	"fmt"
	"math/rand/v2"
	"os/exec"
	"time"
)

// playTrackSimple plays a track from the playlist
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
	a.currentThumb = track.Thumbnail
	a.currentTrack = idx
	a.position = 0
	a.duration = 0
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.updatePlayerInfo()
		a.updateThumbnail(track.Thumbnail)
		a.playlist.SetPlayingIndex(idx)
		a.statusBar.SetText(fmt.Sprintf("[green]▶ " + a.strings.Playing + ": %s", track.Title))
	})

	a.startProgressUpdater()

	go func(expectedCmd *exec.Cmd) {
		expectedCmd.Wait()

		time.Sleep(500 * time.Millisecond)

		a.mu.Lock()
		// Check if still the current process before auto-play
		if a.mpvProcess != expectedCmd {
			// Process was replaced, ignore this callback
			a.mu.Unlock()
			return
		}
		
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

		case ModeShuffle:
			// Shuffle mode: play random song (different from current)
			if len(a.playlistTracks) > 0 {
				if len(a.playlistTracks) == 1 {
					// If only 1 song, repeat it
					nextIdx = 0
				} else {
					// Check if should auto-skip to next
					for {
						nextIdx = rand.IntN(len(a.playlistTracks))
						if nextIdx != idx {
							break
						}
					}
				}
				shouldPlayNext = true
				nextTrack = a.playlistTracks[nextIdx]
			}

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
				a.updatePlayerInfo()
				a.statusBar.SetText("[yellow]" + a.strings.PlaylistFinished)
			})
		}
	}(cmd)
}

// playTrackDirect plays a track directly from results (without playlist)
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
	a.currentThumb = track.Thumbnail
	a.currentTrack = -1
	a.position = 0
	a.duration = 0
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.updatePlayerInfo()
		a.updateThumbnail(track.Thumbnail)
		a.playlist.SetPlayingIndex(-1)
		a.statusBar.SetText(fmt.Sprintf("[green]▶ " + a.strings.PlayingWithoutPlaylist, track.Title))
	})

	a.startProgressUpdater()

	go func(expectedCmd *exec.Cmd) {
		expectedCmd.Wait()

		// Check if still the current process before updating state
		a.mu.Lock()
		if a.mpvProcess != expectedCmd {
			// Process was replaced, ignore this callback
			a.mu.Unlock()
			return
		}
		a.isPlaying = false
		a.mu.Unlock()

		a.app.QueueUpdateDraw(func() {
			a.updatePlayerInfo()
			a.statusBar.SetText("[yellow]" + a.strings.PlaybackFinished)
		})
	}(cmd)
}

// togglePause toggles pause/play
func (a *SimpleApp) togglePause() {
	a.mu.Lock()
	isPlaying := a.isPlaying
	socket := a.mpvSocket
	a.mu.Unlock()

	if !isPlaying || socket == "" {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[red]⚠ " + a.strings.StateError, isPlaying, socket))
		})
		return
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf(`echo '{ "command": ["cycle", "pause"] }' | socat - "%s" 2>&1`, socket))
	output, err := cmd.CombinedOutput()
	if err != nil {
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText(fmt.Sprintf("[red]❌ " + a.strings.Error, err, string(output)))
		})
		return
	}

	a.mu.Lock()
	a.isPaused = !a.isPaused
	isPaused := a.isPaused
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.updatePlayerInfo()
		if isPaused {
			a.statusBar.SetText("[yellow]⏸ " + a.strings.Paused)
		} else {
			a.statusBar.SetText("[green]▶ " + a.strings.Playing)
		}
	})
}

// stopPlayback stops playback completely
func (a *SimpleApp) stopPlayback() {
	a.mu.Lock()
	if a.mpvProcess != nil && a.mpvProcess.Process != nil {
		a.mpvProcess.Process.Kill()
		a.mpvProcess = nil
	}
	a.isPlaying = false
	a.isPaused = false
	a.currentTrack = -1
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.updatePlayerInfo()
		a.updateThumbnail("")
		a.playlist.SetPlayingIndex(-1)
		a.statusBar.SetText("[red]⏹ " + a.strings.Stopped)
	})
}

// playNext skips to the next track
func (a *SimpleApp) playNext() {
	a.mu.Lock()
	currentIsPlaying := a.isPlaying
	currentTrack := a.currentTrack
	playlistLen := len(a.playlistTracks)

	if playlistLen == 0 {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ " + a.strings.PlaylistEmpty)
		})
		return
	}

	if !currentIsPlaying {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ " + a.strings.NothingPlaying)
		})
		return
	}

	// If playing in direct mode (currentTrack < 0), start from beginning of playlist
	if currentTrack < 0 {
		track := a.playlistTracks[0]
		a.skipAutoPlay = true
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[cyan]▶ " + a.strings.EnteringPlaylist)
		})
		go a.playTrackSimple(track, 0)
		return
	}

	var next int

	if a.playlistMode == ModeShuffle {
		// Shuffle: choose different song from current
		if playlistLen == 1 {
			next = 0 // Only 1 song
		} else {
			for {
				next = rand.IntN(playlistLen)
				if next != currentTrack {
					break
				}
			}
		}
	} else {
		next = currentTrack + 1
		if next >= playlistLen {
			if a.playlistMode == ModeRepeatAll {
				next = 0
			} else {
				a.mu.Unlock()
				a.app.QueueUpdateDraw(func() {
					a.statusBar.SetText("[yellow]" + a.strings.AlreadyLastSong)
				})
				return
			}
		}
	}

	track := a.playlistTracks[next]
	a.skipAutoPlay = true
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.statusBar.SetText(fmt.Sprintf("[green]▶ " + a.strings.SkippingTo, next+1, playlistLen, track.Title))
	})

	go a.playTrackSimple(track, next)
}

// playPrevious goes to previous track
func (a *SimpleApp) playPrevious() {
	a.mu.Lock()
	playlistLen := len(a.playlistTracks)
	
	if playlistLen == 0 {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ " + a.strings.PlaylistEmpty)
		})
		return
	}

	if !a.isPlaying {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]⚠ " + a.strings.NothingPlaying)
		})
		return
	}

	// If playing in direct mode (currentTrack < 0), start from end of playlist
	if a.currentTrack < 0 {
		lastIdx := playlistLen - 1
		track := a.playlistTracks[lastIdx]
		a.skipAutoPlay = true
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[cyan]▶ " + a.strings.EnteringPlaylist)
		})
		go a.playTrackSimple(track, lastIdx)
		return
	}

	prev := a.currentTrack - 1
	if prev < 0 {
		if a.playlistMode == ModeRepeatAll {
			prev = len(a.playlistTracks) - 1
		} else {
			a.mu.Unlock()
			a.app.QueueUpdateDraw(func() {
				a.statusBar.SetText("[yellow]" + a.strings.AlreadyFirstSong)
			})
			return
		}
	}
	track := a.playlistTracks[prev]
	a.skipAutoPlay = true
	a.mu.Unlock()

	go a.playTrackSimple(track, prev)
}

// toggleMode toggles between audio and video mode
func (a *SimpleApp) toggleMode() {
	if a.playMode == ModeAudio {
		a.playMode = ModeVideo
	} else {
		a.playMode = ModeAudio
	}

	a.app.QueueUpdateDraw(func() {
		a.updatePlayerInfo()
		a.updateModeBadge()
		a.statusBar.SetText(fmt.Sprintf("[cyan]  " + a.strings.ModeChanged, a.playMode.String()))
	})
}
