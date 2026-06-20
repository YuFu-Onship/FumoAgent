package config

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// 指针类型转换 -----------------------------------------------------------------
type PointerConverter struct{}

// Bool 将 *bool 转为 string，默认返回 "false"
func (c PointerConverter) Bool(ptr *bool) string {
	if ptr == nil {
		return "false"
	}
	return strconv.FormatBool(*ptr)
}

// Int 将 *int 转为 string，如果为 nil 则返回传入的 defaultVal
func (c PointerConverter) Int(ptr *int, defaultVal string) string {
	if ptr == nil {
		return defaultVal
	}
	return strconv.Itoa(*ptr)
}

// IntOrZero 将 *int 转为 string，如果为 nil 则返回 "0"
func (c PointerConverter) IntOrZero(ptr *int) string {
	return c.Int(ptr, "0")
}

// 文件操作 -------------------------------------------------------------------
type FileOperation struct {
}

// 保存csv文件, 第一个参数不会重复
func (fo FileOperation) SaveCSV(path string, data []string) {
	if len(data) == 0 {
		return
	}

	// 开始写入
	targetID := data[0]
	var allRecords [][]string
	found := false

	// 1. 读取全部数据
	file, err := os.Open(path)
	if err == nil {
		reader := csv.NewReader(file)
		allRecords, _ = reader.ReadAll() // 读取所有行到内存
		file.Close()
	}

	// 2. 遍历并修改
	for i, record := range allRecords {
		if len(record) > 0 && record[0] == targetID {
			allRecords[i] = data // 找到则覆写
			found = true
			break
		}
	}

	// 3. 如果没找到，则追加到切片
	if !found {
		allRecords = append(allRecords, data)
	}

	// 4. 覆盖写入文件 (os.Create 会清空原文件)
	outFile, err := os.Create(path)
	if err != nil {
		log.Fatal("无法创建文件:", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	if err := writer.WriteAll(allRecords); err != nil { // 一次性写入所有行
		log.Fatal("写入失败:", err)
	}
}
