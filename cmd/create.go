package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/scaffold"
)

var withActions bool

const InitCommitMessage = "helm gitops chart init"
const PRmarkText = "[create-pr]"

func init() {
	createCmd := newCreateCmd()
	createCmd.Flags().BoolVar(&withActions, "actions", false, "also create .github/workflows/ci-test.yaml")
	rootCmd.AddCommand(createCmd)
}

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [NAME]",
		Short: "create a new Helm chart with GitOps scaffold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return scaffold.Create(args[0], withActions, InitCommitMessage, PRmarkText)
		},
	}
}
