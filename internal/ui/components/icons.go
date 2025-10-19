package components

// GetTrackIcon retorna um ícone musical para um item de track
// Usa diferentes caracteres baseado no índice para variedade visual
func GetTrackIcon(idx int) string {
	icons := []string{"♪", "♫", "♬"}
	return icons[idx%len(icons)]
}

// GetPlaylistIcon retorna um ícone musical para um item da playlist
// Usa uma paleta maior de ícones para mais variedade
func GetPlaylistIcon(idx int) string {
	icons := []string{"♪", "♫", "♬", "♩", "▸", "•"}
	return icons[idx%len(icons)]
}
