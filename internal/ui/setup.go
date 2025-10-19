package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// setupUI configura todos os componentes da interface
func (a *SimpleApp) setupUI() {
	a.setupSearchComponents()
	a.setupPlaylistComponent()
	a.setupDetailsComponent()
	a.setupPlayerComponents()
	a.setupStatusBars()
	a.setupHelpModal()
	a.setupLayout()
	a.setupInputHandlers()
}

// setupSearchComponents cria os componentes de busca e resultados
func (a *SimpleApp) setupSearchComponents() {
	a.searchInput = tview.NewInputField().
		SetLabel(" ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(a.theme.Surface0).
		SetFieldTextColor(a.theme.Text)

	a.searchInput.SetBorder(true).
		SetTitle(" Busca ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(a.theme.Blue)

	a.searchInput.SetDoneFunc(a.onSearchDone)

	a.searchResults = tview.NewList().
		ShowSecondaryText(false).
		SetMainTextColor(a.theme.Text).
		SetSelectedTextColor(a.theme.Base).
		SetSelectedBackgroundColor(a.theme.Blue)

	a.searchResults.SetBorder(true).
		SetTitle(" Resultados [0] ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(a.theme.Surface0)

	a.searchResults.SetSelectedFunc(a.onResultSelected)
	a.searchResults.SetChangedFunc(a.onResultChanged)
}

// setupPlaylistComponent cria o componente de playlist
func (a *SimpleApp) setupPlaylistComponent() {
	a.playlist = tview.NewList().
		ShowSecondaryText(false).
		SetMainTextColor(a.theme.Text).
		SetSelectedTextColor(a.theme.Base).
		SetSelectedBackgroundColor(a.theme.Blue)

	a.playlist.SetBorder(true).
		SetTitle(" Playlist [0] ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(a.theme.Surface0)

	a.playlist.SetSelectedFunc(a.onPlaylistSelected)
}

// setupDetailsComponent cria o painel de detalhes
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

// setupPlayerComponents cria os componentes do player
func (a *SimpleApp) setupPlayerComponents() {
	a.thumbnailView = tview.NewImage().
		SetColors(tview.TrueColor).
		SetDithering(tview.DitheringFloydSteinberg)
	
	a.thumbnailView.SetBorder(true).
		SetTitle("  ").
		SetBorderColor(a.theme.Mauve)

	a.playerInfo = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(a.theme.Text)

	a.playerInfo.SetBorder(true).
		SetTitle(" Player ").
		SetBorderColor(a.theme.Mauve)
}

// setupStatusBars cria as barras de status e comandos
func (a *SimpleApp) setupStatusBars() {
	a.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	a.commandBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	a.updatePlayerInfo()
	a.statusBar.SetText("")
	a.updateCommandBar()
}

// setupHelpModal cria o modal de ajuda
func (a *SimpleApp) setupHelpModal() {
	a.helpModal = tview.NewModal().
		SetText(a.getHelpText()).
		AddButtons([]string{"Fechar"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.app.SetRoot(a.getMainLayout(), true)
		})
}

// setupLayout configura o layout principal da aplicação
func (a *SimpleApp) setupLayout() {
	searchPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.searchInput, 3, 0, true).
		AddItem(a.searchResults, 0, 1, false).
		AddItem(a.detailsView, 7, 0, false)

	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(searchPanel, 0, 1, true).
		AddItem(a.playlist, 0, 1, false)

	playerFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.thumbnailView, 20, 0, false).
		AddItem(a.playerInfo, 0, 1, false)

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topFlex, 0, 1, true).
		AddItem(playerFlex, 5, 0, false).
		AddItem(a.statusBar, 1, 0, false).
		AddItem(a.commandBar, 1, 0, false)

	a.app.SetRoot(mainLayout, true).SetFocus(a.searchInput)
}

// getMainLayout retorna o layout principal (usado para voltar do help)
func (a *SimpleApp) getMainLayout() tview.Primitive {
	searchPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.searchInput, 3, 0, true).
		AddItem(a.searchResults, 0, 1, false).
		AddItem(a.detailsView, 7, 0, false)

	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(searchPanel, 0, 1, true).
		AddItem(a.playlist, 0, 1, false)

	playerFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.thumbnailView, 20, 0, false).
		AddItem(a.playerInfo, 0, 1, false)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topFlex, 0, 1, true).
		AddItem(playerFlex, 5, 0, false).
		AddItem(a.statusBar, 1, 0, false).
		AddItem(a.commandBar, 1, 0, false)
}

// setupInputHandlers configura os handlers de entrada globais
func (a *SimpleApp) setupInputHandlers() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		focused := a.app.GetFocus()

		// Ajuda global
		if event.Rune() == '?' && focused != a.searchInput {
			a.app.SetRoot(a.helpModal, true)
			return nil
		}

		// Se está no input de busca
		if focused == a.searchInput {
			if event.Key() == tcell.KeyTab {
				a.app.SetFocus(a.searchResults)
				a.updateCommandBar()
				return nil
			}
			return event
		}

		// Tab para navegar entre painéis
		switch event.Key() {
		case tcell.KeyTab:
			switch focused {
			case a.searchResults:
				a.app.SetFocus(a.playlist)
				a.updateCommandBar()
			case a.playlist:
				a.app.SetFocus(a.searchInput)
				a.updateCommandBar()
			}
			return nil
		}

		// Atalhos de teclado
		return a.handleKeyPress(event, focused)
	})
}

// getHelpText retorna o texto de ajuda
func (a *SimpleApp) getHelpText() string {
	return `ATALHOS DO YOUTUI

NAVEGAÇÃO:
  Tab       Alternar entre painéis
  /         Focar na busca
  ↑/↓       Navegar nas listas
  ?         Mostrar esta ajuda

BUSCA:
  Digite    Texto para buscar
  Enter     Executar busca

RESULTADOS:
  Enter     Tocar faixa diretamente (sem playlist)
  a         Adicionar à playlist

PLAYLIST:
  Enter     Tocar faixa da playlist
  Space     Tocar playlist do início
  d         Remover item
  J         Mover item para baixo
  K         Mover item para cima

PLAYER (Global):
  c/Space   Pause/Play
  s         Stop
  n         Próxima (só funciona tocando da playlist)
  b         Anterior (só funciona tocando da playlist)
  m         Alternar áudio/vídeo
  r         Ciclar repetição
  h         Toggle shuffle

IMPORTANTE:
  n/b só funcionam quando tocando da PLAYLIST!
  Para tocar da playlist: Entre na Playlist e pressione Enter ou Space

GERAL:
  q         Sair
`
}
