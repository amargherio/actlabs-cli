package styles

import "github.com/charmbracelet/lipgloss"

var (
	Red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	Blue   = lipgloss.AdaptiveColor{Light: "#00A2ED", Dark: "#00A2ED"}
	Green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	Yellow = lipgloss.AdaptiveColor{Light: "", Dark: ""}
	White  = lipgloss.Color("#FFFFFF")
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().Foreground(Green).Bold(true).Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(White).PaddingLeft(1).MarginTop(1)
	s.StatusHeader = lg.NewStyle().Foreground(Blue).Bold(true)
	s.Highlight = lg.NewStyle().Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.Foreground(Red)
	s.Help = lg.NewStyle().Foreground(lipgloss.Color("240"))

	return &s
}
