package model

import (
	"encoding/json"
	"fmt"
	"myapp/src/config"
	conv "myapp/src/tools/Conv"
	fumovoice "myapp/src/tools/FumoVoice"
	live2d "myapp/src/tools/Live2D"
	netserver "myapp/src/tools/NetServer"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// 结构 ---------------------------------------------------------
type Message struct {
	msg  string
	name string
}

// 主干 ---------------------------------------------------------
type Trunk struct {
	// 设置信息
	configData config.Config
	RootPath   string

	// 内部变量
	WindowTitle string
	AiUrl       string
	AiModel     string
	AiKey       string
	ServerPort  string

	// 人设
	CharacterTitle string
	CharacterDesc  string

	// 对话消息
	MessageList [][]string

	// 音频播放器
	Player SimplePlayer

	// 语音信息
	VoicePreset any
	VoiceEnable bool

	// voiceVox相关
	VoiceVoxCmd       *exec.Cmd
	VoiceVoxRun       bool
	VoiceVox          fumovoice.VoiceVox
	VoiceVoxDict      map[string]fumovoice.Speaker
	VoiceVoxSpeakerID int
	VoiceVoxToneCount int
	VoiceVoxPlayCount int

	// 语言转换
	LangConv conv.LangConv

	// 语言
	Language      string
	LanguageTable config.LanguageTable
	languagePack  LangPack

	// 颜色
	ColorID  string
	Color    config.ColorPattle
	DarkMode bool

	// 接口信息
	Handler_AI      ModelAI
	Handler_Message ModelMsg
	Handler_Custom  config.ConfigCustom

	// 服务器
	WsServer  netserver.WsServer
	IsConnect bool

	// 插件
	Plugins   map[string]ModelPlugin
	PluginCoe map[string]map[string]string

	// live2d 参数
	Live2D_State    map[string]string
	Live2D_ModelDic map[string][]string
	Live2D_CurModel []string
	Live2D_CurName  string

	// gui api
	Handler_Gui GuiApi

	// pet_render
	PetRenderPath string
}

func New_Trunk() *Trunk {
	// 设置信息
	configData := config.ReadConfigDate()

	// LLM配置
	title := config.API_MODEL_GetCurrentModelTitle_NOREAD(*configData)
	modelDetail := config.API_MODEL_GetModelDetail_NOREAD(title, *configData)

	// 语言
	curLang := config.API_CUSTOM_LANGUAGE_GetCurrentLanguage()
	langTable := config.API_CUSTOM_LANGUAGE_GetLanguageTable(curLang)

	color_id := config.API_CUSTOM_COLOR_GetColorID()
	color_dark := config.API_CUSTOM_COLOR_GetDarkmode()
	color_pattle := config.API_CUSTOM_COLOR_GetColorPattel(color_id, color_dark)

	self := Trunk{
		WindowTitle: RandomWindowTitle(),
		configData:  *configData,

		// 模型
		AiUrl:   modelDetail.Url,
		AiModel: modelDetail.Model,
		AiKey:   modelDetail.APIKey,

		// live2d通信
		ServerPort: "9527",

		// 播放器
		Player: *NewSimplePlayer(),

		// VoiceVox
		VoiceVox: *fumovoice.New_VoiceVox(config.RootPath),

		// 客户端与服务器
		IsConnect: false,

		// 语言转换
		LangConv: *conv.New_LangConv(config.RootPath),

		// 语言
		Language:      curLang,
		LanguageTable: langTable,

		//颜色
		ColorID:  color_id,
		Color:    color_pattle,
		DarkMode: color_dark,

		// pet_render
		PetRenderPath: filepath.Join(config.RootPath, "/pet_render.exe"),
	}

	// 语音
	cv := config.Config_Voice{}
	self.VoicePreset = cv.Get_CurPreset()
	self.VoiceEnable = cv.Get_EnableState()

	// api
	cm := config.Config_Model{}
	_cur_model_preset := cm.Get_CurPreset(cm.Get_CurName())
	self.AiUrl = _cur_model_preset[2]
	self.AiModel = _cur_model_preset[3]
	self.AiKey = _cur_model_preset[4]

	// 人设
	cc := config.Config_Character{}
	cc.Get_CharList()
	characterTitle := cc.Get_CurName()
	characterDesc, _ := cc.Get_Character(characterTitle)
	self.CharacterTitle = characterTitle
	self.CharacterDesc = characterDesc

	// 插件
	self.Plugins = map[string]ModelPlugin{}
	self.PluginCoe = (&config.Config_Plugin{}).Get_Data()

	// 语言
	self.languagePack = self.LanguagePack()

	// live2d参数
	self.Live2D_State = map[string]string{}
	self.Live2D_State["cur_emotion"] = "Default"

	self.Live2D_CurName = (&config.Config_Live2D{}).Get_CurModelName()
	self.Live2D_ModelDic = (&live2d.Live2D{}).ParamsFolder(filepath.Join(config.RootPath, "/live2d"))
	self.Live2D_CurModel = self.Live2D_ModelDic[self.Live2D_CurName]

	// go func() {
	// 	self.VoiceVox = *fumovoice.New_VoiceVox(config.RootPath)
	// 	self.VoiceVox.StartServer()
	// 	speakers := (self.VoiceVox.API_GetSpeakers())
	// 	for i, n := range speakers {
	// 		fmt.Println(i, n)
	// 	}
	// }()
	// self.InputTest()
	return &self
}

func (self *Trunk) ToServer_SendCommand() {
}

func (self *Trunk) ToModel_SendMessage(message string) {
	// fmt.Println(message)
}

// 语音测试
func (self *Trunk) InputTest() {
	go func() {
		for {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				fmt.Println("读取错误:", err)
				continue
			}
			self.PlayVoice(input)

		}
	}()
}

