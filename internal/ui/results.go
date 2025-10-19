package ui

// onResultSelected é chamado quando Enter é pressionado em um resultado
func (a *SimpleApp) onResultSelected(idx int, _ string, _ string, _ rune) {
	a.mu.Lock()
	if idx >= 0 && idx < len(a.tracks) {
		track := a.tracks[idx]
		a.mu.Unlock()
		go a.playTrackDirect(track)
	} else {
		a.mu.Unlock()
	}
}

// onResultChanged é chamado quando a seleção muda nos resultados
func (a *SimpleApp) onResultChanged(idx int, _ string, _ string, _ rune) {
	a.updateSearchDetailsDebounced(idx)
}
