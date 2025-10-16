package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/levirenato/YouTui/internal/search"
)

type item struct {
	title, desc, url string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title + " " + i.desc }

type errMsg error

type playMode int

const (
	playModeVideo playMode = iota
	playModeAudio
)

func (pm playMode) String() string {
	if pm == playModeAudio {
		return "MP3 (Áudio)"
	}
	return "MP4 (Vídeo)"
}

type Model struct {
	ti            textinput.Model
	li            list.Model
	spin          spinner.Model
	loading       bool
	status        string
	results       []search.Result
	mode          playMode
	width         int
	height        int
	currentTitle  string
	playlist      []item
	focusedPanel  int // 0=busca, 1=playlist, 2=lista, 3=visualizador
	selectedIndex int // índice selecionado nos resultados
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Buscar vídeos…"
	ti.Focus()

	li := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	li.Title = "Resultados"
	li.SetShowHelp(false)
	li.SetShowStatusBar(false)

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return Model{
		ti:            ti,
		li:            li,
		spin:          sp,
		mode:          playModeVideo,
		playlist:      []item{},
		focusedPanel:  0,
		selectedIndex: 0,
	}
}

type searchedMsg struct {
	results []search.Result
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) doSearch(q string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()
		res, err := search.SearchVideos(ctx, q, 30)
		if err != nil {
			return errMsg(err)
		}
		return searchedMsg{results: res}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.li.SetSize(msg.Width-4, msg.Height-12)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			// Se estiver no input com foco e tiver texto -> buscar
			if m.ti.Focused() {
				q := m.ti.Value()
				if q == "" {
					m.status = "Digite algo para buscar."
					return m, nil
				}
				m.loading = true
				m.status = "Buscando…"
				return m, tea.Batch(m.spin.Tick, m.doSearch(q))
			}
			// Se estiver na lista de resultados -> tocar
			if m.focusedPanel == 2 && m.selectedIndex < len(m.results) {
				r := m.results[m.selectedIndex]
				m.currentTitle = r.Title
				return m, m.play(r.URL, r.Title)
			}
			// Se estiver na playlist -> tocar
			if m.focusedPanel == 1 && len(m.playlist) > 0 {
				// TODO: implementar seleção na playlist
			}
		case "tab":
			// alterna entre painéis
			m.focusedPanel = (m.focusedPanel + 1) % 4
			if m.focusedPanel == 0 {
				m.ti.Focus()
			} else {
				m.ti.Blur()
			}
		case "m":
			// alterna modo de reprodução
			if m.mode == playModeVideo {
				m.mode = playModeAudio
			} else {
				m.mode = playModeVideo
			}
			m.status = "Modo: " + m.mode.String()
		case "p":
			// adiciona à playlist
			if m.focusedPanel == 2 && m.selectedIndex < len(m.results) {
				r := m.results[m.selectedIndex]
				it := item{
					title: r.Title,
					desc:  r.Author,
					url:   r.URL,
				}
				m.playlist = append(m.playlist, it)
				m.status = fmt.Sprintf("Adicionado: %s", truncate(r.Title, 50))
			}
		case "up", "k":
			if m.focusedPanel == 2 && m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
			if m.focusedPanel == 2 && m.selectedIndex < len(m.results)-1 {
				m.selectedIndex++
			}
		}
	case searchedMsg:
		m.loading = false
		m.status = fmt.Sprintf("Encontrados: %d", len(msg.results))
		m.results = msg.results
		m.selectedIndex = 0
		m.focusedPanel = 2 // muda foco para resultados

		items := make([]list.Item, 0, len(msg.results))
		for _, r := range msg.results {
			items = append(items, item{
				title: fmt.Sprintf("%s  ·  %s", r.Title, r.Duration),
				desc:  r.Author,
				url:   r.URL,
			})
		}
		m.li.SetItems(items)
		if len(items) > 0 {
			m.ti.Blur()
		}
	case errMsg:
		m.loading = false
		m.status = "Erro: " + msg.Error()
	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spin, cmd = m.spin.Update(msg)
			return m, cmd
		}
	}

	// Encaminhar eventos para widgets
	var cmd tea.Cmd
	if m.ti.Focused() {
		m.ti, cmd = m.ti.Update(msg)
	} else {
		m.li, cmd = m.li.Update(msg)
	}
	return m, cmd
}

