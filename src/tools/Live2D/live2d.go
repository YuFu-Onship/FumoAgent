package live2d

import (
	"os"
	"path/filepath"
	"strings"
)

type Live2D struct{}

// 解析live2d文件夹,获取到模型文件路径,"文件夹名称":["文件夹名称","model3.json名称"]
func (self *Live2D) ParamsFolder(folderPath string) map[string][]string {
	result := map[string][]string{}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return result
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subFolderName := entry.Name()
			subFolderPath := filepath.Join(folderPath, subFolderName)

			// 读取子文件夹内部的内容
			subEntries, err := os.ReadDir(subFolderPath)
			if err != nil {
				continue
			}

			// 寻找是否存在 .model3.json 文件
			for _, subEntry := range subEntries {
				if !subEntry.IsDir() && strings.HasSuffix(strings.ToLower(subEntry.Name()), ".model3.json") {
					// 判定为 Live2D 模型文件夹
					result[subFolderName] = []string{subFolderName, subEntry.Name()}
					break
				}
			}
		}
	}

	return result
}
