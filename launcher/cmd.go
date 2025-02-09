package launcher

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var CONFIG config

func getFlagValues(cmd *cobra.Command, display CobraDisplay) (Shortcut, Executable, []string, error) {
	s := Shortcut{}
	e := Executable{}
	params := []string{}
	executableName, err := cmd.Flags().GetString("executable-name")
	if err != nil {
		return s, e, params, err
	}
	if executableName == "" {
		return s, e, params, fmt.Errorf("The executable-name must be provided.")
	}
	shortcutName, err := cmd.Flags().GetString("shortcut-name")
	if err != nil {
		return s, e, params, err
	}
	if shortcutName == "" {
		shortcutName = display.Prompt(fmt.Sprintf("[%s] Shortcut name", executableName))
	}
	params, err = cmd.Flags().GetStringArray("params")
	if err != nil {
		return s, e, params, err
	}
	s, err = CONFIG.GetShortcut(shortcutName)
	if err != nil {
		return s, e, params, err
	}
	e, err = CONFIG.GetExecutable(executableName)
	if err != nil {
		return s, e, params, err
	}
	if s.HasParams() && len(params) == 0 {
		promptResponse := display.Prompt("Params for template, comma separated")
		params = strings.Split(promptResponse, ",")
	}
	return s, e, params, nil
}

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a shortcut",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		useGUI, _ := cmd.Flags().GetBool("use-gui")
		display := NewDisplay(useGUI, args)
		shortcut, executable, params, err := getFlagValues(cmd, display)
		if err != nil {
			display.Error(err.Error())
			panic(err)
		}
		err = RunCommand(shortcut, executable, params)
		if err != nil {
			display.Error(err.Error())
			panic(err)
		}
	},
}

func init() {
	rootCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")
	rootCmd.Flags().StringP("executable-name", "e", "", "The program that should execute your command template.")
	rootCmd.Flags().StringP("shortcut-name", "s", "", "The name of the shortcut.")
	rootCmd.Flags().StringArrayP("params", "p", []string{}, "(optional) The params for the command.")
	CONFIG = ParseConfig()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
