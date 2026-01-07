package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/scaffold"
)

func init() { rootCmd.AddCommand(newCreateCmd()) }

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [NAME]",
		Short: "create a new Helm chart with GitOps scaffold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return scaffold.Create(args[0])
		},
	}
}
