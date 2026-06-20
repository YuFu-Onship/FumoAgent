package Ctrl

// app控制
import (
	"myapp/src/config"
	"myapp/src/model"
	"myapp/src/plugins"
	aichat "myapp/src/tools/AIModel"
	netserver "myapp/src/tools/NetServer"
	"strconv"

	"myapp/src/ui"
)

type LangPack map[string]map[string]string

type AppCtrl struct {
	Gui       *ui.Page
	AiClient  *aichat.AIClient
	trunk     *model.Trunk
	language  string
	langtable LangPack
}

func New_AppCtrl(trunk *model.Trunk) *AppCtrl {
	self := AppCtrl{trunk: trunk}
	self.LoadModel()
	self.LoadClientPlugin()

	self.language = "Chinese"
	self.langtable = self.Language()
	return &self
}

// 运行
func (self *AppCtrl) Run() {
	go self.Gui.Run()
	go self.trunk.WsServer.Run()
	self.LoadPlugin()
	select {}
}

// 加载模块
func (self *AppCtrl) LoadModel() {
	inst_gui := ui.New_Page(self.trunk)
	inst_ai := aichat.New_AIClient(self.trunk)

	self.trunk.Handler_AI = aichat.New_AIClient(self.trunk)
	self.trunk.Handler_Message = config.New_HistoryMsg()
	self.trunk.Handler_Custom = *config.New_ConfigCustom()
	self.trunk.WsServer = *netserver.New_WsServer(self.trunk.ServerPort)
	self.Gui = inst_gui
	self.AiClient = inst_ai
	self.trunk.InitServer()
	self.trunk.Handler_Gui = self.Gui

	self.trunk.WsServer.OnHeart(func() { self.trunk.IsConnect = self.trunk.WsServer.API_GetAlive() })
	self.trunk.EnableRender()
}

// trunk加载gui接口
func (self *AppCtrl) LoadGuiApi() {
}

// 服务端->客户端 插件
func (self *AppCtrl) LoadPlugin() {
	current_time := plugins.New_CurrentTime(self.trunk)
	selective_speech := plugins.New_SelectiveSpeech(self.trunk)
	send_emotion := plugins.New_SendEmotion(self.trunk)
	set_model_rotate := plugins.New_ModelRotate(self.trunk)
	set_model_shake := plugins.New_ModelShake(self.trunk)
	voice_vox_tone := plugins.New_VoiceVoxTone(self.trunk)

	pluginList := []struct {
		key       string
		modelPlug model.ModelPlugin
		guiCtrl   ui.PluginCtrl
	}{
		{current_time.API_GetID(), current_time, current_time},
		{selective_speech.API_GetID(), selective_speech, selective_speech},
		{send_emotion.API_GetID(), send_emotion, send_emotion},
		{set_model_rotate.API_GetID(), set_model_rotate, set_model_rotate},
		{set_model_shake.API_GetID(), set_model_shake, set_model_shake},
		{voice_vox_tone.API_GetID(), voice_vox_tone, voice_vox_tone},
	}

	// 取出GUI中的API接口
	addPluginFunc := self.Gui.API()["add_plugin"].(func(ui.PluginCtrl))

	// 插件启动初始化
	data := self.trunk.PluginCoe
	if data["plugin"] == nil {
		data["plugin"] = map[string]string{}
	}

	// 循环
	for _, p := range pluginList {
		if data["plugin"][p.key] == "" {
			data["plugin"][p.key] = strconv.FormatBool(p.guiCtrl.API_GetEnable())
		} else {
			if value, err := strconv.ParseBool(data["plugin"][p.key]); err == nil {
				p.guiCtrl.API_SetEnable(value)
			}
		}

		if data[p.guiCtrl.API_GetID()] != nil {
			p.modelPlug.API_SetValueString(data[p.guiCtrl.API_GetID()])
		}

		self.trunk.Plugins[p.key] = p.modelPlug
		addPluginFunc(p.guiCtrl)
	}
	self.trunk.PluginCoe = data
	(&config.Config_Plugin{}).Save_Data(data)
}

// 客户端->服务端 插件
func (self *AppCtrl) LoadClientPlugin() {
	// 摸头
	self.trunk.WsServer.RegisterHandler("client_set_onpat", func(args map[string]string) {
		self.language = self.trunk.Language
		self.Gui.API()["add_msg"].(func(string, string))("system", self.langtable["headpat_desc"][self.language])
		self.trunk.Get_AIRes(self.langtable["headpat_toai"][self.language], "system", self.Gui.API()["add_msg"].(func(string, string)))
	})
	// 跳转到聊天页
	self.trunk.WsServer.RegisterHandler("client_set_chat", func(args map[string]string) {
		self.Gui.API_SetPageChat()
	})
	// 客户端关闭
	self.trunk.WsServer.RegisterHandler("client_close", func(args map[string]string) {
		self.trunk.IsConnect = false
	})
}

// 语言
func (self *AppCtrl) Language() LangPack {
	return LangPack{
		"headpat_desc": map[string]string{
			"Chinese":  "你摸了摸对方的头",
			"English":  "You patted them on the head.",
			"Japanese": "相手の頭をなでなでした。",
		},
		"headpat_toai": map[string]string{
			"Chinese":  "对方摸了摸你的头",
			"English":  "User patted you on the head.",
			"Japanese": "あなたの頭をなでなでした。",
		},
	}
}
