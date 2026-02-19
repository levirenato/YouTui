package ui

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
)

func (a *SimpleApp) setStatus(color tcell.Color, msg string) {
	a.statusBar.SetText("[" + colorTag(color) + "]" + msg)
}

func (a *SimpleApp) setStatusf(color tcell.Color, format string, args ...any) {
	a.statusBar.SetText(fmt.Sprintf("["+colorTag(color)+"]"+format, args...))
}

func (a *SimpleApp) playTrackSimple(track Track, idx int) {
	a.mu.Lock()
	if a.stopProgress != nil {
		select {
		case <-a.stopProgress:
		default:
			close(a.stopProgress)
		}
		a.stopProgress = nil
	}
	if a.mpvProcess != nil && a.mpvProcess.Process != nil {
		if err := a.mpvProcess.Process.Signal(os.Signal(syscall.Signal(0))); err == nil {
			if err := a.mpvProcess.Process.Kill(); err != nil {
				fmt.Printf("Warning %s %v", a.strings.warnFailedKillPrevMpv, err)
			}
		}
		a.mpvProcess = nil
	}
	a.mu.Unlock()

	socketPath := fmt.Sprintf("/tmp/mpv-socket-%d", time.Now().UnixNano())

	args := []string{
		"--no-terminal",
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

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		a.app.QueueUpdateDraw(func() {
			a.setStatusf(a.theme.Red, "❌%s %v", a.strings.errorStartMpv, err)
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
		a.setStatusf(a.theme.Green, "▶ %s: %s", a.strings.Playing, track.Title)
	})

	a.startProgressUpdater()

	go func(expectedCmd *exec.Cmd, trackCopy Track, idxCopy int) {
		err := expectedCmd.Wait()
		time.Sleep(500 * time.Millisecond)

		a.mu.Lock()
		if a.mpvProcess != expectedCmd {
			a.mu.Unlock()
			return
		}

		if err != nil {
			a.isPlaying = false
			a.mu.Unlock()

			stderrOutput := stderrBuf.String()
			if strings.Contains(stderrOutput, "403") || strings.Contains(stderrOutput, "HTTP error 403") {
				a.app.QueueUpdateDraw(func() {
					a.updatePlayerInfo()
					a.setStatus(a.theme.Red, "❌"+a.strings.youtubeBlocked)
				})
			} else {
				a.app.QueueUpdateDraw(func() {
					a.updatePlayerInfo()
					a.setStatusf(a.theme.Red, "❌"+a.strings.MpvError, err)
				})
			}
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
			nextTrack = trackCopy
			nextIdx = idxCopy

		case ModeShuffle:
			if len(a.playlistTracks) > 0 {
				if len(a.playlistTracks) == 1 {
					nextIdx = 0
				} else {
					for {
						nextIdx = rand.IntN(len(a.playlistTracks))
						if nextIdx != idxCopy || len(a.playlistTracks) == 1 {
							break
						}
					}
				}
				shouldPlayNext = true
				nextTrack = a.playlistTracks[nextIdx]
			}

		case ModeRepeatAll, ModeNormal:
			if len(a.playlistTracks) > 0 {
				next := idxCopy + 1
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
				a.setStatus(a.theme.Yellow, a.strings.PlaylistFinished)
			})
		}
	}(cmd, track, idx)
}

func (a *SimpleApp) playTrackDirect(track Track) {
	a.cleanup()

	socketPath := fmt.Sprintf("/tmp/mpv-socket-%d", time.Now().UnixNano())

	args := []string{
		"--no-terminal",
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

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		a.app.QueueUpdateDraw(func() {
			a.setStatusf(a.theme.Red, "❌%s %v", a.strings.errorStartMpv, err)
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
		a.setStatusf(a.theme.Green, "▶ "+a.strings.PlayingWithoutPlaylist, track.Title)
	})

	a.startProgressUpdater()

	go func(expectedCmd *exec.Cmd) {
		err := expectedCmd.Wait()

		a.mu.Lock()
		if a.mpvProcess != expectedCmd {
			a.mu.Unlock()
			return
		}
		a.isPlaying = false
		a.mu.Unlock()

		a.app.QueueUpdateDraw(func() {
			a.updatePlayerInfo()
			if err != nil {
				stderrOutput := stderrBuf.String()
				if strings.Contains(stderrOutput, "403") || strings.Contains(stderrOutput, "HTTP error 403") {
					a.setStatus(a.theme.Red, "❌"+a.strings.youtubeBlocked)
				} else {
					a.setStatusf(a.theme.Red, "❌ "+a.strings.MpvError, err)
				}
			} else {
				a.setStatus(a.theme.Yellow, a.strings.PlaybackFinished)
			}
		})
	}(cmd)
}

func (a *SimpleApp) togglePause() {
	a.mu.Lock()
	isPlaying := a.isPlaying
	socket := a.mpvSocket
	a.mu.Unlock()

	if !isPlaying || socket == "" {
		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Yellow, a.strings.nothingPlaying)
		})
		return
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '{ \"command\": [\"cycle\", \"pause\"] }' | socat - \"%s\" 2>&1", socket))
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), a.strings.socatCmdNotFound) {
			a.app.QueueUpdateDraw(func() {
				a.setStatus(a.theme.Red, "❌"+a.strings.socatNotInstalled)
			})
			return
		}

		a.app.QueueUpdateDraw(func() {
			a.setStatusf(a.theme.Red, "❌"+a.strings.errorPause+"%v", err)
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
			a.setStatus(a.theme.Yellow, "⏸ "+a.strings.Paused)
		} else {
			a.setStatus(a.theme.Green, "▶ "+a.strings.Playing)
		}
	})
}

