package styles

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	Red      = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	Blue     = lipgloss.AdaptiveColor{Light: "#00A2ED", Dark: "#00A2ED"}
	Green    = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	Yellow   = lipgloss.AdaptiveColor{Light: "", Dark: ""}
	White    = lipgloss.Color("#FFFFFF")
	OffWhite = lipgloss.Color("#DDD")
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help,
	CursorStyle,
	FocusedStyle,
	BlurredStyle lipgloss.Style
	FocusedButton,
	BlurredButton string
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().Foreground(Blue).Bold(true).Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(White).PaddingLeft(1).MarginTop(1)
	s.StatusHeader = lg.NewStyle().Foreground(Blue).Bold(true)
	s.Highlight = lg.NewStyle().Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.Foreground(Red)
	s.Help = lg.NewStyle().Foreground(lipgloss.Color("240"))

	// General styling
	s.FocusedStyle = lipgloss.NewStyle().Foreground(Blue).MarginRight(1).Bold(true)
	s.BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginRight(1)

	// Cursor styling
	s.CursorStyle = lipgloss.NewStyle().Foreground(White)

	// Button styling
	s.FocusedButton = s.FocusedStyle.Bold(true).Render("[ Submit ]")
	s.BlurredButton = fmt.Sprintf("[ %s ]", s.BlurredStyle.Render("Submit"))

	return &s
}
