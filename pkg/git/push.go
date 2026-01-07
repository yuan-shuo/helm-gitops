package git

import "fmt"

// PushToRemote 推送指定分支到指定 remote
func PushToRemote(remote, branch string) error {
	fmt.Printf("pushing %s to %s\n", branch, remote)
	return run("", "git", "push", remote, branch)
}
