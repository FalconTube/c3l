package cmd

import (
	"fmt"
	"os"

	"github.com/FalconTube/c3l/utils"
	"github.com/pelletier/go-toml/v2"
)

type SystemCmd struct {
	Add    AddSystemCmd    `cmd:"" help:"Add a system prompt to the config."`
	Remove RemoveSystemCmd `cmd:"" help:"Remove a system prompt from the config."`
	List   ListSystemCmd   `cmd:"" help:"List all system prompts."`
}

type AddSystemCmd struct {
	Short string `arg:"" help:"Shorthand notation of the prompt."`
	Long  string `arg:"" help:"Full prompt to be expanded from the shorthand notation."`
	Force bool   `short:"f" help:"If true, will replace prompt, even if it exists."`
}

type RemoveSystemCmd struct {
	Short string `arg:"" help:"Shorthand notation of prompt to be removed."`
}

type ListSystemCmd struct {
}

// --- Plain prompt command

func (c *SystemCmd) Run() error {
	return nil
}

// --- Add prompt

func (c *AddSystemCmd) Run() error {
	systems, err := utils.GetPredefinedSystemsFromToml()
	if err != nil {
		return err
	}

	// Check map nil
	if systems.Entries == nil {
		systems.Entries = make(map[string]string)
	}

	// Only need to check, if not forcing override
	if !c.Force {
		check := systems.Entries[c.Short]
		if check != "" {
			return fmt.Errorf("prompt '%s' already exists. Use '--force' to override it", c.Short)
		}
	}

	// Add new one
	systems.Entries[c.Short] = c.Long

	err = updateConfigWithSystems(systems)
	if err != nil {
		return err
	}
	return nil
}

// --- prompt list

func (c *ListSystemCmd) Run() error {
	utils.Logger.Info("Predefined Prompts:\n")
	prompts, err := utils.GetPredefinedSystemsFromToml()
	if err != nil {
		return err
	}
	b, _ := toml.Marshal(prompts)
	fmt.Println(string(b))
	return nil
}

// --- prompt remove

func (c *RemoveSystemCmd) Run() error {
	prompts, err := utils.GetPredefinedSystemsFromToml()
	if err != nil {
		return err
	}

	// Check if exists, else return
	check := prompts.Entries[c.Short]
	if check == "" {
		utils.Logger.Info("Prompt does not exist in config. Nothing to do...")
		os.Exit(0)
	}

	for k := range prompts.Entries {
		if k == c.Short {
			delete(prompts.Entries, k)
		}
	}
	err = updateConfigWithSystems(prompts)
	if err != nil {
		return err
	}
	return nil
}

func updateConfigWithSystems(systems utils.Systems) error {
	currentConfig, err := utils.ReadConfigAsStruct()
	if err != nil {
		return err
	}
	currentConfig.Systems = systems
	newConfig, err := toml.Marshal(currentConfig)
	if err != nil {
		return err
	}

	configPath, err := utils.GetConfigPath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, newConfig, 0644)
	if err != nil {
		return err
	}
	return nil
}
