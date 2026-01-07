package helm

import (
	"fmt"
	"os"
	"regexp"
)

// BumpVersion 按 level 把 Chart.yaml 的 version 字段 +1
func BumpVersion(level string) error {
	chartFile := "Chart.yaml"
	body, err := os.ReadFile(chartFile)
	if err != nil {
		return fmt.Errorf("read Chart.yaml: %w", err)
	}
	// 1. 找到旧版本号
	re := regexp.MustCompile(`(?m)^version:\s*(.+)$`)
	m := re.FindSubmatch(body)
	if len(m) < 2 {
		return fmt.Errorf("cannot find 'version:' line in Chart.yaml")
	}
	oldVer := string(m[1])

	// 2. 计算新版本
	newVer := bump(oldVer, level)

	// 3. 替换写回
	newBody := re.ReplaceAllString(string(body), "version: "+newVer)
	if err := os.WriteFile(chartFile, []byte(newBody), 0644); err != nil {
		return fmt.Errorf("write Chart.yaml: %w", err)
	}
	fmt.Printf("bump version: %s -> %s\n", oldVer, newVer)
	return nil
}

func bump(v string, level string) string {
	parts := [3]int{}
	fmt.Sscanf(v, "%d.%d.%d", &parts[0], &parts[1], &parts[2])
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
