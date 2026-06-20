package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Config_Plugin struct {
	mu sync.RWMutex
}

// 解析ini
func (self *Config_Plugin) parseINI(input string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	currentSection := "default" // 默认 section
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 过滤空行和注释
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 Section: [section]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			currentSection = strings.TrimSpace(currentSection)
			continue
		}

		// 解析 Key = Value
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			// 初始化内部 map
			if _, exists := result[currentSection]; !exists {
				result[currentSection] = make(map[string]string)
			}
			result[currentSection][key] = val
		}
	}
	return result
}

// 写入ini
func (self *Config_Plugin) writeINI(filename string, data map[string]map[string]string) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 bufio.NewWriter 提高写入效率
	writer := bufio.NewWriter(file)

	for section, keys := range data {
		// 1. 写入 Section 头部，例如: [database]
		_, err := fmt.Fprintf(writer, "[%s]\n", section)
		if err != nil {
			return err
		}

		// 2. 循环写入该 Section 下的所有 Key-Value
		for key, value := range keys {
			_, err = fmt.Fprintf(writer, "%s = %s\n", key, value)
			if err != nil {
				return err
			}
		}

		// 3. 每个 Section 之间留一个空行，美观好看
		_, err = writer.WriteString("\n")
		if err != nil {
			return err
		}
	}

	// 记得刷新缓冲区，确保所有数据写入磁盘
	return writer.Flush()
}

// 生成路径
func (self *Config_Plugin) get_Path() string {
	path := filepath.Join(RootPath, "/data/plugin.ini")
	return path
}

// 得到数据
func (self *Config_Plugin) Get_Data() map[string]map[string]string {
	data := map[string]map[string]string{}
	path := filepath.Join(RootPath, "/data/plugin.ini")
	if text, err := os.ReadFile(path); err == nil {
		data = self.parseINI(string(text))
		return data
	}
	return data
}

// 追加写入数据
func (self *Config_Plugin) Add_Data(name string, new_data map[string]string) {
	data := self.Get_Data()
	data[name] = new_data
	self.writeINI(self.get_Path(), data)
}

// 写入全部数据
func (self *Config_Plugin) Save_Data(data map[string]map[string]string) {
	self.writeINI(self.get_Path(), data)
}

// 得到插件是否启用的数据
func (self *Config_Plugin) Get_EnableData() map[string]bool {
	result := map[string]bool{}

	return result
}
