package launcher

import (
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func stringToAny(s []string) []any {
	sAny := []any{}
	for _, v := range s {
		sAny = append(sAny, v)
	}
	return sAny
}

func RunCommand(s Shortcut, e Executable, params []string) error {
	command := s.CommandTemplate
	if len(params) > 0 {
		command = fmt.Sprintf(command, stringToAny(params)...)
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
