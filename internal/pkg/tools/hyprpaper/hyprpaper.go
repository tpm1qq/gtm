package hyprpaper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ChangeWallpaper(v string, p string, b string) error {
	v, err := formatString(v)
	if err != nil {
		return err
	}
	err = editConfig(v, p, b)
	if err != nil {
		return err
	}
	return nil
}
func formatString(v string) (string, error) {
	if v == "" {
		return "", fmt.Errorf("no wallpaper given")
	} else {
		switch ext := strings.ToLower(filepath.Ext(v)); ext {
		case ".png", ".jpg", ".jpeg", ".webp":
			return v, nil

		default:
			return "", fmt.Errorf("invalid file format")
		}
	}
}
func editConfig(v string, p string, b string) error {
	var path = filepath.Join(p, "hypr", "gtm_hyprpaper.conf")
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %w", err)
	}
	data := string(file)
	if strings.Contains(data, "gtmc") {
		i := strings.Index(data, "gtmc")
		if i == -1 {
			return fmt.Errorf("config doesn't contain gtmc start marker")
		}
		start := i + len("gtmc")
		i = strings.Index(data, "!?")
		if i == -1 {
			return fmt.Errorf("config doesn't contain gtmc end marker")
		}
		end := i
		curr := data[start+1 : end]
		data = strings.ReplaceAll(data, curr, v)
	}
	if strings.Contains(data, "gtmb") {
		i := strings.Index(data, "gtmb")
		if i == -1 {
			return fmt.Errorf("config doesn't contain gtmb start marker")
		}
		start := i + len("gtmb")
		i = strings.Index(data, "?!")
		if i == -1 {
			return fmt.Errorf("config doesn't contain gtmb end marker")
		}
		end := i
		curr := data[start+1 : end]
		data = strings.ReplaceAll(data, curr, b)
	}
	err = os.WriteFile(path, []byte(data), 0664)
	if err != nil {
		return err
	}
	return nil
}
func Reload() error {
	cmd := exec.Command("pkill", "hyprpaper")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pkill error %w", err)
	}
	cmd = exec.Command("hyprpaper")
	return cmd.Start()
}
