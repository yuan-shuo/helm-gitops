// Checkout 自动切换/创建分支
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func Checkout(name string, syncMain bool) error {
	if isProtected(name) {
		return fmt.Errorf("refusing to create protected branch %q", name)
	}

	mainBranch, err := detectMain()
	if err != nil {
		return err
	}

	// 若用户想同步主分支
	if syncMain {
		ok, err := confirm(fmt.Sprintf("Pull latest %s before creating branch %q? (y/N): ", mainBranch, name))
		if err != nil {
			return err
		}
		if ok {
			fmt.Printf("pulling origin/%s...\n", mainBranch)
			if err := run("", "git", "fetch", "origin", mainBranch); err != nil {
				return err
			}
			if err := run("", "git", "pull", "origin", mainBranch); err != nil {
				return err
			}
		}
	}

	// 常规切换/创建
	exists, err := branchExists(name)
	if err != nil {
		return err
	}
	if exists {
		fmt.Printf("switching to %q\n", name)
		return run("", "git", "switch", name)
	}
	fmt.Printf("creating and switching to %q\n", name)
	return run("", "git", "switch", "-c", name)
}

/* 辅助函数 */
func detectMain() (string, error) {
	for _, b := range []string{"main", "master"} {
		if err := run("", "git", "rev-parse", "--verify", "origin/"+b); err == nil {
			return b, nil
		}
	}
	return "", fmt.Errorf("no main/master branch found on origin")
}

func confirm(prompt string) (bool, error) {
	fmt.Print(prompt)
	var ans string
	_, err := fmt.Scanln(&ans)
	if err != nil && err.Error() != "unexpected newline" {
		return false, err
	}
	return strings.ToLower(ans) == "y", nil
}

func branchExists(name string) (bool, error) {
	out, err := exec.Command("git", "branch", "--list", name).Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) != "", nil
}

func isProtected(b string) bool {
	protected := []string{"master", "main"}
	for _, p := range protected {
		if b == p {
			return true
		}
	}
	return false
}
