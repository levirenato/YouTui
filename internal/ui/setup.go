package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/levirenato/YouTui/internal/config"
	"github.com/levirenato/YouTui/internal/search"
	"github.com/rivo/tview"
)

func (a *SimpleApp) setupUI() {
	a.initLanguageFromConfig()

	a.setupSearchComponents()
	a.setupPlaylistComponent()
	a.setupDetailsComponent()
	a.setupPlayerComponents()
	a.setupStatusBars()
	a.setupHelpModal()
	a.setupConfigModal()
	a.setupLayout()
	a.setupInputHandlers()
}

func (a *SimpleApp) setupSearchComponents() {
	a.searchInput = tview.NewInputField().
		SetLabel(" ").
		SetFieldBackgroundColor(a.theme.Surface0).
		SetFieldTextColor(a.theme.Text)

	a.searchInput.SetBorder(true).
		SetTitle(" " + a.strings.Search + " ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(a.theme.Blue)

	a.searchInput.SetDoneFunc(a.onSearchDone)

	a.searchResults = NewCustomList(a.theme)
	a.searchResults.SetTitle(" " + a.strings.Results + " [0] ")
	a.searchResults.SetSelectedFunc(func(idx int) {
		a.onResultSelectedCustom()
	})
}

func (a *SimpleApp) setupPlaylistComponent() {
	a.playlist = NewCustomList(a.theme)
	a.playlist.SetTitle(fmt.Sprintf(" %s [0] ", a.strings.Playlist))
	a.playlist.SetSelectedFunc(func(idx int) {
		a.onPlaylistSelectedCustom()
	})

	a.playlistFooter = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(a.theme.Subtext0)
	a.playlistFooter.SetBackgroundColor(a.theme.Base)
	a.updatePlaylistFooter()

	playlistContainer := a.playlist.GetItem(0)
	a.playlist.Flex.Clear()
	a.playlist.Flex.AddItem(playlistContainer, 0, 1, false)
	a.playlist.Flex.AddItem(a.playlistFooter, 1, 0, false)
}

func (a *SimpleApp) setupDetailsComponent() {
	a.detailsThumb = tview.NewImage().
		SetColors(tview.TrueColor).
		SetDithering(tview.DitheringFloydSteinberg)

	a.detailsText = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(a.theme.Text).
		SetWordWrap(true)

	a.detailsView = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.detailsThumb, 20, 0, false).
		AddItem(a.detailsText, 0, 1, false)

	a.detailsView.SetBorder(true).
		SetTitle(" Detalhes ").
		SetBorderColor(a.theme.Surface0)
}

func (a *SimpleApp) setupPlayerComponents() {
	a.thumbnailView = tview.NewImage().
		SetColors(tview.TrueColor).
		SetDithering(tview.DitheringFloydSteinberg)

	a.thumbnailView.SetBorder(false)

	a.playerInfo = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(a.theme.Text)

	a.playerInfo.SetBorder(false)

	playerContent := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.thumbnailView, 20, 0, false).
		AddItem(a.playerInfo, 0, 1, false)

	a.playerBox = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(playerContent, 0, 1, false)

	a.playerBox.SetBorder(true).
		SetTitle(" Player ").
		SetBorderColor(a.theme.Surface0)
}

func (a *SimpleApp) setupStatusBars() {
	a.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	a.commandBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	a.modeBadge = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)

	a.updatePlayerInfo()
	a.updateModeBadge()
	a.statusBar.SetText("")
	a.updateCommandBar()
}

func (a *SimpleApp) setupHelpModal() {
	a.helpModal = tview.NewModal().
		SetText(a.getHelpText()).
		AddButtons([]string{a.strings.Close}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.app.SetRoot(a.getMainLayout(), true)
		})
}

func (a *SimpleApp) setupConfigModal() {
	a.configModal = tview.NewModal().
		SetText(a.getConfigText()).
		AddButtons([]string{
			a.strings.Language + ": " + GetLanguageName(a.language),
			a.strings.Theme + ": " + a.theme.Name,
			a.strings.Help,
			a.strings.Close,
		}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonIndex {
			case 0:
				a.cycleLanguage()
			case 1:
				a.cycleTheme()
			case 2:
				a.app.SetRoot(a.helpModal, true)
			case 3:
				a.app.SetRoot(a.getMainLayout(), true)
			}
		})
}

