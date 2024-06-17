package setup

import (
	"fmt"
	"github.com/amargherio/actlabs-cli/internal/tui/styles"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

const ACTLABS_APP_ID string = "bee16ca1-a401-40ee-bb6a-34349ebd993e"
const MAX_WIDTH = 80
const (
	AzureAuth state = iota
	EnsureOwner
	CreateRG
	CreateStorage
	AssignBlobDataContribRole
	DeployLocal
	VerifyLocal
)

type state int
type tickMsg struct{}
type frameMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

type Model struct {
	state
	lg           *lipgloss.Renderer
	styles       *styles.Styles
	width        int
	focusIndex   int
	setupParams  []textinput.Model
	spinner      spinner.Model
	progress     progress.Model
	provisioning bool
}

func InitializeModel(tenant string, sub string, location string) *Model {
	m := Model{width: MAX_WIDTH}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = styles.NewStyles(m.lg)
	m.focusIndex = 0
	m.provisioning = false

	m.setupParams = buildTextInputs(m.styles, tenant, sub, location)

	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Line
	m.progress = progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	return &m
}

func buildTextInputs(s *styles.Styles, tenant string, sub string, location string) []textinput.Model {
	var inputs []textinput.Model
	params := map[string]string{
		"tenant":       "What is the ID of the Azure tenant containing the subscription you want to configure for use with ACTLabs? (if you're not sure, leave blank)",
		"subscription": "What is the ID of the Azure subscription you want to configure for use with ACTLabs? (if you're not sure, leave blank)",
		"location":     "What Azure region would you like to use for your ACTLabs resource group and default resources? (default: eastus2)",
	}

	for k, v := range params {
		input := textinput.New()
		input.Cursor.Style = s.CursorStyle
		input.Prompt = v
		input.PromptStyle = lipgloss.NewStyle().Foreground(styles.OffWhite).Bold(true).MarginRight(2)

		switch k {
		case "tenant":
			input.Placeholder = tenant
		case "subscription":
			input.Placeholder = sub
		case "location":
			input.Placeholder = location
		default:
			input.Placeholder = ""
		}

		inputs = append(inputs, input)
	}
	return inputs
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			key := msg.String()

			if key == "enter" && m.focusIndex == len(m.setupParams) {
				// we've completed the form, so let's rock and roll
				fmt.Println("Focus index: ", m.focusIndex)
				m.provisioning = true

				//provisionLabs(&m)

				return m, tea.Quit
			}

			// cycle the index for input
			if key == "up" || key == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.setupParams) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 0
			}

			cmds := make([]tea.Cmd, len(m.setupParams))
			for i := 0; i < len(m.setupParams); i++ {
				if i == m.focusIndex {
					// set focused state for the correct input
					cmds[i] = m.setupParams[i].Focus()
					m.setupParams[i].PromptStyle = m.styles.FocusedStyle
					m.setupParams[i].Cursor.Style = m.styles.CursorStyle
					m.setupParams[i].TextStyle = m.styles.FocusedStyle

					continue
				}

				if i < m.focusIndex {
					// remove the focused state (if any) from the other inputs
					m.setupParams[i].Blur()
					m.setupParams[i].PromptStyle = m.styles.BlurredStyle
					m.setupParams[i].Cursor.Style = m.styles.BlurredStyle
				}

			}

			return m, tea.Batch(cmds...)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}

	// handle character inputs/blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n%s\n\n", m.styles.HeaderText.Render("ACTLabs Environment Setup"))

	// bold and focus the first one since we're starting new
	if m.focusIndex == 0 {
		m.setupParams[0].PromptStyle = m.styles.FocusedStyle
		m.setupParams[0].Focus()
	}

	for i := range m.setupParams {
		fmt.Fprintf(&b, "%s\n", m.setupParams[i].View())
	}

	button := &m.styles.BlurredButton
	if m.focusIndex == len(m.setupParams) {
		button = &m.styles.FocusedButton
	}

	if !m.provisioning {
		fmt.Fprintf(&b, "\n%s\n", *button)
	} else {
		fmt.Fprintf(&b, "\n%s\n", generateStatusMsg(m))
	}

	return b.String()
}

func generateStatusMsg(m Model) string {
	var b strings.Builder
	switch m.state {
	case AzureAuth:
		fmt.Fprintf(&b, "\n%s - %s\n", m.spinner.View(), "Preparing to start resource creation for ACTLabs...")
	case EnsureOwner:
		fmt.Fprintf(&b, "\n%s - %s\n", m.spinner.View(), "Ensuring you have the necessary permissions to create resources...")
	case CreateRG:
		fmt.Fprintf(&b, "\n%s - %s\n", m.spinner.View(), "Creating the ACTLabs resource group...")
	case CreateStorage:
		fmt.Fprintf(&b, "\n%s - %s\n", m.spinner.View(), "Creating the ACTLabs storage account...")
	case AssignBlobDataContribRole:
		fmt.Fprintf(&b, "\n%s - %s\n", m.spinner.View(), "Assigning the Blob Data Contributor role to the ACTLabs storage account...")
	case DeployLocal:
		fmt.Fprintf(&b, "\n%s - %s\n", m.spinner.View(), "Deploying the ACTLabs server components locally...")
	case VerifyLocal:
		fmt.Fprintf(&b, "\n%s - %s\n", m.spinner.View(), "Verifying the ACTLabs server components are running locally...")
	}
	return b.String()
}

func (m Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.setupParams))
	for i := range m.setupParams {
		m.setupParams[i], cmds[i] = m.setupParams[i].Update(msg)
	}
	return tea.Batch(cmds...)
}
