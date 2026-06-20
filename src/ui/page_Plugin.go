package ui

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type PluginModel interface {
	Update(gtx C, parent *PluginPage) D
}

type PluginPage struct {
	style    Style
	MainPage *Page

	plugin_list widget.List
	plugins     []PluginModel
	pluginModel []PluginGener

	language  string
	langTable LangPack

	cur_expand_enum string
}

func New_PluginPage(style Style, mainPage *Page) *PluginPage {
	self := PluginPage{
		style:    style,
		MainPage: mainPage,
	}

	self.language = "Chinese"
	self.langTable = self.Language()

	self.plugin_list = widget.List{}
	self.plugin_list.Axis = layout.Vertical

	return &self
}

// 恢复默认设置
func (self *PluginPage) Default() {}

// 界面更新
func (self *PluginPage) Update(gtx C) D {
	self.language = self.MainPage.trunk.Language

	// 标题栏样式
	titleFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(material.Label(self.style.theme, unit.Sp(40), self.langTable["title"][self.language]).Layout),
			layout.Flexed(1, Flexer()),
		)
	}

	// 插件列表
	list := material.List(self.style.theme, &self.plugin_list)
	list.ScrollbarStyle.Indicator.Color.A = 64
	list.ScrollbarStyle.Indicator.MinorWidth = unit.Dp(4)

	// 修改/保存参数 -------------------------------
	isStable := true
	eleList := []W{}
	for i := range self.pluginModel {
		eleList = append(eleList, func(gtx C) D { return self.pluginModel[i].Update(gtx, self) })
		isStable = self.pluginModel[i].API_IsStable() && isStable
	}
	if !isStable {
		for i := range self.pluginModel {
			// 参数对齐
			self.pluginModel[i].Calibrate()
			// 普通参数与启用参数
			self.MainPage.trunk.PluginCoe[self.pluginModel[i].API_GetID()] = self.pluginModel[i].Plugin.API_GetValueString()
			self.MainPage.trunk.PluginCoe["plugin"][self.pluginModel[i].API_GetID()] = strconv.FormatBool(self.pluginModel[i].Plugin.API_GetEnable())
		}
		self.MainPage.trunk.UpdatePluginCoe(self.MainPage.trunk.PluginCoe)
	}

	// 返回布局 --------------------------------
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(titleFunc),
		layout.Rigid(func(gtx C) D {
			return list.Layout(gtx, len(eleList), func(gtx C, index int) D {
				return layout.Inset{Bottom: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx C) D {
					return eleList[index](gtx)
				})
			})
		}),
	)
}

func (self *PluginPage) Language() LangPack {
	return LangPack{
		"title": LangZone{
			"Chinese":  "插件",
			"English":  "Plugin",
			"Japanese": "プラグイン",
		},
	}
}

// API
func (self *PluginPage) API() map[string]any {
	return map[string]any{
		"add_plugin": func(plugin PluginCtrl) {
			if self.pluginModel == nil {
				self.pluginModel = []PluginGener{}
			}
			new_model := New_PluginGener(self.style, plugin)
			self.pluginModel = append(self.pluginModel, *new_model)
		},
	}
}

// 测试 ------------------------------------------------------------------------------------------
type Plugin_Test struct {
	Is_test_1 bool
	Is_test_2 bool
}

func New_Plugin_Test() *Plugin_Test {
	self := Plugin_Test{}
	return &self
}

func (self *Plugin_Test) Options() []map[string]any {
	result := []map[string]any{
		{"id": "1", "value": self.Is_test_1},
		{"id": "2", "value": self.Is_test_2},
	}

	return result
}

func (self *Plugin_Test) API_SetValue(name string, value any) {
	switch name {
	case "1":
		self.Is_test_1 = value.(bool)
	case "2":
		self.Is_test_1 = value.(bool)
	}
}

