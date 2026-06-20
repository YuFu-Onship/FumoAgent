package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 设置结构 ------------------------------------------------
type Config struct {
	CurrentVoice string `json:"CurrentVoice"`

	Voice_Preset []string `json:"Voice_Preset"`
	Voice_Enable string   `json:"Voice_Enable"`

	CurrentModel string     `json:"CurrentModel"`
	RawModel     [][]string `json:"Models"`
	CurrentChar  string     `json:"CurrentChar"`

	Language string `json:"Language"`
	Color    string `json:"Color"`
	Darkmode string `json:"Darkmode"`

	Live2DCurModel string `json:"Live2DCurModel"`
}

// 语音设置参数 ----------------------------------------------
type Preset_Aquestalk struct {
	Name       string // 标题
	Default    bool   // 是否为默认
	MonoTone   bool   // 1 	机械读
	Type       string // 3 	声音类型
	Engine     string // 2 	引擎
	Speed      int    // 4	语速	默认:100 	50-300
	Volume     int    // 5 	音量	默认:100	0-300
	Pitch      int    // 6	音高	默认:100	20-200
	Accent     int    // 7	顿音	默认:100	0-200
	Quality    int    // 8 	音质	默认:100	0-200
	Intonation int    // 9	夹音	默认:100	50-200
	Note       string // 10	备注
}
type Preset_OpenJTalk struct {
	Name    string
	Default bool
	Speaker string
	Speed   float64 // 语速 1, 0.5--2
	Volume  float64 // 音量 0, -30 -- 30
	Pitch   float64 // 音高 0, -30 -- 30
}
type Preset_VoiceVox struct {
	Name       string
	Default    bool
	Speaker    string
	Speed      float64 // 语速 1 0.50--2.00
	Pitch      float64 // 音高 0 -0.15--0.15
	Intonation float64 // 顿音 1 0.00--2.00
	Volume     float64 // 音量 1 0.00--2.00
}
type Preset_Teto struct {
	Name    string
	Default bool
}

// 文件锁 ----------------------------------------------------------
var configLock sync.RWMutex
var historyLock sync.RWMutex
var errorLogMutex sync.Mutex

// 文件路径
var RootPath string = GetRootPath()
var configPath string = GetConfigPath()
var Path_PinyinTable string = "data/pinyin_table.json"
var Path_HistoryFolder string = "data/history"
var Path_HistoryMessage string = "data/history/history.csv"
var Path_AquestalkEXE string = "apps/aquestalkplayer/AquesTalkPlayer.exe"
var Path_AquestalkPreset string = "apps/aquestalkplayer/AquesTalkPlayer.preset"
var Path_ErrorFolder string = "error"

// 获取根目录路径
func GetRootPath() string {
	exePath, err := os.Executable()
	if err != nil {
		panic("无法获取程序运行路径: " + err.Error())
	}
	if strings.Contains(exePath, "go-build") || strings.Contains(exePath, "Temp") {
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			panic("无法获取源码路径")
		}
		currentDir := filepath.Dir(filepath.Dir(filepath.Dir(file)))
		return currentDir
	}
	exeDir := filepath.Dir(exePath)
	return exeDir
}

// 获取设置文件.json的路径
func GetConfigPath() string {
	return filepath.Join(RootPath, "config.json")
}

// 设置文件操作 ---------------------------------------------------------
// 读取config.json文件
func ReadConfigDate() *Config {
	data, err := os.ReadFile(GetConfigPath())
	if err != nil {
		fmt.Println("[read file failed]:", err)
		return nil
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("[parser json failed]:", err)
		return nil
	}
	return &config
}

// 保存config.json文件
func SaveConfig(config *Config) error {
	configLock.Lock() // 加写锁
	defer configLock.Unlock()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化失败: %v", err)
	}
	path := GetConfigPath()
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}
	return nil
}

// csv文件相关 ------------------------------------------------
// 检测某个文件夹是否存在
func CheckFileExist(path string, ismake bool) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(path, 0755)
		}
	}
}

// 错误相关
func Error_Init() {
	path := Path_ErrorFolder
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Printf("创建文件夹失败: %v\n", err)
			return
		}
	}

}

func Error_Save(errInfo error, name string) {
	if errInfo == nil {
		return
	}

	go func() {
		errorLogMutex.Lock()
		defer errorLogMutex.Unlock()

		now := time.Now()
		fileName := fmt.Sprintf("%d-%d-%d.txt", now.Year(), int(now.Month()), now.Day())
		fullPath := filepath.Join(Path_ErrorFolder, fileName)
		timeStr := now.Format("2006-1-2 15:04")
		content := fmt.Sprintf("[%s] [%s] %v\n", timeStr, name, errInfo)
		fmt.Println(content)

		// 3. 打开文件
		// os.O_APPEND: 追加模式
		// os.O_CREATE: 如果不存在则创建
		// os.O_WRONLY: 只写模式
		f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("无法打开或创建日志文件: %v\n", err)
			return
		}
		defer f.Close()

		if _, err := f.WriteString(content); err != nil {
			log.Printf("写入日志失败: %v\n", err)
		}
	}()
}
