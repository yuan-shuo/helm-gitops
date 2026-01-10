package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/yuan-shuo/helm-gitops/pkg/git"
)

func CreateChart(name string, withActions bool, initCommitMessage string, prMarkText string) error {
	// 1. helm create
	if err := execCommand("helm", "create", name); err != nil {
		return fmt.Errorf("helm create failed: %w", err)
	}
	root := filepath.Join(".", name)

	// 2. 写骨架
	if err := writeChartSkel(root, withActions, initCommitMessage, prMarkText); err != nil {
		return err
	}

	// 3. git init
	if err := git.Init(root, initCommitMessage); err != nil {
		fmt.Println("warning: git init failed:", err)
	} else {
		fmt.Printf("✅  Chart %q created with GitOps scaffold & initial commit.\n", name)
	}
	return nil
}

func CreateEnvRepo(remoteChartUrl string, chartTag string, EnvInitCommitMessage string) error {
	// 确认创建目录名
	repoName := path.Base(strings.TrimSpace(remoteChartUrl))
	root := filepath.Join(".", repoName+"-env")

	// 从 values.yaml 中获取全文
	valuesContent, err := fetchChartRepoToGetValues(remoteChartUrl, chartTag)
	if err != nil {
		return err
	}
	// 从 Chart.yaml 中获取 chart name
	chartName, err := unmarshalChartNameWithContent(remoteChartUrl, chartTag, valuesContent)
	if err != nil {
		return err
	}
	// 写env骨架
	if err := writeEnvSkel(root, valuesContent, remoteChartUrl, chartTag, chartName); err != nil {
		return err
	}

	// 3. git init
	if err := git.Init(root, EnvInitCommitMessage); err != nil {
		fmt.Println("warning: git init failed:", err)
	} else {
		fmt.Printf("✅  a env repo for %q created with GitOps scaffold & initial commit.\n", remoteChartUrl)
	}
	return nil
}

func execCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