func (self *Plugin_Test) GUI_Language() map[string]map[string]string {
	return map[string]map[string]string{
		"Title": {
			"Chinese":  "获取时间",
			"English":  "",
			"Japanese": "",
		},
		"Desc": {
			"Chinese":  "允许ai获取到当前的系统时间",
			"English":  "",
			"Japanese": "",
		},
		"1": {
			"Chinese":  "一",
			"English":  "one",
			"Japanese": "1",
		},
		"2": {
			"Chinese":  "二",
			"English":  "two",
			"Japanese": "2",
		},
		"Detail": {
			"Chinese":  "无",
			"English":  "",
			"Japanese": "",
		},
	}
}

// 选择性说话 ------------------------------------------------------------------------------------
// 语音生成可以选择所有部分,也可以由ai决定只说一部分
type Plugin_SelectiveSpeech struct {
	style       Style
	expandPanel ExpandPanel
	title       string
	desc        string
	prompt      string

	language      string
	languageTable LangPack

	switch_Enable        widget.Bool
	switch_EnableChinese widget.Bool

	switch_EnableEnglish  widget.Bool
	switch_EnableJapanese widget.Bool
}

func New_Plugin_SelectiveSpeech(style Style) *Plugin_SelectiveSpeech {
	p := Plugin_SelectiveSpeech{
		style:       style,
		expandPanel: *New_ExpandPanel(style, 210),
	}
	p.languageTable = p.Language()
	return &p
}

