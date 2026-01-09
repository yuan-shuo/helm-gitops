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

// //go:embed skel/ci-test.yaml
// var ciTestYAML string

// //go:embed skel/auto-pr.yaml
// var autoPrYAML string

// 把 embed 内容写进 chart
func writeSkel(root string, withActions bool, initCommitMessage string, prMarkText string) error {
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
		fmt.Printf(`
[need action] all the action files are ready
  -  please go your github repo:
  -  Settings -> Actions -> General -> Workflow permissions:
  -  ** OPEN ** < Read and write permissions >
  -  ** OPEN ** < Allow GitHub Actions to create and approve pull requests >
		`)
	}
	return nil
}

// func writeToActionsDir(root string, dest string, content string) error {
// 	return utils.WriteFile(filepath.Join(root, ".github", "workflows", dest), content, 0644)
// }