func (a *SimpleApp) stopPlayback() {
	a.mu.Lock()
	if a.mpvProcess != nil && a.mpvProcess.Process != nil {
		if err := a.mpvProcess.Process.Kill(); err != nil {
			fmt.Print("Warning:", a.strings.warnFailedKillMpv, err)
		}
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
		a.setStatus(a.theme.Red, "⏹ "+a.strings.Stopped)
	})
}

func (a *SimpleApp) playNext() {
	a.mu.Lock()
	currentIsPlaying := a.isPlaying
	currentTrack := a.currentTrack
	playlistLen := len(a.playlistTracks)
	mode := a.playlistMode
	a.mu.Unlock()

	if playlistLen == 0 {
		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Yellow, "⚠ "+a.strings.PlaylistEmpty)
		})
		return
	}

	if !currentIsPlaying {
		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Yellow, "⚠ "+a.strings.NothingPlaying)
		})
		return
	}

	if currentTrack < 0 {
		a.mu.Lock()
		track := a.playlistTracks[0]
		a.skipAutoPlay = true
		a.mu.Unlock()

		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Sapphire, "▶ "+a.strings.EnteringPlaylist)
		})
		go a.playTrackSimple(track, 0)
		return
	}

	var next int
	if mode == ModeShuffle {
		if playlistLen == 1 {
			next = 0
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
			if mode == ModeRepeatAll {
				next = 0
			} else {
				a.app.QueueUpdateDraw(func() {
					a.setStatus(a.theme.Yellow, a.strings.AlreadyLastSong)
				})
				return
			}
		}
	}

	a.mu.Lock()
	track := a.playlistTracks[next]
	a.skipAutoPlay = true
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.setStatusf(a.theme.Green, "▶ "+a.strings.SkippingTo, next+1, playlistLen, track.Title)
	})

	go a.playTrackSimple(track, next)
}

func (a *SimpleApp) playPrevious() {
	a.mu.Lock()
	currentTrack := a.currentTrack
	playlistLen := len(a.playlistTracks)
	isPlaying := a.isPlaying
	mode := a.playlistMode
	a.mu.Unlock()

	if playlistLen == 0 {
		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Yellow, "⚠ "+a.strings.PlaylistEmpty)
		})
		return
	}

	if !isPlaying {
		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Yellow, "⚠ "+a.strings.NothingPlaying)
		})
		return
	}

	if currentTrack < 0 {
		lastIdx := playlistLen - 1
		a.mu.Lock()
		track := a.playlistTracks[lastIdx]
		a.skipAutoPlay = true
		a.mu.Unlock()

		a.app.QueueUpdateDraw(func() {
			a.setStatus(a.theme.Sapphire, "▶ "+a.strings.EnteringPlaylist)
		})
		go a.playTrackSimple(track, lastIdx)
		return
	}

	prev := currentTrack - 1
	if prev < 0 {
		if mode == ModeRepeatAll {
			prev = playlistLen - 1
		} else {
			a.app.QueueUpdateDraw(func() {
				a.setStatus(a.theme.Yellow, a.strings.AlreadyFirstSong)
			})
			return
		}
	}

	a.mu.Lock()
	track := a.playlistTracks[prev]
	a.skipAutoPlay = true
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.setStatusf(a.theme.Green, "⏮ Back to: %s", track.Title)
	})

	go a.playTrackSimple(track, prev)
}

func (a *SimpleApp) seekMedia(seconds float64) {
	a.mu.Lock()
	isPlaying := a.isPlaying
	socket := a.mpvSocket
	a.mu.Unlock()

	if !isPlaying || socket == "" {
		return
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf(`echo '{ "command": ["seek", %g, "relative"] }' | socat - "%s" 2>&1`, seconds, socket))
	cmd.Run()

	go func() {
		time.Sleep(150 * time.Millisecond)
		a.updateProgress()
	}()
}

func (a *SimpleApp) toggleMode() {
	a.mu.Lock()
	if a.playMode == ModeAudio {
		a.playMode = ModeVideo
	} else {
		a.playMode = ModeAudio
	}
	newMode := a.playMode
	a.mu.Unlock()

	a.app.QueueUpdateDraw(func() {
		a.updatePlayerInfo()
		a.updateModeBadge()
		a.setStatusf(a.theme.Sapphire, "  "+a.strings.ModeChanged, newMode.String())
	})
}
