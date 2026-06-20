package plugins

import (
	"myapp/src/model"
)

type Mod_ModelShake struct {
	trunk     *model.Trunk
	lang      string
	langTable LangPack

	enable      bool
	last_enable bool
	id          string
}

func New_ModelShake(trunk *model.Trunk) *Mod_ModelShake {
	self := Mod_ModelShake{
		trunk: trunk,
		lang:  "Chinese",
	}
	self.langTable = self.Language()
	self.enable = true
	self.last_enable = self.enable
	self.id = "set_model_shake"
	return &self
}

func (self *Mod_ModelShake) Execute(pm model.PluginMeta) (string, bool) {
	self.trunk.WsServer.SendCommand("server_set_live2d_shake", pm.Args)
	return "", false
}

func (self *Mod_ModelShake) API_SetLanguage(language string) {
	self.lang = language
}

func (self *Mod_ModelShake) Name() string {
	return "set_model_shake"
}

func (self *Mod_ModelShake) Desc() string {
	return self.langTable["Desc"][self.lang]
}
func (self *Mod_ModelShake) Params() string {
	return self.langTable["Params"][self.lang]
}

func (self *Mod_ModelShake) Language() LangPack {
	return LangPack{
		"Desc": LangZone{
			"chinese":  "可以选择让live2d模型摆动身体/摇头",
			"English":  "Option to make the Live2D model sway its body or shake its head",
			"Japanese": "Live2Dモデルに体を揺らしたり、首を振らせるオプション",
		},
		"Params": LangZone{
			"Chinese":  `<plugin name="set_model_shake" args={} />`,
			"English":  `<plugin name="set_model_shake" args={} />`,
			"Japanese": `<plugin name="set_model_shake" args={} />`,
		},
	}
}

// 外部调用 ---------------------------------------------------------------------
func (self *Mod_ModelShake) Options() []map[string]any {
	return []map[string]any{}
}

func (self *Mod_ModelShake) API_SetValue(name string, value any)  {}
func (self *Mod_ModelShake) API_SetValueString(map[string]string) {}

func (self *Mod_ModelShake) API_GetValueString() map[string]string { return map[string]string{} }
func (self *Mod_ModelShake) GUI_Language() map[string]map[string]string {
	return map[string]map[string]string{
		"Title": {
			"Chinese":  "模型摇晃",
			"English":  "Model Shaking",
			"Japanese": "モデル揺らし",
		},
		"Desc": {
			"Chinese":  "允许AI控制模型摇晃或摇头",
			"English":  "Allows the AI to control model shaking or head shaking",
			"Japanese": "AIがモデルの揺れや首振りを制御できるようにします",
		},
		"Detail": {
			"Chinese":  "无",
			"English":  "None",
			"Japanese": "なし",
		},
	}
}

func (self *Mod_ModelShake) API_SetEnable(value bool) {
	self.enable = value
}
func (self *Mod_ModelShake) API_GetEnable() bool { return self.enable }
func (self *Mod_ModelShake) API_IsStable() bool {
	return self.enable == self.last_enable
}
func (self *Mod_ModelShake) Calibrate() {
	self.last_enable = self.enable
}
func (self *Mod_ModelShake) API_GetID() string { return self.id }
