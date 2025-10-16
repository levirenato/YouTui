package ui

import (
	"context"
	"fmt"
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

type playlistMode int

const (
	playlistNormal playlistMode = iota
	playlistShuffle
	playlistRepeatOne
	playlistRepeatAll
)

func (pm playlistMode) String() string {
	switch pm {
	case playlistShuffle:
		return "🔀 Aleatório"
	case playlistRepeatOne:
		return "🔂 Repetir 1"
	case playlistRepeatAll:
		return "🔁 Repetir Todas"
	default:
		return "▶️ Normal"
	}
}

type logLevel int

const (
	logInfo logLevel = iota
	logWarning
	logError
)

type logEntry struct {
	level   logLevel
	message string
	time    time.Time
}

type Model struct {
	ti              textinput.Model
	li              list.Model
	spin            spinner.Model
	loading         bool
	status          string
	results         []search.Result
	mode            playMode
	width           int
	height          int
	currentTitle    string
	playlist        []item
	focusedPanel    int // 0=busca, 1=playlist, 2=lista, 3=visualizador, 4=logs
	selectedIndex   int // índice selecionado nos resultados
	playlistIndex   int // índice selecionado na playlist
	playlistMode    playlistMode
	playlistPlaying bool
	mpvProcess      *exec.Cmd
	cavaProcess     *exec.Cmd
	isPlaying       bool
	nowPlaying      string
	cavaOutput      []string // linhas do cava para visualização
	logs            []logEntry
	notification    string
	notificationLevel logLevel
	showLogs        bool
	program         *tea.Program // referência para enviar mensagens
}

func NewModel() Model  {
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
		playlistIndex: 0,
		playlistMode:  playlistNormal,
		cavaOutput:    []string{},
		logs:          []logEntry{},
		showLogs:      false,
	}
}

// SetProgram define a referência do programa para enviar mensagens assíncronas
func (m *Model) SetProgram(p *tea.Program) {
	m.program = p
}

// addLog adiciona uma entrada de log
func (m *Model) addLog(level logLevel, message string) {
	entry := logEntry{
		level:   level,
		message: message,
		time:    time.Now(),
	}
	m.logs = append(m.logs, entry)
	
	// Mantém apenas os últimos 100 logs
	if len(m.logs) > 100 {
		m.logs = m.logs[1:]
	}
}

// setNotification define uma notificação temporária
func (m *Model) setNotification(level logLevel, message string) {
	m.notification = message
	m.notificationLevel = level
	m.addLog(level, message)
}

type searchedMsg struct {
	results []search.Result
}

type playbackFinishedMsg struct{}

