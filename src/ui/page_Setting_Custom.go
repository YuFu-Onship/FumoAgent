package ui

import (
	"myapp/src/config"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CustomSetting struct {
	style           Style
	mainSettingPage *SettingPage

	langBox  DropDown
	lastLang string

	mainList widget.List

	// 配色
	colorEnum widget.Enum
	colorList []string
	lastColor string
	colorBox  DropDown
}

func New_CustomSetting(style Style) *CustomSetting {
	mainList := widget.List{}
	mainList.Axis = layout.Vertical

	curLang := config.API_CUSTOM_LANGUAGE_GetCurrentLanguage()
	curColor := config.API_CUSTOM_COLOR_GetColorID()
	var langValue string
	switch curLang {
	case "Chinese":
		langValue = "中文"
	case "Japanese":
		langValue = "日本語"
	case "English":
		langValue = "English"
	}

	ts := CustomSetting{
		style:     style,
		colorEnum: widget.Enum{},
		mainList:  mainList,
		langBox:   *New_DropDown(style, []string{"中文", "日本語", "English"}, langValue),
		lastLang:  curLang,
		colorBox:  *New_DropDown(style, []string{"cirno", "satori", "koishi", "teto", "windows"}, curColor),
		lastColor: curColor,
	}
	return &ts
}
func (self *CustomSetting) Title() string {
	return self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Custom_Title
}
func (self *CustomSetting) Default() {

}

func (self *CustomSetting) Update(gtx C, parent *SettingPage) D {
	self.mainSettingPage = parent
	// 逻辑部分
	// 语言
	curLang := self.langBox.API_GetValue()
	var language string
	switch curLang {
	case "中文":
		language = "Chinese"
	case "日本語":
		language = "Japanese"
	case "English":
		language = "English"
	default:
		language = "Chinese"
	}
	if language != self.lastLang {
		self.lastLang = language
		parent.MainPage.trunk.Language = language
		parent.MainPage.trunk.Handler_Custom.API_CUSTOM_LANGUAGE_SetCurrentLanguage(language)
		parent.MainPage.trunk.LanguageTable = parent.MainPage.trunk.Handler_Custom.API_CUSTOM_LANGUAGE_GetLanguageTable(language)
	}

	// 颜色
	curColor := self.colorBox.API_GetValue()
	if curColor != self.lastColor {
		self.lastColor = curColor
		self.style.darkmode.API_SetColor(curColor)
		self.style.darkmode.Apply(self.style.theme)
	}

	// 界面布局
	listElements := []layout.Widget{
		func(gtx C) D {
			return self.build_TieleLabel(gtx, parent.MainPage.trunk.LanguageTable.SettingPage_Custom_ChoiceLanguage)
		},
		func(gtx C) D { return self.langBox.Update(gtx) },
		func(gtx C) D {
			return self.build_TieleLabel(gtx, parent.MainPage.trunk.LanguageTable.SettingPage_Custom_ChoiceColor)
		},
		func(gtx C) D { return self.colorBox.Update(gtx) },
	}

	list := material.List(self.style.theme, &self.mainList)
	list.ScrollbarStyle.Indicator.Color.A = 64
	list.ScrollbarStyle.Indicator.MinorWidth = unit.Dp(4)
	return list.Layout(gtx, len(listElements), func(gtx C, index int) D {
		return layout.Inset{Right: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx C) D { return listElements[index](gtx) })
	})
}

// 语言选择布局
func (self *CustomSetting) layout_Language(gtx C) D {
	return self.langBox.Update(gtx)

}

func (self *CustomSetting) layout_Color(gtx C) D {
	return D{}
}

// 创建标题栏
func (self *CustomSetting) build_TieleLabel(gtx C, title string) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D { return material.Label(self.style.theme, unit.Sp(16), title).Layout(gtx) }),
		layout.Flexed(1, Flexer()),
	)
}

// // 颜色表示栏
type ColorShowBar struct {
	style Style
	color ColorStruct
	title string
	btn   widget.Clickable
}

func New_ColorShowBar(style Style, title string) *ColorShowBar {
	csb := ColorShowBar{
		style: style,
		title: title,
	}
	return &csb
}

func (self *ColorShowBar) Update(gtx C, enum widget.Enum) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			DrawRRect(gtx, gtx.Constraints.Min, self.style.darkmode.currentColor.IdleBg, gtx.Dp(4))
			return D{}
		}),
		layout.Stacked(func(gtx C) D {
			return self.layout_Stack(gtx)
		}),
	)
}

func (self *ColorShowBar) layout_Stack(gtx C) D {
	textLabel := material.Label(self.style.theme, unit.Sp(18), self.title)
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(textLabel.Layout),
	)
}
