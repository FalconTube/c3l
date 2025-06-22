package utils

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type ConfigToml struct {
	Flags
	Prompts
	Systems
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

type Prompts struct {
	Entries map[string]string `toml:"prompts"`
}

type Systems struct {
	Entries map[string]string `toml:"systems"`
}

func ExpandPromptFromToml(predefined string) (string, error) {
	if predefined == "" {
		return "", nil
	}

	prePrompts, err := GetPredefinedPromptsFromToml()
	if err != nil {
		return "", err
	}

	expanded := prePrompts.Entries[predefined]

	// If no expansion found, just return incoming prompt
	if expanded == "" {
		keys := make([]string, 0, len(prePrompts.Entries))
		for k := range prePrompts.Entries {
			keys = append(keys, k)
		}
		Logger.Warnf("could not find predefined prompt \"%s\" in config file.\nAvailable prompts:\n%s", predefined, keys)
		return predefined, nil
	}

	return expanded, nil

}

func ExpandSystemFromToml(predefined string) (string, error) {
	if predefined == "" {
		return "", nil
	}

	preSystems, err := GetPredefinedSystemsFromToml()
	if err != nil {
		return "", err
	}

	expanded := preSystems.Entries[predefined]

	if expanded == "" {
		keys := make([]string, 0, len(preSystems.Entries))
		for k := range preSystems.Entries {
			keys = append(keys, k)
		}
		Logger.Warnf("could not find predefined system prompt \"%s\" in config file.\nAvailable system prompts:\n%s", predefined, keys)
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

func GetPredefinedPromptsFromToml() (Prompts, error) {
	config, err := ReadConfigAsBytes()
	if err != nil {
		return Prompts{}, err
	}

	var customPrompts Prompts
	err = toml.Unmarshal(config, &customPrompts)
	if err != nil {
		return Prompts{}, err
	}
	return customPrompts, nil

}

func GetPredefinedSystemsFromToml() (Systems, error) {
	config, err := ReadConfigAsBytes()
	if err != nil {
		return Systems{}, err
	}

	var customSystems Systems
	err = toml.Unmarshal(config, &customSystems)
	if err != nil {
		return Systems{}, err
	}
	return customSystems, nil

}
