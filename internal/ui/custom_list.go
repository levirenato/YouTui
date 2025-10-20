package ui

import (
	"image"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CustomListItem representa um item com thumbnail inline
type CustomListItem struct {
	flex      *tview.Flex
	thumbnail *tview.Image
	info      *tview.TextView
	index     int
	track     Track
}

// CustomList é uma lista customizada com thumbnails inline
type CustomList struct {
	*tview.Flex
	container     *tview.Flex
	items         []*CustomListItem
	selectedIndex int
	playingIndex  int // Índice do item atualmente tocando (-1 se nenhum)
	theme         *Theme
	mu            sync.Mutex
	onSelected    func(index int)
	visibleStart  int
	visibleHeight int
}

// NewCustomList cria uma nova lista customizada
func NewCustomList(theme *Theme) *CustomList {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	
	// Wrapper é o componente principal (focusável)
	wrapper := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(container, 0, 1, false) // Container NÃO recebe foco
	
	list := &CustomList{
		Flex:          wrapper,
		container:     container,
		items:         []*CustomListItem{},
		selectedIndex: 0,
		playingIndex:  -1, // -1 = nenhum item tocando
		theme:         theme,
		visibleStart:  0,
		visibleHeight: 10, // Altura padrão em items (não linhas)
	}

	// Configura o wrapper como focusável
	wrapper.SetBorder(true).
		SetTitle(" Resultados [0] ").
		SetBorderColor(theme.Surface0)

	// Captura input no wrapper (que é focusável)
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
			// Deixa Tab e Shift+Tab passarem para o handler global
			return event
		}
		
		// Deixa outros eventos (como 'a', 'd', '/', etc) passarem para o handler global
		// Só captura setas e Enter
		return event
	})

	return list
}

// AddItem adiciona um item com thumbnail
func (c *CustomList) AddItem(track Track, index int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Cria thumbnail (não focusável)
	thumb := tview.NewImage().
		SetColors(tview.TrueColor).
		SetDithering(tview.DitheringFloydSteinberg)

	// Cria info text (não focusável)
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetText(formatItemInfo(track, index)).
		SetTextAlign(tview.AlignLeft)

	// Cria flex horizontal [thumbnail | info] (não focusável)
	itemFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(thumb, 20, 0, false).  // false = não focusável
		AddItem(info, 0, 1, false)     // false = não focusável

	item := &CustomListItem{
		flex:      itemFlex,
		thumbnail: thumb,
		info:      info,
		index:     index,
		track:     track,
	}

	c.items = append(c.items, item)
	
	// Renderiza items visíveis
	c.renderVisibleItems()
}

// SetThumbnail define a thumbnail de um item específico
func (c *CustomList) SetThumbnail(index int, img image.Image) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index >= 0 && index < len(c.items) {
		c.items[index].thumbnail.SetImage(img)
	}
}

// Clear limpa todos os items
func (c *CustomList) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.container.Clear()
	c.items = []*CustomListItem{}
	c.selectedIndex = 0
	c.visibleStart = 0
}

// renderVisibleItems renderiza apenas os items visíveis para evitar overflow
func (c *CustomList) renderVisibleItems() {
	c.container.Clear()
	
	// Calcula altura disponível (em linhas)
	_, _, _, height := c.Flex.GetInnerRect()
	if height < 3 {
		height = 30 // Valor padrão se ainda não foi renderizado
	}
	
	// Cada item tem 3 linhas de altura
	itemsPerPage := height / 3
	if itemsPerPage < 1 {
		itemsPerPage = 1
	}
	
	c.visibleHeight = itemsPerPage
	
	// Calcula quais items devem ser visíveis
	end := c.visibleStart + c.visibleHeight
	if end > len(c.items) {
		end = len(c.items)
	}
	
	// Adiciona apenas os items visíveis
	for i := c.visibleStart; i < end; i++ {
		c.container.AddItem(c.items[i].flex, 3, 0, false)
	}
	
	c.updateSelection()
}

// scrollToSelection ajusta o scroll para mostrar o item selecionado
func (c *CustomList) scrollToSelection() {
	// Se o item selecionado está abaixo da área visível
	if c.selectedIndex >= c.visibleStart+c.visibleHeight {
		c.visibleStart = c.selectedIndex - c.visibleHeight + 1
		c.renderVisibleItems()
	}
	
	// Se o item selecionado está acima da área visível
	if c.selectedIndex < c.visibleStart {
		c.visibleStart = c.selectedIndex
		c.renderVisibleItems()
	}
}

