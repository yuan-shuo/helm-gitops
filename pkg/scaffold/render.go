package scaffold

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/yuan-shuo/helm-gitops/pkg/utils"
)

type Values struct {
	ENV_REPO_NAME string
	ENV_REPO_URL  string
	ENV_REPO_TAG  string
	Envs          []string
}

func renderArgoAppSet(tmpl string, v Values) (string, error) {
	t, err := template.New("argo").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, v); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderEnv 为指定环境完成「下载-解压-渲染-构建」全流程（纯 Go 解压版）
func RenderEnv(envName string, remoteChartURL string, chartTag string, ifUseLocalCache bool, renderFileName string) error {

	// 1. 环境目录必须存在
	envDir := filepath.Join(".", envName)
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return fmt.Errorf("environment %q not found", envName)
	}

	helmCache := filepath.Join(envDir, "rendered", "helm", "helm-chart.yaml")
	if ifUseLocalCache {
		if _, err := os.Stat(helmCache); err == nil {
			// 文件存在，只跑 kustomize
			if renderFileName != "" {
				return renderKustomize(envName, "", renderFileName)
			}
			return renderKustomize(envName, "", "local-uname-render-result")
		}
	}

	// 2. 没有 kustomization.yaml 就仅使用 values 渲染 helm chart
	kustFile := filepath.Join(envDir, "kustomization.yaml")
	kustExist := false
	if _, err := os.Stat(kustFile); os.IsNotExist(err) {
		fmt.Printf("[skip warn] kustomization.yaml missing in %q, just use values to render helm chart", envName)
	} else {
		kustExist = true
	}

	valuesFile := filepath.Join(envDir, "values.yaml")
	if _, err := os.Stat(valuesFile); os.IsNotExist(err) {
		return fmt.Errorf("values.yaml missing in %q", envName)
	}

	// 3. 获取 tgz 下载地址并落盘
	tgzURL, err := fetchFirstTgzURL(remoteChartURL, chartTag)
	if err != nil {
		return err
	}
	chartName := strings.TrimSuffix(filepath.Base(tgzURL), filepath.Ext(tgzURL))

	chartsDir := filepath.Join(envDir, "charts")
	if err := os.MkdirAll(chartsDir, 0755); err != nil {
		return err
	}
	tgzPath := filepath.Join(chartsDir, filepath.Base(tgzURL))
	if err := downloadURL(tgzPath, joinAbsoluteURL(remoteChartURL, chartTag, tgzURL)); err != nil {
		return fmt.Errorf("download tgz failed: %w", err)
	}

	// 4. 纯 Go 解压（--strip-components=1）
	chartUnpacked := filepath.Join(chartsDir, "chart")
	if err := os.MkdirAll(chartUnpacked, 0755); err != nil {
		return err
	}
	if err := utils.UntarStripComponents(tgzPath, chartUnpacked, 1); err != nil {
		return fmt.Errorf("untar failed: %w", err)
	}

	// 5. helm 渲染
	renderedHelmDir := filepath.Join(envDir, "rendered", "helm")
	if err := os.MkdirAll(renderedHelmDir, 0755); err != nil {
		return err
	}
	helmOut := filepath.Join(renderedHelmDir, "helm-chart.yaml")
	helmCmd := exec.Command("helm", "template", envName,
		chartUnpacked,
		"--values", valuesFile)
	out, err := helmCmd.Output()
	if err != nil {
		return fmt.Errorf("helm template failed: %w", err)
	}
	if err := os.WriteFile(helmOut, out, 0644); err != nil {
		return err
	}
	// 清理
	_ = os.Remove(chartsDir)

	if kustExist {
		if err := renderKustomize(envDir, chartName, renderFileName); err != nil {
			return err
		}
	}

	return nil
}

// renderKustomize 渲染 kustomize 配置
func renderKustomize(envDir, chartName, renderFileName string) error {
	var kustOut string
	// 6. kustomize 构建
	renderedKustDir := filepath.Join(envDir, "rendered", "kustomize")
	if err := os.MkdirAll(renderedKustDir, 0755); err != nil {
		return err
	}
	if renderFileName != "" {
		kustOut = filepath.Join(renderedKustDir, renderFileName+".yaml")
	} else {
		kustOut = filepath.Join(renderedKustDir, chartName+".yaml")
	}
	kustCmd := exec.Command("kustomize", "build", envDir)
	kustOutBytes, err := kustCmd.CombinedOutput() // 合并 stderr
	if err != nil {
		return fmt.Errorf("kustomize build failed: %w\noutput: %s", err, string(kustOutBytes))
	}

	if err := os.WriteFile(kustOut, kustOutBytes, 0644); err != nil {
		return err
	}

	return nil
}

