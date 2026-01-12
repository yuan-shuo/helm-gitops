package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// // 统一封装 exec，方便后续加日志、dry-run 等
// func Exec(dir, name string, arg ...string) error {
// 	cmd := exec.Command(name, arg...)
// 	cmd.Dir = dir
// 	return cmd.Run()
// }

// untarStripComponents 类似 tar --strip-components=N 的纯 Go 实现
func UntarStripComponents(tgzPath, dst string, strip int) error {
	f, err := os.Open(tgzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// 去掉前 strip 级目录
		parts := strings.Split(filepath.Clean(hdr.Name), string(os.PathSeparator))
		if len(parts) <= strip {
			continue
		}
		relPath := filepath.Join(parts[strip:]...)
		target := filepath.Join(dst, relPath)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			fw, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(fw, tr); err != nil {
				fw.Close()
				return err
			}
			fw.Close()
		}
	}
	return nil
}

// NormalizeToNS 把任意字符串转换成符合 K8s 命名空间规则的字符串：
// 只能包含小写字母、数字、连字符 "-"；长度 1-63；首尾必须是字母或数字。
// 非法字符一律用 "-" 替代，并做长度截断与收尾校正。
func NormalizeToNS(s string) string {
	// 1. 统一小写
	s = strings.ToLower(s)

	// 2. 非法字符 → "-"
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	s = reg.ReplaceAllString(s, "-")

	// 3. 合并连续 "-"
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")

	// 4. 去掉首尾 "-"
	s = strings.Trim(s, "-")

	// 5. 长度 1-63
	if len(s) > 63 {
		s = s[:63]
	}
	if len(s) == 0 {
		return "default"
	}

	// 6. 确保首尾是字母或数字
	if !regexp.MustCompile(`^[a-z0-9]`).MatchString(s[:1]) {
		s = "a" + s[1:]
	}
	if !regexp.MustCompile(`[a-z0-9]$`).MatchString(s[len(s)-1:]) {
		s = s[:len(s)-1] + "0"
	}

	return s
}

// cleanRepoURL 去掉首尾空白，并删除 URL 最右侧的 "/"（如果有）
func CleanRepoURL(raw string) string {
	return strings.TrimRight(strings.TrimSpace(raw), "/")
}

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
