package ui

func (a *SimpleApp) onResultSelectedCustom(idx int) {
	track := a.searchResults.GetCurrentTrack()
	if track != nil {
		go a.playTrackDirect(*track)
	}
}
