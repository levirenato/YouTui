package ui

import (
	"context"
	"fmt"
	"image"
	"time"
)

func (a *SimpleApp) updateSearchDetailsDebounced(idx int) {
	a.detailsLoadingMutex.Lock()

	if a.detailsDebounceTimer != nil {
		a.detailsDebounceTimer.Stop()
	}

	a.detailsDebounceTimer = time.AfterFunc(150*time.Millisecond, func() {
		a.updateSearchDetails(idx)
	})

	a.detailsLoadingMutex.Unlock()
}

func (a *SimpleApp) updateSearchDetails(idx int) {
	a.detailsLoadingMutex.Lock()
	if a.detailsCancelFunc != nil {
		a.detailsCancelFunc()
		a.detailsCancelFunc = nil
	}
	a.detailsLoadingIdx = idx
	a.detailsLoadingMutex.Unlock()

	a.mu.Lock()
	if idx < 0 || idx >= len(a.tracks) {
		a.mu.Unlock()
		a.app.QueueUpdateDraw(func() {
			a.detailsText.SetText("")
			a.detailsThumb.SetImage(nil)
		})
		return
	}

	track := a.tracks[idx]
	title := track.Title
	author := track.Author
	duration := track.Duration
	thumbnailURL := track.Thumbnail
	a.mu.Unlock()

	if title == "" {
		title = a.strings.NoTitle
	}
	if author == "" {
		author = a.strings.Unknown
	}
	if duration == "" {
		duration = "--:--"
	}

	basicDetails := fmt.Sprintf(
		"[yellow::b]%s[-:-:-]\n[cyan]%s:[-] %s\n[green]%s:[-] %s\n\n[gray]%s[-]",
		title,
		a.strings.Channel,
		author,
		a.strings.Duration,
		duration,
		a.strings.PressEnterToPlay,
	)

	a.app.QueueUpdateDraw(func() {
		a.detailsText.SetText(basicDetails)
	})

	if thumbnailURL != "" && a.thumbCache != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		a.detailsLoadingMutex.Lock()
		a.detailsCancelFunc = cancel
		a.detailsLoadingMutex.Unlock()

		go func(url string, ctx context.Context) {
			type result struct {
				img image.Image
				err error
			}
			resultChan := make(chan result, 1)

			go func() {
				img, err := a.thumbCache.GetThumbnailImageWithContext(ctx, url)
				select {
				case resultChan <- result{img: img, err: err}:
				case <-ctx.Done():
				}
			}()

			select {
			case res := <-resultChan:
				if res.err == nil && res.img != nil {
					a.detailsLoadingMutex.Lock()
					currentIdx := a.detailsLoadingIdx
					a.detailsLoadingMutex.Unlock()

					if currentIdx == idx {
						a.app.QueueUpdateDraw(func() {
							a.detailsThumb.SetImage(res.img)
						})
					}
				}
			case <-ctx.Done():
				return
			}

			a.detailsLoadingMutex.Lock()
			a.detailsCancelFunc = nil
			a.detailsLoadingMutex.Unlock()
		}(thumbnailURL, ctx)
	} else {
		a.app.QueueUpdateDraw(func() {
			a.detailsThumb.SetImage(nil)
		})
	}
}

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
