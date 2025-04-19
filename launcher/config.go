package launcher

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Shortcut struct {
	Template string `yaml:"template"`
}

type Executable struct {
	Command []string `yaml:"command"`
}

type config struct {
	Executables map[string]Executable `yaml:"executables"`
	Shortcuts   map[string]Shortcut   `yaml:"shortcuts"`
}

func (c *config) GetShortcut(shortcutName string) (Shortcut, error) {
	shortcut, ok := c.Shortcuts[shortcutName]
	if !ok {
		return Shortcut{}, fmt.Errorf("the shortcut does not exist: %s", shortcutName)
	}
	return shortcut, nil
}

func (c *config) GetExecutable(executableName string) (Executable, error) {
	executable, ok := c.Executables[executableName]
	if !ok {
		return Executable{}, fmt.Errorf("the executable does not exist: %s", executableName)
	}
	return executable, nil
}

func (c *config) AddShortcut(name string, s Shortcut) error {
	val, ok := c.Shortcuts[name]
	if ok {
		return fmt.Errorf("shortcut named '%s' already exists. The template is: '%s'", name, val.Template)
	}
	c.Shortcuts[name] = s
	return c.updateFile()
}

func (c *config) updateFile() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	configYAML, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	os.WriteFile(configFilePath, configYAML, os.ModePerm)
	return nil
}

func (s *Shortcut) HasParams() bool {
	return strings.Contains(s.Template, "%s")
}

func getConfigFilePath() (string, error) {
	configPath, ok := os.LookupEnv("DLAUNCHER_CONFIG_PATH")
	if ok && configPath != "" {
		return configPath, nil
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	configPath = fmt.Sprintf("%s/.config/dlauncher/config.yaml", homedir)
	return configPath, nil
}

func ParseConfig() (config, error) {
	cfg := config{}
	configPath, err := getConfigFilePath()
	if err != nil {
		return cfg, fmt.Errorf("failed to determine config file path: %w", err)
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return cfg, fmt.Errorf("configuration file not found at '%s'", configPath)
		}
		return cfg, fmt.Errorf("failed to read config file '%s': %w", configPath, err)
	}

	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse config file '%s': %w", configPath, err)
	}

	return cfg, nil
}
