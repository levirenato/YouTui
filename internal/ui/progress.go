package ui

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// startProgressUpdater inicia o atualizador de progresso
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

// updateProgress atualiza a posição e duração do player
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
		a.updatePlayerInfo()
	})
}

// updatePlayerInfo atualiza a informação exibida no player
func (a *SimpleApp) updatePlayerInfo() {
	_, _, width, _ := a.playerInfo.GetInnerRect()
	if width <= 0 {
		width = 80
	}

	a.mu.Lock()
	isPlaying := a.isPlaying
	isPaused := a.isPaused
	nowPlaying := a.nowPlaying
	currentTrack := a.currentTrack
	playlistLen := len(a.playlistTracks)
	position := a.position
	duration := a.duration
	
	// Pega autor se disponível
	var author string
	if currentTrack >= 0 && currentTrack < len(a.playlistTracks) {
		author = a.playlistTracks[currentTrack].Author
	}
	a.mu.Unlock()

	// Pega strings traduzidas
	a.mu.Lock()
	str := a.strings
	a.mu.Unlock()
	
	// Linha 1: Título • Autor + Posição [3/10]
	var titleLine string
	if isPlaying {
		// Monta título com autor se disponível
		if author != "" && author != "Desconhecido" {
			titleLine = fmt.Sprintf("[white::b] %s [gray]•[-] [#a6adc8]%s[-]", nowPlaying, author)
		} else {
			titleLine = fmt.Sprintf("[white::b] %s[-:-:-]", nowPlaying)
		}
		
		// Adiciona posição na playlist se tocando da playlist
		if currentTrack >= 0 && playlistLen > 0 {
			// Calcula padding para alinhar à direita
			plainTitle := fmt.Sprintf(" %s", nowPlaying)
			if author != "" && author != "Desconhecido" {
				plainTitle += fmt.Sprintf(" • %s", author)
			}
			posText := fmt.Sprintf("[%d/%d]", currentTrack+1, playlistLen)
			padding := width - len(plainTitle) - len(posText) - 2
			if padding < 0 {
				padding = 0
			}
			titleLine += strings.Repeat(" ", padding) + fmt.Sprintf("[#585b70]%s[-]", posText)
		}
	} else {
		titleLine = "[gray]⏹ " + str.NoTrackPlaying + "[-]"
	}

	// Linha 2: Ícone + Barra + Tempo
	var progressLine string
	if isPlaying && duration > 0 {
		// Ícone play/pause
		icon := "▶"
		iconColor := "green"
		if isPaused {
			icon = "⏸"
			iconColor = "yellow"
		}

		// Calcula barra de progresso
		percentage := position / duration
		if percentage > 1 {
			percentage = 1
		}
		if percentage < 0 {
			percentage = 0
		}

		totalBars := width - 18 // Espaço para ícone + tempo
		if totalBars < 20 {
			totalBars = 20
		}
		if totalBars > 100 {
			totalBars = 100
		}

		filledBars := int(percentage * float64(totalBars))
		emptyBars := totalBars - filledBars

		posMin := int(position / 60)
		posSec := int(position) % 60
		durMin := int(duration / 60)
		durSec := int(duration) % 60

		progressLine = fmt.Sprintf("[%s]%s[-] [blue]%s[gray]%s[-]  [cyan]%02d:%02d[-] [#585b70]/[-] [white]%02d:%02d[-]",
			iconColor, icon,
			strings.Repeat("█", filledBars),
			strings.Repeat("░", emptyBars),
			posMin, posSec,
			durMin, durSec)
	} else {
		totalBars := width - 18
		if totalBars < 20 {
			totalBars = 20
		}
		if totalBars > 100 {
			totalBars = 100
		}
		progressLine = fmt.Sprintf("[gray]⏹ %s  --:-- / --:--[-]", strings.Repeat("░", totalBars))
	}

	a.playerInfo.SetText(fmt.Sprintf("%s\n%s", titleLine, progressLine))
}

// updateModeBadge atualiza a badge de modo estilo Neovim
func (a *SimpleApp) updateModeBadge() {
	a.mu.Lock()
	mode := a.playMode
	strings := a.strings
	a.mu.Unlock()

	var badge string
	if mode == ModeVideo {
		// Badge azul para vídeo (estilo VISUAL do Neovim)
		badge = "[gray]m[-] [black:blue:b]   " + strings.Video + " [-:-:-] "
	} else {
		// Badge verde para áudio (estilo INSERT do Neovim)
		badge = "[gray]m[-] [black:green:b]   " + strings.Audio + " [-:-:-] "
	}

	a.modeBadge.SetText(badge)
}

// updatePlaylistFooter atualiza o footer da playlist com o ícone do modo
func (a *SimpleApp) updatePlaylistFooter() {
	a.mu.Lock()
	mode := a.playlistMode
	strings := a.strings
	a.mu.Unlock()

	var footer string
	switch mode {
	case ModeShuffle:
		footer = "[#94e2d5]  " + strings.Shuffle + "[-]"
	case ModeRepeatOne:
		footer = "[#fab387]󰑘 " + strings.RepeatOne + "[-]"
	case ModeRepeatAll:
		footer = "[#a6e3a1]󰑖 " + strings.RepeatAll + "[-]"
	default:
		footer = "[#585b70]󰑗 " + strings.NoRepeat + "[-]"
	}

	a.playlistFooter.SetText(footer)
}
