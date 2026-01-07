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

// 把 embed 内容写进 chart
func writeSkel(root string) error {
	if err := utils.WriteFile(filepath.Join(root, ".gitignore"), gitIgnore, 0644); err != nil {
		return err
	}
	return utils.WriteFile(filepath.Join(root, ".github", "workflows", "ci-test.yaml"), ciTestYAML, 0644)
}
