package fumovoice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const (
	TargetURL = "http://127.0.0.1:50021"
	SpeakerID = "1" // 默认说话人 ID，可以根据 /speakers 接口获取更多
)

func Generate(text string) string {
	// 1. 静默启动 VOICEVOX 引擎 (不弹窗)
	// 假设你的 exe 叫 voicevox_engine.exe，请替换为实际路径
	cmd := exec.Command("F:/Teto/VOICEVOX/vv-engine/run.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		// CreationFlags: syscall.CREATE_NO_WINDOW,
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("无法启动引擎: %v", err)
	}
	fmt.Println("VOICEVOX 引擎已在后台静默启动...")

	// 确保主程序退出时关闭后台进程
	defer func() {
		fmt.Println("正在关闭后台引擎...")
		_ = cmd.Process.Kill()
	}()

	// 2. 轮询等待服务器就绪
	if !waitFormServer(TargetURL, 15*time.Second) {
		log.Fatal("等待服务器启动超时，请检查路径或端口是否被占用")
	}
	fmt.Println("VOICEVOX 引擎已就绪，开始处理音频...")

	// 4. 执行语音合成核心逻辑
	audioBytes, err := generateVoice(text, SpeakerID)
	if err != nil {
		log.Fatalf("语音合成失败: %v", err)
	}

	// 5. 将字节流保存为本地音频文件
	outputFile := "output.wav"
	err = os.WriteFile(outputFile, audioBytes, 0644)
	if err != nil {
		log.Fatalf("保存音频文件失败: %v", err)
	}

	fmt.Printf("音频生成成功！已保存至: %s\n", outputFile)
	return outputFile
}

// 语音合成核心逻辑：两步请求
func generateVoice(text string, speaker string) ([]byte, error) {
	// 【第一步】：请求 /audio_query 获取合成参数 JSON
	queryURL := fmt.Sprintf("%s/audio_query?text=%s&speaker=%s", TargetURL, url.QueryEscape(text), speaker)
	respQuery, err := http.Post(queryURL, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("请求 audio_query 失败: %w", err)
	}
	defer respQuery.Body.Close()

	if respQuery.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("audio_query 状态码异常: %d", respQuery.StatusCode)
	}

	// 读取返回的配置 JSON 字节
	queryJson, err := io.ReadAll(respQuery.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 query json 失败: %w", err)
	}

	// 【第二步】：将 JSON 作为请求体，请求 /synthesis 获取音频字节流
	synthesisURL := fmt.Sprintf("%s/synthesis?speaker=%s", TargetURL, speaker)
	respAudio, err := http.Post(synthesisURL, "application/json", bytes.NewBuffer(queryJson))
	if err != nil {
		return nil, fmt.Errorf("请求 synthesis 失败: %w", err)
	}
	defer respAudio.Body.Close()

	if respAudio.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("synthesis 状态码异常: %d", respAudio.StatusCode)
	}

	// 读取最终的音频字节 (WAV 格式)
	audioBytes, err := io.ReadAll(respAudio.Body)
	if err != nil {
		return nil, fmt.Errorf("读取音频字节流失败: %w", err)
	}

	return audioBytes, nil
}

// 检查服务是否存活
func waitFormServer(apiURL string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// 访问引擎自带的 version 接口来验证是否完全启动
		resp, err := http.Get(apiURL + "/version")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

// 类 --------------------------------------------------------------------------------------

// voice speaker 语音子参数
type Style struct {
	Name string `json:"name"` // 声线风格名称，如 "ノーマル" (Normal)
	ID   int    `json:"id"`   // 核心：合成语音时真正需要传给引擎的 Style ID
}

// voice speaker 讲者参数
type Speaker struct {
	Name        string  `json:"name"`         // 角色名称，如 "四国めたん"
	SpeakerUUID string  `json:"speaker_uuid"` // 角色 UUID
	Styles      []Style `json:"styles"`       // 该角色拥有的声线风格列表
}

type VoiceVox struct {
	rootPath string
	exePath  string
	wavPath  string

	cmd       *exec.Cmd
	isRun     bool
	targetUrl string

	isStarting bool
	speakers   []Speaker
}

// 创建voicevox工作类,希望是唯一的
func New_VoiceVox(rootPath string) *VoiceVox {
	v := VoiceVox{
		rootPath:  rootPath,
		exePath:   filepath.Join(rootPath, "/apps/voicevox/vv-engine/run.exe"),
		wavPath:   filepath.Join(rootPath, "/data/temp/voiceVox.wav"),
		isRun:     false,
		targetUrl: "http://127.0.0.1:50021",
	}
	return &v
}

// 启动服务器
func (self *VoiceVox) StartServer() bool {
	if self.isRun {
		return true
	}

	if self.isStarting {
		return false
	}
	self.isStarting = true

	defer func() {
		self.isStarting = false
	}()

	if self.exePath != "" {
		self.cmd = exec.Command(self.exePath)
		self.cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}

		// 启动失败
		if err := self.cmd.Start(); err != nil {
			log.Fatalf("无法启动引擎: %v", err)
			return false
		}

		// 启动超时
		if !self.wait_Server() {
			return false
		}

		self.isRun = true
		return true
	}
	self.API_GetSpeakers()
	return false
}

// 得到角色列表
func (self *VoiceVox) API_GetSpeakers() []string {
	if !self.isRun {
		self.StartServer()
	}

	resp, err := http.Get(TargetURL + "/speakers")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var speakers []Speaker
	// 将返回的 JSON 自动解析到结构体切片中
	err = json.NewDecoder(resp.Body).Decode(&speakers)
	if err != nil {
		return nil
	}
	self.speakers = speakers

	var result []string
	for _, n := range speakers {
		fmt.Println(n)
		result = append(result, n.Name)
	}
	return result
}

