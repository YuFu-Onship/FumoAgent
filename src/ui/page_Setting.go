package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Component_SettingChildPage interface {
	Update(gtx C, settingpage *SettingPage) D
	Title() string
}
type SettingRoute struct {
	Key  string
	Icon string
	Name string
	Desc string
	// Page   Component_SettingChildPage
	Create func() Component_SettingChildPage
	Btn    *widget.Clickable
}

type Component_SettingBtn interface {
}

// 定义参数	-----------------------------------------------------------------------
type SettingPage struct {
	style    Style
	MainPage *Page

	SettingList widget.List
	mainOptions []D
	mainList    widget.List

	pageName   string
	jumpBtn    map[string]*widget.Clickable
	childPages map[string]Component_SettingChildPage
	routes     []SettingRoute

	modelSettingPage     ModelSetting
	soundSettingPage     SoundSetting
	characterSettingPage CharacterSetting

	// 页面切换
	placeholderHeight_cur float64
	placeholderHeight_tar float64

	language      string
	languageTable LangPack
}

func New_SettingPage(style Style, mainPage *Page) *SettingPage {
	self := SettingPage{
		style:    style,
		MainPage: mainPage,
		pageName: "main",

		placeholderHeight_cur: 0,
		placeholderHeight_tar: 0,
	}
	self.Init_JumpBtn()
	self.InitChildPages()

	self.language = "Chinese"
	self.languageTable = self.Language()
	return &self
}

func (self *SettingPage) Default() {
	self.pageName = "main"
}

func (self *SettingPage) Init_JumpBtn() {
	self.jumpBtn = make(map[string]*widget.Clickable)
	self.jumpBtn["main"] = new(widget.Clickable)
	self.jumpBtn["voice"] = new(widget.Clickable)
	self.jumpBtn["model"] = new(widget.Clickable)
	self.jumpBtn["character"] = new(widget.Clickable)
	self.jumpBtn["live2d"] = new(widget.Clickable)
	self.jumpBtn["custom"] = new(widget.Clickable)
	self.jumpBtn["info"] = new(widget.Clickable)
}
func (self *SettingPage) InitChildPages() {
	self.routes = []SettingRoute{
		{
			Key: "voice", Icon: "\ue767", Name: "", Desc: "",
			Create: func() Component_SettingChildPage { return New_SoundSetting(self.style) },
			Btn:    new(widget.Clickable),
		},
		{
			Key: "model", Icon: "\uf003", Name: "", Desc: "",
			Create: func() Component_SettingChildPage { return New_ModelSetting(self.style) },
			Btn:    new(widget.Clickable),
		},
		{
			Key: "character", Icon: "\uED59", Name: "", Desc: "",
			Create: func() Component_SettingChildPage { return New_CharacterSetting(self.style) },
			Btn:    new(widget.Clickable),
		},
		{
			Key: "live2d", Icon: "\uE76E", Name: "", Desc: "",
			Create: func() Component_SettingChildPage { return New_Live2DSetting(self.style) },
			Btn:    new(widget.Clickable),
		},
		{
			Key: "custom", Icon: "\ue790", Name: "", Desc: "",
			Create: func() Component_SettingChildPage { return New_CustomSetting(self.style) },
			Btn:    new(widget.Clickable),
		},
		{
			Key: "info", Icon: "\uE946", Name: "", Desc: "",
			Create: func() Component_SettingChildPage { return New_InfoSetting(self.style) },
			Btn:    new(widget.Clickable),
		},
	}
	self.childPages = make(map[string]Component_SettingChildPage)
}

// 进行绘制的主函数 ----------------------------------------------------------------------
func (self *SettingPage) Update(gtx C) D {
	self.language = self.MainPage.trunk.Language
	for i := range self.routes {
		switch self.routes[i].Key {
		default:
		case "voice":
			self.routes[i].Name = self.languageTable["voice_Title"][self.language]
			self.routes[i].Desc = self.languageTable["voice_Desc"][self.language]
		case "model":
			self.routes[i].Name = self.languageTable["model_Title"][self.language]
			self.routes[i].Desc = self.languageTable["model_Desc"][self.language]
		case "character":
			self.routes[i].Name = self.languageTable["character_Title"][self.language]
			self.routes[i].Desc = self.languageTable["character_Desc"][self.language]
		case "live2d":
			self.routes[i].Name = self.languageTable["live2d_Title"][self.language]
			self.routes[i].Desc = self.languageTable["live2d_Desc"][self.language]
		case "custom":
			self.routes[i].Name = self.languageTable["custom_Title"][self.language]
			self.routes[i].Desc = self.languageTable["custom_Desc"][self.language]
		case "info":
			self.routes[i].Name = self.languageTable["info_Title"][self.language]
			self.routes[i].Desc = self.languageTable["info_Desc"][self.language]
		}
	}

	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Rigid(self.layout_TopBox),
					layout.Flexed(1, func(gtx C) D { return D{Size: gtx.Constraints.Min} }),
				)
			}),
			layout.Flexed(1, self.layout_Options),
		)
	})
}

