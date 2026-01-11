package scaffold

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yuan-shuo/helm-gitops/pkg/git"
	"github.com/yuan-shuo/helm-gitops/pkg/utils"
)

//go:embed skel/gitignore
var gitIgnore string

//go:embed skel/auto-test-pr.yaml
var autoTestPrYAML string

//go:embed skel/auto-tag.yaml
var autoTagYAML string

// 把 embed 内容写进 chart
func writeChartSkel(root string, withActions bool, initCommitMessage string, prMarkText string) error {
	// 必写：gitignore
	if err := utils.WriteFile(filepath.Join(root, ".gitignore"), gitIgnore, 0644); err != nil {
		return err
	}
	// 可选：actions
	if withActions {
		// 替换占位符
		workflowContent := strings.ReplaceAll(autoTestPrYAML, "{{INIT_COMMIT_MESSAGE}}", initCommitMessage)
		workflowContent = strings.ReplaceAll(workflowContent, "{{PR_MARK_TEXT}}", prMarkText)

		// auto-test-pr.yaml
		if err := utils.WriteFile(filepath.Join(root, ".github", "workflows", "auto-test-pr.yaml"), workflowContent, 0644); err != nil {
			return err
		}
		// auto-tag.yaml
		if err := utils.WriteFile(filepath.Join(root, ".github", "workflows", "auto-tag.yaml"), autoTagYAML, 0644); err != nil {
			return err
		}
		// 		fmt.Printf(`
		// [need action] all the action files are ready
		//   -  please go your github repo:
		//   -  Settings -> Actions -> General -> Workflow permissions:
		//   -  ** OPEN ** < Read and write permissions >
		//   -  ** OPEN ** < Allow GitHub Actions to create and approve pull requests >
		// 		`)
		fmt.Print(strings.TrimLeft(`
			[need action] all the action files are ready
			- please go your github repo:
			- Settings -> Actions -> General -> Workflow permissions:
			- ** OPEN ** < Read and write permissions >
			- ** OPEN ** < Allow GitHub Actions to create and approve pull requests >
			`, "\n"))
	}
	return nil
}

//go:embed skel/kustomization.yaml
var envKustomizationYAML string

//go:embed skel/patch.yaml
var envPatchYAML string

//go:embed skel/env-readme.md
var envReadme string

func writeEnvSkel(root string, chartValues string, remoteChartUrl string, chartTag string, chartName string, EnvInitCommitMessage string) error {
	// // 必写：gitignore
	// if err := utils.WriteFile(filepath.Join(root, ".gitignore"), gitIgnore, 0644); err != nil {
	// 	return err
	// }

	// 创建开发测试环境仓库
	envDirList := []string{"dev", "test", "staging"}
	devSuffixOfDir := "non-prod"
	devDir := root + "-" + devSuffixOfDir
	for _, env := range envDirList {
		if err := embedEnvFiles(devDir, env, chartName, chartValues, remoteChartUrl, chartTag); err != nil {
			return err
		}
	}

	// 创建生产环境仓库
	envDirList = []string{"prod"}
	prodSuffixOfDir := "prod"
	prodDir := root + "-" + prodSuffixOfDir
	for _, env := range envDirList {
		if err := embedEnvFiles(prodDir, env, chartName, chartValues, remoteChartUrl, chartTag); err != nil {
			return err
		}
	}

	// 各仓库根目录下其他文件创建
	if err := embedEnvRootFiles(devDir, gitIgnore, devSuffixOfDir, remoteChartUrl, EnvInitCommitMessage); err != nil {
		return err
	}
	if err := embedEnvRootFiles(prodDir, gitIgnore, prodSuffixOfDir, remoteChartUrl, EnvInitCommitMessage); err != nil {
		return err
	}

	// 必写：readme.md
	// envReadmeContent := strings.ReplaceAll(envReadme, "{{REMOTE_HELM_CHART_REPO}}", remoteChartUrl)
	// // 为 dev 环境创建 readme.md
	// devEnvReadmeContent := strings.ReplaceAll(envReadmeContent, "{{REPO_TYPE}}", devSuffixOfDir)
	// if err := utils.WriteFile(filepath.Join(devDir, "README.md"), devEnvReadmeContent, 0644); err != nil {
	// 	return err
	// }
	// // 为 prod 环境创建 readme.md
	// prodEnvReadmeContent := strings.ReplaceAll(envReadmeContent, "{{REPO_TYPE}}", prodSuffixOfDir)
	// if err := utils.WriteFile(filepath.Join(prodDir, "README.md"), prodEnvReadmeContent, 0644); err != nil {
	// 	return err
	// }

	return nil
}

// 为 env 环境创建文件
func embedEnvFiles(root, env, chartName, chartValues, remoteChartUrl, chartTag string) error {
	// 必写：values.yaml
	valuesContent := fmt.Sprintf("# %s/values.yaml\n\n", env) + chartValues
	if err := utils.WriteFile(filepath.Join(root, env, "values.yaml"), valuesContent, 0644); err != nil {
		return err
	}
	// 必写：kustomization.yaml
	// 先搞到一个临时的 kustomization.yaml 全文 然后渲染变量
	kustomizationContent := strings.ReplaceAll(envKustomizationYAML, "{{ENV}}", env)
	// 此部分将 repo 仓库下的 Chart.yaml 中的 name 属性渲染到 kustomization.yaml 中
	kustomizationContent = strings.ReplaceAll(kustomizationContent, "{{CHART_NAME}}", chartName)
	kustomizationContent = strings.ReplaceAll(kustomizationContent, "{{REMOTE_HELM_CHART_REPO}}", remoteChartUrl)
	kustomizationContent = strings.ReplaceAll(kustomizationContent, "{{REMOTE_HELM_CHART_TAG}}", chartTag)
	if err := utils.WriteFile(filepath.Join(root, env, "kustomization.yaml"), kustomizationContent, 0644); err != nil {
		return err
	}
	// 必写：patch.yaml
	patchContent := strings.ReplaceAll(envPatchYAML, "{{ENV}}", env)
	if err := utils.WriteFile(filepath.Join(root, env, "patch.yaml"), patchContent, 0644); err != nil {
		return err
	}

	return nil
}

// 在 env 根目录下创建文件
func embedEnvRootFiles(dir, gitIgnore, suffixOfDir, remoteChartUrl, EnvInitCommitMessage string) error {
	// 必写：gitignore
	if err := utils.WriteFile(filepath.Join(dir, ".gitignore"), gitIgnore, 0644); err != nil {
		return err
	}

	// 必写：readme.md
	envReadmeContent := strings.ReplaceAll(envReadme, "{{REMOTE_HELM_CHART_REPO}}", remoteChartUrl)
	// 为环境创建 readme.md
	envReadmeContent = strings.ReplaceAll(envReadmeContent, "{{REPO_TYPE}}", suffixOfDir)
	if err := utils.WriteFile(filepath.Join(dir, "README.md"), envReadmeContent, 0644); err != nil {
		return err
	}

	// 3. git init
	if err := git.Init(dir, suffixOfDir+EnvInitCommitMessage); err != nil {
		fmt.Println("warning: git init failed:", err)
	} else {
		fmt.Printf("✅  a env repo for %q created with GitOps scaffold & initial commit.\n", remoteChartUrl)
	}

	return nil
}
