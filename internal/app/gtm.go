package app

import (
	"flag"
	"fmt"
	"github.com/tpm1qq/gtm/internal/pkg/tools/hyprland"
	"github.com/tpm1qq/gtm/internal/pkg/tools/rofi"
	"github.com/tpm1qq/gtm/internal/pkg/tools/waybar"
	"os"
	"strings"
)

func RunGTM() {
	var tool string
	var color string
	flag.StringVar(&tool, "tool", "", "which tool's config the user wants to change")
	flag.StringVar(&tool, "t", "", "which tool's config the user wants to change")
	flag.StringVar(&color, "color", "", "change the color of the given tool(s)")
	flag.Parse()

	tool = strings.ToLower(tool)
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
	}
}
