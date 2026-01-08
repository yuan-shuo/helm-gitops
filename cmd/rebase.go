package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/git"
)

var rebaseOnto string

func init() {
	rebaseCmd := newRebaseCmd()
	rebaseCmd.Flags().StringVarP(&rebaseOnto, "branch", "b", "main", "branch to rebase onto (default main)")
	rootCmd.AddCommand(rebaseCmd)
}

func newRebaseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rebase",
		Short: "rebase current branch onto origin/<branch>",
		RunE: func(cmd *cobra.Command, args []string) error {
			return git.RebaseOnto(rebaseOnto)
		},
	}
}
