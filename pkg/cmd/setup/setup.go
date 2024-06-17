package setup

import (
	"fmt"
	"github.com/amargherio/actlabs-cli/internal/tui/setup"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var ResourceGroupName string
var TenantID string
var SubscriptionID string
var Location string

type SetupOptions struct {
	ResourceGroupName string
	TenantID          string
	SubscriptionID    string
	Location          string
	IsInteractive     bool
}

func NewSetupCmd() *cobra.Command {
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Configures the required infrastructure for ACTLabs to run in your sub.",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := parseFlags(cmd, args)
			return setupRun(opts)
		},
	}

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setupCmd.Flags().StringVarP(&ResourceGroupName,
		"resource-group",
		"g",
		"repro-project",
		"The name of the resource group ACTLabs will use when creating resources for labs and other deployments.")

	setupCmd.Flags().StringVarP(&Location,
		"location",
		"l",
		"eastus2",
		"The location for the ACTLabs resource group and default resources. Defaults to eastus2.")

	setupCmd.Flags().StringVar(&TenantID,
		"tenant-id",
		"",
		`The ID of the Azure tenant containing the subscription you want to configure for use with ACTLabs.

If this value is not provided, ACTLabs will attempt to use the value of the AZURE_TENANT_ID environment variable or the tenantId value provided by Azure CLI.`)

	setupCmd.Flags().StringVarP(&SubscriptionID,
		"subscription-id",
		"s",
		"",
		`The ID of the Azure subscription you want to configure for use with ACTLabs.

If this value is not provided, ACTLabs will attempt to use the value of the AZURE_SUBSCRIPTION_ID environment variable or the id value provided by Azure CLI.`)

	setupCmd.Flags().Bool("interactive", true, "Run the ACTLabs environment setup in an interactive mode.")
	setupCmd.Flags().Bool("local", false, "Run the server components of ACTLabs locally in a Docker container.")

	return setupCmd
}

func parseFlags(cmd *cobra.Command, args []string) SetupOptions {
	opts := SetupOptions{}

	interactive, err := cmd.Flags().GetBool("interactive")
	if err != nil {
		fmt.Println("Error getting interactive flag: ", err)
	} else {
		opts.IsInteractive = interactive
	}

	opts.ResourceGroupName = ResourceGroupName
	opts.TenantID = TenantID
	opts.SubscriptionID = SubscriptionID

	return opts
}

func setupRun(opts SetupOptions) error {
	fmt.Println("setup called")

	if !opts.IsInteractive {
		fmt.Println("Running in non-interactive mode.")
		return nil
	} else {
		var (
			_, err = tea.NewProgram(setup.InitializeModel(ResourceGroupName, TenantID, SubscriptionID, Location)).Run()
		)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}
