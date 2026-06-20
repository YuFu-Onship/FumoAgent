package plugins

import (
	"myapp/src/model"
	"time"
)

type INI_CurrentTime struct {
	Enable bool `ini:"enable"`
}

type Mod_CurrentTime struct {
	language      string
	languageTable LangPack
	trunk         *model.Trunk

	enable      bool
	last_enable bool

	id string
}

func New_CurrentTime(trunk *model.Trunk) *Mod_CurrentTime {
	self := Mod_CurrentTime{
		trunk: trunk,
	}
	self.language = "Chinese"
	self.languageTable = self.Language()

	self.enable = true
	self.last_enable = self.enable

	self.id = "current_time"

	return &self
}

// 检测函数
func (self *Mod_CurrentTime) Handle() bool {
	return self.enable
}

// 执行功能（返回型插件：返回时间字符串与 true）
func (self *Mod_CurrentTime) Execute(pm model.PluginMeta) (string, bool) {
	now := time.Now()
	var timeText string

	switch self.language {
	case "English":
		timeText = now.Format("Monday, January 2, 2006, 15:04:05")
	case "Japanese":
		weekdays := []string{"日", "月", "火", "水", "木", "金", "土"}
		weekdayStr := weekdays[now.Weekday()]
		timeText = now.Format("2006年1月2日") + "(" + weekdayStr + ") " + now.Format("15:04:05")
	case "Chinese":
		fallthrough
	default:
		weekdays := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
		weekdayStr := weekdays[now.Weekday()]
		timeText = now.Format("2006年1月2日 ") + weekdayStr + now.Format(" 15:04:05")
	}

	return timeText, true
}

// 获取描述
func (self *Mod_CurrentTime) Name() string { return "current_time" }
func (self *Mod_CurrentTime) Desc() string {
	return self.languageTable["Desc"][self.language]
}

func (self *Mod_CurrentTime) Params() string {
	return self.languageTable["Params"][self.language]
}

// 参数设置
func (self *Mod_CurrentTime) Setting()     {}
func (self *Mod_CurrentTime) LoadSetting() {}

// 设置语言
func (self *Mod_CurrentTime) API_SetLanguage(language string) { self.language = language }

// 内部语言
func (self *Mod_CurrentTime) Language() LangPack {
	lp := LangPack{
		"Desc": LangZone{
			"Chinese":  "获取当前系统的时间与日期。当用户询问有关当前时间、日期、星期几等问题时，你必须调用此插件获取准确信息。使用单标签格式：<plugin name=\"current_time\" args={} />",
			"English":  "Get the current system time and date. You must call this plugin to get accurate information when the user asks about the current time, date, or day of the week. Use single-tag format: <plugin name=\"current_time\" args={} />",
			"Japanese": "現在のシステム時刻と日付を取得します。ユーザーから現在時刻、日付、曜日などについて尋ねられた際、必ずこのプラグインを呼び出して正確な情報を取得してください。単一タグ形式を使用：<plugin name=\"current_time\" args={} />",
		},
		"Params": LangZone{
			"Chinese":  `<plugin name="current_time" args={} />`,
			"English":  `<plugin name="current_time" args={} />`,
			"Japanese": `<plugin name="current_time" args={} />`,
		},
	}
	return lp
}

// 外部调用 ---------------------------------------------------------------------
func (self *Mod_CurrentTime) Options() []map[string]any {
	return []map[string]any{}
}

func (self *Mod_CurrentTime) API_SetValue(name string, value any)   {}
func (self *Mod_CurrentTime) API_SetValueString(map[string]string)  {}
func (self *Mod_CurrentTime) API_GetValueString() map[string]string { return map[string]string{} }

func (self *Mod_CurrentTime) GUI_Language() map[string]map[string]string {
	return map[string]map[string]string{
		"Title": {
			"Chinese":  "获取时间",
			"English":  "Get Time",
			"Japanese": "時間取得",
		},
		"Desc": {
			"Chinese":  "允许AI获取到当前的系统时间",
			"English":  "Allows the AI to get the current system time",
			"Japanese": "AIが現在のシステム時間を取得できるようにします",
		},
		"Detail": {
			"Chinese":  "无",
			"English":  "None",
			"Japanese": "なし",
		},
	}
}

func (self *Mod_CurrentTime) API_SetEnable(value bool) {
	self.enable = value
}
func (self *Mod_CurrentTime) API_GetEnable() bool { return self.enable }
func (self *Mod_CurrentTime) API_IsStable() bool {
	return self.enable == self.last_enable

}
func (self *Mod_CurrentTime) Calibrate() {
	self.last_enable = self.enable
}
func (self *Mod_CurrentTime) API_GetID() string {
	return self.id
}
