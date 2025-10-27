package ui

func (a *SimpleApp) onResultSelectedCustom() {
	track := a.searchResults.GetCurrentTrack()
	if track != nil {
		go a.playTrackDirect(*track)
	}
}
