package config

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// 分散的函数 -------------------------------------------------------------------
// 初始化历史消息文件
func HistoryMsg_Init(rootPath string) {
	historyFilePath := filepath.Join(rootPath, Path_HistoryFolder)
	os.MkdirAll(historyFilePath, 0755)

	path := filepath.Join(rootPath, Path_HistoryMessage)
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			file, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}
			writer := csv.NewWriter(file)
			writer.Write([]string{"Time", "ID", "Content"})
			writer.Flush()
		}
	}
}

// 获取到历史消息文件路径
func HistoryMsg_GetFilePath() string {
	rootPath := RootPath
	historyFilePath := filepath.Join(rootPath, Path_HistoryFolder)
	os.MkdirAll(historyFilePath, 0755)

	path := filepath.Join(rootPath, Path_HistoryMessage)
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			file, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}
			writer := csv.NewWriter(file)
			writer.Write([]string{"Time", "ID", "Content"})
			writer.Flush()
			return ""
		}
	}
	return path
}

// 获取到历史文件内容
func HistoryMsg_GetContent() [][]string {
	file, err := os.Open(HistoryMsg_GetFilePath())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}
	return records
}

// 清空历史消息文件
func HistoryMsg_Clear() {
	rootPath := RootPath
	path := filepath.Join(rootPath, Path_HistoryMessage)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	writer := csv.NewWriter(file)
	writer.Write([]string{"Time", "ID", "Content"})
	writer.Flush()
}

// 添加新的对话
func HistoryMsg_Add(msg []string) error {
	rootPath := RootPath
	filePath := filepath.Join(rootPath, Path_HistoryMessage)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 2. 创建 CSV Writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 3. 写入新行
	if err := writer.Write(msg); err != nil {
		return fmt.Errorf("写入数据失败: %w", err)
	}

	return nil
}

// 类 -----------------------------------------------------------------------------------
type HistoryMsg struct {
	rootPath string
	// folderPath string
	// filePath   string
}

func New_HistoryMsg() *HistoryMsg {
	rootPath := RootPath
	return &HistoryMsg{
		rootPath: rootPath,
		// folderPath: filepath.Join(rootPath, Path_HistoryFolder),
		// filePath:   filepath.Join(rootPath, Path_HistoryMessage),
	}
}

// Init 初始化历史消息目录和文件头
func (self *HistoryMsg) Init() {
	folderPath := filepath.Join(RootPath, "/data/history")
	filePath := filepath.Join(RootPath, "/data/history/history.csv")
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Fatalf("无法创建目录: %v", err)
	}
	// 检查文件是否存在，不存在则创建并写入表头
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		self.createWithHeader()
	}
}

// createWithHeader 内部私有方法：创建文件并写入 CSV 表头
func (self *HistoryMsg) createWithHeader() {
	filePath := filepath.Join(RootPath, "/data/history/history.csv")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("无法创建历史文件: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{"Time", "ID", "Content"})
	writer.Flush()
}

// GetFilePath 获取文件路径（确保目录存在）
func (self *HistoryMsg) GetFilePath() string {
	self.Init() // 确保调用时环境已准备好
	filePath := filepath.Join(RootPath, "/data/history/history.csv")
	return filePath
}

// GetContent 获取所有历史记录
func (self *HistoryMsg) API_GetContent() [][]string {
	file, err := os.Open(self.GetFilePath())
	if err != nil {
		log.Printf("读取文件失败: %v", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("解析 CSV 失败: %v", err)
		return nil
	}
	return records
}

// Clear 清空历史（保留表头）
func (h *HistoryMsg) API_Clear() {
	h.createWithHeader()
}

// Add 添加单条对话
func (self *HistoryMsg) API_Add(msg []string) error {
	folderPath := filepath.Join(RootPath, "/data/history")
	filePath := filepath.Join(RootPath, "/data/history/history.csv")

	// 确保目录存在
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return err
	}

	// 以追加模式打开
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(msg); err != nil {
		return fmt.Errorf("写入数据失败: %w", err)
	}
	writer.Flush()
	return writer.Error()
}

type Config_Message struct{}

// 清空历史聊天记录
func (self *Config_Message) Clear_Msg() {
	filePath := filepath.Join(RootPath, "/data/history/history.csv")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("无法创建历史文件: %v", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Write([]string{"Time", "ID", "Content"})
	writer.Flush()
}

// 获取到当前的历史聊天记录
func (self *Config_Message) Get_Msg() [][]string {
	filePath := filepath.Join(RootPath, "/data/history/history.csv")
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("读取文件失败: %v", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("解析 CSV 失败: %v", err)
		return nil
	}
	return records
}
