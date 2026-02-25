package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ecoker/launchpad/internal/ai"
	"github.com/ecoker/launchpad/internal/scaffold"
	"github.com/ecoker/launchpad/internal/ui"
	"github.com/spf13/cobra"
)

var (
	flagForce bool
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Start a conversation to generate tailored AI instructions",
	Long: `Have a brief conversation about what you're building, then Launchpad
generates customized AI coding instructions for your project.

Set OPENAI_API_KEY in your environment before running.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVarP(&flagForce, "force", "f", false, "Overwrite files in non-empty target")
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Print(ui.Banner)

	// 1. Check for API key (env var, then .env file, then prompt)
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = loadKeyFromDotEnv()
	}
	if apiKey == "" {
		fmt.Println(ui.Warning.Render("No OPENAI_API_KEY found in environment."))
		fmt.Println()
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Paste your OpenAI API key:").
					EchoMode(huh.EchoModePassword).
					Value(&apiKey),
			),
		).Run()
		if err != nil {
			return err
		}
		if apiKey == "" {
			return fmt.Errorf("an OpenAI API key is required — get one at https://platform.openai.com/api-keys")
		}
	}

	// 2. Target directory
	targetDir := ""
	if len(args) > 0 {
		targetDir = args[0]
	}
	if targetDir == "" {
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Where should we set up the project?").
					Placeholder("./my-app").
					Value(&targetDir),
			),
		).Run()
		if err != nil {
			return err
		}
		if targetDir == "" {
			targetDir = "./my-app"
		}
	}

	outputPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}
	projectName := filepath.Base(outputPath)

	// 3. Safety check for non-empty directory
	if !flagForce {
		entries, _ := os.ReadDir(outputPath)
		if len(entries) > 0 {
			force := false
			err := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title("Directory isn't empty. Overwrite existing files?").
						Affirmative("Yes, overwrite").
						Negative("No, abort").
						Value(&force),
				),
			).Run()
			if err != nil {
				return err
			}
			if !force {
				return fmt.Errorf("aborted — directory is not empty")
			}
		}
	}

	// 4. Conversation — natural language with loading spinners
	fmt.Println()
	fmt.Println(ui.Heading.Render("What are you building?"))
	fmt.Println(ui.DimStyle.Render("Describe your project and I'll help you pick the right stack and standards."))
	fmt.Println()

	// Build LLM provider — model is configurable via LAUNCHPAD_MODEL env var.
	var providerOpts []ai.OpenAIOption
	if model := os.Getenv("LAUNCHPAD_MODEL"); model != "" {
		providerOpts = append(providerOpts, ai.WithModel(model))
	}
	provider := ai.NewOpenAIProvider(apiKey, providerOpts...)
	engine := ai.NewEngine(provider)

	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(ui.Accent.Render("You: "))
	firstInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}
	firstInput = strings.TrimSpace(firstInput)
	if firstInput == "" {
		return fmt.Errorf("please describe what you're building")
	}

	fmt.Println()
	spin := ui.NewSpinner("Thinking...")
	reply, err := engine.Chat(ctx, fmt.Sprintf(
		"Project name: %q. What I'm building: %s", projectName, firstInput,
	))
	spin.Stop()
	if err != nil {
		return fmt.Errorf("conversation error: %w", err)
	}
	printLaunchpadReply(reply)

	for !ai.IsReady(reply) {
		fmt.Print(ui.Accent.Render("You: "))
		userInput, readErr := reader.ReadString('\n')
		if readErr != nil {
			return fmt.Errorf("reading input: %w", readErr)
		}
		userInput = strings.TrimSpace(userInput)
		if userInput == "" || strings.EqualFold(userInput, "/done") {
			break
		}

		fmt.Println()
		spin = ui.NewSpinner("Thinking...")
		reply, err = engine.Chat(ctx, userInput)
		spin.Stop()
		if err != nil {
			return fmt.Errorf("conversation error: %w", err)
		}
		printLaunchpadReply(reply)
	}

	// 5. Silent extraction — user never sees this
	spin = ui.NewSpinner("Resolving selection...")
	sel, err := engine.ExtractDecision(ctx)
	spin.Stop()
	if err != nil {
		return fmt.Errorf("extracting decision: %w", err)
	}

	fmt.Println()
	printSelectionSummary(sel)

	// 6. Generate files
	spin = ui.NewSpinner("Generating instruction files...")
	fmt.Println()

	files, err := engine.GenerateFiles(ctx, projectName, sel)
	spin.Stop()
	if err != nil {
		return fmt.Errorf("generation error: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files were generated — try running again with more detail about your project")
	}

	// 6. Write files
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	var created []string
	for _, f := range files {
		fullPath := filepath.Join(outputPath, f.Path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", f.Path, err)
		}
		if err := os.WriteFile(fullPath, []byte(f.Content+"\n"), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", f.Path, err)
		}
		created = append(created, fullPath)
	}

	// 7. Print results
	ui.PrintFileTree(created, outputPath)

	displayPath := ui.DisplayPath(outputPath)
	fmt.Printf("%s Generated %s instruction files in %s\n",
		ui.Success.Render("✔"),
		ui.Accent.Render(fmt.Sprintf("%d", len(created))),
		ui.FileStyle.Render(displayPath),
	)
	fmt.Println()
	fmt.Println(ui.Heading.Render("Next steps:"))
	fmt.Printf("  %s cd %s\n", ui.DimStyle.Render("1."), ui.FileStyle.Render(displayPath))
	fmt.Printf("  %s Review the generated files — tweak anything that doesn't feel right\n", ui.DimStyle.Render("2."))

	// Show scaffold command if available for the selected profile
	if profile := scaffold.FindProfile(sel.ProfileID); profile != nil && profile.ScaffoldCmd != "" {
		scaffoldDisplay := strings.ReplaceAll(profile.ScaffoldCmd, "{{name}}", projectName)
		scaffoldDisplay = strings.ReplaceAll(scaffoldDisplay, "{{module}}", projectName)
		fmt.Printf("  %s Scaffold your project: %s\n", ui.DimStyle.Render("3."), ui.Accent.Render(scaffoldDisplay))
		fmt.Printf("  %s Open Copilot Chat and type %s to start building\n", ui.DimStyle.Render("4."), ui.Accent.Render("/start"))
	} else {
		fmt.Printf("  %s Open Copilot Chat and type %s to bootstrap the project\n", ui.DimStyle.Render("3."), ui.Accent.Render("/start"))
	}

	fmt.Println()
	fmt.Println(ui.DimStyle.Render("Your AI copilot is briefed. Go build something great."))
	fmt.Println()

	return nil
}

func printSelectionSummary(sel *ai.Selection) {
	fmt.Printf("%s %s\n", ui.DimStyle.Render("Profile:"), ui.ProfileID.Render(sel.ProfileID))
	if len(sel.AddonIDs) > 0 {
		fmt.Printf("%s %s\n", ui.DimStyle.Render("Add-ons: "), strings.Join(sel.AddonIDs, ", "))
	}
	if len(sel.AssetIDs) > 0 {
		fmt.Printf("%s %s\n", ui.DimStyle.Render("Assets:  "), strings.Join(sel.AssetIDs, ", "))
	}
	if sel.Rationale != "" {
		fmt.Printf("%s %s\n", ui.DimStyle.Render("Why:     "), sel.Rationale)
	}
	fmt.Println()
}

// printLaunchpadReply displays the AI response, stripping the READY_TO_GENERATE token.
func printLaunchpadReply(reply string) {
	display := strings.ReplaceAll(reply, "READY_TO_GENERATE", "")
	display = strings.ReplaceAll(display, "READY TO GENERATE", "")
	display = strings.TrimSpace(display)
	fmt.Print(ui.DimStyle.Render("Launchpad: "))
	fmt.Println(display)
	fmt.Println()
}

// loadKeyFromDotEnv reads OPENAI_API_KEY or KEY from a .env file in the current directory.
// Handles common formats: quoted values, `export` prefix, inline comments.
func loadKeyFromDotEnv() string {
	data, err := os.ReadFile(".env")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		// Strip optional `export ` prefix
		line = strings.TrimPrefix(line, "export ")
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		// Strip surrounding quotes (single or double)
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		// Strip inline comments (only after unquoted values)
		if ci := strings.Index(val, " #"); ci != -1 {
			val = strings.TrimSpace(val[:ci])
		}
		if key == "OPENAI_API_KEY" || key == "KEY" {
			return val
		}
	}
	return ""
}