type cavaUpdateMsg struct {
	lines []string
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Cleanup para ao sair - CRÍTICO para matar processos mpv
func (m *Model) Cleanup() {
	if m.mpvProcess != nil && m.mpvProcess.Process != nil {
		m.mpvProcess.Process.Kill()
		m.mpvProcess = nil
	}
	if m.cavaProcess != nil && m.cavaProcess.Process != nil {
		m.cavaProcess.Process.Kill()
		m.cavaProcess = nil
	}
	m.isPlaying = false
	m.playlistPlaying = false
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

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.li.SetSize(msg.Width-4, msg.Height-12)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// CRÍTICO: Limpar processos antes de sair
			m.Cleanup()
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
			if m.focusedPanel == 1 && m.playlistIndex < len(m.playlist) {
				it := m.playlist[m.playlistIndex]
				m.currentTitle = it.title
				return m, m.play(it.url, it.title)
			}
		case "tab":
			// alterna entre painéis
			maxPanels := 4
			if m.showLogs {
				maxPanels = 5
			}
			m.focusedPanel = (m.focusedPanel + 1) % maxPanels
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
			} else if m.focusedPanel == 1 && m.playlistIndex > 0 {
				m.playlistIndex--
			}
		case "down", "j":
			if m.focusedPanel == 2 && m.selectedIndex < len(m.results)-1 {
				m.selectedIndex++
			} else if m.focusedPanel == 1 && m.playlistIndex < len(m.playlist)-1 {
				m.playlistIndex++
			}
		case "d", "x":
			// remove da playlist
			if m.focusedPanel == 1 && len(m.playlist) > 0 && m.playlistIndex < len(m.playlist) {
				m.playlist = append(m.playlist[:m.playlistIndex], m.playlist[m.playlistIndex+1:]...)
				if m.playlistIndex >= len(m.playlist) && m.playlistIndex > 0 {
					m.playlistIndex--
				}
				m.status = "Item removido da playlist"
			}
		case "K":
			// move item para cima na playlist
			if m.focusedPanel == 1 && m.playlistIndex > 0 {
				m.playlist[m.playlistIndex], m.playlist[m.playlistIndex-1] = m.playlist[m.playlistIndex-1], m.playlist[m.playlistIndex]
				m.playlistIndex--
				m.status = "Item movido para cima"
			}
		case "J":
			// move item para baixo na playlist
			if m.focusedPanel == 1 && m.playlistIndex < len(m.playlist)-1 {
				m.playlist[m.playlistIndex], m.playlist[m.playlistIndex+1] = m.playlist[m.playlistIndex+1], m.playlist[m.playlistIndex]
				m.playlistIndex++
				m.status = "Item movido para baixo"
			}
		case "s":
			// para reprodução
			if m.isPlaying {
				m.stopPlayback()
				m.setNotification(logInfo, "Reprodução parada")
				m.status = "Reprodução parada"
			}
		case "l":
			// alterna visualização de logs
			m.showLogs = !m.showLogs
			if m.showLogs {
				m.setNotification(logInfo, "Painel de logs ativado")
			}
		case "r":
			// alterna modo de playlist
			switch m.playlistMode {
			case playlistNormal:
				m.playlistMode = playlistShuffle
			case playlistShuffle:
				m.playlistMode = playlistRepeatOne
			case playlistRepeatOne:
				m.playlistMode = playlistRepeatAll
			case playlistRepeatAll:
				m.playlistMode = playlistNormal
			}
			m.setNotification(logInfo, "Modo de playlist: "+m.playlistMode.String())
			m.status = "Modo de playlist: " + m.playlistMode.String()
		case " ":
			// inicia playlist completa
			if len(m.playlist) > 0 && !m.playlistPlaying {
				m.playlistPlaying = true
				m.playlistIndex = 0
				it := m.playlist[0]
				m.setNotification(logInfo, "Iniciando playlist")
				return m, m.play(it.url, it.title)
			}
		}
	case searchedMsg:
		m.loading = false
		m.status = fmt.Sprintf("Encontrados: %d", len(msg.results))
		m.results = msg.results
		m.selectedIndex = 0
		m.focusedPanel = 2 // muda foco para resultados
		m.addLog(logInfo, fmt.Sprintf("Busca concluída: %d resultados", len(msg.results)))

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
		errorMsg := msg.Error()
		m.status = "Erro: " + errorMsg
		m.setNotification(logError, errorMsg)
	case playbackFinishedMsg:
		m.isPlaying = false
		m.nowPlaying = ""
		m.mpvProcess = nil
		if m.cavaProcess != nil {
			m.cavaProcess.Process.Kill()
			m.cavaProcess = nil
		}
		m.addLog(logInfo, "Reprodução finalizada")
		
		// Avança para próxima música se playlist estiver ativa
		if m.playlistPlaying && len(m.playlist) > 0 {
			switch m.playlistMode {
			case playlistNormal, playlistRepeatAll:
				m.playlistIndex++
				if m.playlistIndex >= len(m.playlist) {
					if m.playlistMode == playlistRepeatAll {
						m.playlistIndex = 0
					} else {
						m.playlistPlaying = false
						m.status = "Playlist finalizada"
						return m, nil
					}
				}
			case playlistRepeatOne:
				// Mantém o mesmo índice
			case playlistShuffle:
				// TODO: implementar shuffle
				m.playlistIndex = (m.playlistIndex + 1) % len(m.playlist)
			}
			
			// Toca próximo item
			if m.playlistIndex < len(m.playlist) {
				it := m.playlist[m.playlistIndex]
				m.addLog(logInfo, fmt.Sprintf("Tocando: %s", it.title))
				return m, m.play(it.url, it.title)
			}
		} else {
			m.status = "Reprodução finalizada"
		}
	case cavaUpdateMsg:
		m.cavaOutput = msg.lines
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

