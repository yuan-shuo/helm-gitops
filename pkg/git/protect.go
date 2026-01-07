package git

import (
	"fmt"
	"os/exec"
	"strings"
)

var protected = []string{"master", "main"}

// CurrentBranch 返回当前分支名
func CurrentBranch() (string, error) {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// IsProtected 判断分支是否受保护
func IsProtected(branch string) bool {
	for _, p := range protected {
		if branch == p {
			return true
		}
	}
	return false
}

// ErrProtected 返回统一的保护分支错误
func ErrProtected(branch string) error {
	return fmt.Errorf("you are on protected branch %q, use 'helm gitops checkout <dev-branch>' first", branch)
}