// 不同的界面
func (self *SettingPage) layout_Options(gtx C) D {
	var offsetY int
	self.placeholderHeight_cur = easing(gtx, self.placeholderHeight_cur, self.placeholderHeight_tar, 0.2)
	offsetY = int(self.placeholderHeight_cur)
	trans := op.Offset(image.Pt(0, offsetY)).Push(gtx.Ops)
	dism := func(gtx C) D {
		if self.pageName == "main" {
			return self.layout_MainPage(gtx)
		}
		page, ok := self.childPages[self.pageName]
		if !ok {
			// 如果 map 里没有，去 routes 找对应的工厂函数
			for _, r := range self.routes {
				if r.Key == self.pageName {
					page = r.Create()
					self.childPages[r.Key] = page
					break
				}
			}
		}

		if page != nil {
			return page.Update(gtx, self)
		}
		return D{}
	}(gtx)
	trans.Pop()
	return dism
}

// 跳转到对应的子界面
func (self *SettingPage) layout_MainPage(gtx C) D {
	var items []layout.Widget
	for _, route := range self.routes {
		r := route
		items = append(items, func(gtx C) D {
			if route.Btn.Clicked(gtx) {
				self.pageName = route.Key
				self.placeholderHeight_cur = float64(gtx.Dp(200))
			}
			return self.build_DescBtn(gtx, r.Icon, r.Name, r.Desc, r.Btn)
		})
	}

	list := material.List(self.style.theme, &self.mainList)
	list.ScrollbarStyle.Indicator.Color.A = 64
	list.ScrollbarStyle.Indicator.MinorWidth = unit.Dp(4)

	self.mainList.Axis = layout.Vertical
	self.mainList.Alignment = layout.Middle
	return list.Layout(gtx, len(items), func(gtx C, i int) D {
		return layout.Inset{Bottom: unit.Dp(10), Right: unit.Dp(8)}.Layout(gtx, items[i])
	})
}

// 顶部栏更新 -----------------------------------------------------------------------------
func (self *SettingPage) update_TopBarWidgets() []layout.Widget {
	widgets := []layout.Widget{
		func(gtx C) D {
			return self.build_TopBarBtn(gtx, self.MainPage.trunk.LanguageTable.SettingPage_Title, self.jumpBtn["main"])
		},
	}

	if self.pageName != "main" {
		if currentPage, ok := self.childPages[self.pageName]; ok {
			widgets = append(widgets, func(gtx C) D {
				return self.build_TopBarBtn(gtx, currentPage.Title(), self.jumpBtn[self.pageName])
			})
		}
	}
	return widgets
}

// 顶部导航栏 ---------------------------------------------------------------------------
func (self *SettingPage) layout_TopBox(gtx C) D {
	if self.jumpBtn["main"].Clicked(gtx) && self.pageName != "main" {
		self.pageName = "main"
		// self.placeholderHeight_cur = float64(gtx.Constraints.Max.Y) / 2
		self.placeholderHeight_cur = float64(unit.Dp(200))
		gtx.Execute(op.InvalidateCmd{})
	}

	btnWidgets := self.update_TopBarWidgets()
	children := []layout.FlexChild{}
	for index, ele := range btnWidgets {
		if index > 0 {
			children = append(children, layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: unit.Dp(4), Left: unit.Dp(4)}.Layout(gtx, func(gtx C) D { return material.Label(self.style.theme, unit.Sp(16), "\uf08f").Layout(gtx) })
			}))
		}
		children = append(children, layout.Rigid(ele))
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		children...,
	)
}

