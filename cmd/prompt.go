package cmd

import (
	"fmt"
	"os"

	"github.com/FalconTube/c3l/utils"
	"github.com/pelletier/go-toml/v2"
)

type PromptCmd struct {
	Add    AddPromptCmd    `cmd:"" help:"Add a shorthand notation prompt to the config."`
	Remove RemovePromptCmd `cmd:"" help:"Remove a shorthand notation prompt from the config."`
	List   ListPromptCmd   `cmd:"" help:"List all shorthand notation prompts available in the config."`
}

type AddPromptCmd struct {
	Short string `arg:"" help:"Shorthand notation of the prompt."`
	Long  string `arg:"" help:"Full prompt to be expanded from the shorthand notation."`
	Force bool   `short:"f" help:"If true, will replace prompt, even if it exists."`
}

type RemovePromptCmd struct {
	Short string `arg:"" help:"Shorthand notation of prompt to be removed."`
}

type ListPromptCmd struct {
}

// --- Plain prompt command

func (c *PromptCmd) Run() error {
	return nil
}

// --- Add prompt

func (c *AddPromptCmd) Run() error {
	prompts, err := utils.GetPredefinedFromToml()
	if err != nil {
		return err
	}
	// Only need to check, if not forcing override
	if !c.Force {
		check := prompts.Prompts[c.Short]
		if check != "" {
			return fmt.Errorf("prompt '%s' already exists. Use '--force' to override it", c.Short)
		}
	}

	// Add new one
	prompts.Prompts[c.Short] = c.Long
	err = updateConfig(prompts)
	if err != nil {
		return err
	}
	return nil
}

// --- prompt list

func (c *ListPromptCmd) Run() error {
	utils.Logger.Info("Predefined Prompts:\n")
	prompts, err := utils.GetPredefinedFromToml()
	if err != nil {
		return err
	}
	b, _ := toml.Marshal(prompts)
	fmt.Println(string(b))
	return nil
}

// --- prompt remove

func (c *RemovePromptCmd) Run() error {
	prompts, err := utils.GetPredefinedFromToml()
	if err != nil {
		return err
	}

	// Check if exists, else return
	check := prompts.Prompts[c.Short]
	if check == "" {
		utils.Logger.Info("Prompt does not exist in config. Nothing to do...")
		os.Exit(0)
	}

	for k := range prompts.Prompts {
		if k == c.Short {
			delete(prompts.Prompts, k)
		}
	}
	err = updateConfig(prompts)
	if err != nil {
		return err
	}
	return nil
}

func updateConfig(prompts utils.ExpandPrompts) error {
	currentConfig, err := utils.ReadConfigAsStruct()
	if err != nil {
		return err
	}
	currentConfig.ExpandPrompts = prompts
	newConfig, _ := toml.Marshal(currentConfig)

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
