package scaffold

import (
	"fmt"
	"strings"

	"github.com/yuan-shuo/helm-gitops/pkg/utils"
	"gopkg.in/yaml.v2"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/storage/memory"
)

type chartMeta struct {
	Name string `yaml:"name"`
}

type values struct {
	FullnameOverride string `yaml:"fullnameOverride"`
}

// 输入 tag 如 "v0.1.2"，返回该 tag 下含 kustomization.yaml 的一级子目录
// 要求必须给出有效 tag；空串或找不到都直接报错
func listKustomizeEnvDirs(repoURL string, tag string) ([]string, error) {
	repoURL = strings.TrimRight(strings.TrimSpace(repoURL), "/")
	if tag == "" {
		return nil, fmt.Errorf("tag is required and must not be empty")
	}

	// 组装明确引用：refs/tags/<tag>
	refName := plumbing.NewTagReferenceName(tag)

	r, err := gogit.Clone(memory.NewStorage(), nil, &gogit.CloneOptions{
		URL:           repoURL,
		Depth:         1,
		ReferenceName: refName,
		SingleBranch:  true,
	})
	if err != nil {
		// go-git 返回的错误里会包含 "reference not found" 等字样，可直接外抛
		return nil, fmt.Errorf("clone repo with tag %q failed: %w", tag, err)
	}

	ref, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("get HEAD after clone failed: %w", err)
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("get commit failed: %w", err)
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("get tree failed: %w", err)
	}

	var dirs []string
	for _, e := range tree.Entries {
		if e.Mode != filemode.Dir {
			continue
		}
		subTree, err := r.TreeObject(e.Hash)
		if err != nil {
			continue
		}
		for _, f := range subTree.Entries {
			if f.Name == "kustomization.yaml" || f.Name == "kustomization.yml" {
				dirs = append(dirs, e.Name)
				break
			}
		}
	}
	return dirs, nil
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
