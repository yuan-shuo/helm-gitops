package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yuan-shuo/helm-gitops/pkg/helm"
)

func init() {
	envVersionCmd := newEnvVersionCmd()
	rootCmd.AddCommand(envVersionCmd)
}

func newEnvVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "env-version",
		Short: "print all the environment version",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 打印所有环境的版本
			if err := helm.PrintAllEnvVersions(); err != nil {
				return err
			}
			return nil
		},
	}
}
