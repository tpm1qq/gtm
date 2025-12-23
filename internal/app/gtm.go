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
	var tools core.ToolList
	var theme string
	var global bool
	var color string
	var wallpaper string
	var settings = make(map[core.ToolName]core.ToolSettings)
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

	// per tool functionality using the cli *can* be achieved currently by just running the tool multiple times using different tool flags each time.
	flag.BoolVar(&global, "global", false, "global flag; apply config changes to all tools at the same time")
	flag.Var(&tools, "tools", "which tool's config the user wants to change")
	flag.StringVar(&color, "color", "", "change the color of the given tool(s)")
	flag.StringVar(&wallpaper, "wallpaper", "", "change the current wallpaper")
	flag.StringVar(&theme, "theme", "", "use a gtm theme")

	flag.Parse()

	//filling the settings this way might not be the most efficient.
	//should work fine though under the assumption that no theme means
	//we only get a single color from the commandline, no matter how many tools are selected
	if len(theme) > 0 {
		themePath := filepath.Join(themeDir, theme+".yaml")
		v, err := t.LoadTheme(themePath)
		if err != nil {
			return fmt.Errorf("error loading theme, %w", err)
		}
		if v.Global.Enabled {
			selection = all
			settings[core.ToolHyprland] = core.ToolSettings{
				Color: v.Global.Color,
			}
			settings[core.ToolHyprpaper] = core.ToolSettings{
				Wallpaper: v.Global.Wallpaper,
			}
			settings[core.ToolWaybar] = core.ToolSettings{
				Color: v.Global.Color,
			}
			settings[core.ToolRofi] = core.ToolSettings{
				Color: v.Global.Color,
			}
		} else {
			if v.Tools.Hyprland.Enabled {
				settings[core.ToolHyprland] = core.ToolSettings{
					Color: v.Tools.Hyprland.Color,
				}
				selection = append(selection, core.ToolHyprland)
			}
			if v.Tools.Hyprpaper.Enabled && bkgDirSet {
				settings[core.ToolHyprpaper] = core.ToolSettings{
					Wallpaper: v.Tools.Hyprpaper.Wallpaper,
				}
				selection = append(selection, core.ToolHyprpaper)
			}
			if v.Tools.Waybar.Enabled {
				settings[core.ToolWaybar] = core.ToolSettings{
					Color: v.Tools.Waybar.Color,
				}
				selection = append(selection, core.ToolWaybar)
			}
			if v.Tools.Rofi.Enabled {
				settings[core.ToolRofi] = core.ToolSettings{
					Color: v.Tools.Rofi.Color,
				}
				selection = append(selection, core.ToolRofi)
			}
		}
	} else {
		settings[core.ToolHyprland] = core.ToolSettings{
			Color: color,
		}
		settings[core.ToolHyprpaper] = core.ToolSettings{
			Wallpaper: wallpaper,
		}
		settings[core.ToolWaybar] = core.ToolSettings{
			Color: color,
		}
		settings[core.ToolRofi] = core.ToolSettings{
			Color: color,
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
			conf.CurrentSettings.CurrentTheme.Hyprland.Color = settings[core.ToolHyprland].Color
		case core.ToolWaybar:
			conf.CurrentSettings.CurrentTheme.Waybar.Color = settings[core.ToolWaybar].Color
		case core.ToolRofi:
			conf.CurrentSettings.CurrentTheme.Rofi.Color = settings[core.ToolRofi].Color
		case core.ToolHyprpaper:
			conf.CurrentSettings.CurrentTheme.Hyprpaper.Wallpaper = settings[core.ToolHyprpaper].Wallpaper
		}
	}
	if err := c.WriteData(conf, confPath); err != nil {
		return err
	}
	return nil
}
