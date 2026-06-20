package plugins

import (
	"myapp/src/model"
	"strconv"
)

type Mod_ModelRotate struct {
	trunk     *model.Trunk
	lang      string
	langTable LangPack

	enable      bool
	last_enable bool

	id string
}

func New_ModelRotate(trunk *model.Trunk) *Mod_ModelRotate {
	self := Mod_ModelRotate{
		trunk: trunk,
		lang:  "Chinese",
	}
	self.langTable = self.Language()
	self.enable = true
	self.last_enable = self.enable

	self.id = "set_model_rotate"
	return &self
}

func (self *Mod_ModelRotate) Execute(pm model.PluginMeta) (string, bool) {
	weekStr, ok := pm.Args["week"].(string)
	if !ok {
		return "", false
	}
	value, err := strconv.ParseInt(weekStr, 10, 64)
	if err != nil {
		return "", false
	}

	newArgs := model.PluginMeta{
		Name: "server_set_live2d_rotate",
		Args: map[string]any{
			"week": value,
		},
	}

	self.trunk.WsServer.SendCommand(newArgs.Name, newArgs.Args)
	return "", false
}

func (self *Mod_ModelRotate) API_SetLanguage(language string) {
	self.lang = language
}

func (self *Mod_ModelRotate) Name() string {
	return "set_model_rotate"
}

func (self *Mod_ModelRotate) Desc() string {
	return self.langTable["Desc"][self.lang]
}
func (self *Mod_ModelRotate) Params() string {
	return self.langTable["Params"][self.lang]
}

func (self *Mod_ModelRotate) Language() LangPack {
	return LangPack{
		"Desc": LangZone{
			"chinese":  "可以选择让live2d模型旋转几周",
			"English":  "Option to make the Live2D model spin for several rotations",
			"Japanese": "Live2Dモデルを数周回転させるオプション",
		},
		"Params": LangZone{
			"Chinese":  `<plugin name="set_model_rotate" args={"week":"1"} />`,
			"English":  `<plugin name="set_model_rotate" args={"week":"1"} />`,
			"Japanese": `<plugin name="set_model_rotate" args={"week":"1"} />`,
		},
	}
}

// 外部调用 ---------------------------------------------------------------------
func (self *Mod_ModelRotate) Options() []map[string]any {
	return []map[string]any{}
}

func (self *Mod_ModelRotate) API_SetValue(name string, value any)   {}
func (self *Mod_ModelRotate) API_SetValueString(map[string]string)  {}
func (self *Mod_ModelRotate) API_GetValueString() map[string]string { return map[string]string{} }

func (self *Mod_ModelRotate) GUI_Language() map[string]map[string]string {
	return map[string]map[string]string{
		"Title": {
			"Chinese":  "模型旋转",
			"English":  "Model Rotation",
			"Japanese": "モデル回転",
		},
		"Desc": {
			"Chinese":  "允许AI控制模型旋转",
			"English":  "Allows the AI to control model rotation",
			"Japanese": "AIがモデルの回転を制御できるようにします",
		},
		"Detail": {
			"Chinese":  "无",
			"English":  "None",
			"Japanese": "なし",
		},
	}
}

func (self *Mod_ModelRotate) API_SetEnable(value bool) { self.enable = value }
func (self *Mod_ModelRotate) API_GetEnable() bool      { return self.enable }
func (self *Mod_ModelRotate) API_IsStable() bool {
	return self.enable == self.last_enable
}
func (self *Mod_ModelRotate) Calibrate() {
	self.last_enable = self.enable
}
func (self *Mod_ModelRotate) API_GetID() string { return self.id }
