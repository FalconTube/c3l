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

type AddAnyCmd struct {
	Short string
	Long  string
	Force bool
}

type RemoveAnyCmd struct {
	Short string
}

type ListAnyCmd struct {
}

// --- Plain prompt command

func (c *PromptCmd) Run() error {
	return nil
}

// --- Add prompt

func (c *AddPromptCmd) Run() error {
	cmd := AddAnyCmd{Short: c.Short, Long: c.Long, Force: c.Force}
	err := cmd.addPromptCmd(utils.PromptType)
	return err
}

func (c *AddAnyCmd) addPromptCmd(e utils.ExpandType) error {

	pre, err := utils.GetPredefinedFromToml()
	if err != nil {
		return err
	}
	var entries map[string]string
	switch e {
	case utils.PromptType:
		entries = pre.ExpandPrompts.Entries
	case utils.SystemType:
		entries = pre.ExpandSystems.Entries
	}
	utils.Logger.Info("initial entries", entries)

	// Check map nil
	if entries == nil {
		entries = make(map[string]string)
	}
	// Only need to check, if not forcing override
	if !c.Force {
		check := entries[c.Short]
		if check != "" {
			return fmt.Errorf("%s '%s' already exists. Use '--force' to override it", e, c.Short)
		}
	}

	// Add new one
	entries[c.Short] = c.Long
	utils.Logger.Info("entries", entries)

	switch e {
	case utils.PromptType:

		new := utils.ExpandPrompts{Entries: entries}
		utils.Logger.Info("new: ", new)
		err = updateConfigWithPrompts(new)
		if err != nil {
			return err
		}
	case utils.SystemType:
		new := utils.ExpandSystems{Entries: entries}
		err = updateConfigWithSystems(new)
		if err != nil {
			return err
		}
	}
	return nil

}

// --- prompt list

func (c *ListPromptCmd) Run() error {
	cmd := ListAnyCmd{}
	err := cmd.listCmd(utils.PromptType)
	return err
}

func (c *ListAnyCmd) listCmd(e utils.ExpandType) error {
	utils.Logger.Infof("Predefined %s:\n", e)
	pre, err := utils.GetPredefinedFromToml()
	if err != nil {
		return err
	}
	var prompts interface{}
	switch e {
	case utils.PromptType:
		prompts = pre.ExpandPrompts
	case utils.SystemType:
		prompts = pre.ExpandSystems
	}
	b, _ := toml.Marshal(prompts)
	fmt.Println(string(b))
	return nil
}

// --- prompt remove

func (c *RemovePromptCmd) Run() error {
	cmd := RemoveAnyCmd{}
	err := cmd.removeCmd(utils.PromptType)
	return err
}

func (c *RemoveAnyCmd) removeCmd(e utils.ExpandType) error {

	pre, err := utils.GetPredefinedFromToml()
	if err != nil {
		return err
	}

	// Check if exists, else return
	var check string
	switch e {
	case utils.PromptType:
		check = pre.ExpandPrompts.Entries[c.Short]
	case utils.SystemType:
		check = pre.ExpandSystems.Entries[c.Short]
	}
	if check == "" {
		utils.Logger.Info("Prompt does not exist in config. Nothing to do...")
		os.Exit(0)
	}

	switch e {
	case utils.PromptType:
		ex := pre.ExpandPrompts
		for k := range ex.Entries {
			if k == c.Short {
				delete(ex.Entries, k)
			}
		}
		err = updateConfigWithPrompts(ex)
		if err != nil {
			return err
		}
	case utils.SystemType:
		ex := pre.ExpandSystems
		for k := range ex.Entries {
			if k == c.Short {
				delete(ex.Entries, k)
			}
		}
		err = updateConfigWithSystems(ex)
		if err != nil {
			return err
		}
	}

	return nil

}

func updateConfigWithPrompts(prompts utils.ExpandPrompts) error {
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
