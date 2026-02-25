package cli

import (
	"fmt"

	"github.com/ecoker/launchpad/internal/scaffold"
	"github.com/ecoker/launchpad/internal/ui"
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

		fmt.Println(ui.Heading.Render("  ★ Canonical stacks (coherence-first philosophy):"))
		for _, p := range scaffold.Profiles {
			if p.Tier != 1 {
				continue
			}
			layerTag := fmt.Sprintf("[%s]", p.Layer)
			fmt.Printf("    %s  %s  %s\n", ui.ProfileID.Render(p.ID), ui.DimStyle.Render(layerTag), ui.ProfileDesc.Render(p.Summary))
			if p.ScaffoldCmd != "" {
				fmt.Printf("    %s  %s\n", ui.DimStyle.Render("  scaffold:"), ui.DimStyle.Render(p.ScaffoldCmd))
			}
		}
		fmt.Println()

		fmt.Println(ui.Heading.Render("  Additional supported stacks:"))
		for _, p := range scaffold.Profiles {
			if p.Tier == 1 {
				continue
			}
			layerTag := fmt.Sprintf("[%s]", p.Layer)
			fmt.Printf("    %s  %s  %s\n", ui.ProfileID.Render(p.ID), ui.DimStyle.Render(layerTag), ui.ProfileDesc.Render(p.Summary))
			if p.ScaffoldCmd != "" {
				fmt.Printf("    %s  %s\n", ui.DimStyle.Render("  scaffold:"), ui.DimStyle.Render(p.ScaffoldCmd))
			}
		}
		fmt.Println()

		fmt.Println(ui.Heading.Render("  Specialized add-ons:"))
		for _, a := range scaffold.Addons {
			fmt.Printf("    %s  %s\n", ui.ProfileID.Render(a.ID), ui.ProfileDesc.Render(a.Summary))
		}
		fmt.Println()
		fmt.Println(ui.DimStyle.Render("  UI stacks automatically include frontend-craft, a default palette,"))
		fmt.Println(ui.DimStyle.Render("  and font pairing. No opt-in needed."))
		fmt.Println()

		return nil
	},
}
