package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	System  SystemSettings  `yaml:"system"`
	Current CurrentSettings `yaml:"current"`
}

type SystemSettings struct {
	ToolsDir      string `yaml:"toolsdir"`
	ConfigDir     string `yaml:"configdir"`
	BackgroundDir string `yaml:"backgrounddir"`
}

type CurrentSettings struct {
	// Do I really need this? Wouldn't just having the CurrentTheme in here suffice?
	// perhaps add something related to the setup functionality here in case I need to save something related to that later, might as well just leave this here for the time being
	CurrentTheme CurrentTheme `yaml:"tools"`
}

type CurrentTheme struct {
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

func LoadConfig(p string) (Config, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return Config{}, fmt.Errorf("error reading file: %w", err)
	}
	var r Config
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return Config{}, fmt.Errorf("error yaml.Unmarshal: %w", err)
	}
	return r, nil
}