// 播放当前预设的声音
func (self *Trunk) PlayVoice(content string) {
	if !self.VoiceEnable {
		return
	}

	reg := regexp.MustCompile(`[^。？！\.?!\n]+[。？！\.?!]?`)
	sentences := reg.FindAllString(content, -1)

	if len(sentences) == 0 {
		return
	}

	go func() {
		for _, sentence := range sentences {
			trimmed := strings.TrimSpace(sentence)
			if trimmed == "" {
				continue
			}

			// 确保在播放下一句前，先释放上一句的文件/流
			self.Player.ReleaseFile()
			path := ""
			switch data := self.VoicePreset.(type) {
			case config.Preset_Aquestalk:
				path = config.SOUND_GetVoicePath_Aquestalk(trimmed, data)

			case config.Preset_OpenJTalk:
				path = config.SOUND_GetVoicePath_OpenJTalk(trimmed, data)

			case config.Preset_VoiceVox:
				self.VoiceVoxPlayCount += 1
				if !self.VoiceVox.API_GetRunState() {
					self.VoiceVox.StartServer()
					self.VoiceVoxDict = self.VoiceVox.API_GetSpeakerDict()
				}

				self.VoiceVoxSpeakerID = self.VoiceVox.API_CheckID(
					self.VoiceVoxSpeakerID,
					self.VoiceVoxDict[data.Speaker].Styles,
					self.VoiceVoxPlayCount,
					self.VoiceVoxToneCount,
				)
				path, _ = self.VoiceVox.GenerateWAV(
					trimmed,
					strconv.Itoa(self.VoiceVoxSpeakerID),
					data.Speed,
					data.Pitch,
					data.Intonation,
					data.Volume,
				)
			default:
				return
			}
			// 播放音频
			volist, _ := self.Player.ParamsVolume(path)
			self.WsServer.SendCommand("server_set_volume_list", map[string]any{"volist": volist})
			self.Player.PlayFile(path)
		}
	}()
}

