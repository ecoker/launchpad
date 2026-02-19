package cli

import (
	"fmt"

	"github.com/ehrencoker/agent-kit/internal/scaffold"
	"github.com/ehrencoker/agent-kit/internal/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show the template knowledge base used for generation",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(ui.Banner)

		fmt.Println(ui.Heading.Render("Template knowledge base:"))
		fmt.Println()
		fmt.Println(ui.DimStyle.Render("  Launchpad narrows to these repository-defined assets during conversation,"))
		fmt.Println(ui.DimStyle.Render("  then generates instructions from the selected subset."))
		fmt.Println()

		fmt.Println(ui.Heading.Render("  Language & framework profiles:"))
		for _, p := range scaffold.Profiles {
			fmt.Printf("    %s  %s\n", ui.ProfileID.Render(p.ID), ui.ProfileDesc.Render(p.Summary))
		}
		fmt.Println()

		fmt.Println(ui.Heading.Render("  Specialized add-ons:"))
		for _, a := range scaffold.Addons {
			fmt.Printf("    %s  %s\n", ui.ProfileID.Render(a.ID), ui.ProfileDesc.Render(a.Summary))
		}
		fmt.Println()

		return nil
	},
}
