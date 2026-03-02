package styles

import "github.com/charmbracelet/lipgloss"

// Chill color palette - pastels with dark background
var (
	// Primary colors
	colorPastelBlue   = lipgloss.Color("#A8D8EA")
	colorPastelGreen  = lipgloss.Color("#B4F8C8")
	colorPastelPurple = lipgloss.Color("#AA96DA")
	colorPastelPink   = lipgloss.Color("#FFAAA5")
	colorPastelYellow = lipgloss.Color("#FFE5B4")

	// Background colors
	colorDarkBg      = lipgloss.Color("#1E1E2E")
	colorLighterBg   = lipgloss.Color("#313244")
	colorSelectedBg  = lipgloss.Color("#45475A")
	colorHighlightBg = lipgloss.Color("#585B70")

	// Text colors
	colorPrimaryText = lipgloss.Color("#CDD6F4")
	colorMutedText   = lipgloss.Color("#A6ADC8")
	colorSubtleText  = lipgloss.Color("#6C7086")
	colorSuccessText = lipgloss.Color("#A6E3A1")
	colorErrorText   = lipgloss.Color("#F38BA8")
	colorWarningText = lipgloss.Color("#F9E2AF")
)

// Color styles for direct use with .Render()
var (
	PastelBlue   = lipgloss.NewStyle().Foreground(colorPastelBlue)
	PastelGreen  = lipgloss.NewStyle().Foreground(colorPastelGreen)
	PastelPurple = lipgloss.NewStyle().Foreground(colorPastelPurple)
	PastelPink   = lipgloss.NewStyle().Foreground(colorPastelPink)
	PastelYellow = lipgloss.NewStyle().Foreground(colorPastelYellow)

	PrimaryText = lipgloss.NewStyle().Foreground(colorPrimaryText)
	MutedText   = lipgloss.NewStyle().Foreground(colorMutedText)
	SubtleText  = lipgloss.NewStyle().Foreground(colorSubtleText)
)

// Common styles
var (
	// Base container style
	BaseStyle = lipgloss.NewStyle().
			Background(colorDarkBg).
			Foreground(colorPrimaryText)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(colorPastelBlue).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(colorMutedText).
			Italic(true)

	// Menu and list styles
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(colorPrimaryText).
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(colorPastelBlue).
				Background(colorSelectedBg).
				Bold(true).
				PaddingLeft(2)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Foreground(colorPrimaryText).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPastelPurple).
			Padding(0, 1)

	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(colorPrimaryText).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(colorPastelBlue).
				Padding(0, 1)

	// Label styles
	LabelStyle = lipgloss.NewStyle().
			Foreground(colorMutedText).
			MarginBottom(0)

	FocusedLabelStyle = lipgloss.NewStyle().
				Foreground(colorPastelBlue).
				Bold(true).
				MarginBottom(0)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(colorDarkBg).
			Background(colorPastelPurple).
			Padding(0, 2).
			MarginRight(1)

	FocusedButtonStyle = lipgloss.NewStyle().
				Foreground(colorDarkBg).
				Background(colorPastelBlue).
				Bold(true).
				Padding(0, 2).
				MarginRight(1)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(colorSuccessText)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(colorErrorText)

	WarningStyle = lipgloss.NewStyle().
			Foreground(colorWarningText)

	// Help bar style
	HelpStyle = lipgloss.NewStyle().
			Foreground(colorSubtleText).
			MarginTop(1)

	// Divider
	DividerStyle = lipgloss.NewStyle().
			Foreground(colorSubtleText)

	// Box styles for panels
	PanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorLighterBg).
			Padding(1, 2)

	FocusedPanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(colorPastelBlue).
				Padding(1, 2)

	// Checkbox styles
	CheckboxStyle = lipgloss.NewStyle().
			Foreground(colorPrimaryText)

	CheckedStyle = lipgloss.NewStyle().
			Foreground(colorPastelGreen)

	// Branch/commit styles
	BranchStyle = lipgloss.NewStyle().
			Foreground(colorPastelPurple)

	CurrentBranchStyle = lipgloss.NewStyle().
				Foreground(colorPastelGreen).
				Bold(true)

	CommitHashStyle = lipgloss.NewStyle().
			Foreground(colorPastelYellow)

	CommitMsgStyle = lipgloss.NewStyle().
			Foreground(colorPrimaryText)

	// File status styles
	StagedStyle = lipgloss.NewStyle().
			Foreground(colorPastelGreen)

	ModifiedStyle = lipgloss.NewStyle().
			Foreground(colorPastelYellow)

	UntrackedStyle = lipgloss.NewStyle().
			Foreground(colorPastelPink)

	DeletedStyle = lipgloss.NewStyle().
			Foreground(colorErrorText)

	// Error text color for use in panel borders
	ErrorText = colorErrorText
)

// Logo renders the Git-Graft logo
func Logo() string {
	logo := `
  ╔═╗╦╔╦╗  ╔═╗╦═╗╔═╗╔═╗╔╦╗
  ║ ╦║ ║───║ ╦╠╦╝╠═╣╠╣  ║
  ╚═╝╩ ╩   ╚═╝╩╚═╩ ╩╚   ╩ `
	return lipgloss.NewStyle().
		Foreground(colorPastelBlue).
		Bold(true).
		Render(logo)
}

// Divider creates a horizontal divider
func Divider(width int) string {
	return DividerStyle.Render(lipgloss.NewStyle().
		Width(width).
		Render("─────────────────────────────────────────────────"))
}
