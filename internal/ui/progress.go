package ui

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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

	var author string
	if currentTrack >= 0 && currentTrack < len(a.playlistTracks) {
		author = a.playlistTracks[currentTrack].Author
	}
	a.mu.Unlock()

	a.mu.Lock()
	str := a.strings
	a.mu.Unlock()

	var titleLine string
	if isPlaying {
		if author != "" && author != str.Unknown {
			titleLine = fmt.Sprintf("[white::b] %s [gray]•[-] [#a6adc8]%s[-]", nowPlaying, author)
		} else {
			titleLine = fmt.Sprintf("[white::b] %s[-:-:-]", nowPlaying)
		}

		if currentTrack >= 0 && playlistLen > 0 {
			plainTitle := fmt.Sprintf(" %s", nowPlaying)
			if author != "" && author != str.Unknown {
				plainTitle += fmt.Sprintf(" • %s", author)
			}
			posText := fmt.Sprintf("[%d/%d]", currentTrack+1, playlistLen)
			padding := max(width-len(plainTitle)-len(posText)-2, 0)
			titleLine += strings.Repeat(" ", padding) + fmt.Sprintf("[#585b70]%s[-]", posText)
		}
	} else {
		titleLine = "[gray]⏹ " + str.NoTrackPlaying + "[-]"
	}

	var progressLine string
	if isPlaying && duration > 0 {
		icon := "▶"
		iconColor := "green"
		if isPaused {
			icon = "⏸"
			iconColor = "yellow"
		}

		percentage := position / duration
		if percentage > 1 {
			percentage = 1
		}
		if percentage < 0 {
			percentage = 0
		}

		totalBars := max(width-18, 20)
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

func (a *SimpleApp) updateModeBadge() {
	a.mu.Lock()
	mode := a.playMode
	strings := a.strings
	a.mu.Unlock()

	var badge string
	if mode == ModeVideo {
		badge = "[gray]m[-] [black:blue:b]   " + strings.Video + " [-:-:-] "
	} else {
		badge = "[gray]m[-] [black:green:b]   " + strings.Audio + " [-:-:-] "
	}

	a.modeBadge.SetText(badge)
}

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
