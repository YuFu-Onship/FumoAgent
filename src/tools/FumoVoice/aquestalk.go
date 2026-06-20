package fumovoice

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// CharToKata 处理字符->日语片假音
type CharToKata struct {
	PinyinDict map[string]string
}

func NewCharToKata(path string) *CharToKata {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	var dict map[string]string
	if err := json.Unmarshal(file, &dict); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}

	return &CharToKata{PinyinDict: dict}
}

func (c *CharToKata) GenKata(text string) string {
	var sb strings.Builder
	// 匹配中文字符的正则
	re := regexp.MustCompile(`[\p{Han}]`)
	// 拼音配置
	args := pinyin.NewArgs()

	for _, char := range text {
		charStr := string(char)
		if re.MatchString(charStr) {
			// 获取拼音
			py := pinyin.Pinyin(charStr, args)
			if len(py) > 0 && len(py[0]) > 0 {
				pinyinStr := py[0][0]
				if val, ok := c.PinyinDict[pinyinStr]; ok {
					sb.WriteString(val)
					continue
				}
			}
		}
		// 非中文或未匹配到则直接保留原样
		sb.WriteString(charStr)
	}
	return sb.String()
}

// aquestalk preset文件
type AquestalkPreset struct {
	Name       string // 0	预设名
	MonoTone   *bool  // 1 	机械读
	Engine     string // 2 	引擎
	VoiceType  string // 3 	声音类型
	Speed      *int   // 4	语速	默认:100 	50-300
	Volume     *int   // 5 	音量	默认:100	0-300
	Pitch      *int   // 6	音高	默认:100	20-200
	Accent     *int   // 7	顿音	默认:100	0-200
	Quality    *int   // 8 	音质	默认:100	0-200
	Intonation *int   // 9	夹音	默认:100	50-200
	Note       string // 10	备注
}

func buildAquestalkPreset(params *AquestalkPreset) *AquestalkPreset {
	default_MonoTone := defaultBool(params.MonoTone, false)
	default_Engine := defaultStr(params.Engine, "AquesTalk10")
	default_VoiceType := defaultStr(params.VoiceType, "F1E")
	default_Volume := defaultInt(params.Volume, 100)
	default_Note := defaultStr(params.Note, "")

	params.MonoTone = &default_MonoTone
	params.Engine = default_Engine
	params.VoiceType = default_VoiceType
	params.Volume = &default_Volume
	params.Note = default_Note
	return params
}

// 设置默认值
func defaultBool(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}
func defaultInt(v *int, def int) int {
	if v == nil {
		return def
	}
	return *v
}
func defaultStr(v string, def string) string {
	if v == "" {
		return def
	}
	return v
}

// 转化为csv格式
func PresetToCSV(p *AquestalkPreset) []string {
	pp := buildAquestalkPreset(p)
	return []string{
		pp.Name,

		// bool → "true"/"false"
		strconv.FormatBool(*pp.MonoTone),

		pp.Engine,
		pp.VoiceType,

		strconv.Itoa(*pp.Speed),
		strconv.Itoa(*pp.Volume),
		strconv.Itoa(*pp.Pitch),
		strconv.Itoa(*pp.Accent),
		strconv.Itoa(*pp.Quality),
		strconv.Itoa(*pp.Intonation),

		pp.Note,
	}
}

