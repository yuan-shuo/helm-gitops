package git

import "fmt"

// RebaseOnto 获取最新 origin/<onto> 并 rebase 到当前分支
func RebaseOnto(onto string) error {
	fmt.Printf("fetching origin/%s...\n", onto)
	if err := run("", "git", "fetch", "origin", onto); err != nil {
		return err
	}
	fmt.Printf("rebasing current branch onto origin/%s...\n", onto)
	return run("", "git", "rebase", "origin/"+onto)
}
