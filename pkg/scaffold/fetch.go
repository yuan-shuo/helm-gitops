package scaffold

import (
	"fmt"
	"strings"

	"github.com/yuan-shuo/helm-gitops/pkg/utils"
)

func fetchChartRepoToGetValues(chartRepo, tag string) (string, error) {

	fileURL := fmt.Sprintf("%s/raw/%s/values.yaml",
		strings.TrimRight(strings.TrimSpace(chartRepo), "/"),
		tag)
	return utils.GetFromUrlAndCollectBody(fileURL)
}
