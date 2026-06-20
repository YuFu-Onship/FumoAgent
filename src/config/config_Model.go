package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type ModelDetail struct {
	Title   string
	Default string
	Url     string
	Model   string
	APIKey  string
}

// 得到当前的模型名
func API_MODEL_GetCurrentModelTitle() string {
	data := ReadConfigDate()
	title := data.CurrentModel
	return title
}
func API_MODEL_GetCurrentModelTitle_NOREAD(data Config) string {
	title := data.CurrentModel
	return title
}

// 设置当前的模型名
func API_MODEL_SetCurrentModelTitle(title string) {
	data := ReadConfigDate()
	data.CurrentModel = title
	SaveConfig(data)
}

// 得到所有模型的信息
func API_MODEL_GetAllModelDetail() []ModelDetail {
	data := ReadConfigDate()
	models := data.RawModel
	var modelList []ModelDetail
	for _, m := range models {
		modelList = append(modelList, ModelDetail{
			Title:   m[0],
			Default: m[1],
			Url:     m[2],
			Model:   m[3],
			APIKey:  m[4],
		})
	}
	return modelList
}

// 得到模型相关信息: data不传入
func API_MODEL_GetModelDetail(title string) ModelDetail {
	data := ReadConfigDate()
	for _, m := range data.RawModel {
		if m[0] == title {
			return ModelDetail{
				Title:   m[0],
				Default: m[1],
				Url:     m[2],
				Model:   m[3],
				APIKey:  m[4],
			}
		}
	}
	return ModelDetail{}
}

// 得到模型相关信息: data被传入
func API_MODEL_GetModelDetail_NOREAD(title string, data Config) ModelDetail {
	for _, m := range data.RawModel {
		if m[0] == title {
			return ModelDetail{
				Title:   m[0],
				Default: m[1],
				Url:     m[2],
				Model:   m[3],
				APIKey:  m[4],
			}
		}
	}
	return ModelDetail{}
}

// 删除模型
func API_MODEL_DeleteModel(title string) {
	data := ReadConfigDate()
	modelList := data.RawModel
	for i, m := range modelList {
		if title == m[0] {
			if m[1] != "true" {
				data.RawModel = append(modelList[:i], modelList[i+1:]...)
			}
		}
	}
	SaveConfig(data)
}

// 保存新模型
func API_MODEL_SaveNewModelDetail(m ModelDetail) {
	if !CHECK_MODEL_ModelDetailRepaet(m.Title) {
		API_MODEL_SaveModelDetail(m)
	}
}

// 保存某个模型的修改
func API_MODEL_SaveModelDetail(m ModelDetail) {
	s := []string{m.Title, m.Default, m.Url, m.Model, m.APIKey}
	data := ReadConfigDate()
	modelList := data.RawModel
	found := false
	for i, row := range modelList {
		if row[0] == m.Title {
			data.RawModel[i] = s
			found = true
			break
		}
	}
	if !found {
		data.RawModel = append(data.RawModel, s)
	}
	SaveConfig(data)
}

// 检测模型重复
func CHECK_MODEL_ModelDetailRepaet(title string) bool {
	repeat := false
	data := ReadConfigDate()
	ml := data.RawModel
	for _, m := range ml {
		if m[0] == title {
			return true
		}
	}
	return repeat
}

// 语言模型api ------------------------------------------------------------------------------------
type Model_Preset struct {
	Name    string
	Default bool
	Url     string
	Model   string
	Key     string
}

type Config_Model struct {
	mu sync.RWMutex
}

func (self *Config_Model) CheckInstance() bool {
	path := filepath.Join(RootPath, "/data/preset_Model.json")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.WriteFile(path, []byte("[]"), 0644)
		if err != nil {
			return false
		}
	}
	return true
}
func (self *Config_Model) CheckRepeat(data [][]string, name string) bool {
	for _, n := range data {
		if n[0] == name {
			return true
		}
	}
	return false
}

// 得到所有预设
func (self *Config_Model) Get_Preset() ([][]string, error) {
	path := filepath.Join(RootPath, "data/preset_Model.json")
	var data [][]string

	self.mu.Lock()
	defer self.mu.Unlock()

	f, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// 追加预设
func (self *Config_Model) Append_Preset(preset Model_Preset) {
	path := filepath.Join(RootPath, "/data/preset_Model.json")
	self.mu.Lock()
	defer self.mu.Unlock()

	self.CheckInstance()

	file, err := os.Open(path)
	if err != nil {
		return
	}

	// 获取数据
	var data [][]string
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)

	if self.CheckRepeat(data, preset.Name) {
		file.Close()
		return
	}
	file.Close()

	if err != nil {
		return
	}

	// 在获取到的数据中追加预设
	data = append(data, self.ConvertTo_String(preset))

	// 写入数据
	wf, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer wf.Close()
	encoder := json.NewEncoder(wf)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return
	}
}

// 修改指定ID预设
func (self *Config_Model) Change_Preset(preset Model_Preset) {
	data, _ := self.Get_Preset()
	for _, d := range data {
		if d[0] == preset.Name {
			d[2] = preset.Url
			d[3] = preset.Model
			d[4] = preset.Key
		}
	}
	self.write_data(data)
}

// 保存数据
func (self *Config_Model) write_data(data [][]string) {
	path := filepath.Join(RootPath, "/data/preset_Model.json")
	self.mu.Lock()
	defer self.mu.Unlock()

	wf, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer wf.Close()
	encoder := json.NewEncoder(wf)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return
	}
}

// 删除指定预设
func (self *Config_Model) Delete_Preset(name string) {
	path := filepath.Join(RootPath, "/data/preset_Model.json")
	new_data := [][]string{}
	data, _ := self.Get_Preset()

	// 从数据中移除
	for _, p := range data {
		if p[0] != name {
			new_data = append(new_data, p)
		}
	}

	// 写入文件
	wf, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer wf.Close()
	encoder := json.NewEncoder(wf)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(new_data)
	if err != nil {
		return
	}
}

// 将预设转化为字符串列表
func (self *Config_Model) ConvertTo_String(preset Model_Preset) []string {
	isDefault := strconv.FormatBool(preset.Default)
	data := []string{
		preset.Name,
		isDefault,
		preset.Url,
		preset.Model,
		preset.Key,
	}
	return data
}

// 将字符串转化为预设
func (self *Config_Model) ConvertTo_Preset(data []string) Model_Preset {
	var isDefault bool
	if data[1] == "true" {
		isDefault = true
	} else {
		isDefault = false
	}
	return Model_Preset{
		Name:    data[0],
		Default: isDefault,
		Url:     data[2],
		Model:   data[3],
		Key:     data[4],
	}
}

// 获取当前预设ID
func (self *Config_Model) Get_CurName() string {
	data := ReadConfigDate()
	name := data.CurrentModel
	return name
}

// 获取指定预设的字符串
func (self *Config_Model) Get_CurPreset(name string) []string {
	preset, _ := self.Get_Preset()
	for _, p := range preset {
		if name == p[0] {
			return p
		}
	}
	return []string{}
}
