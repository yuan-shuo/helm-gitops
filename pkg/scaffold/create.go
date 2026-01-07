package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yuan-shuo/helm-gitops/pkg/git"
)

func Create(name string, withActions bool) error {
	// 1. helm create
	if err := execCommand("helm", "create", name); err != nil {
		return fmt.Errorf("helm create failed: %w", err)
	}
	root := filepath.Join(".", name)

	// 2. 写骨架
	if err := writeSkel(root, withActions); err != nil {
		return err
	}

	// 3. git init
	if err := git.Init(root); err != nil {
		fmt.Println("warning: git init failed:", err)
	} else {
		fmt.Printf("✅  Chart %q created with GitOps scaffold & initial commit.\n", name)
	}
	return nil
}

func execCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
