package cli

import (
	"github.com/spf13/cobra"
)

// version is set at build time via ldflags.
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "launchpad",
	Short: "AI-powered instruction scaffolder for your projects",
	Long: `Launchpad sets up opinionated AI coding instructions through a
brief conversation about what you're building.

It generates .github/copilot-instructions.md, scoped .instructions.md
files, and AGENTS.md — all tailored to your stack and style.

Powered by OpenAI. Your copilot should write code the way you would.`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
