package app

import (
	"flag"
	"fmt"
	"github.com/tpm1qq/gtm/internal/pkg/core"
	"os"
	"sync"
)

func RunGTM() {
	var tools core.ToolList
	var global bool
	var color string
	var wallpaper string
	var settings core.ToolSettings
	var selection core.ToolList

	flag.BoolVar(&global, "global", false, "global flag; apply config changes to all tools at the same time")
	flag.BoolVar(&global, "g", false, "global flag; apply config changes to all tools at the same time")
	flag.Var(&tools, "package", "which tool's config the user wants to change")
	flag.Var(&tools, "p", "which tool's config the user wants to change")
	flag.StringVar(&color, "color", "", "change the color of the given tool(s)")
	flag.StringVar(&wallpaper, "wallpaper", "", "change the current wallpaper")
	flag.StringVar(&wallpaper, "w", "", "change the current wallpaper")

	flag.Parse()

	settings = core.ToolSettings{
		Color:     color,
		Wallpaper: wallpaper,
	}

	var all = core.ToolList{
		core.ToolHyprland,
		core.ToolWaybar,
		core.ToolRofi,
		core.ToolHyprpaper,
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