// 写入csv
func writePresetsToFile(filename string, presets []*AquestalkPreset) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写表头（建议保留日文，兼容原软件）
	header := []string{
		"プリセット名", "棒読み", "エンジン", "声種",
		"話速", "音量", "高さ", "アクセント",
		"声質", "音程", "メモ",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, p := range presets {
		row := PresetToCSV(p)
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// 追加写入
func AppendToCSV(root_path string, csvRow []string) {
	if len(csvRow) == 0 {
		return
	}

	presetPath := filepath.Join(root_path, `apps\aquestalkplayer\AquesTalkPlayer.preset`)
	targetID := csvRow[0]
	var allRecords [][]string
	found := false

	// 1. 读取全部数据
	file, err := os.Open(presetPath)
	if err == nil {
		reader := csv.NewReader(file)
		allRecords, _ = reader.ReadAll() // 读取所有行到内存
		file.Close()
	}

	// 2. 遍历并修改
	for i, record := range allRecords {
		if len(record) > 0 && record[0] == targetID {
			allRecords[i] = csvRow // 找到则覆写
			found = true
			break
		}
	}

	// 3. 如果没找到，则追加到切片
	if !found {
		allRecords = append(allRecords, csvRow)
	}

	// 4. 覆盖写入文件 (os.Create 会清空原文件)
	outFile, err := os.Create(presetPath)
	if err != nil {
		log.Fatal("无法创建文件:", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	if err := writer.WriteAll(allRecords); err != nil { // 一次性写入所有行
		log.Fatal("写入失败:", err)
	}
	// WriteAll 内部会自动调用 Flush()
}

// Aquestalk 类 ----------------------------------------------------------------------
type Aquestalk struct {
	Path       string
	OutputPath string
	PinyinPath string
}

func NewAquestalk(root_path string) *Aquestalk {
	return &Aquestalk{
		Path:       filepath.Join(root_path, `apps\aquestalkplayer\AquesTalkPlayer.exe`),
		OutputPath: filepath.Join(root_path, `temp\yukkuri_voice\voice.wav`),
		PinyinPath: filepath.Join(root_path, `data\pinyin_table.json`),
	}
}

// Talk 对应 Python 的 talk 方法
// 默认输出
func (a *Aquestalk) Talk(content string) {
	kata := NewCharToKata(a.PinyinPath)
	message := kata.GenKata(content)
	// fmt.Println(message)
	cmd := exec.Command(a.Path, "/T", message)
	// 如果需要捕获输出，可以使用 cmd.CombinedOutput()
	err := cmd.Run()
	if err != nil {
		fmt.Printf("执行失败: %v\n", err)
	}
}

// 预设输出
func (a *Aquestalk) TalkWithCoefficient(presetName string, content string) {
	kata := NewCharToKata(a.PinyinPath)
	message := kata.GenKata(content)
	// fmt.Println(message)
	cmd := exec.Command(a.Path, "/P", presetName, "/T", message)
	// 如果需要捕获输出，可以使用 cmd.CombinedOutput()
	err := cmd.Run()
	if err != nil {
		fmt.Printf("执行失败: %v\n", err)
	}
}
func (self *Aquestalk) ExportWAV(content string, presetName string, outputPath string) string {
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			fmt.Printf("删除旧音频文件失败: %v\n", err)
		}
	} else if !os.IsNotExist(err) {
		fmt.Printf("检查文件状态时出错: %v\n", err)
	}

	// 2. 执行命令
	cmd := exec.Command(self.Path, "/P", presetName, "/T", content, "/W", outputPath)

	if err := cmd.Run(); err != nil {
		fmt.Printf("保存音频失败: %v\n", err)
	}
	return outputPath
}

// SaveSound 对应 Python 的 save_sound 方法
func (a *Aquestalk) Export(content string) string {
	cmd := exec.Command(a.Path, "/T", content, "/W", a.OutputPath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("保存音频失败: %v\n", err)
	}
	return a.OutputPath
}

// func main() {
// 	// 初始化
// 	converter := NewCharToKata(`data\pinyin_table.json`)
// 	voice := NewAquestalk()

// 	message := "哎呀"
// 	kata := converter.GenKata(message)
// 	fmt.Println("转换结果:", kata)

// 	// 播放
// 	voice.Talk(kata)
// }

// 调用方法
// fumoVoice := fumovoice.NewAquestalk(root_path)
// fumoVoice.Talk("hello")

//----------------------------
// package main

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// )

// func main() {
// 	exePath, err := os.Executable()
// 	if err != nil {
// 		panic(err)
// 	}

// 	dir := filepath.Dir(exePath)
// 	fmt.Println("程序所在目录:", dir)
// }
// ---------------------------------
// package main
// import (
// 	"fmt"
// 	"path/filepath"
// 	"runtime"
// )

// func main() {
// 	_, file, _, ok := runtime.Caller(0)
// 	if !ok {
// 		panic("无法获取路径")
// 	}

// 	dir := filepath.Dir(file)
// 	fmt.Println("当前源码文件目录:", dir)
// }
