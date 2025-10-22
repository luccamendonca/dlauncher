package launcher

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var CONFIG config

func runCmdGetFlagValues(cmd *cobra.Command, display CobraDisplay) (Shortcut, Executable, []string, error) {
	s := Shortcut{}
	e := Executable{}
	params := []string{}
	executableName, err := cmd.Flags().GetString("executable-name")
	if err != nil {
		return s, e, params, err
	}
	if executableName == "" {
		return s, e, params, fmt.Errorf("the executable-name must be provided")
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

func addCmdGetFlagValues(cmd *cobra.Command, display CobraDisplay) (string, Shortcut, error) {
	s := Shortcut{}
	name, err := cmd.Flags().GetString("shortcut-name")
	if err != nil {
		return name, s, err
	}
	if name == "" {
		name = display.Prompt("Shortcut name")
	}
	s.Template, err = cmd.Flags().GetString("shortcut-name")
	if err != nil {
		return name, s, err
	}
	if s.Template == "" {
		s.Template = display.Prompt("Shortcut template")
	}

	return name, s, nil
}

func multiRunCmdGetFlagValues(cmd *cobra.Command, display CobraDisplay) (Executable, []string, error) {
	e := Executable{}
	executableName, err := cmd.Flags().GetString("executable-name")
	if err != nil {
		return e, nil, err
	}
	if executableName == "" {
		return e, nil, fmt.Errorf("the executable-name must be provided")
	}

	linksInput, err := cmd.Flags().GetString("links")
	if err != nil {
		return e, nil, err
	}
	if linksInput == "" {
		linksInput = display.PromptMultiline("Enter links")
	}

	e, err = CONFIG.GetExecutable(executableName)
	if err != nil {
		return e, nil, err
	}

	links := strings.Split(strings.TrimSpace(linksInput), "\n")
	var cleanLinks []string
	for _, link := range links {
		link = strings.TrimSpace(link)
		if link != "" {
			cleanLinks = append(cleanLinks, link)
		}
	}

	return e, cleanLinks, nil
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Go and launcht it!",
	Args:  cobra.NoArgs,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a shortcut",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		useGUI, _ := cmd.Flags().GetBool("use-gui")
		display := NewDisplay(useGUI, args)
		shortcut, executable, params, err := runCmdGetFlagValues(cmd, display)
		if err != nil {
			display.Error(err.Error())
			return
		}
		err = RunCommand(shortcut, executable, params)
		if err != nil {
			display.Error(err.Error())
			return
		}
	},
}

var multiRunCmd = &cobra.Command{
	Use:   "multi-run",
	Short: "Opens multiple links in new browser tabs",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		useGUI, _ := cmd.Flags().GetBool("use-gui")
		display := NewDisplay(useGUI, args)
		executable, links, err := multiRunCmdGetFlagValues(cmd, display)
		if err != nil {
			display.Error(err.Error())
			return
		}

		if len(links) == 0 {
			display.Error("No valid links provided")
			return
		}

		err = RunMultipleCommands(links, executable)
		if err != nil {
			display.Error(err.Error())
			return
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a command to your config file",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		useGUI, _ := cmd.Flags().GetBool("use-gui")
		display := NewDisplay(useGUI, args)
		name, shortcut, err := addCmdGetFlagValues(cmd, display)
		if err != nil {
			display.Error(err.Error())
			return
		}
		err = CONFIG.AddShortcut(name, shortcut)
		if err != nil {
			display.Error(err.Error())
			return
		}
		fmt.Printf("Done! Added shortcut '%s' -> '%s'", name, shortcut.Template)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all configured shortcuts",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		useGUI, _ := cmd.Flags().GetBool("use-gui")
		display := NewDisplay(useGUI, args)
		if len(CONFIG.Shortcuts) == 0 {
			display.Info("No shortcuts configured.")
			return
		}
		var sb strings.Builder
		sb.WriteString("Configured shortcuts:\n")
		for name, shortcut := range CONFIG.Shortcuts {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", name, shortcut.Template))
		}
		display.Info(sb.String())
	},
}

func init() {
	runCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")
	runCmd.Flags().StringP("executable-name", "e", "", "The program that should execute your command template.")
	runCmd.Flags().StringP("shortcut-name", "s", "", "The name of the shortcut.")
	runCmd.Flags().StringArrayP("params", "p", []string{}, "(optional) The params for the command.")

	addCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")
	addCmd.Flags().StringP("shortcut-name", "s", "", "The name of the shortcut.")
	addCmd.Flags().StringP("shortcut-template", "t", "", "The template for the shortcut.")

	multiRunCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")
	multiRunCmd.Flags().StringP("executable-name", "e", "", "The browser executable to use (e.g., chrome, firefox).")
	multiRunCmd.Flags().StringP("links", "l", "", "Links separated by newlines to open in new tabs.")

	listCmd.Flags().BoolP("use-gui", "g", false, "Uses GUI instead of CLI")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(multiRunCmd)
	rootCmd.AddCommand(listCmd)

	var err error
	CONFIG, err = ParseConfig()
	if err != nil {
		panic(err)
	}
}

func Execute() {
	if err := runCmd.Execute(); err != nil {
		panic(err)
	}
}
