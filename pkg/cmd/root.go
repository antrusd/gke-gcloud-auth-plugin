package cmd

import (
	"os"

	"github.com/antrusd/gke-gcloud-auth-plugin/pkg/auth"
	"github.com/antrusd/gke-gcloud-auth-plugin/pkg/conf"
	"github.com/spf13/cobra"
)

func GetRootCmd(args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               conf.BinName,
		Short:             "GKE Authentication Plugin",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		Long:              `GKE Authentication Plugin`,
		RunE: func(c *cobra.Command, args []string) error {
			return auth.Gcp(c.Context())
		},
	}

	var appcreds = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(versionCmd())
	rootCmd.Flags().StringVarP(&conf.AppCreds, "credential", "c", appcreds, "override Google application credential")
	rootCmd.SetArgs(args)

	return rootCmd
}
