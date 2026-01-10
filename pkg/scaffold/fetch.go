package scaffold

import (
	"fmt"
	"strings"

	"github.com/yuan-shuo/helm-gitops/pkg/utils"
	"gopkg.in/yaml.v2"
)

type chartMeta struct {
	Name string `yaml:"name"`
}

type values struct {
	FullnameOverride string `yaml:"fullnameOverride"`
}

func fetchChartRepoToGetValues(chartRepo, tag string) (string, error) {

	fileURL := fmt.Sprintf("%s/raw/%s/values.yaml",
		strings.TrimRight(strings.TrimSpace(chartRepo), "/"),
		tag)
	return utils.GetFromUrlAndCollectBody(fileURL)
}

func fetchChartRepoToGetChartYAML(chartRepo, tag string) (string, error) {

	fileURL := fmt.Sprintf("%s/raw/%s/Chart.yaml",
		strings.TrimRight(strings.TrimSpace(chartRepo), "/"),
		tag)
	return utils.GetFromUrlAndCollectBody(fileURL)
}

func unmarshalChartNameWithContent(chartRepo, tag, valueContent string) (string, error) {
	body, err := fetchChartRepoToGetChartYAML(chartRepo, tag)
	if err != nil {
		return "", err
	}

	// 获取 Chart.yaml 中的 name 属性
	var cm chartMeta
	if err := yaml.Unmarshal([]byte(body), &cm); err != nil {
		return "", err
	}

	// 获取 values.yaml 中的 fullnameOverride 属性, 确认 Chart 是否使用名称覆盖
	// 如果 values.yaml 中没有 fullnameOverride 属性, 则返回 Chart.yaml 中的 name 属性
	// 如果 values.yaml 中 fullnameOverride 属性有值, 则返回 values.yaml 中的 fullnameOverride 属性
	var valuesYAML values
	if err := yaml.Unmarshal([]byte(valueContent), &valuesYAML); err != nil {
		return "", err
	}
	if valuesYAML.FullnameOverride != "" {
		return valuesYAML.FullnameOverride, nil
	}

	return cm.Name, nil
}
