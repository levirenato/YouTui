package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Cores principais
	primaryColor   = lipgloss.Color("#FF0000")
	secondaryColor = lipgloss.Color("#282828")
	accentColor    = lipgloss.Color("#00FF00")
	textColor      = lipgloss.Color("#FFFFFF")
	subtleColor    = lipgloss.Color("#666666")

	// Estilo do título principal
	titleStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		PaddingLeft(2).
		PaddingRight(2)

	// Estilo da barra de comandos
	commandBarStyle = lipgloss.NewStyle().
		Background(secondaryColor).
		Foreground(textColor).
		PaddingLeft(1).
		PaddingRight(1).
		MarginTop(1)

	// Estilo para teclas de atalho
	keyStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)

	// Estilo para descrição de comandos
	descStyle = lipgloss.NewStyle().
		Foreground(textColor)

	// Estilo do input de busca
	inputBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		PaddingLeft(1).
		PaddingRight(1).
		MarginTop(1).
		MarginBottom(1)

	// Estilo da lista de resultados
	listBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		PaddingLeft(1).
		PaddingRight(1)

	// Estilo do status
	statusStyle = lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true).
		PaddingLeft(2)

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
	borderColor := secondaryColor
	if focused {
		borderColor = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width).
		Height(height).
		Padding(1)

	// Adiciona título
	titleStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Bold(true)

	panelContent := titleStyle.Render(title) + "\n\n" + content

	return style.Render(panelContent)
}
