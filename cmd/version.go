package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/git"
	"github.com/yuan-shuo/helm-gitops/pkg/helm"
)

var (
	level     string
	tagMode   string
	tagSuffix string
)

func init() {
	versionCmd := newVersionCmd()
	// main -> direct bump on main + tag (no PR) | or | pr -> commit push pr ci-auto-tag
	versionCmd.Flags().StringVarP(&tagMode, "mode", "m", "", "tag mode: main|pr")
	// versionCmd.Flags().StringVar(&bumpLevel, "bump", "", "bump level: patch|minor|major (required)")
	versionCmd.Flags().StringVarP(&level, "level", "l", "", "bump level: patch|minor|major|no")
	versionCmd.Flags().StringVarP(&tagSuffix, "suffix", "s", "", "tag suffix")
	// _ = versionCmd.MarkFlagRequired("bump")
	rootCmd.AddCommand(versionCmd)
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print or bump Chart.yaml version",
		Example: `helm gitops version                # 显示远程分支 Chart 最新版本
helm gitops version --bump patch   # 创建 release 分支并提交 PR`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 确保能获得当前真正的版本号
			if err := git.GoToMainAndPullLatest(); err != nil {
				return err
			}
			cur_version, err := helm.GetVersion()
			if err != nil {
				return err
			}
			// 未指定 bump 级别或未指定 tag 模式时, 仅返回最新版本
			if level == "" || tagMode == "" {
				fmt.Println(cur_version)
				return nil
			}
			// 检查version-tag模式
			switch tagMode {
			case "pr":
				// 使用pr-ci自动tag(需要actions)
				return helm.BumpWithPushAndPR(cur_version, level, PRmarkText, tagSuffix)
			case "main":
				// 在主分支上直接commit tag 然后同时推送分支和tag
				return helm.BumpDirectlyOnDefaultBranch(cur_version, level, PRmarkText, tagSuffix)
			}

			return nil
		},
	}
}