// 得到AI消息回复
func (self *Trunk) Get_AIRes(text string, role string, onMsg func(role string, content string)) string {
	// 将当前的消息添加倒历史消息中
	self.Handler_Message.API_Add([]string{"0", role, text})

	//获取历史消息
	content := self.Handler_Message.API_GetContent()
	start := 0
	if len(content) > 40 {
		start = len(content) - 40
	}
	recentContent := content[start:]
	var rows []string
	for _, line := range recentContent {
		rows = append(rows, strings.Join(line, " "))
	}

	// 设置当前插件内部语言, 决定返回的是什么语言的描述
	pluginMsg := self.languagePack["PluginDesc"][self.Language]

	// 将插件依次添加进去
	for i := range self.Plugins {
		self.Plugins[i].API_SetLanguage(self.Language)
		if self.PluginCoe["plugin"][self.Plugins[i].API_GetID()] == "false" {
			continue
		}
		name := self.Plugins[i].Name()
		desc := self.Plugins[i].Desc()
		params := self.Plugins[i].Params()
		pluginMsg = pluginMsg + name + desc + params
	}

	// 构建历史消息
	hisMsg := fmt.Sprintf(
		"%s%s\n%s%s",
		self.languagePack["PromptDesc_01"][self.Language],
		strings.Join(rows, "\n"),
		self.languagePack["PromptDesc_02"][self.Language],
		text,
	)

	var res string = ""
	oriCont := pluginMsg + hisMsg
	for range 3 {
		res = self.Handler_AI.API_GetAIResponse(oriCont, self.CharacterDesc)

		// 解析并执行插件
		pms, ai_msg := self.ParamsPlugin(res)

		// 当选择性说话插件被关闭
		if self.PluginCoe["plugin"]["selective_speech"] == "false" {
			var text string
			if self.Language == "Chinese" {
				text = self.LangConv.ChineseToKatana(ai_msg)
			} else {
				text = ai_msg
			}
			self.PlayVoice(text)
		}

		// 执行保存与追加到gui
		ai_msg = strings.TrimSpace(ai_msg)
		if ai_msg != "" {
			onMsg("ai", ai_msg)
			self.Handler_Message.API_Add([]string{"0", "ai", ai_msg})
		}

		// 插件处理
		isSkip := true
		for _, pm := range pms {
			pluginReturn, ok := self.ExecutePlugin(pm)
			isSkip = !ok && isSkip
			if ok {
				oriCont += fmt.Sprintf(
					"%s%s%s%s",
					self.languagePack["PluginReturn_01"][self.Language],
					pm.Name,
					self.languagePack["PluginReturn_02"][self.Language],
					pluginReturn,
				)
			}
		}

		if isSkip {
			break
		}
	}
	return ""
}

// 执行插件功能（修改返回值：返回结果和是否为返回型插件）
func (self *Trunk) ExecutePlugin(pm PluginMeta) (string, bool) {
	plugin, err := self.Plugins[pm.Name]
	if !err {
		return "未找到该插件", false
	}
	return plugin.Execute(pm)
}

// 解析并提取插件
func (self *Trunk) ParamsPlugin(content string) ([]PluginMeta, string) {
	pm := []PluginMeta{}

	// 正则匹配: <plugin name="xxx" args={...} />
	// [^"]+ 匹配插件名，({.*?}) 非贪婪匹配 {} 内的简易 JSON 字符串
	var pluginRegex = regexp.MustCompile(`(?s)<plugin\s+name="([^"]+)"\s+args=({.*?})\s*/?>`)
	allMatches := pluginRegex.FindAllStringSubmatch(content, -1)

	if len(allMatches) == 0 {
		return pm, strings.TrimSpace(content)
	}

	var appendix string
	for _, match := range allMatches {
		name := match[1]
		jsonStr := match[2]

		// 直接反序列化为 map[string]string
		var args map[string]any
		err := json.Unmarshal([]byte(jsonStr), &args)
		if err != nil {
			// 如果大模型生成的 JSON 格式有严重错误导致解析失败，跳过该插件
			continue
		}

		meta := PluginMeta{
			Name: name,
			Args: args,
		}
		pm = append(pm, meta)

		// 特殊插件参数保留逻辑
		switch meta.Name {
		case "selective_speech":
			// 此时 Args 的 Value 已经是 string，可以直接安全地判断和取值

			var text string
			if val, ok := meta.Args["text"]; ok {
				text = val.(string)
			}

			appendix = appendix + text
		default:
		}
	}

	// 统一擦除文本中所有的 <plugin ... /> 标签
	cleanContent := pluginRegex.ReplaceAllString(content, "")

	// 追加需要保留的内容
	if appendix != "" {
		cleanContent = cleanContent + appendix
	}

	return pm, strings.TrimSpace(cleanContent)
}

