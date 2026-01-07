package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/git"
	"github.com/yuan-shuo/helm-gitops/pkg/helm"
)

var tagPush bool

func init() {
	tagCmd := newTagCmd()
	tagCmd.Flags().BoolVarP(&tagPush, "push", "p", false, "push tag to origin")
	rootCmd.AddCommand(tagCmd)
}

func newTagCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tag",
		Short: "create a Git tag from Chart.yaml version",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 0. 强制 lint
			if err := helm.Lint(); err != nil {
				return fmt.Errorf("lint check failed, tag aborted: %w", err)
			}
			ver, err := helm.GetVersion()
			if err != nil {
				return err
			}
			if err := git.Tag(ver); err != nil {
				return err
			}
			if tagPush {
				return git.PushTag(ver)
			}
			return nil
		},
	}
}
