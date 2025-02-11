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

func (c *config) AddShortcut(name string, s Shortcut) error {

}

func (c *config) updateFile() error {
	getConfigFilePath
}

func (s *Shortcut) HasParams() bool {
	return strings.Contains(s.Template, "%s")
}

func getConfigFilePath() (string, error) {
	configPath, ok := os.LookupEnv("DLAUNCHER_CONFIG_PATH")
	if configPath != "" && ok {
		return configPath, nil
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath = fmt.Sprintf("%s/.config/dlauncher/config.yaml", homedir)
	return configPath, nil
}

func ParseConfig() config {
	configPath, err := getConfigFilePath()
	if err != nil {
		panic(err)
	}
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	config := config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}