// stopPlayback para a reprodução atual
func (m *Model) stopPlayback() {
	if m.mpvProcess != nil && m.mpvProcess.Process != nil {
		m.mpvProcess.Process.Kill()
		m.mpvProcess = nil
	}
	if m.cavaProcess != nil && m.cavaProcess.Process != nil {
		m.cavaProcess.Process.Kill()
		m.cavaProcess = nil
	}
	m.isPlaying = false
	m.nowPlaying = ""
}

// play executa MPV em background para reproduzir vídeo/áudio do YouTube
func (m *Model) play(videoURL, title string) tea.Cmd {
	return func() tea.Msg {
		// Para reprodução anterior se houver
		if m.mpvProcess != nil {
			m.stopPlayback()
		}

		// Argumentos base do mpv
		args := []string{
			"--no-terminal",
			"--really-quiet",
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

		// Cria comando mpv
		cmd := exec.Command("mpv", args...)
		
		// Redireciona saída para /dev/null
		cmd.Stdout = nil
		cmd.Stderr = nil

		// Inicia processo em background
		if err := cmd.Start(); err != nil {
			return errMsg(fmt.Errorf("erro ao iniciar mpv: %w", err))
		}

		m.mpvProcess = cmd
		m.isPlaying = true
		m.nowPlaying = title
		m.addLog(logInfo, fmt.Sprintf("Reproduzindo: %s", truncate(title, 60)))

		// Inicia cava se modo áudio
		if m.mode == playModeAudio {
			m.startCava()
		}

		// Aguarda finalização em goroutine
		prog := m.program
		go func() {
			cmd.Wait()
			// Envia mensagem de finalização
			if prog != nil {
				prog.Send(playbackFinishedMsg{})
			}
		}()

		return nil
	}
}

// startCava inicia o visualizador de áudio cava
func (m *Model) startCava() {
	// Para cava anterior se houver
	if m.cavaProcess != nil && m.cavaProcess.Process != nil {
		m.cavaProcess.Process.Kill()
	}

	// Cria arquivo de config temporário para cava
	// Por enquanto, apenas simula a visualização
	// TODO: Implementar integração real com cava via pipe
	m.cavaOutput = []string{
		"█▓▒░",
		"██▓▒░",
		"███▓▒░",
	}
	
	m.addLog(logInfo, "Visualizador de áudio iniciado")
}

func (m *Model) View() string {
	if m.width == 0 {
		return "Carregando..."
	}

	// 🎵 NOVO LAYOUT: Player de Música Moderno
	
	// Header com logo
	header := titleStyle.Render("♫ YouTui Music Player")
	
	// Notificação (se houver)
	var notification string
	if m.notification != "" {
		notifColor := successColor
		icon := "ℹ️"
		
		switch m.notificationLevel {
		case logError:
			notifColor = errorColor
			icon = "❌"
		case logWarning:
			notifColor = warningColor
			icon = "⚠️"
		}
		
		notifStyle := lipgloss.NewStyle().
			Foreground(notifColor).
			Bold(true).
			Padding(0, 2)
		
		notification = notifStyle.Render(icon + " " + m.notification)
	}

	// Loading indicator (não invasivo)
	var loadingIndicator string
	if m.loading {
		loadingIndicator = lipgloss.NewStyle().
			Foreground(primaryColor).
			Render(fmt.Sprintf(" %s Buscando...", m.spin.View()))
	}

	// ═══════════════════════════════════════════════
	// PLAYER CENTRAL (Now Playing)
	// ═══════════════════════════════════════════════
	var playerSection string
	if m.isPlaying {
		nowPlayingText := truncate(m.nowPlaying, 60)
		playButton := renderPlayButton(m.isPlaying)
		volumeBar := renderVolumeBar()
		
		playerContent := lipgloss.JoinVertical(lipgloss.Center,
			"",
			lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("♪ TOCANDO AGORA"),
			"",
			lipgloss.NewStyle().Foreground(textColor).Bold(true).Render(nowPlayingText),
			"",
			playButton,
			"",
			volumeBar,
			"",
		)
		
		playerSection = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor).
			Padding(2, 4).
			Width(m.width - 10).
			Align(lipgloss.Center).
			Render(playerContent)
	}

	// ═══════════════════════════════════════════════
	// LAYOUT PRINCIPAL
	// ═══════════════════════════════════════════════
	
	var mainLayout string
	
	if m.isPlaying {
		// Quando tocando: Layout vertical com player no topo
		leftWidth := (m.width - 10) / 2
		rightWidth := m.width - leftWidth - 12
		contentHeight := m.height - 18
		
		// Busca + Playlist (superior)
		searchPanel := m.renderModernSearchPanel(leftWidth, contentHeight/2)
		playlistPanel := m.renderModernPlaylistPanel(rightWidth, contentHeight/2)
		topRow := lipgloss.JoinHorizontal(lipgloss.Top, searchPanel, playlistPanel)
		
		// Resultados + Visualizador (inferior)
		resultsPanel := m.renderModernResultsPanel(leftWidth, contentHeight/2)
		visualizerPanel := m.renderModernVisualizerPanel(rightWidth, contentHeight/2)
		bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, resultsPanel, visualizerPanel)
		
		content := lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)
		mainLayout = lipgloss.JoinVertical(lipgloss.Center, playerSection, "", content)
	} else {
		// Quando não tocando: Layout 2x2 clássico
		leftWidth := (m.width - 10) / 2
		rightWidth := m.width - leftWidth - 12
		panelHeight := (m.height - 10) / 2
		
		searchPanel := m.renderModernSearchPanel(leftWidth, panelHeight)
		playlistPanel := m.renderModernPlaylistPanel(rightWidth, panelHeight)
		topRow := lipgloss.JoinHorizontal(lipgloss.Top, searchPanel, playlistPanel)
		
		resultsPanel := m.renderModernResultsPanel(leftWidth, panelHeight)
		infoPanel := m.renderModernInfoPanel(rightWidth, panelHeight)
		bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, resultsPanel, infoPanel)
		
		mainLayout = lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)
	}

	// ═══════════════════════════════════════════════
	// STATUS BAR + CONTROLES
	// ═══════════════════════════════════════════════
	
	modeIndicator := renderModeIndicator(m.mode.String(), m.playlistMode.String())
	
	statusBar := lipgloss.JoinHorizontal(lipgloss.Left,
		statusStyle.Render(m.status),
		"  ",
		modeIndicator,
		loadingIndicator,
	)
	
	commandBar := m.renderCommandBarExtended()

	// Logs panel (se ativo)
	if m.showLogs {
		logsPanel := m.renderLogsPanel(m.width - 4, 10)
		mainLayout = lipgloss.JoinVertical(lipgloss.Left, mainLayout, "", logsPanel)
	}

	// ═══════════════════════════════════════════════
	// MONTAGEM FINAL
	// ═══════════════════════════════════════════════
	
	result := header
	if notification != "" {
		result += "\n" + notification
	}
	result += "\n\n" + mainLayout + "\n\n" + statusBar + "\n" + commandBar
	
	return result
}

