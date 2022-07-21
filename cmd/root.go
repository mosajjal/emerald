package cmd

import "github.com/spf13/cobra"

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "emailtools",
		Short: "Email Parsing tools",
		Long:  `Emailtools is a swiss army knife of dealing with email-related investigations.`,
	}
)
