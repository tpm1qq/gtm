package hyprpaper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var home, _ = os.UserHomeDir()
var path = filepath.Join(home, ".config", "hypr", "gtm_hyprpaper.conf")
var path_wp = filepath.Join(home, "Backgrounds")

// todo: when config functionality is introduced, add variable for saving default wallpaper path
// currently hardcoding the path is the only way, perhaps use path given in hyprpaper.conf and assume its the default for all wallpapers?
// setup functionality can perhaps give this choice

func Hyprpaper_changeWallpaper(v string) error {
	v, err := formatString(v)
	if err != nil {
		return err
	}
	err = editConfig(v)
	if err != nil {
		return err
	}
	return nil
}
func formatString(v string) (string, error) {
	switch {
	case v == "":
		return "", fmt.Errorf("no wallpaper given")

	case v != "":
		switch ext := strings.ToLower(filepath.Ext(v)); ext {
		case ".png", ".jpg", ".jpeg", ".webp":
			return v, nil

		default:
			return "", fmt.Errorf("invalid file format")
		}

	default:
		return "", fmt.Errorf("wallpaper string not formatted correctly")
	}
}
func editConfig(v string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	data := string(file)
	if strings.Contains(data, "gtmc") {
		i := strings.Index(data, "gtmc")
		if i == -1 {
			return fmt.Errorf("config doesn't contain start marker")
		}
		start := i + len("gtmc")
		i = strings.Index(data, "!?")
		if i == -1 {
			return fmt.Errorf("config doesn't contain end marker")
		}
		end := i
		curr := data[start+1 : end]
		data = strings.ReplaceAll(data, curr, v)
	}
	err = os.WriteFile(path, []byte(data), 0664)
	if err != nil {
		return err
	}
	return nil
}
func Hyprpaper_reload() error {
	cmd := exec.Command("pkill", "hyprpaper")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pkill error %w", err)
	}
	cmd = exec.Command("hyprpaper")
	return cmd.Start()
}
