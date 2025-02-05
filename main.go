package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	"gopkg.in/yaml.v3"
)

type Shortcut struct {
	Kind               string   `yaml:"kind"`
	Shortcut           string   `yaml:"shortcut"`
	CommandTemplate    string   `yaml:"commandTemplate"`
	AllowedExecutables []string `yaml:"allowedExecutables"`
}

type Executable struct {
	Command []string `yaml:"command"`
}

type config struct {
	Executables map[string]Executable `yaml:"executables"`
	Shortcuts   map[string]Shortcut   `yaml:"shortcuts"`
}

func (c *config) getShortcut(shortcutName string) (Shortcut, error) {
	shortcut, ok := c.Shortcuts[shortcutName]
	if !ok {
		return Shortcut{}, fmt.Errorf("the shortcut does not exist: %s", shortcutName)
	}
	return shortcut, nil
}

func (c *config) getExecutable(executableName string) (Executable, error) {
	executable, ok := c.Executables[executableName]
	if !ok {
		return Executable{}, fmt.Errorf("the executable does not exist: %s", executableName)
	}
	return executable, nil
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

func runCommand(c config, executableName, shortcutName string, params []any) error {
	e, err := c.getExecutable(executableName)
	if err != nil {
		return err
	}
	s, err := c.getShortcut(shortcutName)
	if err != nil {
		return err
	}
	command := s.CommandTemplate
	if len(params) > 0 {
		command = fmt.Sprintf(command, params...)
	}
	cmd := exec.Command(e.Command[0], append(e.Command[1:], command)...)
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(currentUser.Gid)
	if err != nil {
		return err
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:         uint32(uid),
			Gid:         uint32(gid),
			NoSetGroups: true,
		},
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Process.Release()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	c := ParseConfig()
	params := []any{"bla"}
	shortcutName := "gs"
	executableName := "chrome"
	// executableName := "firefox"
	err := runCommand(c, executableName, shortcutName, params)
	if err != nil {
		panic(err)
	}
}
