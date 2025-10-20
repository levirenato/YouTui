package components

import (
	"fmt"
	"image"
)

type ItemData struct {
	Title     string
	Author    string
	Duration  string
	Thumbnail image.Image
	Icon      string
}

func FormatItemText(title, author, duration string, idx int) string {
	icon := GetTrackIcon(idx)

	maxTitleLen := 60
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-3] + "..."
	}

	maxAuthorLen := 40
	if len(author) > maxAuthorLen {
		author = author[:maxAuthorLen-3] + "..."
	}

	return fmt.Sprintf("%s [yellow::b]%s[-:-:-]\n   [cyan]%s[-] • [green]%s[-]",
		icon, title, author, duration)
}

func FormatItemWithoutColor(title, author, duration string, idx int) string {
	icon := GetTrackIcon(idx)

	maxTitleLen := 60
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-3] + "..."
	}

	maxAuthorLen := 40
	if len(author) > maxAuthorLen {
		author = author[:maxAuthorLen-3] + "..."
	}

	return fmt.Sprintf("%s %s\n   %s • %s", icon, title, author, duration)
}
