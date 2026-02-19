package ui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PrintFileTree prints the list of created files in a pretty tree.
func PrintFileTree(files []string, rootDir string) {
	fmt.Println(Heading.Render("\nCreated files:"))
	fmt.Println()

	sorted := make([]string, len(files))
	copy(sorted, files)
	sort.Strings(sorted)

	prefix := rootDir + "/"
	for _, f := range sorted {
		rel := strings.TrimPrefix(f, prefix)
		fmt.Printf("  %s %s\n", DimStyle.Render("└─"), FileStyle.Render(rel))
	}
	fmt.Println()
}

// PrintDone prints the success message and next steps.
func PrintDone(profileLabel, targetDir string) {
	fmt.Printf("%s Scaffolded %s into %s\n",
		Success.Render("✔"),
		Accent.Render(profileLabel),
		FileStyle.Render(targetDir),
	)
	fmt.Println()
	fmt.Println(Heading.Render("Next steps:"))
	fmt.Printf("  %s cd %s\n", DimStyle.Render("1."), FileStyle.Render(targetDir))
	fmt.Printf("  %s Review %s — your always-on standards\n", DimStyle.Render("2."), FileStyle.Render(".github/copilot-instructions.md"))
	fmt.Printf("  %s Browse %s — language & framework rules\n", DimStyle.Render("3."), FileStyle.Render(".github/instructions/"))
	fmt.Printf("  %s Edit freely — these are %s opinions now\n", DimStyle.Render("4."), lipgloss.NewStyle().Italic(true).Render("your"))
	fmt.Println()
	fmt.Println(DimStyle.Render("Happy building. Write something beautiful. ✨"))
	fmt.Println()
}

// DisplayPath returns a clean display path: relative if under cwd, absolute otherwise.
func DisplayPath(outputPath string) string {
	cwd, err := filepath.Abs(".")
	if err != nil {
		return outputPath
	}
	rel, err := filepath.Rel(cwd, outputPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return outputPath
	}
	if rel == "" {
		return "."
	}
	return rel
}
