package app

import (
	"flag"
	"fmt"
	"github.com/tpm1qq/gtm/internal/pkg/tools/hyprland"
	"github.com/tpm1qq/gtm/internal/pkg/tools/hyprpaper"
	"github.com/tpm1qq/gtm/internal/pkg/tools/rofi"
	"github.com/tpm1qq/gtm/internal/pkg/tools/waybar"
	"os"
	"strings"
	"sync"
)

func RunGTM() {
	var tool string
	var global bool
	var color string
	var wallpaper string
	flag.BoolVar(&global, "global", false, "global flag; apply config changes to all tools at the same time")
	flag.BoolVar(&global, "g", false, "global flag; apply config changes to all tools at the same time")
	flag.StringVar(&tool, "tool", "", "which tool's config the user wants to change")
	flag.StringVar(&tool, "t", "", "which tool's config the user wants to change")
	flag.StringVar(&color, "color", "", "change the color of the given tool(s)")
	flag.StringVar(&wallpaper, "wallpaper", "", "change the current wallpaper")
	flag.StringVar(&wallpaper, "w", "", "change the current wallpaper")

	flag.Parse()

	tool = strings.ToLower(tool)
	switch {
	case global && tool == "":
		var wg = sync.WaitGroup{}
		var errors = make(chan error, 4)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := hyprland.Hyprland_SetColor(color); err != nil {
				errors <- err
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := waybar.Waybar_SetColor(color); err != nil {
				errors <- err
			}
			if err := waybar.Waybar_reload(); err != nil {
				errors <- err
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := rofi.Rofi_SetColor(color); err != nil {
				errors <- err
			}
		}()
		wg.Wait()
		close(errors)
		for err := range errors {
			if err != nil {
				fmt.Fprintln(os.Stderr, "error setting color:", err)
			}
		}
	case !global && tool != "":
		switch tool {
		case "hyprland":
			if err := hyprland.Hyprland_SetColor(color); err != nil {
				fmt.Fprintln(os.Stderr, "error setting color:", err)
			}
		case "waybar":
			if err := waybar.Waybar_SetColor(color); err != nil {
				fmt.Fprintln(os.Stderr, "error setting color:", err)
			}
			if err := waybar.Waybar_reload(); err != nil {
				fmt.Fprintln(os.Stderr, "error reloading waybar", err)
			}

		case "rofi":
			if err := rofi.Rofi_SetColor(color); err != nil {
				fmt.Fprintln(os.Stderr, "error setting color:", err)
			}
		case "hyprpaper":
			if err := hyprpaper.Hyprpaper_changeWallpaper(wallpaper); err != nil {
				fmt.Fprintln(os.Stderr, "error changing wallpaper:", err)
			}
			if err := hyprpaper.Hyprpaper_reload(); err != nil {
				fmt.Fprintln(os.Stderr, "error reloading hyprpaper:", err)
			}
		}
	case !global && tool == "":
		fmt.Println("neither global nor tool flag set!")
	}
}
