package launcher

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Shortcut struct {
	Template             string   `yaml:"template"`
	SupportedExecutables []string `yaml:"supportedExecutables"`
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

func (s *Shortcut) HasParams() bool {
	return strings.Contains(s.Template, "")
}

func ParseConfig() config {
	config := config{}
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}
