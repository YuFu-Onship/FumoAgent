package plugins

import (
	"myapp/src/model"
	"strconv"
)

type INI_SelectiveSpeech struct {
	enable   bool `ini:"enable"`
	chinese  bool `ini:"chinese"`
	english  bool `ini:"english"`
	japanese bool `ini:"japanese"`
}

type Mod_SelectiveSpeech struct {
	language      string
	languageTable LangPack
	trunk         *model.Trunk

	last_chinese  bool
	last_english  bool
	last_japanese bool

	chinese  bool
	english  bool
	japanese bool

	enable      bool
	last_enable bool

	id string
}

func New_SelectiveSpeech(trunk *model.Trunk) *Mod_SelectiveSpeech {
	self := Mod_SelectiveSpeech{
		trunk: trunk,
	}
	self.language = "Chinese"
	self.languageTable = self.Language()

	self.chinese = true
	self.english = true
	self.japanese = true
	self.last_chinese = self.chinese
	self.last_english = self.english
	self.last_japanese = self.japanese

	self.enable = true
	self.last_enable = self.enable

	self.id = "selective_speech"
	return &self
}

// 检测函数
func (self *Mod_SelectiveSpeech) Handle() bool {
	return true
}

// 执行功能
func (self *Mod_SelectiveSpeech) Execute(pm model.PluginMeta) (string, bool) {
	var text string
	if val, ok := pm.Args["text"]; ok {
		if str, ok := val.(string); ok {
			text = str
		}
	}
	var language string
	if val, ok := pm.Args["language"]; ok {
		language = val.(string)
	}

	if language == "Chinese" {
		text = self.trunk.LangConv.ChineseToKatana(text)
	}
	self.trunk.PlayVoice(text)
	return "", false
}

// 获取描述
func (self *Mod_SelectiveSpeech) Name() string { return "selective_speech" }
func (self *Mod_SelectiveSpeech) Desc() string {
	desc := self.languageTable["Desc"][self.language]
	if self.chinese {
		desc += "<Chinese>"
	}
	if self.english {
		desc += "<English>"
	}
	if self.japanese {
		desc += "<Japanese>"
	}

	return desc
}

func (self *Mod_SelectiveSpeech) Params() string {
	return self.languageTable["Params"][self.language]
}

// 参数设置
func (self *Mod_SelectiveSpeech) Setting()     {}
func (self *Mod_SelectiveSpeech) LoadSetting() {}

// 设置语言
func (self *Mod_SelectiveSpeech) API_SetLanguage(language string) { self.language = language }

// 内部语言
func (self *Mod_SelectiveSpeech) Language() LangPack {
	lp := LangPack{
		"Desc": LangZone{
			"Chinese":  "选择发言内容,可以控制说语音引擎什么,你只能使用的语言参数有:",
			"English":  "Select the speech content to control what the TTS engine says. The available language parameters are:",
			"Japanese": "発言内容を選択し、音声エンジンが話す内容を制御できます。使用可能な言語パラメータは以下の通りです:",
		},
		"Params": LangZone{
			"Chinese":  `<plugin name="selective_speech" args={"text":"需要说的内容","language":"Chinese"} />`,
			"English":  `<plugin name="selective_speech" args={"text":"Content to be spoken","language":"English"} />`,
			"Japanese": `<plugin name="selective_speech" args={"text":"発言する内容","language":"Japanese"} />`,
		},
	}
	return lp
}

// 外部调用 ---------------------------------------------------------------------
func (self *Mod_SelectiveSpeech) Options() []map[string]any {
	return []map[string]any{
		{"id": "isChinese", "value": self.chinese},
		{"id": "isEnglish", "value": self.english},
		{"id": "isJapanese", "value": self.japanese},
	}
}

func (self *Mod_SelectiveSpeech) API_SetValue(name string, value any) {
	switch name {
	case "isChinese":
		self.chinese = value.(bool)
	case "isEnglish":
		self.english = value.(bool)
	case "isJapanese":
		self.japanese = value.(bool)
	}
}
func (self *Mod_SelectiveSpeech) API_SetValueString(data map[string]string) {
	if value, err := strconv.ParseBool(data["isChinese"]); err == nil {
		self.chinese = value
	}
	if value, err := strconv.ParseBool(data["isEnglish"]); err == nil {
		self.english = value
	}
	if value, err := strconv.ParseBool(data["isJapanese"]); err == nil {
		self.japanese = value
	}
}

func (self *Mod_SelectiveSpeech) API_GetValueString() map[string]string {
	return map[string]string{
		"isChinese":  strconv.FormatBool(self.chinese),
		"isEnglish":  strconv.FormatBool(self.english),
		"isJapanese": strconv.FormatBool(self.japanese),
	}
}

func (self *Mod_SelectiveSpeech) GUI_Language() map[string]map[string]string {
	return map[string]map[string]string{
		"Title": {
			"Chinese":  "选择性发言",
			"English":  "Selective Speaking",
			"Japanese": "選択発言",
		},
		"Desc": {
			"Chinese":  "能够让AI自主选择说什么。",
			"English":  "Allows the AI to independently choose what to say.",
			"Japanese": "AIが何を話すかを自主的に選択できるようにします。",
		},
		"Detail": {
			"Chinese":  "能够使AI自主选择可以说什么,而不用直接将整段内容都转化为音频。\n可以选择说话时的语言,英语在一些情况下可能并不能正常工作。",
			"English":  "Enables the AI to selectively choose what to say, instead of converting the entire text into audio directly.\nYou can choose the speaking language; however, English may not work properly under certain circumstances.",
			"Japanese": "すべてのテキストを直接音声に変換するのではなく、AIが話す内容を自主的に選択できるようにします。\n発言時の言語を選択できますが、一部の状況下では英語が正常に機能しない場合があります。",
		},
		"isChinese": {
			"Chinese":  "中文",
			"English":  "Chinese",
			"Japanese": "中国語",
		},
		"isEnglish": {
			"Chinese":  "英语",
			"English":  "English",
			"Japanese": "英語",
		},
		"isJapanese": {
			"Chinese":  "日语",
			"English":  "Japanese",
			"Japanese": "日本語",
		},
	}
}

func (self *Mod_SelectiveSpeech) API_SetEnable(value bool) { self.enable = value }
func (self *Mod_SelectiveSpeech) API_GetEnable() bool      { return self.enable }
func (self *Mod_SelectiveSpeech) API_IsStable() bool {
	result := true
	result = self.enable == self.last_enable && result
	result = self.chinese == self.last_chinese && result
	result = self.english == self.last_english && result
	result = self.japanese == self.last_japanese && result
	return result
}
func (self *Mod_SelectiveSpeech) Calibrate() {
	self.last_enable = self.enable
	self.last_chinese = self.chinese
	self.last_english = self.english
	self.last_japanese = self.japanese
}
func (self *Mod_SelectiveSpeech) API_GetID() string { return self.id }
