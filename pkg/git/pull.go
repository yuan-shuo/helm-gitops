package git

import "fmt"

// PullRebase 如果当前分支有 upstream，则 git pull --rebase
func PullRebase() error {
	if _, err := CurrentBranch(); err != nil {
		return err
	}
	// 看是否有 upstream
	if err := run("", "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"); err != nil {
		// 没有 upstream，跳过
		return nil
	}
	fmt.Println("pulling latest changes (rebase)...")
	return run("", "git", "pull", "--rebase")
}
