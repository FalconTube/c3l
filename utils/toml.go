package utils

import (
	"github.com/pelletier/go-toml"
	"os"
	"path/filepath"
)

type ExpandPrompts struct {
	Prompts map[string]any `toml:"prompts"`
}

func ExpandPromptFromToml(predefined string) string {
	prePrompts, err := getPredefinedFromToml()
	if err != nil {
		Logger.Fatal(err)
	}

	expanded := prePrompts.Prompts[predefined]

	if expanded == nil {
		keys := make([]string, 0, len(prePrompts.Prompts))
		for k := range prePrompts.Prompts {
			keys = append(keys, k)
		}
		Logger.Fatalf("Could not find predefined prompt \"%s\" in config file.\nAvailable prompts:\n%s", predefined, keys)
	}

	expandedString := expanded.(string)
	return expandedString

}
func GetConfigPath() (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFile := filepath.Join(home, ".c3l.toml")
	return configFile, nil
}

func ReadConfigToml() ([]byte, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getPredefinedFromToml() (ExpandPrompts, error) {
	config, err := ReadConfigToml()
	if err != nil {
		return ExpandPrompts{}, err
	}

	var customPrompts ExpandPrompts
	err = toml.Unmarshal(config, &customPrompts)
	if err != nil {
		return ExpandPrompts{}, err
	}
	return customPrompts, nil

}