// 语言
func (self *Trunk) LanguagePack() LangPack {
	return LangPack{
		"PluginDesc": LangZone{
			"Chinese":  "调用插件时，必须使用单标签格式：<plugin name=\"插件名\" args={...} />。请确保 args 属性是一个合法的 JSON 对象。你可以调用的插件功能有：",
			"English":  "When invoking a plugin, you must use the single-tag format: <plugin name=\"plugin_name\" args={...} />. Ensure that the args attribute is a valid JSON object. Available plugin functions:",
			"Japanese": "プラグインを呼び出す際は、必ず単一タグ形式 <plugin name=\"プラグイン名\" args={...} /> を使用してください。args 属性は有効な JSON オブジェクトである必要があります。利用可能なプラグイン機能は以下の通りです:",
		},
		"PromptDesc_01": LangZone{
			"Chinese":  "以下是历史消息记录,对方是user,你的角色是ai: ",
			"English":  "Below is the historical message log. The counterparty is 'user', and your role is 'ai': ",
			"Japanese": "以下は過去のメッセージ履歴です。相手は'user'、あなたの役割は'ai'です: ",
		},
		"PromptDesc_02": LangZone{
			"Chinese":  "当前来自user的对话是:",
			"English":  "The current conversation from user is:",
			"Japanese": "ユーザーからの現在の発言は以下の通りです:",
		},
		"PluginReturn_01": LangZone{
			"Chinese":  "你刚刚调用了插件<",
			"English":  "",
			"Japanese": "",
		},
		"PluginReturn_02": LangZone{
			"Chinese":  ">,返回结果是<",
			"English":  "",
			"Japanese": "",
		},
	}
}

// 服务端与客户端初始化
func (self *Trunk) InitServer() {
	// 发送live2d模型路径
	self.WsServer.RegisterHandler("client_get_live2d_model", func(map[string]string) {
		if len(self.Live2D_CurModel) == 0 {
			return
		}
		self.WsServer.SendCommand("server_set_live2d_model", map[string]any{
			"FolderPath": self.Live2D_CurModel[0],
			"JsonName":   self.Live2D_CurModel[1],
		})
	})

	// 发送live2d模型参数(scale,x,y)
	self.WsServer.RegisterHandler("client_get_live2d_trans", func(map[string]string) {
		if len(self.Live2D_CurModel) == 0 {
			return
		}
		cl := config.Config_Live2D{}
		coe := cl.Get_TarCoe(self.Live2D_CurName)
		s, _ := coe["scale"].(float64)
		x, _ := coe["x"].(float64)
		y, _ := coe["y"].(float64)
		self.WsServer.SendCommand("server_set_live2d_trans", map[string]any{
			"scale":  s,
			"transx": x,
			"transy": y,
		})
	})
}

// 插件参数的更新与保存
func (self *Trunk) UpdatePluginCoe(data map[string]map[string]string) {
	self.PluginCoe = data
	cp := config.Config_Plugin{}
	cp.Save_Data(data)
}

// 启动渲染器
func (self *Trunk) EnableRender() {
	if self.IsConnect {
		return
	}

	// 检查 pet_render.exe 是否已经在运行
	if isProcessRunning("pet_render.exe") {
		// 客户端已在运行，不需要再次启动
		// IsConnect 的最终状态由 WebSocket 心跳回调接管
		return
	}

	self.IsConnect = true
	go func() {
		cmd := exec.Command(self.PetRenderPath)
		cmd.CombinedOutput()
		self.IsConnect = false
	}()
}

// isProcessRunning 检查指定名称的进程是否正在运行
func isProcessRunning(name string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist", "/FI", "IMAGENAME eq "+name, "/NH")
	default:
		cmd = exec.Command("pgrep", "-x", strings.TrimSuffix(name, ".exe"))
	}
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		return strings.Contains(string(output), name)
	}
	return len(strings.TrimSpace(string(output))) > 0
}
