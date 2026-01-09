package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

// // 统一封装 exec，方便后续加日志、dry-run 等
// func Exec(dir, name string, arg ...string) error {
// 	cmd := exec.Command(name, arg...)
// 	cmd.Dir = dir
// 	return cmd.Run()
// }

func WriteFile(name, content string, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}
	return os.WriteFile(name, []byte(content), perm)
}

func Run(dir string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Stdin, cmd.Stdout, cmd.Stderr = nil, nil, os.Stderr
	return cmd.Run()
}
