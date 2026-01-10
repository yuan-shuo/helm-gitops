package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/scaffold"
)

var (
	remoteChartUrl string // Helm chart 目录
	chartTag       string // 从远程 helm 仓库选择需要拉取的 chart 版本/tag
)

const EnvInitCommitMessage = "env for helm gitops chart init"

func init() {
	envCreateCmd := newEnvCreateCmd()
	envCreateCmd.Flags().StringVarP(&remoteChartUrl, "remote", "r", "", "path to the Helm chart directory")
	envCreateCmd.Flags().StringVarP(&chartTag, "tag", "t", "", "chart tag to pull")
	_ = envCreateCmd.MarkFlagRequired("remote")
	rootCmd.AddCommand(envCreateCmd)
}

func newEnvCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-env",
		Short: "create a new environment repository from an existing remote Helm chart",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 这里把 chartDir/mode 传下去即可，withActions 等变量按你实际来
			return scaffold.CreateEnvRepo(remoteChartUrl, chartTag, EnvInitCommitMessage)
		},
	}
}