func (self *Plugin_SelectiveSpeech) Update(gtx C, parent *PluginPage) D {
	self.language = parent.MainPage.trunk.Language
	headFunc := func(gtx C) D {
		gtx.Constraints.Min.Y = gtx.Dp(50)
		gtx.Constraints.Max.Y = gtx.Dp(100)
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: unit.Dp(12), Right: unit.Dp(40)}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(material.Label(self.style.theme, unit.Sp(20), self.languageTable["Title"][self.language]).Layout),
						layout.Rigid(material.Label(self.style.theme, unit.Sp(12), self.languageTable["Desc"][self.language]).Layout),
					)
				})
			}),
			layout.Flexed(1, Flexer()),
		)
	}
	panelFunc := func(gtx C) D {
		return self.layout_panel(gtx)
	}
	d, _ := self.expandPanel.Update(gtx, headFunc, panelFunc)
	return d
}
func (self *Plugin_SelectiveSpeech) layout_panel(gtx C) D {
	switchFunc := func(gtx C, lable material.LabelStyle, switchWidget *widget.Bool) D {
		s := material.Switch(self.style.theme, switchWidget, "")
		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(lable.Layout),
				layout.Flexed(1, Flexer()),
				layout.Rigid(s.Layout),
			)
		})
	}

	label_1 := material.Label(self.style.theme, unit.Sp(18), self.languageTable["Option_Enable"][self.language])
	label_2 := material.Label(self.style.theme, unit.Sp(18), self.languageTable["Option_Lang"][self.language])
	label_cn := material.Label(self.style.theme, unit.Sp(14), self.languageTable["Option_CN"][self.language])
	label_en := material.Label(self.style.theme, unit.Sp(14), self.languageTable["Option_EN"][self.language])
	label_jp := material.Label(self.style.theme, unit.Sp(14), self.languageTable["Option_JP"][self.language])

	leftFunc := func(gtx C) D {
		return BorderBox(
			self.style,
			self.style.darkmode.currentColor.IdleBg,
			self.style.darkmode.currentColor.Fg,
			self.languageTable["Option_Title"][self.language],
		).Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Dp(150)
			gtx.Constraints.Max.X = gtx.Dp(150)
			gtx.Constraints.Min.Y = gtx.Dp(150)
			gtx.Constraints.Max.Y = gtx.Dp(150)
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D { return switchFunc(gtx, label_1, &self.switch_Enable) }),
				layout.Rigid(func(gtx C) D { return label_2.Layout(gtx) }),
				layout.Rigid(Spacer(0, gtx.Dp(8))),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(Spacer(gtx.Dp(12), 0)),
						layout.Rigid(func(gtx C) D { return switchFunc(gtx, label_cn, &self.switch_EnableChinese) }),
					)
				}),
				layout.Rigid(Spacer(0, gtx.Dp(8))),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(Spacer(gtx.Dp(12), 0)),
						layout.Rigid(func(gtx C) D { return switchFunc(gtx, label_en, &self.switch_EnableEnglish) }),
					)
				}),
				layout.Rigid(Spacer(0, gtx.Dp(8))),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(Spacer(gtx.Dp(12), 0)),
						layout.Rigid(func(gtx C) D { return switchFunc(gtx, label_jp, &self.switch_EnableJapanese) }),
					)
				}),
			)
		})
	}

	rightFunc := func(gtx C) D {
		l := BorderBox(
			self.style,
			self.style.darkmode.currentColor.IdleBg,
			self.style.darkmode.currentColor.Fg,
			self.languageTable["Detail_Title"][self.language],
		).Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Dp(150)
			gtx.Constraints.Max.Y = gtx.Dp(150)
			label := material.Label(self.style.theme, unit.Sp(14), self.languageTable["Detail_Desc"][self.language])
			return layout.Flex{}.Layout(gtx, layout.Flexed(1, label.Layout))
		})
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D { return l }),
		)

	}

	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(leftFunc),
			layout.Rigid(Spacer(gtx.Sp(16), 0)),
			layout.Rigid(rightFunc),
		)
	})
}
func (self *Plugin_SelectiveSpeech) Language() LangPack {
	return LangPack{
		"Title": LangZone{
			"Chinese":  "选择性发言",
			"English":  "Selective Speaking",
			"Japanese": "選択的発言",
		},
		"Desc": LangZone{
			"Chinese":  "能够让AI自主选择说什么,而不用直接讲述所有的内容。",
			"English":  "Enables the AI to choose what to say independently, instead of outputting all content directly.",
			"Japanese": "AIが自分で判断して話す内容を選べるようにする機能です。",
		},
		"Detail_Title": LangZone{
			"Chinese":  "详情",
			"English":  "Details",
			"Japanese": "詳細",
		},
		"Detail_Desc": LangZone{
			"Chinese":  "能够使ai自主选择可以说什么,而不用直接将整段内容都转化为音频。\n可以选择说话时的语言,英语在一些情况下可能并不能正常工作。",
			"English":  "Allows the AI to choose what to say instead of converting the entire text into audio.\nYou can select the speech language, though English may not work properly under certain circumstances.",
			"Japanese": "テキスト全体をそのまま音声に変換するのではなく、AIが話す内容を自主的に選択できるようになります。\n発話時の言語を選択できますが、一部の環境では英語が正常に機能しない場合があります。",
		},
		"Option_Title": LangZone{
			"Chinese":  "选项",
			"English":  "Options",
			"Japanese": "オプション",
		},
		"Option_Enable": LangZone{
			"Chinese":  "启用",
			"English":  "Enable",
			"Japanese": "有効化",
		},
		"Option_Lang": LangZone{
			"Chinese":  "语言",
			"English":  "Language",
			"Japanese": "言語",
		},
		"Option_CN": LangZone{
			"Chinese":  "中文",
			"English":  "Chinese",
			"Japanese": "中国語",
		},
		"Option_EN": LangZone{
			"Chinese":  "英语",
			"English":  "English",
			"Japanese": "英語",
		},
		"Option_JP": LangZone{
			"Chinese":  "日语",
			"English":  "Japanese",
			"Japanese": "日本語",
		},
	}
}

type Plugin_SendEmotion struct {
	style       Style
	expandPanel ExpandPanel
	title       string
	desc        string

	language      string
	languageTable LangPack

	switch_Enable widget.Bool
}

