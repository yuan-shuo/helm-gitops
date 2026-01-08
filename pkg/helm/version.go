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

// BumpString 仅计算新版本号（不读写文件）
func BumpString(oldVer, level string) string {
	parts := [3]int{}
	fmt.Sscanf(oldVer, "%d.%d.%d", &parts[0], &parts[1], &parts[2])
	switch level {
	case "major":
		parts[0]++
		parts[1], parts[2] = 0, 0
	case "minor":
		parts[1]++
		parts[2] = 0
	default: // patch
		parts[2]++
	}
	return fmt.Sprintf("%d.%d.%d", parts[0], parts[1], parts[2])
}

// BumpVersionAndSave 按 level +1 并写回 Chart.yaml
func BumpVersionAndSave(level string) (string, error) {
	old, err := GetVersion()
	if err != nil {
		return "", err
	}
	newVer := BumpString(old, level)
	body, err := os.ReadFile("Chart.yaml")
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`(?m)^version:\s*.+$`)
	newBody := re.ReplaceAllString(string(body), "version: "+newVer)
	if err := os.WriteFile("Chart.yaml", []byte(newBody), 0644); err != nil {
		return "", err
	}
	return newVer, nil
}