// Executa MPV para reproduzir vídeo/áudio do YouTube.
// Estratégia atualizada para 2025: usa formatos que não requerem PO Token.
// Suporta modo vídeo (MP4) e modo áudio (MP3).
func (m Model) play(videoURL, title string) tea.Cmd {
	return func() tea.Msg {
		// Argumentos base do mpv
		args := []string{
			"--force-window=no",
			"--quiet",
			"--script-opts=ytdl_hook-ytdl_path=yt-dlp",
			fmt.Sprintf("--title=%s", title),
		}

		// Se modo áudio, adiciona flags específicas
		if m.mode == playModeAudio {
			args = append(args,
				"--no-video",
				"--ytdl-format=bestaudio",
			)
		}

		args = append(args, videoURL)

		// Tentativa 1: mpv com yt-dlp usando formatos padrão
		cmd := exec.Command("mpv", args...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd.Run(); err == nil {
			return nil // sucesso
		}

		// Tentativa 2: fallback com formato específico
		args2 := []string{
			"--force-window=no",
			"--quiet",
			"--script-opts=ytdl_hook-ytdl_path=yt-dlp",
			fmt.Sprintf("--title=%s", title),
		}

		if m.mode == playModeAudio {
			args2 = append(args2, "--no-video", "--ytdl-format=140") // m4a audio
		} else {
			args2 = append(args2, "--ytdl-format=18") // 360p MP4
		}

		args2 = append(args2, videoURL)

		cmd2 := exec.Command("mpv", args2...)
		cmd2.Stdin, cmd2.Stdout, cmd2.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd2.Run(); err == nil {
			return nil
		}

		fmt.Fprintf(os.Stderr, "Erro ao reproduzir.\n")
		fmt.Fprintf(os.Stderr, "Certifique-se de ter mpv e yt-dlp instalados e atualizados.\n")
		return nil
	}
}

func (m Model) View() string {
	if m.width == 0 {
		return "Carregando..."
	}

	// Título principal
	header := titleStyle.Render("🎥 YouTui - YouTube Terminal Interface")

	// Barra de comandos
	commandBar := renderCommandBar(m.mode.String())

	// Se estiver carregando
	if m.loading {
		content := fmt.Sprintf("\n%s %s\n", m.spin.View(), m.status)
		return header + "\n" + content + "\n" + commandBar
	}

	// Calcula dimensões dos painéis
	panelWidth := (m.width - 6) / 2
	panelHeight := (m.height - 8) / 2

	// Painel 1: Busca (superior esquerdo)
	searchPanel := m.renderSearchPanel(panelWidth, panelHeight)

	// Painel 2: Playlist (superior direito)
	playlistPanel := m.renderPlaylistPanel(panelWidth, panelHeight)

	// Painel 3: Lista de resultados (inferior esquerdo)
	resultsPanel := m.renderResultsPanel(panelWidth, panelHeight)

	// Painel 4: Visualizador/Info (inferior direito)
	infoPanel := m.renderInfoPanel(panelWidth, panelHeight)

	// Monta linha superior
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, searchPanel, playlistPanel)

	// Monta linha inferior
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, resultsPanel, infoPanel)

	// Monta grid completo
	grid := lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)

	// Status
	statusView := statusStyle.Render(m.status)

	return header + "\n\n" + grid + "\n" + statusView + "\n" + commandBar
}

