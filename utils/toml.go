package utils

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type ConfigToml struct {
	Flags
	ExpandPrompts
	ExpandSystems
}

type Flags struct {
	Think      bool   `short:"t" help:"If true, uses thinking mode, if applicable in model. If false, adds '/no_think' to prompt" negatable:""`
	Silent     bool   `short:"s" help:"If true, prints response to stdout (default: false)" negatable:""`
	Replace    bool   `short:"r" help:"If true, put Ollama output on clipboard" negatable:""`
	Model      string `short:"m" help:"Ollama model to use. Available models: https://ollama.com/library" default:"qwen3:0.6b"`
	Notify     bool   `short:"n" help:"If true, display tray notification when finished." negatable:"" default:"false"`
	Expand     bool   `short:"e" help:"Expand given prompt into long version, as defined in $HOME/.c3l.toml " negatable:"" `
	System     string `help:"(Optional) System prompt to pass to model. Can be expanded into long version." `
	OllamaHost string `help:"IP Address for the Ollama server." env:"OLLAMA_HOST" default:"127.0.0.1:11434"`
}

type ExpandPrompts struct {
	Entries map[string]string `toml:"prompts"`
}

type ExpandSystems struct {
	Entries map[string]string `toml:"systems"`
}

type PredefinedExpansions struct {
	ExpandPrompts ExpandPrompts
	ExpandSystems ExpandSystems
}

type ExpandType int

const (
	PromptType ExpandType = iota
	SystemType
)

func ExpandAnyFromToml(predefined string, expandType ExpandType) (string, error) {
	if predefined == "" {
		return "", nil
	}

	pre, err := GetPredefinedFromToml()
	if err != nil {
		return "", err
	}
	var entries map[string]string
	switch expandType {
	case PromptType:
		entries = pre.ExpandPrompts.Entries
	case SystemType:
		entries = pre.ExpandSystems.Entries
	}
	expanded := entries[predefined]
	// If no expansion found, just return incoming prompt
	if expanded == "" {
		keys := make([]string, 0, len(entries))
		for k := range entries {
			keys = append(keys, k)
		}
		Logger.Warnf("could not find predefined key \"%s\" in config file.\nAvailable values:\n%s", predefined, keys)
		return predefined, nil
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
	err = toml.Unmarshal(raw, &config)
	if err != nil {
		return ConfigToml{}, err
	}

	return config, nil
}

func GetPredefinedFromToml() (PredefinedExpansions, error) {
	config, err := ReadConfigAsBytes()
	if err != nil {
		return PredefinedExpansions{}, err
	}
	var prompts ExpandPrompts
	err = toml.Unmarshal(config, &prompts)
	if err != nil {
		return PredefinedExpansions{}, err
	}
	var systems ExpandSystems
	err = toml.Unmarshal(config, &systems)
	if err != nil {
		return PredefinedExpansions{}, err
	}
	pre := PredefinedExpansions{ExpandPrompts: prompts, ExpandSystems: systems}
	return pre, nil

}
