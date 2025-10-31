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

	a.playerBox.SetBackgroundColor(a.theme.Base)
	a.thumbnailView.SetBackgroundColor(a.theme.Base)
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
			titleLine = fmt.Sprintf("["+colorTag(a.theme.Text)+"::b] %s ["+colorTag(a.theme.Subtext0)+"]•[-] ["+colorTag(a.theme.Subtext1)+"]%s[-]", nowPlaying, author)
		} else {
			titleLine = fmt.Sprintf("["+colorTag(a.theme.Text)+"::b] %s[-:-:-]", nowPlaying)
		}

		if currentTrack >= 0 && playlistLen > 0 {
			plainTitle := fmt.Sprintf(" %s", nowPlaying)
			if author != "" && author != str.Unknown {
				plainTitle += fmt.Sprintf(" • %s", author)
			}
			posText := fmt.Sprintf("[%d/%d]", currentTrack+1, playlistLen)
			padding := max(width-len(plainTitle)-len(posText)-2, 0)
			titleLine += strings.Repeat(" ", padding) + fmt.Sprintf("["+colorTag(a.theme.Surface2)+"]%s[-]", posText)
		}
	} else {
		titleLine = "[" + colorTag(a.theme.Subtext0) + "]⏹ " + str.NoTrackPlaying + "[-]"
	}

	var progressLine string
	if isPlaying && duration > 0 {
		icon := "▶"
		iconColor := colorTag(a.theme.Green)
		if isPaused {
			icon = "⏸"
			iconColor = colorTag(a.theme.Yellow)
		}

		percentage := position / duration
		if percentage > 1 {
			percentage = 1
		}
		if percentage < 0 {
			percentage = 0
		}

		totalBars := min(max(width-18, 20), 100)

		filledBars := int(percentage * float64(totalBars))
		emptyBars := totalBars - filledBars

		posMin := int(position / 60)
		posSec := int(position) % 60
		durMin := int(duration / 60)
		durSec := int(duration) % 60

		progressLine = fmt.Sprintf("[%s]%s[-] ["+colorTag(a.theme.Blue)+"]%s["+colorTag(a.theme.Surface1)+"]%s[-]  ["+colorTag(a.theme.Sapphire)+"]%02d:%02d[-] ["+colorTag(a.theme.Surface2)+"]/[-] ["+colorTag(a.theme.Text)+"]%02d:%02d[-]",
			iconColor, icon,
			strings.Repeat("█", filledBars),
			strings.Repeat("░", emptyBars),
			posMin, posSec,
			durMin, durSec)
	} else {
		totalBars := min(max(width-18, 20), 100)
		progressLine = fmt.Sprintf("["+colorTag(a.theme.Subtext0)+"]⏹ %s  --:-- / --:--[-]", strings.Repeat("░", totalBars))
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
		badge = "[" + colorTag(a.theme.Subtext0) + "]m[-] [" + colorTag(a.theme.Crust) + ":" + colorTag(a.theme.Blue) + ":b] 󰗃 " + strings.Video + " [-:-:-] "
	} else {
		badge = "[" + colorTag(a.theme.Subtext0) + "]m[-] [" + colorTag(a.theme.Crust) + ":" + colorTag(a.theme.Green) + ":b]  " + strings.Audio + " [-:-:-] "
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
		footer = "[" + colorTag(a.theme.Teal) + "]  " + strings.Shuffle + "[-]"
	case ModeRepeatOne:
		footer = "[" + colorTag(a.theme.Peach) + "]󰑘 " + strings.RepeatOne + "[-]"
	case ModeRepeatAll:
		footer = "[" + colorTag(a.theme.Green) + "]󰑖 " + strings.RepeatAll + "[-]"
	default:
		footer = "[" + colorTag(a.theme.Surface2) + "]󰑗 " + strings.NoRepeat + "[-]"
	}

	a.playlistFooter.SetText(footer)
}
