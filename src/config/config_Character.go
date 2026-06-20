package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type CharDetail struct {
	Title   string
	Default string
	Desc    string
}

// 得到当前的角色id
func API_CHARACTER_GetCurrentCharacter() string {
	data := ReadConfigDate()
	title := data.CurrentChar
	return title
}
func API_CHARACTER_GetCurrentCharacter_NOREAD(data *Config) string {
	title := data.CurrentChar
	return title
}

// 设置当前角色id
func API_CHARACTER_SetCurrentCharacter(title string) {
	data := ReadConfigDate()
	data.CurrentChar = title
	SaveConfig(data)
}
func API_CHARACTER_SetCurrentCharacter_NOREAD(data *Config, title string) {
	data.CurrentChar = title
	SaveConfig(data)
}

type CharacterConfig struct {
}

func New_CharacterConfig() *CharacterConfig {
	return &CharacterConfig{}
}

// func (self *CharacterConfig)API

type Character_Preset struct {
}

// 人设设置 ---------------------------------------------------------------------------
type Config_Character struct{}

// 获取所有预设
func (self *Config_Character) Get_Preset() [][]string {
	preset := [][]string{}
	charList := self.Get_CharList()
	for _, name := range charList {
		desc, _ := self.Get_Character(name)
		preset = append(preset, []string{name, desc})
	}
	return preset
}

// 获取到角色列表
func (self *Config_Character) Get_CharList() []string {
	folder_path := filepath.Join(RootPath, "/data/character")
	entries, err := os.ReadDir(folder_path)
	if err != nil {
		return []string{}
	}

	items := []string{}
	for _, r := range entries {
		if r.IsDir() {
			continue
		}

		file_name := r.Name()
		if filepath.Ext(file_name) == ".txt" {
			item_name := strings.TrimSuffix(file_name, ".txt")
			items = append(items, item_name)
		}
	}
	return items
}

// 检测是否重复
func (self *Config_Character) check_repeat(name string) bool {
	items := self.Get_CharList()
	return slices.Contains(items, name)
}

// 创建新的角色
func (self *Config_Character) Save_Character(name string, desc string) bool {
	if self.check_repeat(name) {
		return false
	}
	file_name := name + ".txt"
	file_path := filepath.Join(RootPath, "/data/character", file_name)
	content := fmt.Sprintf("Name: %s\r\nDescription: %s\r\n", name, desc)
	err := os.WriteFile(file_path, []byte(content), 0644)
	if err != nil {
		return false
	}
	return true
}

// 删除角色文件
func (self *Config_Character) Delete_Character(name string) bool {
	if !self.check_repeat(name) {
		return true
	}
	file_name := name + ".txt"
	file_path := filepath.Join(RootPath, "/data/character", file_name)
	err := os.Remove(file_path)
	if err != nil {
		return false
	}
	return true
}

// 读取指定角色文件
func (self *Config_Character) Get_Character(name string) (string, bool) {
	fileName := name + ".txt"
	filePath := filepath.Join(RootPath, "/data/character", fileName)

	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", false
	}
	return string(contentBytes), true
}

// 获取当前角色ID
func (self *Config_Character) Get_CurName() string {
	data := ReadConfigDate()
	return data.CurrentChar
}

// 设置当前角色ID
func (self *Config_Character) Set_CurName(name string) {
	data := ReadConfigDate()
	data.CurrentChar = name
	SaveConfig(data)
}
