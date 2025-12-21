package theme

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Theme struct {
	Global GlobalBlock `yaml:"global"`
	Tools  ToolsBlock  `yaml:"tools"`
}

type GlobalBlock struct {
	Enabled   bool   `yaml:"enabled"`
	Wallpaper string `yaml:"wallpaper"`
	Color     string `yaml:"color"`
}

type ToolsBlock struct {
	Hyprland  HyprlandBlock  `yaml:"hyprland"`
	Hyprpaper HyprpaperBlock `yaml:"hyprpaper"`
	Waybar    WaybarBlock    `yaml:"waybar"`
	Rofi      RofiBlock      `yaml:"rofi"`
}

type HyprlandBlock struct {
	Enabled bool   `yaml:"enabled"`
	Color   string `yaml:"color"`
}

type HyprpaperBlock struct {
	Enabled   bool   `yaml:"enabled"`
	Wallpaper string `yaml:"wallpaper"`
}

type WaybarBlock struct {
	Enabled bool   `yaml:"enabled"`
	Color   string `yaml:"color"`
}

type RofiBlock struct {
	Enabled bool   `yaml:"enabled"`
	Color   string `yaml:"color"`
}

func LoadTheme(p string) (Theme, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return Theme{}, fmt.Errorf("error reading file %w", err)
	}
	var r Theme
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return Theme{}, fmt.Errorf("error yaml.Unmarshal %w", err)
	}
	return r, nil
}
