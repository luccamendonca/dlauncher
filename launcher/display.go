package launcher

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/ncruces/zenity"
)

type CobraDisplay interface {
	Prompt(msg string) string
	Error(msg string)
	Info(msg string)
	Debug(params any)
	Panic(err error)
}

type DisplayCLI struct {
	args []string
}
type DisplayGUI struct {
	args []string
}

func NewDisplay(useGUI bool, args []string) CobraDisplay {
	if useGUI {
		return DisplayGUI{args}
	}
	return DisplayCLI{args}
}

// DisplayCLI
func (cli DisplayCLI) Prompt(msg string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	i := 0
	for {
		if i == 0 {
			msg = msg + "\n"
		} else {
			msg = msg + " "
		}
		fmt.Fprint(os.Stderr, msg)
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
func (cli DisplayCLI) Error(msg string) {
	cli.Info(msg)
}
func (cli DisplayCLI) Info(msg string) {
	fmt.Println(msg)
}
func (cli DisplayCLI) Debug(params any) {
	repr.Println(params)
}
func (cli DisplayCLI) Panic(err error) {
	cli.Error(err.Error())
	panic(err)
}

// DisplayGUI
func (gui DisplayGUI) Prompt(msg string) string {
	resp, err := zenity.Entry(msg)
	if err != nil {
		zenity.Error(err.Error())
		os.Exit(1)
	}
	return resp
}
func (gui DisplayGUI) Error(msg string) {
	zenity.Error(msg)
}
func (gui DisplayGUI) Info(msg string) {
	zenity.Info(msg)
}
func (gui DisplayGUI) Debug(params any) {
	zenity.Info(repr.String(params))
}
func (cli DisplayGUI) Panic(err error) {
	panic(err)
}