func (a *SimpleApp) setupLayout() {
	searchPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.searchInput, 3, 0, true).
		AddItem(a.searchResults.Flex, 0, 1, true)

	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(searchPanel, 0, 1, true).
		AddItem(a.playlist.Flex, 0, 1, true)

	statusBarFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.statusBar, 0, 1, false).
		AddItem(a.modeBadge, 18, 0, false)

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topFlex, 0, 1, true).
		AddItem(a.playerBox, 5, 0, false).
		AddItem(statusBarFlex, 1, 0, false).
		AddItem(a.commandBar, 1, 0, false)

	a.app.SetRoot(mainLayout, true).SetFocus(a.searchInput)
}

func (a *SimpleApp) getMainLayout() tview.Primitive {
	searchPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.searchInput, 3, 0, true).
		AddItem(a.searchResults.Flex, 0, 1, true)

	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(searchPanel, 0, 1, true).
		AddItem(a.playlist.Flex, 0, 1, true)

	statusBarFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.statusBar, 0, 1, false).
		AddItem(a.modeBadge, 18, 0, false)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topFlex, 0, 1, true).
		AddItem(a.playerBox, 5, 0, false).
		AddItem(statusBarFlex, 1, 0, false).
		AddItem(a.commandBar, 1, 0, false)
}

func (a *SimpleApp) setupInputHandlers() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		focused := a.app.GetFocus()

		if event.Key() == tcell.KeyCtrlQ {
			a.cleanup()
			a.app.Stop()
			return nil
		}

		if event.Key() == tcell.KeyCtrlC {
			a.app.SetRoot(a.configModal, true)
			return nil
		}

		if event.Rune() == '?' && focused != a.searchInput {
			a.app.SetRoot(a.helpModal, true)
			return nil
		}

		if focused == a.searchInput {
			if event.Key() == tcell.KeyTab {
				a.app.SetFocus(a.searchResults.Flex)
				a.updateCommandBar()
				return nil
			}
			return event
		}

		switch event.Key() {
		case tcell.KeyTab:
			switch focused {
			case a.searchResults.Flex:
				a.app.SetFocus(a.playlist.Flex)
				a.updateCommandBar()
			case a.playlist.Flex:
				a.app.SetFocus(a.playerBox)
				a.updateCommandBar()
			case a.playerBox:
				a.app.SetFocus(a.searchInput)
				a.updateCommandBar()
			}
			return nil
		}

		return a.handleKeyPress(event, focused)
	})
}

func (a *SimpleApp) getHelpText() string {
	return a.strings.HelpTitle + "\n\n" +
		a.strings.HelpNavigation + ":\n" + a.strings.HelpNavigationText + "\n\n" +
		a.strings.HelpSearch + ":\n" + a.strings.HelpSearchText + "\n\n" +
		a.strings.HelpResults + ":\n" + a.strings.HelpResultsText + "\n\n" +
		a.strings.HelpPlaylist + ":\n" + a.strings.HelpPlaylistText + "\n\n" +
		a.strings.HelpPlayer + ":\n" + a.strings.HelpPlayerText + "\n  \n" +
		a.strings.HelpGlobal + ":\n" + a.strings.HelpGlobalText + "\n\n" +
		a.strings.HelpIcons + ":\n" + a.strings.HelpIconsText
}

func (a *SimpleApp) getConfigText() string {
	return a.strings.ConfigText
}

