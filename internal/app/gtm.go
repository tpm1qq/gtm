package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	c "github.com/tpm1qq/gtm/internal/pkg/config"
	"github.com/tpm1qq/gtm/internal/pkg/core"
	t "github.com/tpm1qq/gtm/internal/pkg/theme"
)

func RunGTM() error {

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("error loading theme, %w", err)

	}
	themeDir := filepath.Join(cfgDir, "gtm", "themes")
	confPath := filepath.Join(cfgDir, "gtm", "gtm.yaml")
	conf, err := c.GetData(confPath)
	var bkgDirSet bool
	if err != nil {
		switch {
		case conf == c.Data{}:
			return fmt.Errorf("error reading conf data, %w", err)
		case conf.Paths == c.Paths{}:
			fmt.Fprintln(os.Stderr, "path error; no paths set", err)
		case conf.Paths.BackgroundPath == "" && conf.Paths.ToolsPath != "":
			fmt.Println("Warning: BackgroundPath not set!")
			bkgDirSet = false
		case conf.Paths.BackgroundPath != "" && conf.Paths.ToolsPath == "":
			return fmt.Errorf("config error; ToolsPath not set! %w", err)
		}
	} else {
		bkgDirSet = true
	}
	// todo change current settings in conf when calling tools api
	var tools core.ToolList
	var theme string
	var global bool
	var color string
	var wallpaper string
	var settings core.ToolSettings
	var selection core.ToolList
	var all = core.ToolList{
		core.ToolHyprland,
		core.ToolWaybar,
		core.ToolRofi,
		core.ToolHyprpaper,
	}
	var allNoBkg = core.ToolList{
		core.ToolHyprland,
		core.ToolWaybar,
		core.ToolRofi,
	}

	// this still just uses a global color for everything despite giving the illusion of supporting multiple colors and seperate settings per tool.
	// the config and themes track settings seperatly but the application logic still uses a single settings struct made from either the color flag or the hyprland tool color field.
	// per tool functionality *can* be achieved currently by just running the tool multiple times using different tool flags each time.
	flag.BoolVar(&global, "global", false, "global flag; apply config changes to all tools at the same time")
	flag.Var(&tools, "tools", "which tool's config the user wants to change")
	flag.StringVar(&color, "color", "", "change the color of the given tool(s)")
	flag.StringVar(&wallpaper, "wallpaper", "", "change the current wallpaper")
	flag.StringVar(&theme, "theme", "", "use a gtm theme")

	flag.Parse()

	if len(theme) > 0 {
		themePath := filepath.Join(themeDir, theme+".yaml")
		v, err := t.LoadTheme(themePath)
		if err != nil {
			return fmt.Errorf("error loading theme, %w", err)
		}
		if v.Global.Enabled {
			selection = all
			settings = core.ToolSettings{
				Color:     v.Global.Color,
				Wallpaper: v.Global.Wallpaper,
			}
		} else {
			settings = core.ToolSettings{
				Color:     v.Tools.Hyprland.Color,
				Wallpaper: v.Tools.Hyprpaper.Wallpaper,
			}
			if v.Tools.Hyprland.Enabled {
				selection = append(selection, core.ToolHyprland)
			}
			if v.Tools.Hyprpaper.Enabled && bkgDirSet {
				selection = append(selection, core.ToolHyprpaper)
			}
			if v.Tools.Waybar.Enabled {
				selection = append(selection, core.ToolWaybar)
			}
			if v.Tools.Rofi.Enabled {
				selection = append(selection, core.ToolRofi)
			}
		}
	} else {
		settings = core.ToolSettings{
			Color:     color,
			Wallpaper: wallpaper,
		}

		switch {
		case global && len(tools) == 0:
			if bkgDirSet {
				selection = all
			} else {
				selection = allNoBkg
			}
		case !global && len(tools) > 0:
			selection = tools
		case !global && len(tools) == 0:
			return fmt.Errorf("neither global nor tool flag set %w", err)
		default:
			return fmt.Errorf("flag error, %w", err)
		}
	}

	var wg = sync.WaitGroup{}
	var errs = make(chan error, len(selection))
	for _, v := range selection {
		curr := v
		wg.Go(func() {
			if err := core.ApplyChanges(curr, settings, conf.Paths); err != nil {
				errs <- err
			}
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			fmt.Fprintln(os.Stderr, "error setting color", err)
		}
	}
	//updating config with current settings
	for _, v := range selection {
		switch v {
		case core.ToolHyprland:
			conf.CurrentSettings.CurrentTheme.Hyprland.Color = settings.Color
		case core.ToolWaybar:
			conf.CurrentSettings.CurrentTheme.Waybar.Color = settings.Color
		case core.ToolRofi:
			conf.CurrentSettings.CurrentTheme.Rofi.Color = settings.Color
		case core.ToolHyprpaper:
			conf.CurrentSettings.CurrentTheme.Hyprpaper.Wallpaper = settings.Wallpaper
		}
	}
	if err := c.WriteData(conf, confPath); err != nil {
		return err
	}
	return nil
}
