package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Paths struct {
	BackgroundPath string
	ToolsPath      string
}

type Data struct {
	Paths           Paths
	CurrentSettings currentSettings
}

type config struct {
	System  systemSettings  `yaml:"system"`
	Current currentSettings `yaml:"current"`
}

type systemSettings struct {
	ToolsDir      string `yaml:"toolsdir"`
	BackgroundDir string `yaml:"backgrounddir"`
}

type currentSettings struct {
	// Do I really need this? Wouldn't just having the CurrentTheme in here suffice?
	// perhaps add something related to the setup functionality here in case I need to save something related to that later, might as well just leave this here for the time being
	CurrentTheme currentTheme `yaml:"tools"`
}

type currentTheme struct {
	Hyprland  hyprlandBlock  `yaml:"hyprland"`
	Hyprpaper hyprpaperBlock `yaml:"hyprpaper"`
	Waybar    waybarBlock    `yaml:"waybar"`
	Rofi      rofiBlock      `yaml:"rofi"`
}

type hyprlandBlock struct {
	Enabled bool   `yaml:"enabled"`
	Color   string `yaml:"color"`
}

type hyprpaperBlock struct {
	Enabled   bool   `yaml:"enabled"`
	Wallpaper string `yaml:"wallpaper"`
}

type waybarBlock struct {
	Enabled bool   `yaml:"enabled"`
	Color   string `yaml:"color"`
}

type rofiBlock struct {
	Enabled bool   `yaml:"enabled"`
	Color   string `yaml:"color"`
}

func loadConfig(p string) (config, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return config{}, fmt.Errorf("error reading file %w", err)
	}
	var r config
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return config{}, fmt.Errorf("error yaml.Unmarshal %w", err)
	}
	return r, nil
}
func getPaths(c config) (Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Paths{}, fmt.Errorf("error getting home dir %w", err)
	}
	var p Paths
	var errs []error
	curr := c.System.ToolsDir
	if curr == "" {
		errs = append(errs, fmt.Errorf("toolsDir is emtpy"))
	}
	curr, found := strings.CutPrefix(curr, "~/")
	if found {
		curr = filepath.Join(home, curr)
	}
	p.ToolsPath, err = filepath.Abs(curr)
	if err != nil {
		return Paths{}, fmt.Errorf("error getting absolute path %w", err)
	}
	curr = c.System.BackgroundDir
	if curr == "" {
		errs = append(errs, fmt.Errorf("backgroundDir empty"))
	}
	curr, found = strings.CutPrefix(curr, "~/")
	if found {
		curr = filepath.Join(home, curr)
	}
	p.BackgroundPath, err = filepath.Abs(curr)
	if err != nil {
		return Paths{}, fmt.Errorf("error getting absolute path %w", err)
	}
	return p, errors.Join(errs...)
}
func GetData(path string) (Data, error) {
	var errs []error
	c, err := loadConfig(path)
	if err != nil {
		return Data{}, fmt.Errorf("error loading config %w", err)
	}
	p, err := getPaths(c)
	if err != nil {
		errs = append(errs, fmt.Errorf("error loading paths %w", err))
	}
	var d Data
	d.CurrentSettings = c.Current
	d.Paths = p
	return d, errors.Join(errs...)
}
func WriteData(data Data, path string) error {
	//this is bad
	cfg, err := loadConfig(path)
	if err != nil {
		return fmt.Errorf("error loading config %w", err)
	}
	cfg.Current = data.CurrentSettings
	v, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, v, 0644); err != nil {
		return err
	}
	return nil
}