func Mew_Plugin_SendEmotion(style Style) *Plugin_SendEmotion {
	p := Plugin_SendEmotion{
		style:       style,
		expandPanel: *New_ExpandPanel(style, 210),
	}
	p.languageTable = p.Language()
	return &p
}
func (self *Plugin_SendEmotion) Update(gtx C, parent *PluginPage) D {
	self.language = parent.MainPage.trunk.Language
	headFunc := func(gtx C) D {
		gtx.Constraints.Min.Y = gtx.Dp(50)
		gtx.Constraints.Max.Y = gtx.Dp(100)
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: unit.Dp(12), Right: unit.Dp(40)}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(material.Label(self.style.theme, unit.Sp(20), self.languageTable["Title"][self.language]).Layout),
						layout.Rigid(material.Label(self.style.theme, unit.Sp(12), self.languageTable["Desc"][self.language]).Layout),
					)
				})
			}),
			layout.Flexed(1, Flexer()),
		)
	}
	d, _ := self.expandPanel.Update(gtx, headFunc, self.layout_Panel)
	return d
}

func (self *Plugin_SendEmotion) layout_Panel(gtx C) D {
	switchFunc := func(gtx C, lable material.LabelStyle, switchWidget *widget.Bool) D {
		s := material.Switch(self.style.theme, switchWidget, "")
		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(lable.Layout),
				layout.Flexed(1, Flexer()),
				layout.Rigid(s.Layout),
			)
		})
	}

	label_1 := material.Label(self.style.theme, unit.Sp(18), self.languageTable["Option_Enable"][self.language])

	leftFunc := func(gtx C) D {
		return BorderBox(
			self.style,
			self.style.darkmode.currentColor.IdleBg,
			self.style.darkmode.currentColor.Fg,
			self.languageTable["Option_Title"][self.language],
		).Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Dp(150)
			gtx.Constraints.Max.X = gtx.Dp(150)
			gtx.Constraints.Min.Y = gtx.Dp(150)
			gtx.Constraints.Max.Y = gtx.Dp(150)
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D { return switchFunc(gtx, label_1, &self.switch_Enable) }),
				layout.Rigid(Spacer(0, gtx.Dp(8))),
			)
		})
	}

	rightFunc := func(gtx C) D {
		l := BorderBox(
			self.style,
			self.style.darkmode.currentColor.IdleBg,
			self.style.darkmode.currentColor.Fg,
			self.languageTable["Detail_Title"][self.language],
		).Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Dp(150)
			gtx.Constraints.Max.Y = gtx.Dp(150)
			label := material.Label(self.style.theme, unit.Sp(14), self.languageTable["Detail_Desc"][self.language])
			return layout.Flex{}.Layout(gtx, layout.Flexed(1, label.Layout))
		})
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D { return l }),
		)

	}

	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(leftFunc),
			layout.Rigid(Spacer(gtx.Sp(16), 0)),
			layout.Rigid(rightFunc),
		)
	})
}

func (self *Plugin_SendEmotion) Language() LangPack {
	return LangPack{
		"Title": LangZone{
			"Chinese":  "发送表情",
			"English":  "Send Emotion",
			"Japanese": "感情送信",
		},

		"Desc": LangZone{
			"Chinese":  "允许 AI 控制模型做出相应的表情。",
			"English":  "Allows the AI to control the model's facial expressions.",
			"Japanese": "AI がモデルの表情を制御できるようにします。",
		},

		"Detail_Title": LangZone{
			"Chinese":  "详情",
			"English":  "Details",
			"Japanese": "詳細",
		},

		"Detail_Desc": LangZone{
			"Chinese":  "能够让 AI 发送表情指令，使模型做出对应的表情变化。",
			"English":  "Enables the AI to send emotion commands and trigger corresponding facial expressions on the model.",
			"Japanese": "AI が感情コマンドを送信し、モデルに対応する表情変化を行わせます。",
		},

		"Option_Title": LangZone{
			"Chinese":  "选项",
			"English":  "Options",
			"Japanese": "オプション",
		},

		"Option_Enable": LangZone{
			"Chinese":  "启用",
			"English":  "Enable",
			"Japanese": "有効化",
		},
	}
}
