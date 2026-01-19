package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/scaffold"
)

var (
	needLint   bool   // 是否执行测试
	renderDest string // 渲染目标目录
)

func init() {
	renderHelmChartCmd := newRenderHelmChartCmd()
	renderHelmChartCmd.Flags().BoolVarP(&needLint, "lint", "l", false, "lint Helm chart before render")
	renderHelmChartCmd.Flags().StringVarP(&renderDest, "dest", "d", "rendered", "render destination directory")
	rootCmd.AddCommand(renderHelmChartCmd)
}

func newRenderHelmChartCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "r",
		Short:   "render Helm chart to yaml file",
		Example: `helm gitops r -l`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scaffold.RenderHelmChart(needLint, renderDest)
		},
	}
}
