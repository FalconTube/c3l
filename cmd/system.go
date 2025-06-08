package cmd

import (
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
	cmd := AddAnyCmd{Short: c.Short, Long: c.Long, Force: c.Force}
	err := cmd.addPromptCmd(utils.SystemType)
	return err

}

// --- prompt list

func (c *ListSystemCmd) Run() error {
	cmd := ListAnyCmd{}
	err := cmd.listCmd(utils.SystemType)
	return err
}

// --- prompt remove

func (c *RemoveSystemCmd) Run() error {
	cmd := RemoveAnyCmd{}
	err := cmd.removeCmd(utils.SystemType)
	return err
}
func updateConfigWithSystems(systems utils.ExpandSystems) error {
	currentConfig, err := utils.ReadConfigAsStruct()
	if err != nil {
		return err
	}
	currentConfig.ExpandSystems = systems
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
