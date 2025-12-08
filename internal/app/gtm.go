package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/tpm1qq/gtm/internal/pkg/core"
	t "github.com/tpm1qq/gtm/internal/pkg/theme"
)

func RunGTM() {
	var homeDir, _ = os.UserHomeDir()
	var cfgDir = filepath.Join(homeDir, ".config", "gtm", "themes")
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

	flag.BoolVar(&global, "global", false, "global flag; apply config changes to all tools at the same time")
	flag.BoolVar(&global, "g", false, "global flag; apply config changes to all tools at the same time")
	flag.Var(&tools, "tools", "which tool's config the user wants to change")
	flag.StringVar(&color, "color", "", "change the color of the given tool(s)")
	flag.StringVar(&wallpaper, "wallpaper", "", "change the current wallpaper")
	flag.StringVar(&theme, "theme", "", "use a gtm theme")

	flag.Parse()

	if len(theme) > 0 {
		themePath := filepath.Join(cfgDir, theme+".yaml")
		v, err := t.LoadTheme(themePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error loading theme:", err)
			return
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
			if v.Tools.Hyprpaper.Enabled {
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
			selection = all
		case !global && len(tools) > 0:
			selection = tools
		case !global && len(tools) == 0:
			fmt.Println("neither global nor tool flag set!")
			return
		default:
			fmt.Println("flag error")
			return
		}
	}

	var wg = sync.WaitGroup{}
	var errs = make(chan error, len(selection))
	for _, v := range selection {
		curr := v
		wg.Go(func() {
			if err := core.ApplyChanges(curr, settings); err != nil {
				errs <- err
			}
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			fmt.Fprintln(os.Stderr, "error setting color:", err)
		}
	}
}
