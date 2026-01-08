package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/git"
)

var syncMain bool

func init() {
	checkoutCmd := newCheckoutCmd()
	checkoutCmd.Flags().BoolVarP(&syncMain, "sync-main", "s", false, "pull latest main/master before checkout")
	rootCmd.AddCommand(checkoutCmd)
}

func newCheckoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkout <branch-name>",
		Short: "switch to or create a new development branch",
		Example: `helm gitops checkout feature/foo
helm gitops checkout feature/foo -s    # 先同步主分支`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return git.Checkout(args[0], syncMain)
		},
	}
}
