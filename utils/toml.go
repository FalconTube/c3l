package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type ConfigToml struct {
	Flags
	ExpandPrompts
}

type Flags struct {
	Think      bool   `short:"t" help:"If true, uses thinking mode, if applicable in model. If false, adds '/no_think' to prompt" negatable:""`
	Print      bool   `short:"p" help:"If true, prints response to stdout (default: true)" negatable:""`
	Replace    bool   `short:"r" help:"If true, put Ollama output on clipboard" negatable:""`
	Model      string `short:"m" help:"Ollama model to use. Available models: https://ollama.com/library" default:"qwen3:0.6b"`
	Notify     bool   `short:"n" help:"If true, display tray notification when finished." negatable:"" default:"false"`
	Expand     bool   `short:"e" help:"Expand given prompt into long version, as defined in $HOME/.c3l.toml " negatable:"" `
	OllamaHost string `help:"IP Address for the Ollama server." env:"OLLAMA_HOST" default:"127.0.0.1:11434"`
}

type ExpandPrompts struct {
	Prompts map[string]string `toml:"prompts"`
}

func ExpandPromptFromToml(predefined string) (string, error) {
	prePrompts, err := GetPredefinedFromToml()
	if err != nil {
		return "", err
	}

	expanded := prePrompts.Prompts[predefined]

	if expanded == "" {
		keys := make([]string, 0, len(prePrompts.Prompts))
		for k := range prePrompts.Prompts {
			keys = append(keys, k)
		}
		return "", fmt.Errorf("Could not find predefined prompt \"%s\" in config file.\nAvailable prompts:\n%s", predefined, keys)
	}

	return expanded, nil

}
func GetConfigPath() (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFile := filepath.Join(home, ".c3l.toml")
	return configFile, nil
}

func ReadConfigAsBytes() ([]byte, error) {
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

func ReadConfigAsStruct() (ConfigToml, error) {
	config := ConfigToml{}
	raw, err := ReadConfigAsBytes()
	if err != nil {
		return ConfigToml{}, err
	}
	toml.Unmarshal(raw, config)

	return config, nil
}

func GetPredefinedFromToml() (ExpandPrompts, error) {
	config, err := ReadConfigAsBytes()
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
