package git

import (
	"fmt"
	"os/exec"
)

// Tag 创建轻量标签 v{x}
func Tag(ver string) error {
	name := "v" + ver
	fmt.Printf("creating tag %s\n", name)
	cmd := exec.Command("git", "tag", name)
	return cmd.Run()
}

// PushTag 推送单个标签
func PushTag(ver string) error {
	name := "v" + ver
	fmt.Printf("pushing tag %s\n", name)
	cmd := exec.Command("git", "push", "origin", name)
	return cmd.Run()
}
