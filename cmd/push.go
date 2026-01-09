package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/git"
	"github.com/yuan-shuo/helm-gitops/pkg/helm"
)

// var pushRemote string

func init() {
	pushCmd := newPushCmd()
	// pushCmd.Flags().StringVar(&pushRemote, "remote", "origin", "remote name to push")
	rootCmd.AddCommand(pushCmd)
}

func newPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "push current development branch to remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cur, err := git.CurrentBranch()
			if err != nil {
				return err
			}
			if git.IsProtected(cur) {
				return git.ErrProtected(cur)
			}
			// 0. 先同步
			// if err := git.PullRebase(); err != nil {
			// 	return fmt.Errorf("cannot pull latest changes: %w", err)
			// }
			// 0. 强制 lint
			if err := helm.Lint(); err != nil {
				return fmt.Errorf("lint check failed, push aborted: %w", err)
			}
			return git.PushHead()
		},
	}
}
