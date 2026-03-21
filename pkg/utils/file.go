package utils

import (
	"fmt"
	"os"
	"strings"
)

// RenderFileWithPlaceholders 渲染文件内容，替换占位符
// content: 原始文件内容
// placeholders: 占位符映射，key 为占位符，value 为替换内容
// 返回渲染后的内容
func RenderFileWithPlaceholders(content string, placeholders map[string]string) string {
	result := content
	for placeholder, value := range placeholders {
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// AppendFile 将内容追加到文件末尾
// 如果文件不存在则返回错误
func AppendFile(name, content string, perm os.FileMode) error {
	// 检查文件是否存在
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return fmt.Errorf("file %q not found", name)
	}

	// 以追加模式打开文件
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}
