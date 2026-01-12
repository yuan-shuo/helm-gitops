package scaffold

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

// 将 url 的文件下载到指定本地路径（自动创建目录）
func downloadURL(dstPath, url string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("mkdir failed: %w", err)
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create file failed: %w", err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("copy body failed: %w", err)
	}
	return nil
}

// 基于 remote 仓库 + tag 获取 index.yaml，并返回第一个 .tgz 的 URL（相对或绝对）
func fetchFirstTgzURL(repo, tag string) (string, error) {
	repo = strings.TrimRight(strings.TrimSpace(repo), "/")
	// 构造 index.yaml 裸文件地址（GitHub/Gitee 均可）
	indexURL := fmt.Sprintf("%s/raw/%s/index.yaml", repo, tag)

	body, err := utils.GetFromUrlAndCollectBody(indexURL)
	if err != nil {
		return "", fmt.Errorf("get index.yaml failed: %w", err)
	}

	// 最小结构：只关心 entries[*][0].urls[0]
	var idx struct {
		Entries map[string][]struct {
			URLs []string `yaml:"urls"`
		} `yaml:"entries"`
	}
	if err := yaml.Unmarshal([]byte(body), &idx); err != nil {
		return "", fmt.Errorf("unmarshal index.yaml failed: %w", err)
	}

	// 取任意 chart 的第一个 URL（若有多版本/多 chart，可自行过滤）
	for _, versions := range idx.Entries {
		if len(versions) > 0 && len(versions[0].URLs) > 0 {
			return versions[0].URLs[0], nil
		}
	}
	return "", fmt.Errorf("no tgz url found in index.yaml")
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
				dirs = append(dirs, utils.NormalizeToNS(e.Name))
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
