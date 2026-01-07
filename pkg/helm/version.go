package helm

import (
	"fmt"
	"os"
	"regexp"
)

// GetVersion 读取 Chart.yaml 的 version 字段
func GetVersion() (string, error) {
	body, err := os.ReadFile("Chart.yaml")
	if err != nil {
		return "", fmt.Errorf("read Chart.yaml: %w", err)
	}
	re := regexp.MustCompile(`(?m)^version:\s*(.+)$`)
	m := re.FindSubmatch(body)
	if len(m) < 2 {
		return "", fmt.Errorf("no version field in Chart.yaml")
	}
	return string(m[1]), nil
}
