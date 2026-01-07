// Checkout 自动切换/创建分支
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func Checkout(name string) error {
	// 0. 保护分支检测（可选）
	if isProtected(name) {
		return fmt.Errorf("refusing to create protected branch %q", name)
	}

	// 1. 看分支是否存在
	exists, err := branchExists(name)
	if err != nil {
		return err
	}

	// 2. 切换 or 创建
	if exists {
		fmt.Printf("branch %q exists, switching\n", name)
		return run("", "git", "switch", name) // dir="" 表示当前目录
	}
	fmt.Printf("creating and switching to %q\n", name)
	return run("", "git", "switch", "-c", name) // dir="" 表示当前目录
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
