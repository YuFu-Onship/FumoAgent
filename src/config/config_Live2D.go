package config

import (
	"encoding/json"
	"io"
	live2d "myapp/src/tools/Live2D"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type Config_Live2D struct {
	mu sync.RWMutex
}

// 获取当前的live2d id
func (self *Config_Live2D) Get_CurModelName() string {
	name := ""
	data := ReadConfigDate()
	name = data.Live2DCurModel
	if name == "" {
		modelList := (&live2d.Live2D{}).ParamsFolder(filepath.Join(RootPath, "live2d"))
		if len(modelList) > 0 {
			for folderName := range modelList {
				name = folderName
				break
			}
			data.Live2DCurModel = name
			SaveConfig(data)
		}
	}
	return name
}

// 获取所有的可用模型
func (self *Config_Live2D) Get_ModelList() []string {
	result := []string{}
	modelList := (&live2d.Live2D{}).ParamsFolder(filepath.Join(RootPath, "live2d"))
	for i := range modelList {
		result = append(result, i)
	}
	slices.Sort(result)
	return result
}

// 添加新模型
func (self *Config_Live2D) Add_Model() {

}

// 检测zip是否合规
func (self *Config_Live2D) CheckZipForModelFile(archivePath string) bool {
	ext := strings.ToLower(filepath.Ext(archivePath))
	var checker ZipHandle
	var err error

	switch ext {
	case ".zip":
		checker, err = NewZipChecker(archivePath)
	case ".7z":
		checker, err = NewSevenZipChecker(archivePath)
	case ".rar":
		checker, err = NewRarChecker(archivePath)
	default:
		return false
	}

	if err != nil {
		return false
	}
	defer checker.Close()

	return checker.Check()
}

// 解压压缩包
func (self *Config_Live2D) Unzip(zipFilePath string, targetDir string) {
	ext := strings.ToLower(filepath.Ext(zipFilePath))
	var err error

	var zip_handle ZipHandle
	switch ext {
	case ".zip":
		zip_handle, err = NewZipChecker(zipFilePath)
	case ".7z":
		zip_handle, err = NewSevenZipChecker(zipFilePath)
	case ".rar":
		zip_handle, err = NewRarChecker(zipFilePath)
	default:
		return
	}

	if err != nil {
		return
	}

	zip_handle.Extract(targetDir)
}

// 保存当前live2d id
func (self *Config_Live2D) Save_CurName(name string) {
	data := ReadConfigDate()
	data.Live2DCurModel = name
	SaveConfig(data)
}

// 得到指定ID的参数
func (self *Config_Live2D) Get_TarCoe(id string) map[string]any {
	data := self.Read_file()
	for _, item := range data {
		if m, ok := item.(map[string]any); ok {
			if m["id"] == id {
				return m
			}
		}
	}
	new_preset := map[string]any{
		"id":    id,
		"scale": 1,
		"x":     0,
		"y":     0,
	}
	self.Save_Live2DCoe(id, 1, 0, 0)
	return new_preset
}

// 保存当前live2d 模型参数 [id,scale,x,y]
func (self *Config_Live2D) Save_Live2DCoe(id string, scale float32, x int, y int) bool {
	err := self.init_file()
	if err != nil {
		return false
	}
	list := self.Read_file()

	newData := map[string]any{
		"id":    id,
		"scale": scale,
		"x":     x,
		"y":     y,
	}

	found := false
	for i, item := range list {
		// 将 any 类型断言为 map[string]any
		if m, ok := item.(map[string]any); ok {
			if m["id"] == id {
				list[i] = newData
				found = true
				break
			}
		}
	}
	if !found {
		list = append(list, newData)
	}

	// 转化为字节流
	newBytes, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return false
	}

	filePath := filepath.Join(RootPath, "data/preset_Live2d.json")
	err = os.WriteFile(filePath, newBytes, 0666)
	if err != nil {
		return false
	}

	return true
}

// Delete_Live2DCoe 根据 id 删除指定的 live2d 模型参数
func (self *Config_Live2D) Delete_Live2DCoe(id string) bool {
	err := self.init_file()
	if err != nil {
		return false
	}

	// 1. 读取当前的所有数据
	list := self.Read_file()

	// 用于标记是否真的找到了并删除了数据
	found := false

	// 2. 遍历并过滤切片
	// 注意：在遍历中删除切片元素，最安全的方法是从后往前遍历，或者构建一个新切片
	var newList []any
	for _, item := range list {
		if m, ok := item.(map[string]any); ok {
			// 如果当前项的 id 与要删除的 id 相同，则跳过（实现删除效果）
			if m["id"] == id {
				found = true
				continue
			}
		}
		// 不是要删除的 id，保留到新切片中
		newList = append(newList, item)
	}

	// 3. 如果根本没找到这个 id，说明不需要修改文件，直接返回 true 或 false（取决于你的业务逻辑）
	if !found {
		return true // 或者返回 false，表示“未找到无需删除”
	}

	// 4. 将过滤后的新切片序列化为 JSON
	newBytes, err := json.MarshalIndent(newList, "", "    ")
	if err != nil {
		return false
	}

	// 5. 覆盖写入文件
	filePath := filepath.Join(RootPath, "data/preset_Live2d.json")
	err = os.WriteFile(filePath, newBytes, 0666)
	if err != nil {
		return false
	}

	return true
}

// 文件初始化,检测文件是否存在,以及写入信息
func (self *Config_Live2D) init_file() error {
	filePath := filepath.Join(RootPath, "data/preset_Live2d.json")
	content := "[]"

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.Size() == 0 {
		_, err := file.WriteString(content)
		if err != nil {
			return err
		}
	}
	return nil
}

// 读取设置文件
func (self *Config_Live2D) Read_file() []any {
	var list []any
	filePath := filepath.Join(RootPath, "data/preset_Live2d.json")
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return list
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return list
	}

	if len(bytes) > 0 && string(bytes) != "[]" {
		if err := json.Unmarshal(bytes, &list); err != nil {
			list = []any{}
		}
	}
	return list
}

// 得到文件夹中的JsonName
func (self *Config_Live2D) Get_Model3Name(folderPath string) string {
	info, err := os.Stat(folderPath)
	if err != nil {
		return ""
	}
	if !info.IsDir() {
		return ""
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if strings.HasSuffix(name, ".model3.json") {
			return name
		}
	}

	return ""
}
