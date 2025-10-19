package components

import (
	"fmt"
	"image"
)

// ItemData contém os dados de um item para exibição
type ItemData struct {
	Title     string
	Author    string
	Duration  string
	Thumbnail image.Image
	Icon      string
}

// FormatItemText formata o texto de um item expandido
// Retorna o texto formatado para exibir título, autor e duração
func FormatItemText(title, author, duration string, idx int) string {
	icon := GetTrackIcon(idx)
	
	// Limita o tamanho do título para não quebrar o layout
	maxTitleLen := 60
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-3] + "..."
	}
	
	// Limita o tamanho do autor
	maxAuthorLen := 40
	if len(author) > maxAuthorLen {
		author = author[:maxAuthorLen-3] + "..."
	}
	
	return fmt.Sprintf("%s [yellow::b]%s[-:-:-]\n   [cyan]%s[-] • [green]%s[-]", 
		icon, title, author, duration)
}

// FormatItemWithoutColor formata o item sem cores (para fallback)
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
