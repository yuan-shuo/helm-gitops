package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/scaffold"
)

var withActions bool

const ChartInitCommitMessage = "helm gitops chart init"
const PRmarkText = "[create-pr]"

func init() {
	chartCreateCmd := newChartCreateCmd()
	chartCreateCmd.Flags().BoolVar(&withActions, "actions", false, "also create .github/workflows/ci-test.yaml")
	rootCmd.AddCommand(chartCreateCmd)
}

func newChartCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [NAME]",
		Short: "create a new Helm chart with GitOps scaffold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return scaffold.CreateChart(args[0], withActions, ChartInitCommitMessage, PRmarkText)
		},
	}
}
