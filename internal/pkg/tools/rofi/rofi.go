package rofi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SetColor(v string, p string) error {
	v, err := formatString(v)
	if err != nil {
		return err
	}
	err = editConfig(v, p)
	if err != nil {
		return err
	}
	return nil
}
func formatString(v string) (string, error) {
	switch {
	case v == "":
		return "", fmt.Errorf("no color value given")

	case strings.HasPrefix(v, "#") && len(v) == 7:
		v = strings.ToLower(v)
		v = v + "ff"
		return v, nil

	default:
		return "", fmt.Errorf("color not formatted correctly")
	}
}
func editConfig(v string, p string) error {
	var path = filepath.Join(p, "rofi", "gtm_rofi.rasi")
	file, err := os.ReadFile(path)
	data := string(file)
	if err != nil {
		return err
	}
	if strings.Contains(data, "gtmc") {
		i := strings.Index(data, "gtmc")
		if i == -1 {
			return fmt.Errorf("config doesn't contain color value")
		}
		curr := data[i+5 : i+14]
		data = strings.ReplaceAll(data, curr, v)
	}
	err = os.WriteFile(path, []byte(data), 0664)
	if err != nil {
		return err
	}
	return nil
}
