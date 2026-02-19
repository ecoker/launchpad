package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Cyan    = lipgloss.Color("86")
	Magenta = lipgloss.Color("205")
	Green   = lipgloss.Color("82")
	Yellow  = lipgloss.Color("220")
	Red     = lipgloss.Color("196")
	Blue    = lipgloss.Color("75")
	Dim     = lipgloss.Color("241")
	White   = lipgloss.Color("255")

	// Styles
	Bold = lipgloss.NewStyle().Bold(true)

	Heading = lipgloss.NewStyle().Bold(true).Foreground(Cyan)

	Accent = lipgloss.NewStyle().Bold(true).Foreground(Magenta)

	Success = lipgloss.NewStyle().Foreground(Green)

	Warning = lipgloss.NewStyle().Foreground(Yellow)

	Error = lipgloss.NewStyle().Bold(true).Foreground(Red)

	DimStyle = lipgloss.NewStyle().Foreground(Dim)

	FileStyle = lipgloss.NewStyle().Foreground(Blue).Underline(true)

	ProfileID = lipgloss.NewStyle().Bold(true).Foreground(Cyan)

	ProfileDesc = lipgloss.NewStyle().Foreground(Dim)
)

func buildBanner() string {
	cyanBold := lipgloss.NewStyle().Bold(true).Foreground(Cyan)
	magentaBold := lipgloss.NewStyle().Bold(true).Foreground(Magenta)
	dim := lipgloss.NewStyle().Foreground(Dim)

	top := cyanBold.Render("   \u250c\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2510")
	mid1 := cyanBold.Render("   \u2502           ") + magentaBold.Render("\U0001f680 launchpad") + cyanBold.Render("                   \u2502")
	mid2 := cyanBold.Render("   \u2502   ") + dim.Render("AI-powered coding instruction setup") + cyanBold.Render("    \u2502")
	bot := cyanBold.Render("   \u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2518")

	return "\n" + top + "\n" + mid1 + "\n" + mid2 + "\n" + bot + "\n"
}

var Banner = buildBanner()
