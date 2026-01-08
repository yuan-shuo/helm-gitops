package scaffold

import (
	_ "embed"
	"path/filepath"

	"github.com/yuan-shuo/helm-gitops/pkg/utils"
)

//go:embed skel/gitignore
var gitIgnore string

//go:embed skel/ci-test.yaml
var ciTestYAML string

//go:embed skel/auto-pr.yaml
var autoPrYAML string

// 把 embed 内容写进 chart
func writeSkel(root string, withActions bool) error {
	// 必写：gitignore
	if err := utils.WriteFile(filepath.Join(root, ".gitignore"), gitIgnore, 0644); err != nil {
		return err
	}
	// 可选：actions
	if withActions {
		// 可选：ci-test.yaml
		if err := utils.WriteFile(filepath.Join(root, ".github", "workflows", "ci-test.yaml"), ciTestYAML, 0644); err != nil {
			return err
		}
		// 可选：auto-pr.yaml
		if err := utils.WriteFile(filepath.Join(root, ".github", "workflows", "auto-pr.yaml"), autoPrYAML, 0644); err != nil {
			return err
		}
		// 可选：auto-tag.yaml
		if err := utils.WriteFile(filepath.Join(root, ".github", "workflows", "auto-tag.yaml"), autoPrYAML, 0644); err != nil {
			return err
		}
	}
	return nil
}
