package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/FalconTube/c3l/utils"
	"github.com/pelletier/go-toml/v2"
)

type ConfigCmd struct {
	Init InitConfigCmd `cmd:"" help:"Create default config at $HOME/.c3l.toml"`
	List ListConfigCmd `cmd:"" help:"List contents of config"`
}

type InitConfigCmd struct {
	Force bool `short:"f" help:"If true, will replace current config with default config, if file exists."`
}

type ListConfigCmd struct {
}

// --- Plain config command

func (c *ConfigCmd) Run() error {
	return nil
}

// --- config init command

func (c *InitConfigCmd) Run() error {
	configPath, err := utils.GetConfigPath()
	if err != nil {
		utils.Logger.Fatal(err)
	}
	// If we are not forcing a write, check if file already exists
	if c.Force == false {
		_, err = os.Stat(configPath)
		if os.IsNotExist(err) {
			utils.Logger.Infof("Config file does not exist yet, will create default one now...")
		} else {
			utils.Logger.Fatalf("Config already exists at %s", configPath)
		}
	}

	err = writeConfigToFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func writeConfigToFile(filename string) error {
	defaultFlags := utils.Flags{Model: "qwen3:0.6b", OllamaHost: "127.0.0.1:11434"}
	defaultPrompts := utils.ExpandPrompts{Prompts: map[string]string{"example": "Using the '--expand' flag with the word 'example' as the prompt will expand into this string."}}
	defaultConfig := utils.ConfigToml{Flags: defaultFlags, ExpandPrompts: defaultPrompts}
	content, err := toml.Marshal(defaultConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

// --- config list command

func (c *ListConfigCmd) Run() error {
	configPath, err := utils.GetConfigPath()
	if err != nil {
		return err
	}
	config, err := utils.ReadConfigAsBytes()
	if err != nil {
		return err
	}
	utils.Logger.Infof("Content of default file %s\n", configPath)

	fmt.Printf(string(config))
	return nil
}
