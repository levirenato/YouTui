# YouTui Themes

## Built-in Themes

YouTui comes with 4 beautiful Catppuccin themes:

- **🌻 Latte** - Light mode (catppuccin-latte)
- **🪴 Frappé** - Dark cool (catppuccin-frappe)
- **🌺 Macchiato** - Dark warm (catppuccin-macchiato)
- **🌿 Mocha** - Dark deep (catppuccin-mocha) [default]

## Switching Themes

1. Press **Ctrl+C** to open settings
2. Select **Theme** button
3. Cycle through themes by pressing it multiple times
4. Theme is saved automatically to `~/.config/youtui/youtui.conf`

## Configuration File

Location: `~/.config/youtui/youtui.conf`

Example:
```toml
[theme]
active = "catppuccin-mocha"
```

Available values:
- `catppuccin-latte`
- `catppuccin-frappe`
- `catppuccin-macchiato`
- `catppuccin-mocha`
- `custom` (requires custom_path)

## Custom Themes

To use a custom theme:

1. Create a TOML file with your theme colors (see `custom-theme.toml.example`)
2. Edit `~/.config/youtui/youtui.conf`:

```toml
[theme]
active = "custom"
custom_path = "/path/to/your/theme.toml"
```

### Custom Theme Format

```toml
name = "My Theme"

rosewater = "#f5e0dc"
flamingo = "#f2cdcd"
pink = "#f5c2e7"
mauve = "#cba6f7"
red = "#f38ba8"
maroon = "#eba0ac"
peach = "#fab387"
yellow = "#f9e2af"
green = "#a6e3a1"
teal = "#94e2d5"
sky = "#89dceb"
sapphire = "#74c7ec"
blue = "#89b4fa"
lavender = "#b4befe"

text = "#cdd6f4"
subtext1 = "#bac2de"
subtext0 = "#a6adc8"

overlay2 = "#9399b2"
overlay1 = "#7f849c"
overlay0 = "#6c7086"

surface2 = "#585b70"
surface1 = "#45475a"
surface0 = "#313244"

base = "#1e1e2e"
mantle = "#181825"
crust = "#11111b"
```

## Tips

- **Light mode**: Use Latte theme for comfortable daytime viewing
- **Dark modes**: Choose based on preference (Frappé is cooler, Macchiato warmer, Mocha deepest)
- **Contrast**: All themes have WCAG AA compliant contrast ratios
