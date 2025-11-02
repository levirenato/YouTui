package ui

import (
	"image"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CustomListItem struct {
	flex      *tview.Flex
	thumbnail *tview.Image
	info      *tview.TextView
	index     int
	track     Track
}

type CustomList struct {
	*tview.Flex
	container     *tview.Flex
	items         []*CustomListItem
	selectedIndex int
	playingIndex  int
	theme         *Theme
	mu            sync.Mutex
	onSelected    func(index int)
	visibleStart  int
	visibleHeight int
}

func NewCustomList(theme *Theme) *CustomList {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBackgroundColor(theme.Base)

	wrapper := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(container, 0, 1, false)
	wrapper.SetBackgroundColor(theme.Base)

	list := &CustomList{
		Flex:          wrapper,
		container:     container,
		items:         []*CustomListItem{},
		selectedIndex: 0,
		playingIndex:  -1,
		theme:         theme,
		visibleStart:  0,
		visibleHeight: 10,
	}

	wrapper.SetBorder(true).
		SetTitle(" ").
		SetBorderColor(theme.Surface0)

	wrapper.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			list.SelectPrevious()
			return nil
		case tcell.KeyDown:
			list.SelectNext()
			return nil
		case tcell.KeyEnter:
			if list.onSelected != nil {
				list.onSelected(list.selectedIndex)
			}
			return nil
		case tcell.KeyTab, tcell.KeyBacktab:
			return event
		}
		return event
	})

	list.renderVisibleItems()
	return list
}

func (c *CustomList) AddItem(track Track, index int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	thumb := tview.NewImage().
		SetColors(tview.TrueColor).
		SetDithering(tview.DitheringFloydSteinberg)
	thumb.SetBackgroundColor(c.theme.Base)

	info := tview.NewTextView().
		SetDynamicColors(true).
		SetText(formatItemInfo(track, index, c.theme)).
		SetTextAlign(tview.AlignLeft)
	info.SetBackgroundColor(c.theme.Base)
	info.SetTextColor(c.theme.Text)

	itemFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(thumb, 20, 0, false).
		AddItem(info, 0, 1, false)
	itemFlex.SetBackgroundColor(c.theme.Base)

	item := &CustomListItem{
		flex:      itemFlex,
		thumbnail: thumb,
		info:      info,
		index:     index,
		track:     track,
	}

	c.items = append(c.items, item)
	c.renderVisibleItems()
}

func (c *CustomList) SetThumbnail(index int, img image.Image) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if index >= 0 && index < len(c.items) {
		c.items[index].thumbnail.SetImage(img)
	}
}

func (c *CustomList) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.container.Clear()
	c.items = []*CustomListItem{}
	c.selectedIndex = 0
	c.visibleStart = 0
}

func (c *CustomList) renderVisibleItems() {
	c.container.Clear()

	_, _, _, wrapperHeight := c.GetInnerRect()

	availableHeight := wrapperHeight
	if availableHeight <= 0 {
		availableHeight = 10
	}

	const itemHeight = 3
	itemsPerPage := max(availableHeight/itemHeight, 1)

	c.visibleHeight = itemsPerPage

	end := min(c.visibleStart+c.visibleHeight, len(c.items))

	if len(c.items) == 0 {
		spacer := tview.NewBox().SetBackgroundColor(c.theme.Base)
		c.container.AddItem(spacer, 0, 1, false)
	} else {
		itemsRendered := 0
		for i := c.visibleStart; i < end; i++ {
			c.container.AddItem(c.items[i].flex, itemHeight, 0, false)
			itemsRendered++
		}

		remainingHeight := availableHeight - (itemsRendered * itemHeight)
		if remainingHeight > 0 {
			spacer := tview.NewBox().SetBackgroundColor(c.theme.Base)
			c.container.AddItem(spacer, remainingHeight, 0, false)
		}
	}

	c.updateSelection()
}

