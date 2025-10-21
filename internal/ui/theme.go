package ui

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gdamore/tcell/v2"
)

type Theme struct {
	Name      string
	ID        string
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

type ThemeTOML struct {
	Name      string `toml:"name"`
	Rosewater string `toml:"rosewater"`
	Flamingo  string `toml:"flamingo"`
	Pink      string `toml:"pink"`
	Mauve     string `toml:"mauve"`
	Red       string `toml:"red"`
	Maroon    string `toml:"maroon"`
	Peach     string `toml:"peach"`
	Yellow    string `toml:"yellow"`
	Green     string `toml:"green"`
	Teal      string `toml:"teal"`
	Sky       string `toml:"sky"`
	Sapphire  string `toml:"sapphire"`
	Blue      string `toml:"blue"`
	Lavender  string `toml:"lavender"`
	Text      string `toml:"text"`
	Subtext1  string `toml:"subtext1"`
	Subtext0  string `toml:"subtext0"`
	Overlay2  string `toml:"overlay2"`
	Overlay1  string `toml:"overlay1"`
	Overlay0  string `toml:"overlay0"`
	Surface2  string `toml:"surface2"`
	Surface1  string `toml:"surface1"`
	Surface0  string `toml:"surface0"`
	Base      string `toml:"base"`
	Mantle    string `toml:"mantle"`
	Crust     string `toml:"crust"`
}

var CatppuccinLatte = Theme{
	Name:      "🌻 Latte",
	ID:        "catppuccin-latte",
	Rosewater: tcell.NewHexColor(0xdc8a78),
	Flamingo:  tcell.NewHexColor(0xdd7878),
	Pink:      tcell.NewHexColor(0xea76cb),
	Mauve:     tcell.NewHexColor(0x8839ef),
	Red:       tcell.NewHexColor(0xd20f39),
	Maroon:    tcell.NewHexColor(0xe64553),
	Peach:     tcell.NewHexColor(0xfe640b),
	Yellow:    tcell.NewHexColor(0xdf8e1d),
	Green:     tcell.NewHexColor(0x40a02b),
	Teal:      tcell.NewHexColor(0x179299),
	Sky:       tcell.NewHexColor(0x04a5e5),
	Sapphire:  tcell.NewHexColor(0x209fb5),
	Blue:      tcell.NewHexColor(0x1e66f5),
	Lavender:  tcell.NewHexColor(0x7287fd),
	Text:      tcell.NewHexColor(0x303030),
	Subtext1:  tcell.NewHexColor(0x404040),
	Subtext0:  tcell.NewHexColor(0x505050),
	Overlay2:  tcell.NewHexColor(0x606060),
	Overlay1:  tcell.NewHexColor(0x707070),
	Overlay0:  tcell.NewHexColor(0x808080),
	Surface2:  tcell.NewHexColor(0x909090),
	Surface1:  tcell.NewHexColor(0xa0a0a0),
	Surface0:  tcell.NewHexColor(0xb0b0b0),
	Base:      tcell.NewHexColor(0xeff1f5),
	Mantle:    tcell.NewHexColor(0xe6e9ef),
	Crust:     tcell.NewHexColor(0xdce0e8),
}

var CatppuccinFrappe = Theme{
	Name:      "🪴 Frappé",
	ID:        "catppuccin-frappe",
	Rosewater: tcell.NewHexColor(0xf2d5cf),
	Flamingo:  tcell.NewHexColor(0xeebebe),
	Pink:      tcell.NewHexColor(0xf4b8e4),
	Mauve:     tcell.NewHexColor(0xca9ee6),
	Red:       tcell.NewHexColor(0xe78284),
	Maroon:    tcell.NewHexColor(0xea999c),
	Peach:     tcell.NewHexColor(0xef9f76),
	Yellow:    tcell.NewHexColor(0xe5c890),
	Green:     tcell.NewHexColor(0xa6d189),
	Teal:      tcell.NewHexColor(0x81c8be),
	Sky:       tcell.NewHexColor(0x99d1db),
	Sapphire:  tcell.NewHexColor(0x85c1dc),
	Blue:      tcell.NewHexColor(0x8caaee),
	Lavender:  tcell.NewHexColor(0xbabbf1),
	Text:      tcell.NewHexColor(0xc6d0f5),
	Subtext1:  tcell.NewHexColor(0xb5bfe2),
	Subtext0:  tcell.NewHexColor(0xa5adce),
	Overlay2:  tcell.NewHexColor(0x949cbb),
	Overlay1:  tcell.NewHexColor(0x838ba7),
	Overlay0:  tcell.NewHexColor(0x737994),
	Surface2:  tcell.NewHexColor(0x626880),
	Surface1:  tcell.NewHexColor(0x51576d),
	Surface0:  tcell.NewHexColor(0x414559),
	Base:      tcell.NewHexColor(0x303446),
	Mantle:    tcell.NewHexColor(0x292c3c),
	Crust:     tcell.NewHexColor(0x232634),
}

var CatppuccinMacchiato = Theme{
	Name:      "🌺 Macchiato",
	ID:        "catppuccin-macchiato",
	Rosewater: tcell.NewHexColor(0xf4dbd6),
	Flamingo:  tcell.NewHexColor(0xf0c6c6),
	Pink:      tcell.NewHexColor(0xf5bde6),
	Mauve:     tcell.NewHexColor(0xc6a0f6),
	Red:       tcell.NewHexColor(0xed8796),
	Maroon:    tcell.NewHexColor(0xee99a0),
	Peach:     tcell.NewHexColor(0xf5a97f),
	Yellow:    tcell.NewHexColor(0xeed49f),
	Green:     tcell.NewHexColor(0xa6da95),
	Teal:      tcell.NewHexColor(0x8bd5ca),
	Sky:       tcell.NewHexColor(0x91d7e3),
	Sapphire:  tcell.NewHexColor(0x7dc4e4),
	Blue:      tcell.NewHexColor(0x8aadf4),
	Lavender:  tcell.NewHexColor(0xb7bdf8),
	Text:      tcell.NewHexColor(0xcad3f5),
	Subtext1:  tcell.NewHexColor(0xb8c0e0),
	Subtext0:  tcell.NewHexColor(0xa5adcb),
	Overlay2:  tcell.NewHexColor(0x939ab7),
	Overlay1:  tcell.NewHexColor(0x8087a2),
	Overlay0:  tcell.NewHexColor(0x6e738d),
	Surface2:  tcell.NewHexColor(0x5b6078),
	Surface1:  tcell.NewHexColor(0x494d64),
	Surface0:  tcell.NewHexColor(0x363a4f),
	Base:      tcell.NewHexColor(0x24273a),
	Mantle:    tcell.NewHexColor(0x1e2030),
	Crust:     tcell.NewHexColor(0x181926),
}

var CatppuccinMocha = Theme{
	Name:      "🌿 Mocha",
	ID:        "catppuccin-mocha",
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

func GetAllThemes() []Theme {
	return []Theme{
		CatppuccinLatte,
		CatppuccinFrappe,
		CatppuccinMacchiato,
		CatppuccinMocha,
	}
}

func GetThemeByID(id string) *Theme {
	themes := GetAllThemes()
	for i := range themes {
		if themes[i].ID == id {
			return &themes[i]
		}
	}
	return &CatppuccinMocha
}

func LoadCustomTheme(path string) (*Theme, error) {
	var themeTOML ThemeTOML
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, &themeTOML); err != nil {
		return nil, err
	}

	theme := &Theme{
		Name:      themeTOML.Name,
		ID:        "custom",
		Rosewater: parseHexColor(themeTOML.Rosewater),
		Flamingo:  parseHexColor(themeTOML.Flamingo),
		Pink:      parseHexColor(themeTOML.Pink),
		Mauve:     parseHexColor(themeTOML.Mauve),
		Red:       parseHexColor(themeTOML.Red),
		Maroon:    parseHexColor(themeTOML.Maroon),
		Peach:     parseHexColor(themeTOML.Peach),
		Yellow:    parseHexColor(themeTOML.Yellow),
		Green:     parseHexColor(themeTOML.Green),
		Teal:      parseHexColor(themeTOML.Teal),
		Sky:       parseHexColor(themeTOML.Sky),
		Sapphire:  parseHexColor(themeTOML.Sapphire),
		Blue:      parseHexColor(themeTOML.Blue),
		Lavender:  parseHexColor(themeTOML.Lavender),
		Text:      parseHexColor(themeTOML.Text),
		Subtext1:  parseHexColor(themeTOML.Subtext1),
		Subtext0:  parseHexColor(themeTOML.Subtext0),
		Overlay2:  parseHexColor(themeTOML.Overlay2),
		Overlay1:  parseHexColor(themeTOML.Overlay1),
		Overlay0:  parseHexColor(themeTOML.Overlay0),
		Surface2:  parseHexColor(themeTOML.Surface2),
		Surface1:  parseHexColor(themeTOML.Surface1),
		Surface0:  parseHexColor(themeTOML.Surface0),
		Base:      parseHexColor(themeTOML.Base),
		Mantle:    parseHexColor(themeTOML.Mantle),
		Crust:     parseHexColor(themeTOML.Crust),
	}

	return theme, nil
}

func parseHexColor(hex string) tcell.Color {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	
	var r, g, b uint8
	if len(hex) == 6 {
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		return tcell.NewRGBColor(int32(r), int32(g), int32(b))
	}
	
	return tcell.ColorDefault
}
