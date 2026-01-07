package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/git"
	"github.com/yuan-shuo/helm-gitops/pkg/helm"
)

var (
	bumpLevel string
	doPush    bool
	commitMsg string
	noBump    bool
)

func init() {
	commitCmd := newCommitCmd()
	commitCmd.Flags().StringVarP(&commitMsg, "message", "m", "", "commit message (required)")
	_ = commitCmd.MarkFlagRequired("message")
	commitCmd.Flags().StringVarP(&bumpLevel, "bump", "b", "patch", "version bump level: patch|minor|major")
	commitCmd.Flags().BoolVar(&noBump, "no-bump", false, "skip version bump")
	commitCmd.Flags().BoolVarP(&doPush, "push", "p", false, "push after commit")
	rootCmd.AddCommand(commitCmd)
}

func newCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commit",
		Short: "bump version, git add & commit",
		Example: `helm gitops commit -m "feat: add metrics" -b minor
helm gitops commit -m "fix: foo" --push`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 0. 保护分支检测
			if cur, err := git.CurrentBranch(); err == nil && git.IsProtected(cur) {
				return git.ErrProtected(cur)
			}

			// 1. 可选 bump
			if !noBump {
				if err := helm.BumpVersion(bumpLevel); err != nil {
					return err
				}
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
