package plugins

import (
	"myapp/src/config"
	"myapp/src/model"
	fumovoice "myapp/src/tools/FumoVoice"
)

type Mod_VoiceVoxTone struct {
	language  string
	langTable map[string]map[string]string
	trunk     *model.Trunk

	enable      bool
	last_enable bool

	id string
}

func New_VoiceVoxTone(trunk *model.Trunk) *Mod_VoiceVoxTone {
	self := &Mod_VoiceVoxTone{trunk: trunk}
	self.language = "Chinese"
	self.langTable = self.Language()
	self.id = "voice_vox_tone"

	self.enable = true
	self.last_enable = true
	return self
}

func (self *Mod_VoiceVoxTone) Execute(pm model.PluginMeta) (string, bool) {
	switch data := self.trunk.VoicePreset.(type) {
	case config.Preset_VoiceVox:
		self.trunk.VoiceVoxToneCount = 0
		self.trunk.VoiceVoxPlayCount = 0
		self.trunk.VoiceVoxToneCount += 1
		if !self.trunk.VoiceVox.API_GetRunState() {
			self.trunk.VoiceVox.StartServer()
			self.trunk.VoiceVoxDict = self.trunk.VoiceVox.API_GetSpeakerDict()
		}
		speaker := data.Speaker
		styles := self.trunk.VoiceVoxDict[speaker].Styles

		tone := pm.Args["tone"].(string)
		id := self.find_tone_id(styles, tone)
		self.trunk.VoiceVoxSpeakerID = id
	}

	return "", false
}

// 根据语气找到ID
func (self *Mod_VoiceVoxTone) find_tone_id(styles []fumovoice.Style, tone string) int {
	for i := range styles {
		if styles[i].Name == tone {
			return styles[i].ID
		}
	}
	return 1
}

// 插件功能描述
func (self *Mod_VoiceVoxTone) Desc() string {
	tones := ""
	switch data := self.trunk.VoicePreset.(type) {
	case config.Preset_VoiceVox:
		if !self.trunk.VoiceVox.API_GetRunState() {
			self.trunk.VoiceVox.StartServer()
			self.trunk.VoiceVoxDict = self.trunk.VoiceVox.API_GetSpeakerDict()
		}
		speaker := data.Speaker
		styles := self.trunk.VoiceVoxDict[speaker].Styles
		for i := range styles {
			tones = tones + styles[i].Name + "><"
		}
	}

	desc := self.langTable["Desc"][self.language]
	return desc + "<" + tones + ">"
}

// 插件参数描述
func (self *Mod_VoiceVoxTone) Params() string {
	return self.langTable["Params"][self.language]
}

// 语言索引
func (self *Mod_VoiceVoxTone) Language() map[string]map[string]string {
	return map[string]map[string]string{
		"Desc": {
			"Chinese":  "必须先使用语音插件才能使用该插件。设置当前的说话的语气，不设置则使用默认语气，你可以使用的语气有：",
			"English":  "This plugin requires the voice plugin to be enabled first. Sets the current speaking tone; if not set, the default tone will be used. Available tones are:",
			"Japanese": "このプラグインを使用するには、まず音声プラグインを有効にする必要があります。現在の話し方のトーンを設定します。設定しない場合はデフォルトのトーンが使用されます。利用可能なトーンは以下の通りです：",
		},
		"Params": {
			"Chinese":  `<plugin name="voice_vox_tone" args={tone=""} />`,
			"English":  `<plugin name="voice_vox_tone" args={tone=""} />`,
			"Japanese": `<plugin name="voice_vox_tone" args={tone=""} />`,
		},
	}
}

// 外部调用 ---------------------------------------------------------------------
func (self *Mod_VoiceVoxTone) Options() []map[string]any {
	return []map[string]any{}
}
func (self *Mod_VoiceVoxTone) API_SetLanguage(lang string) { self.language = lang }

func (self *Mod_VoiceVoxTone) API_SetValue(name string, value any)   {}
func (self *Mod_VoiceVoxTone) API_SetValueString(map[string]string)  {}
func (self *Mod_VoiceVoxTone) API_GetValueString() map[string]string { return map[string]string{} }

func (self *Mod_VoiceVoxTone) API_SetEnable(value bool) { self.enable = value }
func (self *Mod_VoiceVoxTone) API_GetEnable() bool      { return self.enable }

func (self *Mod_VoiceVoxTone) API_IsStable() bool { return self.enable == self.last_enable }
func (self *Mod_VoiceVoxTone) Calibrate()         { self.last_enable = self.enable }

func (self *Mod_VoiceVoxTone) API_GetID() string { return self.id }
func (self *Mod_VoiceVoxTone) Name() string      { return self.id }

func (self *Mod_VoiceVoxTone) GUI_Language() map[string]map[string]string {
	return map[string]map[string]string{
		"Title": {
			"Chinese":  "VoiceVox语气选择",
			"English":  "VoiceVox Tone Selection",
			"Japanese": "VoiceVoxトーン選択",
		},
		"Desc": {
			"Chinese":  "允许AI指定VoiceVox语音语气",
			"English":  "Allow AI to specify the VoiceVox speech tone",
			"Japanese": "AIがVoiceVoxの音声トーンを指定できるようにします",
		},
		"Detail": {
			"Chinese":  "当VoiceVox当前人设下拥有其他语气时调用",
			"English":  "Called when the current VoiceVox character style has alternative tones available",
			"Japanese": "VoiceVoxの現在のキャラクターに他のトーンが存在する場合に呼び出されます",
		},
	}
}
