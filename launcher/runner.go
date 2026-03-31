//go:build linux || darwin

package launcher

import (
	"fmt"
	"os/exec"
	"syscall"
)

func stringToAny(s []string) []any {
	sAny := []any{}
	for _, v := range s {
		sAny = append(sAny, v)
	}
	return sAny
}

func prepareCommand(name string, arg ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	return cmd, nil
}

func RunCommand(s Shortcut, e Executable, params []string) error {
	command := s.Template
	if len(params) > 0 {
		command = fmt.Sprintf(command, stringToAny(params)...)
	}
	cmd, err := prepareCommand(e.Command[0], append(e.Command[1:], command)...)
	if err != nil {
		return err
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

// RunMultipleCommands opens multiple links in separate browser tabs
func RunMultipleCommands(links []string, e Executable) error {
	for _, link := range links {
		cmd, err := prepareCommand(e.Command[0], append(e.Command[1:], link)...)
		if err != nil {
			return fmt.Errorf("failed to prepare command for link %s: %v", link, err)
		}
		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to start command for link %s: %v", link, err)
		}
		err = cmd.Process.Release()
		if err != nil {
			return fmt.Errorf("failed to release process for link %s: %v", link, err)
		}
	}
	return nil
}
