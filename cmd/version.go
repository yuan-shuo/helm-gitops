package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/git"
	"github.com/yuan-shuo/helm-gitops/pkg/helm"
)

var bumpLevel string

func init() {
	versionCmd := newVersionCmd()
	versionCmd.Flags().StringVar(&bumpLevel, "bump", "", "bump level: patch|minor|major (required)")
	// _ = versionCmd.MarkFlagRequired("bump")
	rootCmd.AddCommand(versionCmd)
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print or bump Chart.yaml version",
		Example: `helm gitops version                # 显示当前版本
helm gitops version --bump patch   # 创建 release 分支并提交 PR`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 无 bump → 仅打印
			if bumpLevel == "" {
				v, err := helm.GetVersion()
				if err != nil {
					return err
				}
				fmt.Println(v)
				return nil
			}

			// 有 bump → 一键毕业流程
			return runGraduate(bumpLevel)
		},
	}
}

// runGraduate 复用现有子命令逻辑
func runGraduate(level string) error {
	oldVer, err := helm.GetVersion()
	if err != nil {
		return err
	}
	newVer := helm.BumpString(oldVer, level)

	// 1. 创建 release 分支（复用 checkout）
	releaseBranch := "release/v" + newVer
	checkoutCmd := newCheckoutCmd()
	if err := checkoutCmd.RunE(nil, []string{releaseBranch}); err != nil {
		return err
	}

	// 2. 改版本号（复用 BumpVersionAndSave）
	if _, err := helm.BumpVersionAndSave(level); err != nil {
		return err
	}

	// 3. commit + push + PR（复用 commit 命令）
	// 1.保护分支检测
	if cur, err := git.CurrentBranch(); err == nil && git.IsProtected(cur) {
		return git.ErrProtected(cur)
	}
	// 2.添加到缓存区
	if err := git.Add("."); err != nil {
		return err
	}
	// 3.提交带有PR标记的代码
	if err := git.Commit(git.AddPRMarkToCommitMsg("bump: v" + newVer)); err != nil {
		return err
	}
	// 4.语法检查
	if err := helm.Lint(); err != nil {
		return fmt.Errorf("lint check failed, push aborted: %w", err)
	}
	// 5.提交更改
	return git.PushHead()

	// commitCmd := newCommitCmd()
	// commitCmd.SetArgs([]string{
	// 	"-m", "bump: v" + newVer,
	// 	"--push",
	// 	"--pr",
	// })
	// return commitCmd.Execute()

}
