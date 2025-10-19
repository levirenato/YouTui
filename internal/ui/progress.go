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

	var status string
	if a.isPlaying {
		icon := "▶"
		if a.isPaused {
			icon = "⏸"
		}
		status = fmt.Sprintf("[green::b]%s[-:-:-] [white]%s[-]", icon, a.nowPlaying)
	} else {
		status = "[gray]⏹ Nenhuma faixa tocando[-]"
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

		totalBars := width - 20
		if totalBars < 20 {
			totalBars = 20
		}
		if totalBars > 100 {
			totalBars = 100
		}

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
		totalBars := width - 20
		if totalBars < 20 {
			totalBars = 20
		}
		if totalBars > 100 {
			totalBars = 100
		}
		progress = "[gray]" + strings.Repeat("░", totalBars) + " --:--/--:--[-]"
	}

	modeInfo := fmt.Sprintf("[cyan]%s[-] [yellow]|[-] [magenta]%s[-]", a.playMode.String(), a.playlistMode.String())

	a.playerInfo.SetText(fmt.Sprintf("%s\n%s\n%s", status, progress, modeInfo))
}
