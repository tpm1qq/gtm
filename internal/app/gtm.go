package app

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tpm1qq/gtm/internal/pkg/tools"
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
		if err := tools.Hyprland_SetColor(color); err != nil {
			fmt.Fprintln(os.Stderr, "error setting color:", err)
		}
		/* case "waybar":
			if err := tools.Waybar_SetColor(color); err != nil {
				fmt.Fprintln(os.Stderr, "error setting color:", err)
		case "rofi":
			if err := tools.Rofi_SetColor(color); err != nil {
				fmt.Fprintln(os.Stderr, "error setting color:", err)
				os.Exit(1)*/
	}
}
