package ui

// onResultSelectedCustom é chamado quando Enter é pressionado na CustomList
func (a *SimpleApp) onResultSelectedCustom(idx int) {
	track := a.searchResults.GetCurrentTrack()
	if track != nil {
		go a.playTrackDirect(*track)
	}
}
