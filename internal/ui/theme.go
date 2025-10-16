package ui

import "github.com/gdamore/tcell/v2"

type Theme struct {
	Name string
	
	Rosewater tcell.Color
	Flamingo  tcell.Color
	Pink      tcell.Color
	Mauve     tcell.Color
	Red       tcell.Color
	Maroon    tcell.Color
	Peach     tcell.Color
	Yellow    tcell.Color
	Green     tcell.Color
	Teal      tcell.Color
	Sky       tcell.Color
	Sapphire  tcell.Color
	Blue      tcell.Color
	Lavender  tcell.Color
	Text      tcell.Color
	Subtext1  tcell.Color
	Subtext0  tcell.Color
	Overlay2  tcell.Color
	Overlay1  tcell.Color
	Overlay0  tcell.Color
	Surface2  tcell.Color
	Surface1  tcell.Color
	Surface0  tcell.Color
	Base      tcell.Color
	Mantle    tcell.Color
	Crust     tcell.Color
}

var CatppuccinMocha = Theme{
	Name:      "Catppuccin Mocha",
	Rosewater: tcell.NewHexColor(0xf5e0dc),
	Flamingo:  tcell.NewHexColor(0xf2cdcd),
	Pink:      tcell.NewHexColor(0xf5c2e7),
	Mauve:     tcell.NewHexColor(0xcba6f7),
	Red:       tcell.NewHexColor(0xf38ba8),
	Maroon:    tcell.NewHexColor(0xeba0ac),
	Peach:     tcell.NewHexColor(0xfab387),
	Yellow:    tcell.NewHexColor(0xf9e2af),
	Green:     tcell.NewHexColor(0xa6e3a1),
	Teal:      tcell.NewHexColor(0x94e2d5),
	Sky:       tcell.NewHexColor(0x89dceb),
	Sapphire:  tcell.NewHexColor(0x74c7ec),
	Blue:      tcell.NewHexColor(0x89b4fa),
	Lavender:  tcell.NewHexColor(0xb4befe),
	Text:      tcell.NewHexColor(0xcdd6f4),
	Subtext1:  tcell.NewHexColor(0xbac2de),
	Subtext0:  tcell.NewHexColor(0xa6adc8),
	Overlay2:  tcell.NewHexColor(0x9399b2),
	Overlay1:  tcell.NewHexColor(0x7f849c),
	Overlay0:  tcell.NewHexColor(0x6c7086),
	Surface2:  tcell.NewHexColor(0x585b70),
	Surface1:  tcell.NewHexColor(0x45475a),
	Surface0:  tcell.NewHexColor(0x313244),
	Base:      tcell.NewHexColor(0x1e1e2e),
	Mantle:    tcell.NewHexColor(0x181825),
	Crust:     tcell.NewHexColor(0x11111b),
}

func (t *Theme) GetFocusedBorder() tcell.Color {
	return t.Blue
}

func (t *Theme) GetUnfocusedBorder() tcell.Color {
	return t.Surface0
}

func (t *Theme) GetBackground() tcell.Color {
	return t.Base
}

func (t *Theme) GetPrimaryText() tcell.Color {
	return t.Text
}

func (t *Theme) GetSecondaryText() tcell.Color {
	return t.Subtext1
}

func (t *Theme) GetSuccess() tcell.Color {
	return t.Green
}

func (t *Theme) GetError() tcell.Color {
	return t.Red
}

func (t *Theme) GetWarning() tcell.Color {
	return t.Yellow
}

func (t *Theme) GetAccent() tcell.Color {
	return t.Mauve
}

func (t *Theme) GetHighlight() tcell.Color {
	return t.Sapphire
}