// renderSearchPanel renderiza o painel de busca (superior esquerdo)
func (m Model) renderSearchPanel(width, height int) string {
	title := "🔍 Busca"
	if m.focusedPanel == 0 {
		title = "🔍 Busca [ATIVO]"
	}

	content := m.ti.View()
	if !m.ti.Focused() {
		content = lipgloss.NewStyle().Foreground(subtleColor).Render(content)
	}

	// Adiciona instruções
	instructions := lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true).
		Render("\nDigite e pressione Enter para buscar")

	panelContent := content + instructions

	return createPanel(title, panelContent, width, height, m.focusedPanel == 0)
}

// renderPlaylistPanel renderiza o painel de playlist (superior direito)
func (m Model) renderPlaylistPanel(width, height int) string {
	title := "📋 Playlist"
	if m.focusedPanel == 1 {
		title = "📋 Playlist [ATIVO]"
	}

	var content strings.Builder
	if len(m.playlist) == 0 {
		content.WriteString(lipgloss.NewStyle().
			Foreground(subtleColor).
			Italic(true).
			Render("Nenhum item na playlist\n\nPressione 'p' para adicionar"))
	} else {
		for i, item := range m.playlist {
			line := fmt.Sprintf("%d. %s\n", i+1, truncate(item.title, width-8))
			content.WriteString(line)
		}
	}

	return createPanel(title, content.String(), width, height, m.focusedPanel == 1)
}

// renderResultsPanel renderiza o painel de resultados (inferior esquerdo)
func (m Model) renderResultsPanel(width, height int) string {
	title := "📺 Resultados"
	if m.focusedPanel == 2 {
		title = "📺 Resultados [ATIVO]"
	}

	var content string
	if len(m.results) == 0 {
		content = lipgloss.NewStyle().
			Foreground(subtleColor).
			Italic(true).
			Render("Nenhum resultado ainda\n\nFaça uma busca para começar")
	} else {
		// Renderiza lista simplificada
		items := []string{}
		maxItems := height - 6
		if maxItems > len(m.results) {
			maxItems = len(m.results)
		}

		for i := 0; i < maxItems; i++ {
			r := m.results[i]
			prefix := "  "
			if i == m.selectedIndex {
				prefix = "▶ "
			}

			line := fmt.Sprintf("%s%s", prefix, truncate(r.Title, width-10))
			if r.Duration != "" {
				line += fmt.Sprintf(" [%s]", r.Duration)
			}

			// Destaca item selecionado
			if i == m.selectedIndex && m.focusedPanel == 2 {
				line = lipgloss.NewStyle().
					Foreground(accentColor).
					Bold(true).
					Render(line)
			}

			items = append(items, line)
		}
		content = strings.Join(items, "\n")
	}

	return createPanel(title, content, width, height, m.focusedPanel == 2)
}

// renderInfoPanel renderiza o painel de informações/visualizador (inferior direito)
func (m Model) renderInfoPanel(width, height int) string {
	title := "🎵 Controles"
	if m.focusedPanel == 3 {
		title = "🎵 Controles [ATIVO]"
	}

	var content strings.Builder

	// Modo atual
	content.WriteString(lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Render(fmt.Sprintf("Modo: %s\n\n", m.mode.String())))

	// Estatísticas
	content.WriteString(fmt.Sprintf("Resultados: %d\n", len(m.results)))
	content.WriteString(fmt.Sprintf("Playlist: %d itens\n\n", len(m.playlist)))

	// Atalhos
	content.WriteString(lipgloss.NewStyle().
		Foreground(subtleColor).
		Render("Atalhos:\n"))
	content.WriteString("Tab - Trocar painel\n")
	content.WriteString("m - Mudar modo\n")
	content.WriteString("p - Adicionar à playlist\n")
	content.WriteString("Enter - Reproduzir\n")

	return createPanel(title, content.String(), width, height, m.focusedPanel == 3)
}

// truncate trunca uma string se for maior que o limite
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
