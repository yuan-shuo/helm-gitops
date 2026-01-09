package git

import (
	"fmt"

	"github.com/yuan-shuo/helm-gitops/pkg/utils"
)

func Add(path string) error {
	// cmd := exec.Command("git", "add", path)
	// cmd.Stdout, cmd.Stderr = nil, nil
	// return cmd.Run()
	return run("", "git", "add", path)
}

func Commit(msg string) error {
	// 空工作区直接返回 nil，不提交
	if err := run("", "git", "diff-index", "--quiet", "HEAD", "--"); err == nil {
		fmt.Println("nothing to commit, skipping")
		return nil
	}
	// if err := exec.Command("git", "diff-index", "--quiet", "HEAD", "--").Run(); err == nil {
	// 	fmt.Println("nothing to commit, skipping")
	// 	return nil
	// }
	// cmd := exec.Command("git", "commit", "-m", msg)
	// cmd.Stdout, cmd.Stderr = nil, nil
	// return cmd.Run()
	return run("", "git", "commit", "-m", msg)
}

func AddPRMarkToCommitMsg(msg string, prMarkText string) string {
	return fmt.Sprintf("%s %s", msg, prMarkText)
}

func PushHead() error {
	// cmd := exec.Command("git", "push", "--set-upstream", "origin", "HEAD")
	// cmd.Stdout, cmd.Stderr = nil, os.Stderr
	// return cmd.Run()
	return run("", "git", "push", "--set-upstream", "origin", "HEAD")
}

// Init 在 dir 执行 git init + 初始 commit
func Init(dir string, initCommitMessage string) error {
	if err := run(dir, "git", "init"); err != nil {
		return err
	}
	if err := run(dir, "git", "add", "."); err != nil {
		return err
	}
	return run(dir, "git", "commit", "-m", initCommitMessage)
}

// git tag {{version}}, no 'v' prefix added in this function
func Tag(version string) error {
	// 已存在则直接返回成功
	if err := run("", "git", "rev-parse", "-q", "--verify", "refs/tags/"+version); err == nil {
		return nil
	}
	return run("", "git", "tag", version)
}

func run(dir string, name string, arg ...string) error {
	// cmd := exec.Command(name, arg...)
	// cmd.Dir = dir
	// cmd.Stdin, cmd.Stdout, cmd.Stderr = nil, nil, os.Stderr
	return utils.Run(dir, name, arg...)
}
