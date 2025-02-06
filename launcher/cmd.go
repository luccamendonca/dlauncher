package launcher

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dlanucher",
	Short: "go and launch it",
}
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a shortcut",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := ParseConfig()
		useGUI, _ := cmd.Flags().GetBool("use-gui")
		display := NewDisplay(useGUI, args)
		shortcutName, _ := cmd.Flags().GetString("shortcut-name")
		if shortcutName == "" {
			shortcutName = display.Prompt("Shortcut name")
		}
		shortcutParams, _ := cmd.Flags().GetStringArray("params")
		hasParams, _ := cmd.Flags().GetBool("has-params")
		if hasParams && len(shortcutParams) == 0 {
			shortcutParams = []string{display.Prompt("Params for template")}
		}
		executableName, _ := cmd.Flags().GetString("executable-name")
		RunCommand(c, executableName, shortcutName, shortcutParams)
	},
}

func init() {
	rootCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")

	runCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")
	runCmd.Flags().StringP("executable-name", "e", "", "The program that should execute your command template.")
	runCmd.Flags().StringP("shortcut-name", "s", "", "The name of the shortcut.")
	runCmd.Flags().BoolP("has-params", "a", false, "Tells the launcher whether or not you'll provide params for the command")
	runCmd.Flags().StringArrayP("params", "p", []string{}, "(optional) The params for the command. In GUI mode, this will only work if --has-params is set.")

	rootCmd.AddCommand(runCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
