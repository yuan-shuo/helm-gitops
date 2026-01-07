package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "helm-gitops",
	Short: "Helm GitOps utilities",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
