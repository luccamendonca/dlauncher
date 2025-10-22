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
	PromptMultiline(msg string) string
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
func (cli DisplayCLI) PromptMultiline(msg string) string {
	fmt.Fprintf(os.Stderr, "%s (enter empty line when done):\n", msg)
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
func (cli DisplayCLI) Error(msg string) {
	cli.Info(fmt.Sprintf("error: %s", msg))
}
func (cli DisplayCLI) Info(msg string) {
	fmt.Println(msg)
}
func (cli DisplayCLI) Debug(params any) {
	repr.Println(params)
}
func (cli DisplayCLI) Panic(err error) {
	panic(err)
}

// DisplayGUI
func (gui DisplayGUI) Prompt(msg string) string {
	resp, err := zenity.Entry(
		msg,
		zenity.CancelLabel(""),
		zenity.OKLabel(""),
	)
	if err != nil {
		zenity.Error(err.Error())
		os.Exit(1)
	}
	return resp
}
func (gui DisplayGUI) PromptMultiline(msg string) string {
	resp, err := zenity.Entry(
		msg+"\n(separate links with newlines)",
		zenity.CancelLabel(""),
		zenity.OKLabel(""),
	)
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
