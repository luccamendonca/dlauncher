package launcher

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Shortcut struct {
	Template             string   `yaml:"template"`
	SupportedExecutables []string `yaml:"supportedExecutables,omitempty"`
}

type Executable struct {
	Command []string `yaml:"command"`
}

type Config struct {
	Executables map[string]Executable `yaml:"executables"`
	Shortcuts   map[string]Shortcut   `yaml:"shortcuts"`
}

func (c *Config) GetShortcut(shortcutName string) (Shortcut, error) {
	shortcut, ok := c.Shortcuts[shortcutName]
	if !ok {
		return Shortcut{}, fmt.Errorf("the shortcut does not exist: %s", shortcutName)
	}
	return shortcut, nil
}

func (c *Config) GetExecutable(executableName string) (Executable, error) {
	executable, ok := c.Executables[executableName]
	if !ok {
		return Executable{}, fmt.Errorf("the executable does not exist: %s", executableName)
	}
	return executable, nil
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

func ParseConfig() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("config file not found at %s.  Run 'dlauncher init' to create a default config file.", configPath)
	}

	config := Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func CreateDefaultConfig() error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configDir := strings.TrimSuffix(configPath, "/config.yaml")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	defaultConfig := Config{
		Executables: map[string]Executable{
			"chrome": {
				Command: []string{"/usr/bin/google-chrome-stable", "--new-tab"},
			},
			"firefox": {
				Command: []string{"/usr/bin/firefox", "--new-tab", "--url"},
			},
		},
		Shortcuts: map[string]Shortcut{
			"any": {
				Template: "%s",
			},
			"blank": {
				Template: "about:blank",
			},
			"google": {
				Template: "https://www.google.com/search?q=%s",
			},
		},
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
