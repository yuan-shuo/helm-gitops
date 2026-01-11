package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/scaffold"
)

var (
	remoteEnvRepoUrl string // 非生产环境远程仓库
	createArgoMode   string // 创建 argo 模式
	envRepoTag       string // 图表 tag
	argoCreateDryRun bool   //  dry run
)

func init() {
	argoCreateCmd := newArgoCreateCmd()
	argoCreateCmd.Flags().StringVarP(&remoteEnvRepoUrl, "remote", "r", "", "path to the Helm chart directory")
	argoCreateCmd.Flags().StringVarP(&envRepoTag, "tag", "t", "", "chart tag to pull")
	argoCreateCmd.Flags().StringVarP(&createArgoMode, "mode", "m", "", "create argo mode, non-prod or prod")
	argoCreateCmd.Flags().BoolVarP(&argoCreateDryRun, "dry-run", "d", false, "dry run")

	// 三项均为必填项
	_ = argoCreateCmd.MarkFlagRequired("remote")
	_ = argoCreateCmd.MarkFlagRequired("tag")
	_ = argoCreateCmd.MarkFlagRequired("mode")

	rootCmd.AddCommand(argoCreateCmd)
}

func newArgoCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "create-argo",
		Short:   "create a new argo yaml from an existing remote env repo",
		Example: `helm gitops create-argo -r https://github.com/yuan-shuo/helm-charts.git -t v0.1.0 -m non-prod --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scaffold.CreateArgoYaml(remoteEnvRepoUrl, envRepoTag, createArgoMode, argoCreateDryRun)
		},
	}
}
