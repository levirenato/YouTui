package ui

// onResultSelectedCustom is called when Enter is pressed on CustomList
func (a *SimpleApp) onResultSelectedCustom(idx int) {
	track := a.searchResults.GetCurrentTrack()
	if track != nil {
		go a.playTrackDirect(*track)
	}
}
