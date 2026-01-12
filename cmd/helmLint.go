package cmd

import (
	"github.com/yuan-shuo/helm-gitops/pkg/helm"

	"github.com/spf13/cobra"
)

func init() { rootCmd.AddCommand(newLintCmd()) }

func newLintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lint",
		Short: "helm lint + unittest",
		RunE: func(cmd *cobra.Command, args []string) error {
			return helm.Lint()
		},
	}
}
