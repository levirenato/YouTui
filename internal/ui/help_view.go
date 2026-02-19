package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tabEntry struct {
	label string
	text  string
}

type HelpView struct {
	Flex      *tview.Flex
	tabBar    *tview.TextView
	content   *tview.TextView
	tabs      []tabEntry
	activeTab int
	theme     *Theme
	app       *tview.Application
	onClose   func()
}

func NewHelpView(s Strings, theme *Theme, app *tview.Application, onClose func()) *HelpView {
	hv := &HelpView{
		theme:   theme,
		app:     app,
		onClose: onClose,
	}
	hv.buildTabs(s)
	hv.build(s)
	return hv
}

func (hv *HelpView) buildTabs(s Strings) {
	type rawTab struct {
		label string
		text  string
	}
	sections := []rawTab{
		{s.HelpNavigation, s.HelpNavigationText},
		{s.HelpSearch, s.HelpSearchText},
		{s.HelpResults, s.HelpResultsText},
		{s.HelpPlaylist, s.HelpPlaylistText},
		{s.HelpPlayer, s.HelpPlayerText},
		{s.HelpGlobal, s.HelpGlobalText},
		{s.HelpIcons, s.HelpIconsText},
	}

	var allBuilder strings.Builder
	for _, sec := range sections {
		allBuilder.WriteString("## " + sec.label + "\n")
		allBuilder.WriteString(sec.text + "\n\n")
	}

	hv.tabs = make([]tabEntry, 0, len(sections)+1)
	hv.tabs = append(hv.tabs, tabEntry{label: s.HelpTabAll, text: allBuilder.String()})
	for _, sec := range sections {
		label := sec.label
		if idx := strings.Index(label, " ("); idx >= 0 {
			label = label[:idx]
		}
		hv.tabs = append(hv.tabs, tabEntry{
			label: label,
			text:  "## " + sec.label + "\n" + sec.text,
		})
	}
}

func (hv *HelpView) build(s Strings) {
	blue := colorTag(hv.theme.Blue)
	sub0 := colorTag(hv.theme.Subtext0)

	hv.tabBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	hv.tabBar.SetBackgroundColor(hv.theme.Surface0)

	hv.content = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)
	hv.content.SetBackgroundColor(hv.theme.Base)
	hv.content.SetTextColor(hv.theme.Text)

	tabsWord := "Abas"
	if s.HelpTabAll == "All" {
		tabsWord = "Tabs"
	}
	hintText := fmt.Sprintf(
		"[%s]Esc[-] %s  [%s]←/→ h/l[-] %s  [%s]j/k[-] Scroll",
		sub0, s.Close, blue, tabsWord, blue,
	)

	hintView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(hintText)
	hintView.SetBackgroundColor(hv.theme.Surface0)

	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(hv.tabBar, 1, 0, false).
		AddItem(hv.content, 0, 1, false).
		AddItem(hintView, 1, 0, false)

	innerFlex.SetBorder(true).
		SetTitle(" " + s.HelpTitle + " ").
		SetBorderColor(hv.theme.Blue).
		SetBackgroundColor(hv.theme.Base)

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox().SetBackgroundColor(hv.theme.Base), 0, 1, false).
		AddItem(innerFlex, 0, 5, true).
		AddItem(tview.NewBox().SetBackgroundColor(hv.theme.Base), 0, 1, false)

	hv.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBackgroundColor(hv.theme.Base), 0, 1, false).
		AddItem(centerFlex, 0, 4, true).
		AddItem(tview.NewBox().SetBackgroundColor(hv.theme.Base), 0, 1, false)

	hv.Flex.SetBackgroundColor(hv.theme.Base)

	hv.renderTabBar()
	hv.renderContent()

	hv.content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			hv.onClose()
			return nil
		case tcell.KeyLeft:
			hv.switchTab(-1)
			return nil
		case tcell.KeyRight:
			hv.switchTab(+1)
			return nil
		}
		switch event.Rune() {
		case 'h':
			hv.switchTab(-1)
			return nil
		case 'l':
			hv.switchTab(+1)
			return nil
		}
		return event
	})
}

func (hv *HelpView) renderTabBar() {
	blue := colorTag(hv.theme.Blue)
	sub := colorTag(hv.theme.Subtext0)
	base := colorTag(hv.theme.Surface0)

	var sb strings.Builder
	sb.WriteString(" ")
	for i, tab := range hv.tabs {
		if i == hv.activeTab {
			sb.WriteString(fmt.Sprintf("[%s:%s:b] %s [-:-:-]  ", blue, base, tab.label))
		} else {
			sb.WriteString(fmt.Sprintf("[%s] %s [-]  ", sub, tab.label))
		}
	}
	hv.tabBar.SetText(sb.String())
}

func (hv *HelpView) renderContent() {
	src := hv.tabs[hv.activeTab].text

	blue := colorTag(hv.theme.Blue)
	yellow := colorTag(hv.theme.Yellow)
	dim := colorTag(hv.theme.Subtext1)

	var sb strings.Builder
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(line, "## ") {
			title := strings.TrimPrefix(line, "## ")
			sb.WriteString(fmt.Sprintf("\n[%s::b]  %s  [-:-:-]\n", yellow, title))
			continue
		}
		if strings.HasPrefix(line, "  ") {
			trimmed := strings.TrimLeft(line, " ")
			if idx := strings.Index(trimmed, "  "); idx > 0 {
				key := trimmed[:idx]
				desc := strings.TrimLeft(trimmed[idx:], " ")
				sb.WriteString(fmt.Sprintf("  [%s]%-22s[-][%s]%s[-]\n", blue, key, dim, desc))
				continue
			}
		}
		if line == "" {
			sb.WriteString("\n")
			continue
		}
		sb.WriteString(line + "\n")
	}

	hv.content.SetText(sb.String())
	hv.content.ScrollToBeginning()
}

func (hv *HelpView) switchTab(delta int) {
	hv.activeTab = (hv.activeTab + delta + len(hv.tabs)) % len(hv.tabs)
	hv.renderTabBar()
	hv.renderContent()
}

func (hv *HelpView) FocusContent() {
	hv.app.SetFocus(hv.content)
}
