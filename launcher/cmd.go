package launcher

import (
	"github.com/spf13/cobra"
)

var CONFIG config

func getFlagValues(cmd *cobra.Command, args []string) (Shortcut, Executable, []string, error) {
	s := Shortcut{}
	e := Executable{}
	params := []string{}
	useGUI, _ := cmd.Flags().GetBool("use-gui")
	display := NewDisplay(useGUI, args)
	shortcutName, err := cmd.Flags().GetString("shortcut-name")
	if err != nil {
		return s, e, params, err
	}
	if shortcutName == "" {
		shortcutName = display.Prompt("Shortcut name")
	}
	params, err = cmd.Flags().GetStringArray("params")
	if err != nil {
		return s, e, params, err
	}
	executableName, err := cmd.Flags().GetString("executable-name")
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
		params = []string{display.Prompt("Param for template")}
	}
	return s, e, params, nil
}

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a shortcut",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		shortcut, executable, params, err := getFlagValues(cmd, args)
		if err != nil {
			panic(err)
		}
		err = RunCommand(shortcut, executable, params)
		if err != nil {
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
