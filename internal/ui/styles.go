package ui

import (
	"strings"
	
	"github.com/charmbracelet/lipgloss"
)

var (
	// Paleta de cores moderna (inspirada em Spotify/Apple Music)
	primaryColor   = lipgloss.Color("#1DB954") // verde Spotify
	secondaryColor = lipgloss.Color("#1ed760") // verde claro
	accentColor    = lipgloss.Color("#FF006E") // rosa vibrante
	bgColor        = lipgloss.Color("#121212") // preto suave
	surfaceColor   = lipgloss.Color("#282828") // cinza escuro
	textColor      = lipgloss.Color("#FFFFFF") // branco
	subtleColor    = lipgloss.Color("#B3B3B3") // cinza médio
	dimColor       = lipgloss.Color("#535353") // cinza escuro
	errorColor     = lipgloss.Color("#FF4444") // vermelho erro
	warningColor   = lipgloss.Color("#FFB84D") // laranja warning
	successColor   = lipgloss.Color("#1DB954") // verde sucesso

	// Estilos principais
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1)

	// Player principal - destaque máximo
	nowPlayingStyle = lipgloss.NewStyle().
			Background(surfaceColor).
			Foreground(textColor).
			Bold(true).
			Padding(1, 3).
			MarginTop(1).
			MarginBottom(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Width(50).
			Align(lipgloss.Center)

	// Painel moderno com sombra
	panelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(dimColor).
			Padding(1, 2).
			MarginRight(1)

	// Painel focado - destaque com cor primária
	focusedPanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				MarginRight(1)

	// Barra de status moderna
	statusStyle = lipgloss.NewStyle().
			Background(surfaceColor).
			Foreground(subtleColor).
			Padding(0, 2).
			MarginTop(1)

	// Controles de mídia
	controlButtonStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(surfaceColor).
				Padding(0, 2).
				Bold(true)

	controlButtonActiveStyle = lipgloss.NewStyle().
					Foreground(primaryColor).
					Background(surfaceColor).
					Padding(0, 2).
					Bold(true)

	// Lista de músicas
	listItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(2)

	listItemSelectedStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				PaddingLeft(1)

	// Tags e badges
	badgeStyle = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(primaryColor).
			Padding(0, 1).
			Bold(true)

	// Teclas de atalho
	keyStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			PaddingRight(1)

	descStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			PaddingRight(2)

	// Barra de comandos inferior
	commandBarStyle = lipgloss.NewStyle().
			Background(surfaceColor).
			Foreground(subtleColor).
			Padding(1, 2)

	// Progress bar
	progressFilledStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(dimColor)

	// Visualizador de áudio
	visualizerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center)

	// Input de busca
	inputBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		PaddingLeft(1).
		PaddingRight(1).
		MarginBottom(1)

	// Estilo da lista de resultados
	listBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		PaddingLeft(1).
		PaddingRight(1)

	// Estilo do modo de reprodução
	modeStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		PaddingLeft(2).
		PaddingRight(2)

	// Estilo de erro
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)
)

func renderCommandBar(mode string) string {
	commands := []string{
		keyStyle.Render("[Tab]") + descStyle.Render(" Trocar painel"),
		keyStyle.Render("[Enter]") + descStyle.Render(" Buscar/Reproduzir"),
		keyStyle.Render("[m]") + descStyle.Render(" Modo: ") + modeStyle.Render(mode),
		keyStyle.Render("[p]") + descStyle.Render(" Add Playlist"),
		keyStyle.Render("[q]") + descStyle.Render(" Sair"),
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Left, commands...)
	return commandBarStyle.Render(bar)
}

// createPanel cria um painel com borda e título
func createPanel(title, content string, width, height int, focused bool) string {
	borderColor := dimColor
	borderStyle := lipgloss.RoundedBorder()
	
	if focused {
		borderColor = primaryColor
		borderStyle = lipgloss.ThickBorder()
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColor).
		Width(width).
		Height(height).
		Padding(1)

	// Adiciona título com ícone
	titleStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Bold(true).
		Background(surfaceColor).
		Padding(0, 1)

	panelContent := titleStyle.Render(title) + "\n\n" + content

	return style.Render(panelContent)
}

// renderProgressBar renderiza uma barra de progresso visual
func renderProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}
	
	percentage := float64(current) / float64(total)
	filled := int(float64(width) * percentage)
	empty := width - filled
	
	bar := progressFilledStyle.Render(strings.Repeat("█", filled))
	bar += progressEmptyStyle.Render(strings.Repeat("░", empty))
	
	return bar
}

// renderPlayButton renderiza botão de play/pause
func renderPlayButton(isPlaying bool) string {
	if isPlaying {
		return controlButtonActiveStyle.Render("⏸ PAUSAR")
	}
	return controlButtonStyle.Render("▶ PLAY")
}

// renderModeIndicator renderiza indicador de modo visual
func renderModeIndicator(mode, playlistMode string) string {
	modeIcon := "🎵"
	if mode == "MP4 (Vídeo)" {
		modeIcon = "🎬"
	}
	
	modeBadge := badgeStyle.Render(" " + modeIcon + " " + mode + " ")
	playlistBadge := lipgloss.NewStyle().
		Foreground(textColor).
		Background(dimColor).
		Padding(0, 1).
		Bold(true).
		Render(" " + playlistMode + " ")
	
	return lipgloss.JoinHorizontal(lipgloss.Center, modeBadge, " ", playlistBadge)
}

// renderVolumeBar renderiza indicador de volume (visual apenas)
func renderVolumeBar() string {
	return lipgloss.NewStyle().
		Foreground(subtleColor).
		Render("🔊 ████████░░")
}

// renderVisualizerBars renderiza barras do visualizador de áudio
func renderVisualizerBars(lines []string, width int) string {
	if len(lines) == 0 {
		// Visualização padrão animada
		bars := []string{
			"    ▁▂▃▅▆▇█▇▆▅▃▂▁    ",
			"  ▂▄▅▇█▇▅▄▂    ▂▄▅▇  ",
			"▁▃▅▇█▇▅▃▁  ▁▃▅▇█▇▅▃▁",
		}
		
		var result strings.Builder
		for _, bar := range bars {
			result.WriteString(visualizerStyle.Render(bar) + "\n")
		}
		return result.String()
	}
	
	var result strings.Builder
	for _, line := range lines {
		result.WriteString(visualizerStyle.Render(line) + "\n")
	}
	return result.String()
}
