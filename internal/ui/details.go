package ui

func (a *SimpleApp) updateThumbnail(thumbnailURL string) {
	if thumbnailURL == "" || a.thumbCache == nil {
		a.thumbnailView.SetImage(nil)
		return
	}

	go func() {
		img, err := a.thumbCache.GetThumbnailImage(thumbnailURL)
		if err != nil {
			return
		}

		a.app.QueueUpdateDraw(func() {
			a.thumbnailView.SetImage(img)
		})
	}()
}
