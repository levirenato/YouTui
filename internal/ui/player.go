package ui

import (
	"fmt"
	"math/rand/v2"
	"os/exec"
	"time"
)

// playTrackSimple reproduz uma faixa da playlist
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
		a.playlist.SetPlayingIndex(idx) // Destaca o item tocando em verde
		a.statusBar.SetText(fmt.Sprintf("[green]▶ Tocando: %s", track.Title))
	})

	a.startProgressUpdater()

	go func(expectedCmd *exec.Cmd) {
		expectedCmd.Wait()

		time.Sleep(500 * time.Millisecond)

		a.mu.Lock()
		// Verifica se ainda é o processo atual antes de fazer auto-play
		if a.mpvProcess != expectedCmd {
			// Processo foi substituído, ignora este callback
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
			// Modo shuffle: toca música aleatória (diferente da atual)
			if len(a.playlistTracks) > 0 {
				if len(a.playlistTracks) == 1 {
					// Se só tem 1 música, repete ela
					nextIdx = 0
				} else {
					// Escolhe uma música diferente da atual
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
				a.statusBar.SetText("[yellow]Playlist finalizada")
			})
		}
	}(cmd)
}

// playTrackDirect reproduz uma faixa diretamente (sem playlist)
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
		a.playlist.SetPlayingIndex(-1) // Limpa destaque (não está tocando da playlist)
		a.statusBar.SetText(fmt.Sprintf("[green]▶ Tocando: %s (sem playlist)", track.Title))
	})

	a.startProgressUpdater()

	go func(expectedCmd *exec.Cmd) {
		expectedCmd.Wait()

		// Verifica se ainda é o processo atual antes de atualizar estado
		a.mu.Lock()
		if a.mpvProcess != expectedCmd {
			// Processo foi substituído, ignora este callback
			a.mu.Unlock()
			return
		}
		a.isPlaying = false
		a.mu.Unlock()

		a.app.QueueUpdateDraw(func() {
			a.updatePlayerInfo()
			a.statusBar.SetText("[yellow]Reprodução finalizada")
		})
	}(cmd)
}

// togglePause alterna entre pausar e retomar
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
		a.updatePlayerInfo()
		if isPaused {
			a.statusBar.SetText("[yellow]⏸ Pausado")
		} else {
			a.statusBar.SetText("[green]▶ Tocando")
		}
	})
}

// stopPlayback para a reprodução completamente
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
		a.playlist.SetPlayingIndex(-1) // Limpa destaque ao parar
		a.statusBar.SetText("[red]⏹ Parado")
	})
}

// playNext pula para a próxima faixa
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

	// Se está tocando em modo direto (currentTrack < 0), começa do início da playlist
	if currentTrack < 0 {
		track := a.playlistTracks[0]
		a.skipAutoPlay = true
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[cyan]▶ Entrando na playlist...")
		})
		go a.playTrackSimple(track, 0)
		return
	}

	var next int

	if a.playlistMode == ModeShuffle {
		// Shuffle: escolhe música diferente da atual
		if playlistLen == 1 {
			next = 0 // Só tem 1 música
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

// playPrevious volta para a faixa anterior
func (a *SimpleApp) playPrevious() {
	a.mu.Lock()
	playlistLen := len(a.playlistTracks)
	
	if playlistLen == 0 {
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

	// Se está tocando em modo direto (currentTrack < 0), começa do final da playlist
	if a.currentTrack < 0 {
		lastIdx := playlistLen - 1
		track := a.playlistTracks[lastIdx]
		a.skipAutoPlay = true
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[cyan]▶ Entrando na playlist...")
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

// toggleMode alterna entre modo áudio e vídeo
func (a *SimpleApp) toggleMode() {
	if a.playMode == ModeAudio {
		a.playMode = ModeVideo
	} else {
		a.playMode = ModeAudio
	}

	a.app.QueueUpdateDraw(func() {
		a.updatePlayerInfo()
		a.updateModeBadge()
		a.statusBar.SetText(fmt.Sprintf("[cyan]  Modo: %s", a.playMode.String()))
	})
}
