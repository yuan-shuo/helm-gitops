package scaffold

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

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

func writeEnvSkel(root string, chartValues string, remoteChartUrl string, chartTag string) error {
	// 必写：gitignore
	if err := utils.WriteFile(filepath.Join(root, ".gitignore"), gitIgnore, 0644); err != nil {
		return err
	}
	envDirList := []string{"dev", "test", "staging", "prod"}
	for _, env := range envDirList {
		// 必写：values.yaml
		valuesContent := fmt.Sprintf("# %s/values.yaml\n\n", env) + chartValues
		if err := utils.WriteFile(filepath.Join(root, env, "values.yaml"), valuesContent, 0644); err != nil {
			return err
		}
		// 必写：kustomization.yaml
		kustomizationContent := strings.ReplaceAll(envKustomizationYAML, "{{ENV}}", env)
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

	}

	// 必写：readme.md
	envReadmeContent := strings.ReplaceAll(envReadme, "{{REMOTE_HELM_CHART_REPO}}", remoteChartUrl)
	if err := utils.WriteFile(filepath.Join(root, "README.md"), envReadmeContent, 0644); err != nil {
		return err
	}

	return nil
}
