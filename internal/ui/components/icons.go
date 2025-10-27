// Package components
package components

func GetTrackIcon(idx int) string {
	icons := []string{"♪", "♫", "♬"}
	return icons[idx%len(icons)]
}

func GetPlaylistIcon(idx int) string {
	icons := []string{"♪", "♫", "♬", "♩", "▸", "•"}
	return icons[idx%len(icons)]
}