// renderModernSearchPanel renderiza o painel de busca moderno
func (m *Model) renderModernSearchPanel(width, height int) string {
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

// renderModernPlaylistPanel renderiza o painel de playlist moderno
func (m *Model) renderModernPlaylistPanel(width, height int) string {
	title := "📋 Playlist"
	if m.focusedPanel == 1 {
		title = "📋 Playlist [ATIVO]"
	}

	var content strings.Builder
	
	// Mostra modo de playlist
	modeLabel := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Render(m.playlistMode.String())
	content.WriteString(modeLabel + "\n\n")
	
	if len(m.playlist) == 0 {
		content.WriteString(lipgloss.NewStyle().
			Foreground(subtleColor).
			Italic(true).
			Render("Nenhum item na playlist\n\nPressione 'p' para adicionar"))
	} else {
		maxItems := height - 8
		if maxItems > len(m.playlist) {
			maxItems = len(m.playlist)
		}

		for i := 0; i < maxItems; i++ {
			item := m.playlist[i]
			prefix := "  "
			if i == m.playlistIndex && m.playlistPlaying {
				prefix = "▶ "
			}

			line := fmt.Sprintf("%s%d. %s", prefix, i+1, truncate(item.title, width-12))

			// Destaca item selecionado ou tocando
			if (i == m.playlistIndex && m.focusedPanel == 1) || (i == m.playlistIndex && m.playlistPlaying) {
				line = lipgloss.NewStyle().
					Foreground(accentColor).
					Bold(true).
					Render(line)
			}

			content.WriteString(line + "\n")
		}

		// Adiciona instruções
		if m.focusedPanel == 1 {
			instructions := lipgloss.NewStyle().
				Foreground(subtleColor).
				Render("\nSpace: play | d: remover | J/K: mover")
			content.WriteString(instructions)
		}
	}

	return createPanel(title, content.String(), width, height, m.focusedPanel == 1)
}

// renderModernResultsPanel renderiza o painel de resultados moderno
func (m *Model) renderModernResultsPanel(width, height int) string {
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

// renderModernInfoPanel renderiza o painel de info moderno
func (m *Model) renderModernInfoPanel(width, height int) string {
	title := "🎵 Visualizador"
	if m.focusedPanel == 3 {
		title = "🎵 Visualizador [ATIVO]"
	}

	var content strings.Builder

	// Se estiver tocando, mostra informações
	if m.isPlaying {
		content.WriteString(lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Render("▶ TOCANDO\n\n"))

		content.WriteString(lipgloss.NewStyle().
			Foreground(textColor).
			Render(truncate(m.nowPlaying, width-6) + "\n\n"))

		// Visualizador de áudio (cava)
		if m.mode == playModeAudio && len(m.cavaOutput) > 0 {
			for _, line := range m.cavaOutput {
				content.WriteString(line + "\n")
			}
		} else if m.mode == playModeAudio {
			content.WriteString(lipgloss.NewStyle().
				Foreground(subtleColor).
				Render("🎵 Áudio reproduzindo...\n"))
		} else {
			content.WriteString(lipgloss.NewStyle().
				Foreground(subtleColor).
				Render("🎬 Vídeo reproduzindo...\n"))
		}

		content.WriteString("\n" + lipgloss.NewStyle().
			Foreground(subtleColor).
			Render("Pressione 's' para parar"))
	} else {
		// Modo atual
		content.WriteString(lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Render(fmt.Sprintf("Modo: %s\n\n", m.mode.String())))

		// Estatísticas
		content.WriteString(fmt.Sprintf("Resultados: %d\n", len(m.results)))
		content.WriteString(fmt.Sprintf("Playlist: %d itens\n\n", len(m.playlist)))

		// Atalhos principais
		content.WriteString(lipgloss.NewStyle().
			Foreground(subtleColor).
			Render("Atalhos:\n"))
		content.WriteString("Tab - Trocar painel\n")
		content.WriteString("m - Mudar modo\n")
		content.WriteString("p - Add playlist\n")
		content.WriteString("s - Parar música\n")
	}

	return createPanel(title, content.String(), width, height, m.focusedPanel == 3)
}

// renderModernVisualizerPanel renderiza visualizador de áudio moderno
func (m *Model) renderModernVisualizerPanel(width, height int) string {
	title := "🎵 Visualizador de Áudio"
	if m.focusedPanel == 3 {
		title = "🎵 Visualizador [ATIVO]"
	}

	var content strings.Builder

	if m.isPlaying && m.mode == playModeAudio {
		// Renderiza visualizador de áudio
		visualizer := renderVisualizerBars(m.cavaOutput, width-10)
		content.WriteString(visualizer)
		content.WriteString("\n")
		content.WriteString(lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Render("♫ Áudio ao vivo"))
	} else if m.isPlaying {
		content.WriteString(lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Render("🎬 Modo Vídeo\n\n"))
		content.WriteString(lipgloss.NewStyle().
			Foreground(subtleColor).
			Render("O vídeo está sendo reproduzido\nno mpv em janela externa"))
	} else {
		content.WriteString(lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true).
			Render("Visualizador inativo\n\nInicie uma música para ver\na visualização de áudio"))
	}

	return createPanel(title, content.String(), width, height, m.focusedPanel == 3)
}

// truncate trunca uma string se for maior que o limite
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// renderLogsPanel renderiza o painel de logs
func (m *Model) renderLogsPanel(width, height int) string {
	title := "📝 Logs"
	if m.focusedPanel == 4 {
		title = "📝 Logs [ATIVO]"
	}

	var content strings.Builder
	
	if len(m.logs) == 0 {
		content.WriteString(lipgloss.NewStyle().
			Foreground(subtleColor).
			Italic(true).
			Render("Nenhum log ainda"))
	} else {
		// Mostra os últimos logs que cabem no painel
		maxLogs := height - 4
		startIdx := 0
		if len(m.logs) > maxLogs {
			startIdx = len(m.logs) - maxLogs
		}
		
		for i := startIdx; i < len(m.logs); i++ {
			log := m.logs[i]
			timeStr := log.time.Format("15:04:05")
			
			levelStyle := lipgloss.NewStyle()
			var levelIcon string
			
			switch log.level {
			case logError:
				levelStyle = levelStyle.Foreground(lipgloss.Color("#FF0000"))
				levelIcon = "❌"
			case logWarning:
				levelStyle = levelStyle.Foreground(lipgloss.Color("#FFAA00"))
				levelIcon = "⚠️ "
			default:
				levelStyle = levelStyle.Foreground(accentColor)
				levelIcon = "ℹ️ "
			}
			
			line := fmt.Sprintf("[%s] %s %s",
				timeStr,
				levelIcon,
				truncate(log.message, width-20))
			
			content.WriteString(levelStyle.Render(line) + "\n")
		}
	}

	return createPanel(title, content.String(), width, height, m.focusedPanel == 4)
}

// renderCommandBarExtended renderiza barra de comandos estendida
func (m Model) renderCommandBarExtended() string {
	commands := []string{
		keyStyle.Render("[Tab]") + descStyle.Render(" Painel"),
		keyStyle.Render("[m]") + descStyle.Render(" Modo"),
		keyStyle.Render("[r]") + descStyle.Render(" ") + m.playlistMode.String(),
		keyStyle.Render("[Space]") + descStyle.Render(" Play Playlist"),
		keyStyle.Render("[s]") + descStyle.Render(" Parar"),
		keyStyle.Render("[l]") + descStyle.Render(" Logs"),
		keyStyle.Render("[q]") + descStyle.Render(" Sair"),
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Left, commands...)
	return commandBarStyle.Render(bar)
}
