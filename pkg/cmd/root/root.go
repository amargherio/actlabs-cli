package root

import (
	"github.com/amargherio/actlabs-cli/pkg/cmd/setup"
	"github.com/spf13/cobra"
)

var cfgFile string

//// Execute adds all child commands to the root command and sets flags appropriately.
//// This is called by main.main(). It only needs to happen once to the rootCmd.
//func Execute() {
//	err := rootCmd.Execute()
//	if err != nil {
//		os.Exit(1)
//	}
//}

func NewRootCmd() (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "actlabs <command> <subcommand> [flags]",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Version: "0.1.0",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.actlabs-cli.yaml)")

	rootCmd.AddCommand(setup.NewSetupCmd())

	return rootCmd, nil
}
