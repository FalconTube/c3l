package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/FalconTube/c3l/utils"
	"github.com/pelletier/go-toml/v2"
)

type ConfigCmd struct {
	Init InitCmd `cmd:"" help:"Create default config at $HOME/.c3l.toml"`
	List ListCmd `cmd:"" help:"List contents of config"`
}

type InitCmd struct {
	Force bool `short:"f" help:"If true, will replace current config with default config, if file exists."`
}

type ListCmd struct {
}

// --- Plain config command

func (c *ConfigCmd) Run() error {
	return nil
}

// --- config init command

func (c *InitCmd) Run() error {
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
	defaultFlags := Flags{Model: "qwen3:0.6b", OllamaHost: "127.0.0.1:11434"}
	f, err := toml.Marshal(defaultFlags)
	if err != nil {
		return err
	}
	defaultPrompts := utils.ExpandPrompts{Prompts: map[string]any{"example": "Using the '--expand' flag with the word 'example' as the prompt will expand into this string."}}
	p, err := toml.Marshal(defaultPrompts)
	if err != nil {
		return err
	}
	content := fmt.Sprintf("%s\n%s", f, p)
	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

// --- config list command

func (c *ListCmd) Run() error {
	configPath, err := utils.GetConfigPath()
	if err != nil {
		return err
	}
	config, err := utils.ReadConfigToml()
	if err != nil {
		return err
	}
	utils.Logger.Infof("Content of default file %s\n", configPath)

	fmt.Printf(string(config))
	return nil
}
