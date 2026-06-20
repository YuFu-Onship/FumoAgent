package plugins

import (
	"fmt"
	"myapp/src/model"
)

type Mod_SendEmotion struct {
	trunk     *model.Trunk
	lang      string
	langTable LangPack

	emotions string
	cur_emo  string

	enable      bool
	last_enable bool

	id string
}

func New_SendEmotion(trunk *model.Trunk) *Mod_SendEmotion {
	self := Mod_SendEmotion{trunk: trunk}
	self.lang = "Chinese"
	self.langTable = self.Language()
	self.cur_emo = "Default"
	self.enable = true
	self.last_enable = self.enable

	self.id = "send_emotion"

	// 服务端接收到到该指令后就执行内部的语句
	self.trunk.WsServer.RegisterHandler("client_set_all_emo",
		func(args map[string]string) { self.emotions = args["emotion"] },
	)
	self.trunk.WsServer.RegisterHandler("client_set_cur_emo",
		func(args map[string]string) {
			self.lang = self.trunk.Language
			self.cur_emo = (args["emotion"])
			if self.cur_emo == "Default" {
				return
			}
			desc := fmt.Sprintf(self.langTable["cueEmoSysDesc"][self.lang], self.cur_emo)
			self.trunk.Handler_Gui.API()["add_msg"].(func(string, string))("system", desc)
			self.trunk.Handler_Message.API_Add([]string{"0", "system", desc})
		},
	)

	self.trunk.WsServer.OnConnect(self.get_emotions)
	return &self
}

func (self *Mod_SendEmotion) Handle() bool { return true }
func (self *Mod_SendEmotion) Execute(pm model.PluginMeta) (string, bool) {
	self.get_emotions()
	self.trunk.WsServer.SendCommand("server_set_emotion", pm.Args)
	return "", false
}

// 获取描述
func (self *Mod_SendEmotion) Name() string { return "send_emotion" }
func (self *Mod_SendEmotion) Desc() string {
	desc := ""
	desc += self.langTable["CurEmoDesc"][self.lang] + self.cur_emo + ">"
	desc += self.langTable["Desc"][self.lang] + self.emotions
	return desc
}

func (self *Mod_SendEmotion) Params() string {
	return self.langTable["Params"][self.lang]
}

// 参数设置
func (self *Mod_SendEmotion) Setting()     {}
func (self *Mod_SendEmotion) LoadSetting() {}

// 设置语言
func (self *Mod_SendEmotion) API_SetLanguage(language string) { self.lang = language }

// 向客户获取当前的可用表情
func (self *Mod_SendEmotion) get_emotions() {
	self.trunk.WsServer.SendCommand("server_get_emotions", map[string]any{})
}

// 内部语言
func (self *Mod_SendEmotion) Language() LangPack {
	lp := LangPack{
		"Desc": LangZone{
			"Chinese":  "选择发送的表情，你可以在 live2d 模型上做出的表情有：",
			"English":  "Select an emotion to send. The emotions you can display on the Live2D model are:",
			"Japanese": "送信する表情を選択します。Live2Dモデルで表現できる表情は以下の通りです：",
		},
		"Params": LangZone{
			"Chinese":  `<plugin name="send_emotion" args={"emotion":"表情"} />`,
			"English":  `<plugin name="send_emotion" args={"emotion":"emotion_name"} />`,
			"Japanese": `<plugin name="send_emotion" args={"emotion":"表情の名称"} />`,
		},
		"CurEmoDesc": LangZone{
			"Chinese":  "当前的表情是:<",
			"English":  "Current emotion is:<",
			"Japanese": "現在の表情は:<",
		},
		"cueEmoSysDesc": LangZone{
			"Chinese":  "刚刚做了 <%s> 表情",
			"English":  "Just made a <%s> expression",
			"Japanese": "さっき <%s> 表情をしました",
		},
	}
	return lp
}

// 外部调用 ---------------------------------------------------------------------
func (self *Mod_SendEmotion) Options() []map[string]any { return []map[string]any{} }

func (self *Mod_SendEmotion) API_SetValue(name string, value any)  {}
func (self *Mod_SendEmotion) API_SetValueString(map[string]string) {}

func (self *Mod_SendEmotion) API_GetValueString() map[string]string { return map[string]string{} }

func (self *Mod_SendEmotion) GUI_Language() map[string]map[string]string {
	return map[string]map[string]string{
		"Title": {
			"Chinese":  "发送表情",
			"English":  "Send Expression",
			"Japanese": "表情送信",
		},
		"Desc": {
			"Chinese":  "允许 AI 控制模型做出相应的表情。",
			"English":  "Allows the AI to control the model to make corresponding expressions.",
			"Japanese": "AIがモデルを制御して対応する表情を作れるようにします。",
		},
		"Detail": {
			"Chinese":  "能够让 AI 发送表情指令，使模型做出对应的表情变化。",
			"English":  "Enables the AI to send expression commands, causing the model to change its expression accordingly.",
			"Japanese": "AIが表情コマンドを送信し、それに応じてモデルの表情を変化させることができます。",
		},
	}
}
func (self *Mod_SendEmotion) API_SetEnable(value bool) { self.enable = value }
func (self *Mod_SendEmotion) API_GetEnable() bool      { return self.enable }
func (self *Mod_SendEmotion) API_IsStable() bool {
	return self.enable == self.last_enable
}
func (self *Mod_SendEmotion) Calibrate() {
	self.last_enable = self.enable
}
func (self *Mod_SendEmotion) API_GetID() string { return self.id }
