package config

// 语言 ---------------------------------------------------------------------------------------------
func API_CUSTOM_LANGUAGE_GetCurrentLanguage() string {
	data := ReadConfigDate()
	return data.Language
}
func API_CUSTOM_LANGUAGE_GetLanguageTable(curLang string) LanguageTable {
	switch curLang {
	case "Chinese":
		return LanguageChinese()
	case "English":
		return LanguageEnglish()
	case "Japanese":
		return LanguageJapanese()
	default:
		return LanguageChinese()
	}
}

func API_CUSTOM_LANGUAGE_SetCurrentLanguage(lang string) {
	data := ReadConfigDate()
	data.Language = lang
	SaveConfig(data)
}

// 颜色 ---------------------------------------------------------------------------------------------
// 获得当前的颜色id
func API_CUSTOM_COLOR_GetColorID() string {
	data := ReadConfigDate()
	id := data.Color
	return id
}

// 获得当前的暗色模式
func API_CUSTOM_COLOR_GetDarkmode() bool {
	data := ReadConfigDate()
	isdark := data.Darkmode
	if isdark == "true" {
		return true
	} else {
		return false
	}
}

// 保存当前的颜色id
func API_CUSTOM_COLOR_SetColorID(id string) {
	data := ReadConfigDate()
	data.Color = id
	SaveConfig(data)
}

// 保存当前的暗色模式
func API_CUSTOM_COLOR_SetDarkmode(value bool) {
	data := ReadConfigDate()
	if value {
		data.Darkmode = "true"
	} else {
		data.Darkmode = "false"
	}
	SaveConfig(data)
}

// 得到当前色板
func API_CUSTOM_COLOR_GetColorPattel(id string, isdark bool) ColorPattle {
	list := InitColorList()
	if isdark {
		return list[id].Dark
	} else {
		return list[id].Light
	}
}

// 个性化设置 ---------------------------------------------------------------------------------------
type ConfigCustom struct {
}

func New_ConfigCustom() *ConfigCustom {
	self := ConfigCustom{}
	return &self
}
func (self *ConfigCustom) API_CUSTOM_LANGUAGE_GetCurrentLanguage() string {
	data := ReadConfigDate()
	return data.Language
}
func (self *ConfigCustom) API_CUSTOM_LANGUAGE_GetLanguageTable(curLang string) LanguageTable {
	switch curLang {
	case "Chinese":
		return LanguageChinese()
	case "English":
		return LanguageEnglish()
	case "Japanese":
		return LanguageJapanese()
	default:
		return LanguageChinese()
	}
}

func (self *ConfigCustom) API_CUSTOM_LANGUAGE_SetCurrentLanguage(lang string) {
	switch lang {
	case "Chinese", "English", "Japanese":
		data := ReadConfigDate()
		data.Language = lang
		SaveConfig(data)
	default:
	}
}

// 获得当前的颜色id
func (self *ConfigCustom) API_CUSTOM_COLOR_GetColorID() string {
	data := ReadConfigDate()
	id := data.Color
	return id
}

// 获得当前的暗色模式
func (self *ConfigCustom) API_CUSTOM_COLOR_GetDarkmode() bool {
	data := ReadConfigDate()
	isdark := data.Darkmode
	if isdark == "true" {
		return true
	} else {
		return false
	}
}

// 保存当前的颜色id
func (self *ConfigCustom) API_CUSTOM_COLOR_SetColorID(id string) {
	data := ReadConfigDate()
	data.Color = id
	SaveConfig(data)
}

// 保存当前的暗色模式
func (self *ConfigCustom) API_CUSTOM_COLOR_SetDarkmode(value bool) {
	data := ReadConfigDate()
	if value {
		data.Darkmode = "true"
	} else {
		data.Darkmode = "false"
	}
	SaveConfig(data)
}
