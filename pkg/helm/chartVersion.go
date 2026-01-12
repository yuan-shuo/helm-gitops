package helm

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuan-shuo/helm-gitops/pkg/git"
	"github.com/yuan-shuo/helm-gitops/pkg/utils"
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
	// 不 bump 版本号
	if level == "no" {
		return oldVer
	}
	// 正常 bump 版本号逻辑
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
func BumpVersionAndSave(newVer string) (string, error) {
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

func BumpWithPushAndPR(curVersion string, level string, PRmarkText string, tagSuffix string) error {

	newVer := BumpString(curVersion, level)
	newTagByVerWithSuffix := addSuffixToTag(newVer, tagSuffix)

	// 1. 创建 release 分支（复用 checkout）
	releaseBranch := "release/v" + newTagByVerWithSuffix
	if err := git.SwitchtoBranchByAutoCreate(releaseBranch); err != nil {
		return err
	}

	// 提交带有 PR 标记的 commit 并执行 lint
	commitMsg := git.AddPRMarkToCommitMsg("pr-bump: v"+newTagByVerWithSuffix, PRmarkText)
	if err := changeChartVersionAndCommitWithLint(newVer, commitMsg, true); err != nil {
		return err
	}
	if err := git.PushHead(); err != nil {
		return err
	}
	// 任务完成提示
	fmt.Printf("created release branch %q and pushed to remote successfully\n", releaseBranch)

	// 询问是否清理
	// err = git.DeleteBranch(releaseBranch)
	// if err != nil {
	// 	return err
	// }
	// return nil
	// confirm逻辑在此处表现一般, 会卡死, 直接强制删
	return git.ForceDeleteBranch(releaseBranch)

}

func packChartAndIndex(chartDir string) error {
	// 1. 清理旧 tgz（只在 chartDir 根目录）
	entries, _ := os.ReadDir(chartDir)
	for _, e := range entries {
		if e.Type().IsRegular() && strings.HasSuffix(e.Name(), ".tgz") {
			_ = os.Remove(filepath.Join(chartDir, e.Name()))
		}
	}

	// 2. 打包 + 生成索引
	if err := execCommand("helm", "package", chartDir); err != nil {
		return fmt.Errorf("helm package failed: %w", err)
	}
	if err := execCommand("helm", "repo", "index", chartDir); err != nil {
		return fmt.Errorf("helm repo index failed: %w", err)
	}
	return nil
}

func BumpDirectlyOnDefaultBranch(curVersion string, level string, PRmarkText string, tagSuffix string) error {

	newVer := BumpString(curVersion, level)
	newTagByVerWithSuffix := addSuffixToTag(newVer, tagSuffix)

	// 提交不带 PR 标记的 commit 并执行 lint
	commitMsg := "main-bump: v" + newTagByVerWithSuffix
	if err := changeChartVersionAndCommitWithLint(newVer, commitMsg, false); err != nil {
		return err
	}

	tag := "v" + newTagByVerWithSuffix

	// 打 tag
	if err := git.Tag(tag); err != nil {
		return err
	}

	// 同时推送HEAD和tag
	return pushBranchAndTagTogether(tag)
}

func pushBranchAndTagTogether(tagName string) error {
	return utils.Run("", "git", "push", "origin", "HEAD", tagName)
}

func changeChartVersionAndCommitWithLint(newVer string, commitMsg string, protectBranch bool) error {
	// 2. 改版本号（复用 BumpVersionAndSave）
	if _, err := BumpVersionAndSave(newVer); err != nil {
		return err
	}

	// 3. commit + push + PR（复用 commit 命令）
	// 1.保护分支检测
	if protectBranch {
		if cur, err := git.CurrentBranch(); err == nil && git.IsProtected(cur) {
			return git.ErrProtected(cur)
		}
	}

	if err := packChartAndIndex("."); err != nil {
		return err
	}

	// 2.添加到缓存区
	if err := git.Add("."); err != nil {
		return err
	}
	// 3.提交带有PR标记的代码
	if err := git.Commit(commitMsg); err != nil {
		return err
	}
	// 4.语法检查
	if err := Lint(); err != nil {
		return fmt.Errorf("lint check failed, push aborted: %w", err)
	}

	return nil
}

func addSuffixToTag(tag string, suffix string) string {
	if suffix == "" {
		return tag
	}
	return tag + "-" + suffix
}