func (a *SimpleApp) cycleLanguage() {
	languages := GetAllLanguages()

	currentIdx := 0
	for i, lang := range languages {
		if lang == a.language {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(languages)
	nextLang := languages[nextIdx]

	a.applyLanguage(nextLang)
	a.refreshUI()

	langName := GetLanguageName(a.language)
	a.setStatus(a.theme.Green, "✓ "+fmt.Sprintf(a.strings.LanguageChanged, langName))
}

func (a *SimpleApp) cycleTheme() {
	themes := GetAllThemes()

	currentIdx := 0
	for i, theme := range themes {
		if theme.ID == a.theme.ID {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(themes)
	newTheme := themes[nextIdx]
	a.theme = &newTheme

	cfg, _ := config.LoadConfig()
	cfg.Theme.Active = a.theme.ID
	cfg.Theme.CustomPath = ""
	config.SaveConfig(cfg)

	a.applyTheme()
	a.refreshUI()

	a.setStatus(a.theme.Green, "✓ "+fmt.Sprintf(a.strings.ThemeChanged, a.theme.Name))
}

func (a *SimpleApp) applyTheme() {
	tview.Styles.PrimitiveBackgroundColor = a.theme.Base
	tview.Styles.ContrastBackgroundColor = a.theme.Surface0
	tview.Styles.MoreContrastBackgroundColor = a.theme.Surface1
	tview.Styles.BorderColor = a.theme.Surface0
	tview.Styles.TitleColor = a.theme.Text
	tview.Styles.GraphicsColor = a.theme.Blue
	tview.Styles.PrimaryTextColor = a.theme.Text
	tview.Styles.SecondaryTextColor = a.theme.Subtext1
	tview.Styles.TertiaryTextColor = a.theme.Subtext0
	tview.Styles.InverseTextColor = a.theme.Base
	tview.Styles.ContrastSecondaryTextColor = a.theme.Subtext0

	a.searchInput.SetFieldBackgroundColor(a.theme.Surface0).
		SetFieldTextColor(a.theme.Text).
		SetBorderColor(a.theme.Blue)

	a.searchResults.SetTheme(a.theme)
	a.playlist.SetTheme(a.theme)

	a.playerBox.SetBorderColor(a.theme.Surface0)

	a.statusBar.SetBackgroundColor(a.theme.Base)
	a.statusBar.SetTextColor(a.theme.Text)

	a.commandBar.SetBackgroundColor(a.theme.Base)
	a.commandBar.SetTextColor(a.theme.Subtext1)

	a.modeBadge.SetBackgroundColor(a.theme.Base)
	a.modeBadge.SetTextColor(a.theme.Mauve)

	a.playlistFooter.SetBackgroundColor(a.theme.Base)
	a.playlistFooter.SetTextColor(a.theme.Subtext0)

	a.playerInfo.SetBackgroundColor(a.theme.Base)
	a.playerInfo.SetTextColor(a.theme.Text)

	a.detailsText.SetBackgroundColor(a.theme.Base)
	a.detailsText.SetTextColor(a.theme.Text)

	a.updateCommandBar()
}

func (a *SimpleApp) refreshUI() {
	a.searchInput.SetBorder(true).SetTitle(" " + a.strings.Search + " ")
	a.searchResults.SetTitle(" " + a.strings.Results + " [0] ")

	count := len(a.playlistTracks)
	a.playlist.SetTitle(fmt.Sprintf(" %s [%d] ", a.strings.Playlist, count))
	a.playlist.SetTitleColor(a.theme.Subtext0)

	a.playerBox.SetTitle(" " + a.strings.Player + " ")

	a.helpModal.SetText(a.getHelpText())
	a.helpModal.ClearButtons().AddButtons([]string{a.strings.Close})

	a.configModal.SetText(a.getConfigText())
	a.configModal.ClearButtons().AddButtons([]string{
		a.strings.Language + ": " + GetLanguageName(a.language),
		a.strings.Theme + ": " + a.theme.Name,
		a.strings.Help,
		a.strings.Close,
	})

	a.updateCommandBar()
	a.updatePlaylistFooter()
	a.updateModeBadge()
	a.updatePlayerInfo()

	a.app.SetRoot(a.configModal, true)
}

func (a *SimpleApp) initLanguageFromConfig() {
	cfg, _ := config.LoadConfig()

	lang := parseLanguage(cfg.UI.Language)
	if lang == "" {
		lang = LanguagePT
	}

	a.applyLanguage(lang)
}

func parseLanguage(s string) Language {
	s = strings.ToLower(strings.TrimSpace(s))
	switch {
	case s == "en", strings.HasPrefix(s, "en-"), strings.HasPrefix(s, "en_"):
		return LanguageEN
	case s == "pt", strings.HasPrefix(s, "pt-"), strings.HasPrefix(s, "pt_"), s == "br":
		return LanguagePT
	default:
		return LanguagePT
	}
}

func (a *SimpleApp) applyLanguage(lang Language) {
	a.language = lang
	a.strings = GetStrings(lang)

	cfg, _ := config.LoadConfig()
	cfg.UI.Language = string(lang)
	_ = config.SaveConfig(cfg)

	search.SetTexts(search.Texts{
		EmptyQuery:       a.strings.EmptyQuery,
		NoResultsFor:     a.strings.NoResultsFor,
		YtDlpNotFound:    a.strings.YtDlpNotFound,
		YtDlpStartFailed: a.strings.YtDlpStartFailed,
		YtDlpError:       a.strings.YtDlpError,
		UnknownDate:      a.strings.UnknownDate,
		NoDescription:    a.strings.NoDescription,
	})
}
