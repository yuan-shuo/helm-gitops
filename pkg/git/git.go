package git

import (
	"fmt"
	"os/exec"
)

func Add(path string) error {
	cmd := exec.Command("git", "add", path)
	cmd.Stdout, cmd.Stderr = nil, nil
	return cmd.Run()
}

func Commit(msg string) error {
	// 空工作区直接返回 nil，不提交
	if err := exec.Command("git", "diff-index", "--quiet", "HEAD", "--").Run(); err == nil {
		fmt.Println("nothing to commit, skipping")
		return nil
	}
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout, cmd.Stderr = nil, nil
	return cmd.Run()
}

func PushHead() error {
	cmd := exec.Command("git", "push", "origin", "HEAD")
	cmd.Stdout, cmd.Stderr = nil, nil
	return cmd.Run()
}

// Init 在 dir 执行 git init + 初始 commit
func Init(dir string) error {
	if err := run(dir, "git", "init"); err != nil {
		return err
	}
	if err := run(dir, "git", "add", "."); err != nil {
		return err
	}
	return run(dir, "git", "commit", "-m", "helm gitops chart init")
}

func run(dir string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	return cmd.Run()
}
