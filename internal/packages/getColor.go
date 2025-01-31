package packages

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
)

type Colors struct {
	Colors []string `json:"colors"`
}

// GetRandomColor fn selects a color randomly from colors.json
func GetRandomColor() (string, error) {
	filePath := filepath.Join("internal", "data", "colors.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	var colors Colors
	err = json.Unmarshal(data, &colors)
	if err != nil {
		return "", err
	}

	randomIndex := rand.Intn(len(colors.Colors))

	return colors.Colors[randomIndex], nil
}