func (c *CustomList) scrollToSelection() {
	if c.selectedIndex >= c.visibleStart+c.visibleHeight {
		c.visibleStart = c.selectedIndex - c.visibleHeight + 1
		c.renderVisibleItems()
	}
	if c.selectedIndex < c.visibleStart {
		c.visibleStart = c.selectedIndex
		c.renderVisibleItems()
	}
}

func (c *CustomList) SelectNext() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.selectedIndex < len(c.items)-1 {
		c.selectedIndex++
		c.scrollToSelection()
		c.updateSelection()
	}
}

func (c *CustomList) SelectPrevious() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.selectedIndex > 0 {
		c.selectedIndex--
		c.scrollToSelection()
		c.updateSelection()
	}
}

func (c *CustomList) GetCurrentItem() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.selectedIndex
}

func (c *CustomList) SetCurrentIndex(idx int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if idx >= 0 && idx < len(c.items) {
		c.selectedIndex = idx
		c.scrollToSelection()
		c.updateSelection()
	}
}

func (c *CustomList) GetCurrentTrack() *Track {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.selectedIndex >= 0 && c.selectedIndex < len(c.items) {
		return &c.items[c.selectedIndex].track
	}
	return nil
}

func (c *CustomList) IsFocused(app *tview.Application) bool {
	focused := app.GetFocus()
	return focused == c.Flex
}

func (c *CustomList) SetSelectedFunc(handler func(int)) {
	c.onSelected = handler
}

func (c *CustomList) updateSelection() {
	for i, item := range c.items {
		switch i {
		case c.selectedIndex:
			item.flex.SetBackgroundColor(c.theme.Blue)
			item.info.SetTextColor(c.theme.Crust)
			item.info.SetBackgroundColor(c.theme.Blue)
			item.info.SetText(formatItemInfoPlain(item.track, item.index))
		case c.playingIndex:
			item.flex.SetBackgroundColor(c.theme.Green)
			item.info.SetTextColor(c.theme.Crust)
			item.info.SetBackgroundColor(c.theme.Green)
			item.info.SetText(formatItemInfoPlain(item.track, item.index))
		default:
			item.flex.SetBackgroundColor(c.theme.Base)
			item.info.SetTextColor(c.theme.Text)
			item.info.SetBackgroundColor(c.theme.Base)
			item.info.SetText(formatItemInfo(item.track, item.index, c.theme))
		}
	}
}

func (c *CustomList) SetTitle(title string) *CustomList {
	c.Flex.SetTitle(title)
	return c
}

func (c *CustomList) SetBorderColor(color tcell.Color) *CustomList {
	c.Flex.SetBorderColor(color)
	return c
}

func (c *CustomList) SetPlayingIndex(idx int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.playingIndex = idx
	c.updateSelection()
}

func formatItemInfo(track Track, index int, theme *Theme) string {
	icons := []string{"♪", "♫", "♬"}
	icon := icons[index%len(icons)]
	title := track.Title
	if len(title) > 50 {
		title = title[:47] + "..."
	}
	return icon + " [" + colorTag(theme.Yellow) + "::b]" + title + "[-:-:-]\n" +
		"[" + colorTag(theme.Green) + "]⏱ " + track.Duration + "[-] " +
		"[" + colorTag(theme.Sapphire) + "]• " + track.Author + "[-]"
}

func formatItemInfoPlain(track Track, index int) string {
	icons := []string{"♪", "♫", "♬"}
	icon := icons[index%len(icons)]
	title := track.Title
	if len(title) > 50 {
		title = title[:47] + "..."
	}
	return icon + " " + title + "\n" +
		"⏱ " + track.Duration + " • " + track.Author
}

func (c *CustomList) SetTheme(theme *Theme) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.theme = theme
	c.SetBorderColor(theme.Surface0)
	c.SetBackgroundColor(theme.Base)
	c.container.SetBackgroundColor(theme.Base)
	for _, item := range c.items {
		item.flex.SetBackgroundColor(theme.Base)
		item.thumbnail.SetBackgroundColor(theme.Base)
		item.info.SetBackgroundColor(theme.Base)
		item.info.SetTextColor(theme.Text)
	}
	c.renderVisibleItems()
}
