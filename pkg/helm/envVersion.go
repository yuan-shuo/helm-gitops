package helm

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// 最小结构只拿 version
type kustomizeFile struct {
	HelmCharts []struct {
		Version string `yaml:"version"`
	} `yaml:"helmCharts"`
}

// PrintAllEnvVersions 遍历当前目录下所有含 kustomization.yaml 的子目录并打印版本
func PrintAllEnvVersions() error {
	return filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 只处理文件，且文件名精确匹配
		if !d.IsDir() && d.Name() == "kustomization.yaml" {
			dir := filepath.Dir(path)
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var kf kustomizeFile
			if err := yaml.Unmarshal(data, &kf); err != nil {
				fmt.Fprintf(os.Stderr, "warn: %s parse fail: %v\n", path, err)
				return nil
			}
			if len(kf.HelmCharts) == 0 {
				fmt.Fprintf(os.Stderr, "warn: %s no helmCharts\n", path)
				return nil
			}
			fmt.Printf("%s: %s\n", filepath.Base(dir), kf.HelmCharts[0].Version)
		}
		return nil
	})
}
