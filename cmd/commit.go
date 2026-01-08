package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/git"
)

var (
	doPush    bool
	commitMsg string
)

func init() {
	commitCmd := newCommitCmd()
	commitCmd.Flags().StringVarP(&commitMsg, "message", "m", "", "commit message (required)")
	_ = commitCmd.MarkFlagRequired("message")
	commitCmd.Flags().BoolVarP(&doPush, "push", "p", false, "push after commit")
	rootCmd.AddCommand(commitCmd)
}

func newCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "commit",
		Short:   "git add & commit",
		Example: `helm gitops commit -m "fix: foo" --push`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 0. 保护分支检测
			if cur, err := git.CurrentBranch(); err == nil && git.IsProtected(cur) {
				return git.ErrProtected(cur)
			}

			// 0. 先同步
			if err := git.PullRebase(); err != nil {
				return fmt.Errorf("cannot pull latest changes: %w", err)
			}

			if err := git.Add("."); err != nil {
				return err
			}
			if err := git.Commit(commitMsg); err != nil {
				return err
			}
			if doPush {
				return git.PushHead()
			}
			return nil
		},
	}
}
