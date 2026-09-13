package main

import "github.com/charmbracelet/lipgloss"

var (
	colorX        = lipgloss.Color("#04B575")
	colorO        = lipgloss.Color("#EE6FF8")
	colorMuted    = lipgloss.Color("#383838") 
	colorSubtle   = lipgloss.Color("#767676") 
	colorText     = lipgloss.Color("#F9FAFA") 

	xStyle = lipgloss.NewStyle().Foreground(colorX).Bold(true)
	oStyle = lipgloss.NewStyle().Foreground(colorO).Bold(true)

	appLayout = lipgloss.NewStyle().Margin(1, 2)

	titleStyle = lipgloss.NewStyle().
		Foreground(colorO).
		Bold(true).
		MarginBottom(1)

	cellStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorMuted).
		Width(7).
		Height(3).
		Align(lipgloss.Center, lipgloss.Center)

	cursorStyle = cellStyle.Copy().
		BorderForeground(colorText).
		Background(lipgloss.Color("#1A1A24"))

	boardContainer = lipgloss.NewStyle().MarginRight(4)

	statusStyle = lipgloss.NewStyle().
		Foreground(colorSubtle).
		MarginTop(1)
		
	activeTurnStyle = lipgloss.NewStyle().
		Foreground(colorText).
		Bold(true)

	identityStyle = lipgloss.NewStyle().
		Foreground(colorSubtle).
		MarginBottom(1)

	// Keep the chat box simple and reliable
	chatContainer = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderTop(false).BorderRight(false).BorderBottom(false).
		BorderForeground(colorMuted).
		PaddingLeft(4).
		Width(38).
		Height(17) 

	chatHeaderStyle = lipgloss.NewStyle().
		Foreground(colorSubtle).
		MarginBottom(1)

	chatSenderXStyle = lipgloss.NewStyle().Foreground(colorX).Bold(true)
	chatSenderOStyle = lipgloss.NewStyle().Foreground(colorO).Bold(true)
	chatTextStyle    = lipgloss.NewStyle().Foreground(colorText).Width(34) // Force text wrap safely

	chatInputBox = lipgloss.NewStyle().
		MarginTop(1).
		Foreground(colorText)
)
