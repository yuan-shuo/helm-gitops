package utils

import (
	"fmt"
	"io"
	"net/http"
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

func GetFromUrlAndCollectBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("[fetch failed] GET %s -> %d", url, resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
