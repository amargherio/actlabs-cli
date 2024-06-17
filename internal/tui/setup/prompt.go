package setup

import (
	cmd "github.com/amargherio/actlabs-cli/pkg/cmd/setup"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const ACTLABS_APP_ID string = "bee16ca1-a401-40ee-bb6a-34349ebd993e"
const (
	ensureOwner state = iota
	createRG
	createStorage
	assignBlobDataContribRole
	deployLocal
	verifyLocal
)
const MAX_WIDTH = 80

type state int
type tickMsg struct{}
type Model struct {
	state
	setupParams 			  *huh.Form
	lg *lipgloss.Renderer
	//styles *Styles
	width int
	rgName                    string
	tenantID                  string
	subID                     string
	location                  string
}

func initializeModel(options *cmd.SetupOptions) Model {
	m := Model{width: MAX_WIDTH}
	m.lg = lipgloss.DefaultRenderer()
	m.rgName = options.ResourceGroupName
	m.tenantID = options.TenantID
	m.subID = options.SubscriptionID
	m.location = options.Location

	m.setupParams = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What would you like to name the resource group?").
				Prompt(">>").
				Validate(validateRGName).
				Value(&m.rgName),

			huh.NewInput().
				Title(`What is the ID of the Azure tenant containing the subscription you want to configure for use with ACTLabs?

(if you don't know, leave it blank and we'll work it out`).
				Prompt(">>").
				Value(&m.tenantID),

			huh.NewInput().
				Title(`What is the ID of the Azure subscription you want to configure for use with ACTLabs?

(if you don't know, leave it blank and we'll work it out`).
				Prompt(">>").
				Value(&m.subID),

			huh.NewInput().
				Title("What is the location for the ACTLabs resource group and default resources?").
				Prompt(">>").
				Value(&m.location),
		)


	ti := textinput.New()
	ti.Placeholder =
	return labsSetup{
		rgName:                    huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("What would you like to name the resource group?").Prompt(">>").Validate(validateRGName).,
		tenantID:                  "",
		subID:                     "",
		location:                  "",
		ensureOwner:               false,
		createRG:                  false,
		createStorage:             false,
		assignBlobDataContribRole: false,
		deployLocal:               false,
		verifyLocal:               false,
	}

}

func validateRGName(s string) error {
	// todo
	return nil
}

func InteractiveSetup(opts *cmd.SetupOptions) {
	var _, err = tea.NewProgram(initializeModel(opts)).Run()
	if err != nil {
		log.Error(err)
	}
}
