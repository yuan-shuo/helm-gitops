package helm

import (
	"fmt"
	"os"
	"os/exec"
)

// Lint 运行 helm lint + helm unittest
func Lint() error {
	fmt.Println("running helm lint...")
	if err := execCommand("helm", "lint", "."); err != nil {
		return fmt.Errorf("helm lint failed: %w", err)
	}

	fmt.Println("running helm unittest...")
	if err := execCommand("helm", "unittest", "."); err != nil {
		if _, err2 := exec.LookPath("helm-unittest"); err2 != nil {
			return fmt.Errorf("helm unittest not found; install with: helm plugin install https://github.com/helm-unittest/helm-unittest --verify=false")
		}
		return fmt.Errorf("helm unittest failed: %w", err)
	}

	fmt.Println("local tests passed")
	return nil
}

func execCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 对应原 Makefile 的 dep-up
func DependencyUpdate() error { return nil }

// 对应原 Makefile 的 package
func Package() (string, error) { return "", nil }
