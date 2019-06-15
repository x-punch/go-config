package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Configor represents configuration loader
type Configor struct {
	Options *Options
}

// NewConfigor create configor instance
func NewConfigor(opts ...Option) *Configor {
	if files := os.Getenv("CONFIG_FILES"); files != "" {
		opts = append(opts, Files(strings.Split(files, ",")))
	}
	if prefix := os.Getenv("CONFIG_PREFIX"); prefix != "" {
		opts = append(opts, Prefix(prefix))
	}
	if showLog := os.Getenv("CONFIG_SHOW_LOG"); showLog != "" {
		opts = append(opts, ShowLog(true))
	}
	return &Configor{Options: newOptions(opts...)}
}

// Load will unmarshal configurations to struct from files that you provide
func (c *Configor) Load(config interface{}) error {
	for _, file := range c.getConfigFiles(c.Options.Files...) {
		if c.Options.ShowLog {
			fmt.Printf("[Config]Decode file: %s\n", file)
		}
		if _, err := toml.DecodeFile(file, config); err != nil {
			return fmt.Errorf("Failed to decode %s: %+v", file, err)
		}
	}
	if err := c.ApplyEnvOverrides(c.Options.Prefix, config); err != nil {
		return fmt.Errorf("Failed to apply env args: %+v", err)
	}
	if c.Options.ShowLog {
		fmt.Printf("[Config]Loaded:  %+v\n", config)
	}
	return nil
}

func (c *Configor) getConfigFiles(files ...string) []string {
	var validFiles []string
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			if c.Options.ShowLog {
				fmt.Printf("[Config]File not found: %s\n", file)
			}
		} else {
			if c.Options.ShowLog {
				fmt.Printf("[Config]Found file: %s\n", file)
			}
			validFiles = append(validFiles, file)
		}
	}
	return validFiles
}