// 创建导航栏按钮
func (self *SettingPage) build_TopBarBtn(gtx C, desc string, clickable *widget.Clickable) D {
	label := material.Label(self.style.theme, unit.Sp(30), desc)
	btn := material.Button(self.style.theme, clickable, "")
	btn.Background = color.NRGBA{A: 0}
	gtx.Constraints.Min.Y = gtx.Dp(40)
	gtx.Constraints.Max.Y = gtx.Dp(40)
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return btn.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Center.Layout(gtx,
				func(gtx C) D {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx C) D { gtx.Constraints.Max.X = gtx.Dp(5); return D{Size: gtx.Constraints.Max} }),
						layout.Rigid(label.Layout),
						layout.Rigid(func(gtx C) D { gtx.Constraints.Max.X = gtx.Dp(5); return D{Size: gtx.Constraints.Max} }),
					)
				},
			)
		}),
	)
}

// 描述按钮
func (self *SettingPage) build_DescBtn(gtx C, icon string, title string, desc string, btn *widget.Clickable) D {

	button := material.Button(self.style.theme, btn, "")
	button.Background = self.style.darkmode.currentColor.IdleBg
	button.Color = self.style.darkmode.currentColor.IdleFg
	button.CornerRadius = unit.Dp(4)

	gtx.Constraints.Min.Y = gtx.Dp(60)
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return button.Layout(gtx)
		}),
		layout.Expanded(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, Flexer()),
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return material.Label(self.style.theme, unit.Sp(16), "\uf08f").Layout(gtx)
					})
				}),
				layout.Rigid(Spacer(gtx.Dp(16), 0)),
			)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				// 图标
				layout.Rigid(Spacer(gtx.Dp(15), gtx.Dp(60))),
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return material.Label(self.style.theme, unit.Sp(24), icon).Layout(gtx)
					})
				}),
				layout.Rigid(Spacer(gtx.Dp(15), gtx.Dp(60))),
				// 标题和描述
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: unit.Dp(40), Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(material.Label(self.style.theme, unit.Sp(22), title).Layout),
							layout.Rigid(material.Label(self.style.theme, unit.Sp(14), desc).Layout),
						)
					})
				}),
				layout.Flexed(1, Flexer()),
			)
		}),
	)
}

// API
func (self *SettingPage) API() map[string]any {
	return map[string]any{}
}

// 语言
func (self *SettingPage) Language() LangPack {
	l := LangPack{
		"voice_Title": LangZone{
			"Chinese":  "语音",
			"English":  "Voice",
			"Japanese": "音声",
		},
		"voice_Desc": LangZone{
			"Chinese":  "启用语音,添加语音,选择语音",
			"English":  "Enable voice, Add voice, Select voice",
			"Japanese": "音声を有効にする,音声を追加する,音声を選択する",
		},
		"model_Title": LangZone{
			"Chinese":  "AI模型",
			"English":  "AI Model",
			"Japanese": "AIモデル",
		},
		"model_Desc": LangZone{
			"Chinese":  "添加模型,加载模型",
			"English":  "Add model, Load model",
			"Japanese": "モデルを追加する,モデルを読み込む",
		},
		"character_Title": LangZone{
			"Chinese":  "人设",
			"English":  "Character",
			"Japanese": "キャラクター設定",
		},
		"character_Desc": LangZone{
			"Chinese":  "添加人设,选择人设",
			"English":  "Add character, Select character",
			"Japanese": "キャラクターを追加する,キャラクターを選択する",
		},
		"live2d_Title": LangZone{
			"Chinese":  "Live2D",
			"English":  "Live2D",
			"Japanese": "Live2D",
		},
		"live2d_Desc": LangZone{
			"Chinese":  "添加模型,选择模型",
			"English":  "Add model, Select model",
			"Japanese": "モデルを追加する,モデルを選択する",
		},
		"custom_Title": LangZone{
			"Chinese":  "自定义",
			"English":  "Custom",
			"Japanese": "カスタム",
		},
		"custom_Desc": LangZone{
			"Chinese":  "语言,配色",
			"English":  "Language, Color theme",
			"Japanese": "言語,配色",
		},
		"info_Title": LangZone{
			"Chinese":  "关于",
			"English":  "About",
			"Japanese": "情報",
		},
		"info_Desc": LangZone{
			"Chinese":  "版本,开发者,鸣谢",
			"English":  "Version, Developer, Credits",
			"Japanese": "バージョン,開発者,クレジット",
		},
	}
	return l
}
