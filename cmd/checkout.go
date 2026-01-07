package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/git"
)

func init() { rootCmd.AddCommand(newCheckoutCmd()) }

func newCheckoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkout <branch-name>",
		Short: "switch to or create a new development branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return git.Checkout(args[0])
		},
	}
}
