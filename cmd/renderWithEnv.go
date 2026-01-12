package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/scaffold"
)

var (
	renderRemoteChartUrl string // Helm chart 目录
	renderChartTag       string // 从远程 helm 仓库选择需要拉取的 chart 版本/tag
	envName              string // 指定要渲染的目录
	ifUseLocalCache      bool   // 是否直接利用本地目录渲染
	renderFileName       string // 本地渲染文件名称
)

func init() {
	renderEnvCmd := newRenderEnvCmd()
	renderEnvCmd.Flags().StringVarP(&renderRemoteChartUrl, "remote", "r", "", "path to the Helm chart directory")
	renderEnvCmd.Flags().StringVarP(&renderChartTag, "tag", "t", "", "chart tag to pull")
	renderEnvCmd.Flags().StringVarP(&envName, "env", "e", "", "environment name to render")
	renderEnvCmd.Flags().BoolVarP(&ifUseLocalCache, "use-local-cache", "l", false, "use local cache to render")
	renderEnvCmd.Flags().StringVarP(&renderFileName, "render-file-name", "n", "", "render file name")
	// _ = renderEnvCmd.MarkFlagRequired("remote")
	// _ = renderEnvCmd.MarkFlagRequired("tag")
	_ = renderEnvCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(renderEnvCmd)
}

func newRenderEnvCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "render-env",
		Short:   "create a new environment repository from an existing remote Helm chart",
		Example: `helm gitops render-env -e prod -r https://github.com/yuan-shuo/helm-charts.git -t v0.1.0`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if renderRemoteChartUrl != "" && renderChartTag == "" {
				return fmt.Errorf("flag --tag is required when --remote is provided")
			}
			return scaffold.RenderEnv(envName, renderRemoteChartUrl, renderChartTag, ifUseLocalCache, renderFileName)
		},
	}
}
