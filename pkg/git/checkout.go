// Checkout 自动切换/创建分支
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func Checkout(name string, syncMain bool) error {

	// 若用户想同步主分支
	if syncMain {
		// 1. 检测主分支
		mainBranch, err := DetectMain()
		if err != nil {
			return err
		}
		// // 2. 确认是否同步主分支
		// ok, err := confirm(fmt.Sprintf("Pull latest %s before creating branch %q? (y/N): ", mainBranch, name))
		// if err != nil {
		// 	return err
		// }
		//
		// if ok {
		fmt.Printf("switching to latest %s...(switch to main and pull main:latest)\n", mainBranch)
		// // 1. 跳到主分支
		// if err := run("", "git", "switch", mainBranch); err != nil {
		// 	return err
		// }
		// // 2. 拉最新（fast-forward only，拒绝 merge）
		// if err := run("", "git", "pull", "--ff-only", "origin", mainBranch); err != nil {
		// 	return err
		// }
		if err := GoToMainAndPullLatest(); err != nil {
			return err
		}
		// 3. 从最新主分支创建新分支
		if err := SwitchtoBranchByAutoCreate(name); err != nil {
			return err
		}
		// fmt.Printf("creating %q from latest %s...\n", name, mainBranch)
		// return run("", "git", "switch", "-c", name)
		// }
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

func ForceDeleteBranch(name string) error {

	mainBranch, err := DetectMain()
	if err != nil {
		return err
	}
	// 先切换回主分支
	if err := SwitchtoBranchByAutoCreate(mainBranch); err != nil {
		return err
	}
	// 确认分支是否存在
	exists, err := branchExists(name)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("branch %q does not exist", name)
	}
	// 2. 删除分支
	fmt.Printf("deleting %q branch\n", name)
	return run("", "git", "branch", "-D", name)
}

// func DeleteBranch(name string) error {
// 	// 1. 确认是否删除
// 	ok, err := confirm(fmt.Sprintf("Delete branch %q? (y/N): ", name))
// 	if err != nil {
// 		return err
// 	}
// 	//
// 	if !ok {
// 		return fmt.Errorf("delete branch %q cancelled", name)
// 	}

// 	mainBranch, err := DetectMain()
// 	if err != nil {
// 		return err
// 	}
// 	// 先切换回主分支
// 	if err := SwitchtoBranchByAutoCreate(mainBranch); err != nil {
// 		return err
// 	}
// 	// 确认分支是否存在
// 	exists, err := branchExists(name)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("branch %q does not exist", name)
// 	}
// 	// 2. 删除分支
// 	fmt.Printf("deleting %q branch\n", name)
// 	return run("", "git", "branch", "-D", name)
// }

func SwitchtoBranchByAutoCreate(name string) error {
	// 1. 先尝试直接切换（本地已存在）
	if err := run("", "git", "switch", name); err == nil {
		return nil // 成功切换，直接返回
	}
	// 2. 本地没有，再走创建并切换
	fmt.Printf("creating %q branch\n", name)
	return run("", "git", "switch", "-c", name)
}

// GoToMainAndPullLatest 跳到主分支并拉取主分支最新代码
func GoToMainAndPullLatest() error {
	// 1. 跳到主分支
	mainBranch, err := DetectMain()
	if err != nil {
		return err
	}
	if err := run("", "git", "switch", mainBranch); err != nil {
		return err
	}
	// 2. 拉最新（fast-forward only，拒绝 merge）
	if err := run("", "git", "pull", "--ff-only", "origin", mainBranch); err != nil {
		return err
	}
	return nil
}

// DetectMain 检测主分支名称并返回（默认 main 或 master）
func DetectMain() (string, error) {
	for _, b := range []string{"main", "master"} {
		if err := run("", "git", "rev-parse", "--verify", "origin/"+b); err == nil {
			return b, nil
		}
	}
	return "", fmt.Errorf("no main/master branch found on origin")
}

// 失败的函数, 会偶尔出现卡死情况: 在git push后进行confirm会卡死.26-01-09 | go1.25.5
// func confirm(prompt string) (bool, error) {
// 	fmt.Print(prompt)
// 	var ans string
// 	_, err := fmt.Scanln(&ans)
// 	if err != nil && err.Error() != "unexpected newline" {
// 		return false, err
// 	}
// 	return strings.ToLower(ans) == "y", nil
// }

func branchExists(name string) (bool, error) {
	out, err := exec.Command("git", "branch", "--list", name).Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) != "", nil
}

// func isProtected(b string) bool {
// 	protected := []string{"master", "main"}
// 	for _, p := range protected {
// 		if b == p {
// 			return true
// 		}
// 	}
// 	return false
// }