func (self *VoiceVox) API_GetSpeakerDict() map[string]Speaker {
	if !self.isRun {
		self.StartServer()
	}

	resp, err := http.Get(TargetURL + "/speakers")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	result := map[string]Speaker{}
	var speakers []Speaker
	err = json.NewDecoder(resp.Body).Decode(&speakers)
	if err != nil {
		return nil
	}

	for _, s := range speakers {
		result[s.Name] = s
	}
	return result
}

// 验证服务器是否启动
func (self *VoiceVox) wait_Server() bool {
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		// 访问引擎自带的 version 接口来验证是否完全启动
		resp, err := http.Get(self.targetUrl + "/version")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

// 获得是否运行的信息
func (self *VoiceVox) API_GetRunState() bool {
	return self.isRun
}

// 检验ID
func (self *VoiceVox) API_CheckID(id int, styles []Style, playCount int, toneCount int) int {
	if playCount == toneCount {
		return id
	}
	if len(styles) == 0 {
		return 1
	}
	return styles[0].ID
}

// 获取到当前语者的索引
func (self *VoiceVox) API_GetSpeakIndex(speaker string) int {
	if self.speakers == nil {
		self.API_GetSpeakers()
	}
	for i, s := range self.speakers {
		if s.Name == speaker {
			return i
		}
	}
	return 0
}

// 生成语音文件, 返回wav路径
func (self *VoiceVox) GenerateWAV(
	text string,
	speakerID string,
	speed float64,
	pitch float64,
	intonation float64,
	volume float64,
) (string, error) {
	audioBytes, err := self.generateVoiceByte(text, speakerID, speed, pitch, intonation, volume)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(self.wavPath, audioBytes, 0644)
	if err != nil {
		return "", err
	}
	return self.wavPath, nil
}

// 生成语音音频字节
func (self *VoiceVox) generateVoiceByte(
	text string,
	speakerID string,
	speed float64,
	pitch float64,
	intonation float64,
	volume float64,
) ([]byte, error) {
	// 【第一步】：请求 /audio_query 获取合成参数 JSON
	queryURL := fmt.Sprintf("%s/audio_query?text=%s&speaker=%s", self.targetUrl, url.QueryEscape(text), speakerID)
	respQuery, err := http.Post(queryURL, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer respQuery.Body.Close()

	if respQuery.StatusCode != http.StatusOK {
		return nil, err
	}

	// 读取返回的配置 JSON 字节
	queryJson, err := io.ReadAll(respQuery.Body)
	if err != nil {
		return nil, err
	}

	// 【第二步】：将 JSON 作为请求体，请求 /synthesis 获取音频字节流
	var queryMap map[string]any
	if err := json.Unmarshal(queryJson, &queryMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query json: %w", err)
	}

	queryMap["speedScale"] = speed
	queryMap["pitchScale"] = pitch
	queryMap["intonationScale"] = intonation
	queryMap["volumeScale"] = volume

	modifiedJson, err := json.Marshal(queryMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified json: %w", err)
	}

	synthesisURL := fmt.Sprintf("%s/synthesis?speaker=%s", TargetURL, speakerID)
	respAudio, err := http.Post(synthesisURL, "application/json", bytes.NewBuffer(modifiedJson))
	if err != nil {
		return nil, err
	}
	defer respAudio.Body.Close()

	if respAudio.StatusCode != http.StatusOK {
		return nil, err
	}

	// 读取最终的音频字节 (WAV 格式)
	audioBytes, err := io.ReadAll(respAudio.Body)
	if err != nil {
		return nil, err
	}

	return audioBytes, nil
}

// 主程序
// func main() {

// 	// Generate()
// 	StartServer()
// 	defer func() {
// 		if cmd != nil && cmd.Process != nil {
// 			fmt.Println("\n正在关闭后台服务...")
// 			_ = cmd.Process.Kill()
// 		}
// 	}()

// 	// 3. 【核心修复】阻塞等待服务器完全初始化完毕
// 	fmt.Println("正在连接信道，请稍候...")
// 	if !waitFormServer(TargetURL, 15*time.Second) {
// 		log.Fatal("错误：VOICEVOX 引擎启动超时，请检查路径或端口占用情况。")
// 	}
// 	fmt.Println("信道连接成功！")

// 	// 4. 服务就绪后，现在可以安全地获取角色信息了
// 	speakers, err := GetSpeaker()
// 	if err != nil {
// 		log.Fatalf("获取角色失败: %v", err)
// 	}

// 	// 打印获取到的第一个角色的名字测试一下
// 	// if len(speakers) > 0 {
// 	// 	fmt.Printf("成功获取到角色信息！第一个角色是: %s\n", speakers[0].Name)
// 	// }

// 	for _, s := range speakers {
// 		fmt.Println(s)
// 		for _, e := range s.Styles {
// 			fmt.Println(e)
// 		}
// 	}

// 	// 5. 开启你的长期终端交互信道
// 	fmt.Println("--------------------------------------------------")
// 	fmt.Println("进入终端交互模式，请输入文本测试（或输入 exit 退出）：")

// 	scanner := bufio.NewScanner(os.Stdin)
// 	for {
// 		fmt.Print("\n请输入 > ")
// 		if !scanner.Scan() {
// 			break
// 		}
// 		text := scanner.Text()
// 		if text == "exit" {
// 			break
// 		}

// 		wavfile := Generate(text)
// 		go func() {
// 			PlayVoice(wavfile)
// 		}()

// 		// 	// 在这里可以调用你之前的 Generate 逻辑
// 		// 	fmt.Printf("收到了输入: %s（准备请求语音合成...）\n", text)
// 	}
// }