func joinAbsoluteURL(prefix, tag, rawURL string) string {
	if strings.HasPrefix(rawURL, "http") {
		return rawURL
	}
	prefix = strings.TrimRight(prefix, "/")
	return fmt.Sprintf("%s/raw/%s/%s", prefix, tag, rawURL)
}

// // RenderEnv 为指定环境完成「下载-解压-渲染-构建」全流程
// func RenderEnv(envName string, remoteChartURL string, chartTag string) error {
// 	// 1. 环境目录必须存在
// 	envDir := filepath.Join(".", envName)
// 	if _, err := os.Stat(envDir); os.IsNotExist(err) {
// 		return fmt.Errorf("environment %q not found", envName)
// 	}

// 	// 2. 没有 kustomization.yaml 就仅使用 values 渲染 helm chart
// 	kustFile := filepath.Join(envDir, "kustomization.yaml")
// 	kustExist := false
// 	if _, err := os.Stat(kustFile); os.IsNotExist(err) {
// 		fmt.Printf("[skip warn] kustomization.yaml missing in %q, just use values to render helm chart", envName)
// 	} else {
// 		kustExist = true
// 	}

// 	valuesFile := filepath.Join(envDir, "values.yaml")
// 	if _, err := os.Stat(valuesFile); os.IsNotExist(err) {
// 		return fmt.Errorf("values.yaml missing in %q", envName)
// 	}

// 	// 3. 获取 tgz 下载地址并落盘
// 	tgzURL, err := fetchFirstTgzURL(remoteChartURL, chartTag)
// 	if err != nil {
// 		return err
// 	}
// 	chartName := strings.TrimSuffix(filepath.Base(tgzURL), filepath.Ext(tgzURL)) // test-nor-0.1.3

// 	chartsDir := filepath.Join(envDir, "charts")
// 	if err := os.MkdirAll(chartsDir, 0755); err != nil {
// 		return err
// 	}
// 	tgzPath := filepath.Join(chartsDir, filepath.Base(tgzURL))
// 	if err := downloadURL(tgzPath, tgzURL); err != nil {
// 		return fmt.Errorf("download tgz failed: %w", err)
// 	}

// 	// 4. 解压到 charts/ 目录（剥掉顶层）
// 	if err := os.MkdirAll(filepath.Join(chartsDir, "chart"), 0755); err != nil {
// 		return err
// 	}
// 	cmd := exec.Command("tar", "-zxvf", tgzPath, "--strip-components=1", "-C", filepath.Join(chartsDir, "chart"))
// 	if err := cmd.Run(); err != nil {
// 		return fmt.Errorf("untar failed: %w", err)
// 	}

// 	// 5. helm 渲染
// 	renderedHelmDir := filepath.Join(envDir, "rendered", "helm")
// 	if err := os.MkdirAll(renderedHelmDir, 0755); err != nil {
// 		return err
// 	}
// 	helmOut := filepath.Join(renderedHelmDir, chartName+".yaml")
// 	helmCmd := exec.Command("helm", "template", envName,
// 		filepath.Join(chartsDir, "chart"),
// 		"--values", filepath.Join(envDir, "values.yaml"))
// 	out, err := helmCmd.Output()
// 	if err != nil {
// 		return fmt.Errorf("helm template failed: %w", err)
// 	}
// 	if err := os.WriteFile(helmOut, out, 0644); err != nil {
// 		return err
// 	}

// 	if kustExist {
// 		// 6. kustomize 构建
// 		renderedKustDir := filepath.Join(envDir, "rendered", "kustomize")
// 		if err := os.MkdirAll(renderedKustDir, 0755); err != nil {
// 			return err
// 		}
// 		kustOut := filepath.Join(renderedKustDir, chartName+".yaml")
// 		kustCmd := exec.Command("kustomize", "build", envDir)
// 		kustOutBytes, err := kustCmd.Output()
// 		if err != nil {
// 			return fmt.Errorf("kustomize build failed: %w", err)
// 		}
// 		if err := os.WriteFile(kustOut, kustOutBytes, 0644); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
