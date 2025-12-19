package core

import (
	"fmt"
	"strings"

	"github.com/tpm1qq/gtm/internal/pkg/config"
	"github.com/tpm1qq/gtm/internal/pkg/tools/hyprland"
	"github.com/tpm1qq/gtm/internal/pkg/tools/hyprpaper"
	"github.com/tpm1qq/gtm/internal/pkg/tools/rofi"
	"github.com/tpm1qq/gtm/internal/pkg/tools/waybar"
)

type ToolName string
type ToolList []ToolName

const (
	ToolHyprland  ToolName = "hyprland"
	ToolWaybar    ToolName = "waybar"
	ToolRofi      ToolName = "rofi"
	ToolHyprpaper ToolName = "hyprpaper"
)

type ToolSettings struct {
	Color     string
	Wallpaper string
}

func (t *ToolList) Set(v string) error {
	//split string here to support multiple flags in one command
	for s := range strings.SplitSeq(v, ",") {
		tool := strings.TrimSpace(s)
		if tool == "" {
			continue
		}
		*t = append(*t, ToolName(strings.ToLower(tool)))
	}
	return nil
}

func (t *ToolList) String() string {
	var res []string
	for _, e := range *t {
		res = append(res, string(e))
	}
	return strings.Join(res, ",")
}

func ApplyChanges(n ToolName, s ToolSettings, p config.Paths) error {
	switch n {
	case ToolHyprland:
		if s.Color != "" {
			if err := hyprland.Hyprland_SetColor(s.Color, p.ToolsPath); err != nil {
				return fmt.Errorf("error setting hyprland color: %w", err)
			}
		}
	case ToolWaybar:
		if s.Color != "" {
			if err := waybar.Waybar_SetColor(s.Color, p.ToolsPath); err != nil {
				return fmt.Errorf("error setting waybar color: %w", err)
			}
		}
		if err := waybar.Waybar_reload(); err != nil {
			return fmt.Errorf("error reloading waybar: %w", err)
		}

	case ToolRofi:
		if s.Color != "" {
			if err := rofi.Rofi_SetColor(s.Color, p.ToolsPath); err != nil {
				return fmt.Errorf("error setting rofi color: %w", err)
			}
		}
	case ToolHyprpaper:
		if err := hyprpaper.Hyprpaper_changeWallpaper(s.Wallpaper, p.ToolsPath, p.BackgroundPath); err != nil {
			return fmt.Errorf("error applying hyprpaper change: %w", err)
		}
		if err := hyprpaper.Hyprpaper_reload(); err != nil {
			return fmt.Errorf("error reloading hyprpaper: %w", err)
		}
	}
	return nil
}