// SelectNext seleciona o próximo item
func (c *CustomList) SelectNext() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.selectedIndex < len(c.items)-1 {
		c.selectedIndex++
		c.scrollToSelection()
		c.updateSelection()
	}
}

// SelectPrevious seleciona o item anterior
func (c *CustomList) SelectPrevious() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.selectedIndex > 0 {
		c.selectedIndex--
		c.scrollToSelection()
		c.updateSelection()
	}
}

// GetCurrentItem retorna o índice do item selecionado
func (c *CustomList) GetCurrentItem() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.selectedIndex
}

// SetCurrentIndex define o índice do item selecionado
func (c *CustomList) SetCurrentIndex(idx int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if idx >= 0 && idx < len(c.items) {
		c.selectedIndex = idx
		c.scrollToSelection()
		c.updateSelection()
	}
}

// GetCurrentTrack retorna o track do item selecionado
func (c *CustomList) GetCurrentTrack() *Track {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.selectedIndex >= 0 && c.selectedIndex < len(c.items) {
		return &c.items[c.selectedIndex].track
	}
	return nil
}

// IsFocused verifica se esta lista está focada na aplicação
func (c *CustomList) IsFocused(app *tview.Application) bool {
	focused := app.GetFocus()
	return focused == c.Flex
}

// SetSelectedFunc define o callback quando Enter é pressionado
func (c *CustomList) SetSelectedFunc(handler func(int)) {
	c.onSelected = handler
}

// updateSelection atualiza a aparência visual da seleção
func (c *CustomList) updateSelection() {
	for i, item := range c.items {
		if i == c.selectedIndex {
			// Item selecionado - fundo azul + TEXTO PRETO para contraste
			item.flex.SetBackgroundColor(c.theme.Blue)
			item.info.SetTextColor(c.theme.Crust)
			item.info.SetBackgroundColor(c.theme.Blue)
			// Atualiza texto sem cores inline (texto preto puro)
			item.info.SetText(formatItemInfoPlain(item.track, item.index))
		} else if i == c.playingIndex {
			// Item tocando - fundo verde + TEXTO PRETO para contraste
			item.flex.SetBackgroundColor(c.theme.Green)
			item.info.SetTextColor(c.theme.Crust)
			item.info.SetBackgroundColor(c.theme.Green)
			// Atualiza texto sem cores inline (texto preto puro)
			item.info.SetText(formatItemInfoPlain(item.track, item.index))
		} else {
			// Item normal - fundo padrão com texto colorido
			item.flex.SetBackgroundColor(c.theme.Base)
			item.info.SetTextColor(c.theme.Text)
			item.info.SetBackgroundColor(c.theme.Base)
			// Atualiza texto COM cores inline (amarelo, verde, ciano)
			item.info.SetText(formatItemInfo(item.track, item.index))
		}
	}
}

// SetTitle atualiza o título da lista
func (c *CustomList) SetTitle(title string) *CustomList {
	c.Flex.SetTitle(title)
	return c
}

// SetBorderColor atualiza a cor da borda
func (c *CustomList) SetBorderColor(color tcell.Color) *CustomList {
	c.Flex.SetBorderColor(color)
	return c
}

// SetPlayingIndex define qual item está atualmente tocando
func (c *CustomList) SetPlayingIndex(idx int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.playingIndex = idx
	c.updateSelection()
}

// formatItemInfo formata as informações do item com cores
func formatItemInfo(track Track, index int) string {
	icons := []string{"♪", "♫", "♬"}
	icon := icons[index%len(icons)]
	
	// Trunca título se muito longo
	title := track.Title
	if len(title) > 50 {
		title = title[:47] + "..."
	}
	
	return icon + " [yellow::b]" + title + "[-:-:-]\n" +
		"[green]⏱ " + track.Duration + "[-] [cyan]• " + track.Author + "[-]"
}

// formatItemInfoPlain formata as informações sem cores inline (para itens selecionados/tocando)
func formatItemInfoPlain(track Track, index int) string {
	icons := []string{"♪", "♫", "♬"}
	icon := icons[index%len(icons)]
	
	// Trunca título se muito longo
	title := track.Title
	if len(title) > 50 {
		title = title[:47] + "..."
	}
	
	// SEM markup de cores - texto puro que herda a cor do SetTextColor
	return icon + " " + title + "\n" +
		"⏱ " + track.Duration + " • " + track.Author
}
