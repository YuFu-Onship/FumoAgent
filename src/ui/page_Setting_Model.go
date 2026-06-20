package ui

import (
	"image"
	"image/color"
	"myapp/src/config"

	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// 模型设置 -----------------------------------------------------------------------
type ModelSetting struct {
	style    Style
	mainList widget.List

	expandPanelList      widget.List
	newModel_expandPanel ExpandPanel

	newModelSaveBtn    widget.Clickable
	newModelRestoreBtn widget.Clickable

	urlEditor   *widget.Editor
	modelEditor *widget.Editor
	editor_key  *widget.Editor
	titleEditor *widget.Editor

	modelEnum widget.Enum
	modelList []ModelPresetPanel

	mainSettingPage *SettingPage
}

func New_ModelSetting(style Style) *ModelSetting {
	self := ModelSetting{
		style:    style,
		mainList: widget.List{},

		expandPanelList: widget.List{},

		newModel_expandPanel: *New_ExpandPanel(style, 265),

		titleEditor:        &widget.Editor{SingleLine: true},
		urlEditor:          &widget.Editor{SingleLine: true},
		modelEditor:        &widget.Editor{SingleLine: true},
		editor_key:         &widget.Editor{SingleLine: true},
		newModelSaveBtn:    widget.Clickable{},
		newModelRestoreBtn: widget.Clickable{},

		modelEnum: widget.Enum{Value: config.API_MODEL_GetCurrentModelTitle()},
	}

	mf := config.Config_Model{}
	list, _ := mf.Get_Preset()

	self.modelList = []ModelPresetPanel{}
	for _, p := range list {
		self.modelList = append(self.modelList, *New_ModelPresetPanel(style, mf.ConvertTo_Preset(p)))
	}

	return &self
}

func (self *ModelSetting) Default() {
	self.titleEditor.SetText("")
	self.urlEditor.SetText("")
	self.modelEditor.SetText("")
	self.editor_key.SetText("")
	self.newModel_expandPanel.API_SetPanelState(false)
	for _, e := range self.modelList {
		e.Default()
	}
}
func (self *ModelSetting) Title() string {
	return self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_Title
}
func (self *ModelSetting) Update(gtx C, mainSettingPage *SettingPage) D {
	self.mainSettingPage = mainSettingPage
	listElements := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Label(self.style.theme, unit.Sp(16), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Title).Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D { size := gtx.Constraints.Min; return D{Size: size} }),
			)
		},
		func(gtx C) D { return self.layout_NewModel(gtx) },
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Label(self.style.theme, unit.Sp(16), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_ChoiceAPI_Title).Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D { size := gtx.Constraints.Min; return D{Size: size} }),
			)
		},
	}

	for i := range self.modelList {
		modelPtr := &self.modelList[i]
		listElements = append(listElements, func(gtx C) D { return modelPtr.Update(gtx, self, &self.modelEnum) })
	}

	if self.modelEnum.Update(gtx) {
	}

	self.mainList.Axis = layout.Vertical
	self.mainList.Alignment = layout.Middle
	list := material.List(self.style.theme, &self.mainList)
	list.ScrollbarStyle.Indicator.Color.A = 64
	list.ScrollbarStyle.Indicator.MinorWidth = unit.Dp(4)
	return list.Layout(gtx, len(listElements), func(gtx C, index int) D {
		return layout.Inset{Right: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx C) D { return listElements[index](gtx) })
	})
}

// 布局 -----------------------------------------------------------------------------
func (self *ModelSetting) layout_NewModel(gtx C) D {
	// 恢复按钮
	if self.newModelRestoreBtn.Clicked(gtx) {
		self.titleEditor.SetText("")
		self.urlEditor.SetText("")
		self.modelEditor.SetText("")
		self.editor_key.SetText("")
	}

	// 保存按钮
	if self.newModelSaveBtn.Clicked(gtx) {
		name := strings.Trim(self.titleEditor.Text(), " \n\t")
		url := self.urlEditor.Text()
		model := self.modelEditor.Text()
		key := self.editor_key.Text()

		result := true
		result = (self.tool_FormatString(name) != "") && result
		result = (self.tool_FormatString(url) != "") && result
		result = (self.tool_FormatString(model) != "") && result
		result = (self.tool_FormatString(key) != "") && result

		if result {
			mf := config.Config_Model{}
			mf.Append_Preset(config.Model_Preset{
				Name:    name,
				Default: false,
				Url:     url,
				Model:   model,
				Key:     key,
			})
			self.tool_ReloadModelDetail()
		}
	}

	// 展开按钮样式
	headFunc := func(gtx C) D {
		icon_label := material.Label(self.style.theme, unit.Sp(20), "\uf8aa")
		text_label := material.Label(self.style.theme, unit.Sp(18), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Btn)
		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(Spacer(gtx.Dp(15), 0)),
				layout.Rigid(icon_label.Layout),
				layout.Rigid(Spacer(gtx.Dp(5), 0)),
				layout.Rigid(text_label.Layout),
				layout.Flexed(1, Flexer()),
			)
		})
	}

	// 保存按钮样式
	saveBtnFunc := func(gtx C) layout.Widget {
		saveBtnStyle := func(gtx C) D {
			icon_label := material.Label(self.style.theme, unit.Sp(20), "\ue74e")
			text_label := material.Label(self.style.theme, unit.Sp(18), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Save)
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(Spacer(gtx.Dp(15), 0)),
					layout.Rigid(icon_label.Layout),
					layout.Rigid(Spacer(gtx.Dp(5), 0)),
					layout.Rigid(text_label.Layout),
					layout.Flexed(1, Flexer()),
				)
			})
		}
		return self.build_Button(
			saveBtnStyle,
			&self.newModelSaveBtn,
			layout.Inset{},
			color.NRGBA{A: 0},
			gtx.Constraints.Max.X, gtx.Dp(40),
		)
	}

	// 恢复按钮样式
	restoreBtnFunc := func(gtx C) layout.Widget {
		restoreBtnStyle := func(gtx C) D {
			icon_label := material.Label(self.style.theme, unit.Sp(20), "\ue777")
			text_label := material.Label(self.style.theme, unit.Sp(18), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Restore)
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(Spacer(gtx.Dp(15), 0)),
					layout.Rigid(icon_label.Layout),
					layout.Rigid(Spacer(gtx.Dp(5), 0)),
					layout.Rigid(text_label.Layout),
					layout.Flexed(1, func(gtx C) D { return D{Size: gtx.Constraints.Max} }),
				)
			})
		}
		return self.build_Button(
			restoreBtnStyle,
			&self.newModelRestoreBtn,
			layout.Inset{},
			color.NRGBA{A: 0},
			gtx.Constraints.Max.X, gtx.Dp(40),
		)
	}
	// 按钮栏
	btnFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D { return saveBtnFunc(gtx)(gtx) }),
			layout.Flexed(1, func(gtx C) D { return restoreBtnFunc(gtx)(gtx) }),
		)
	}

	// 标题栏
	titleEditorFunc := func(gtx C) D {
		return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx C) D {
			return BorderBox(
				self.style,
				self.style.darkmode.currentColor.Fg,
				self.style.darkmode.currentColor.IdleBg,
				self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Name,
			).Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D { return material.Editor(self.style.theme, self.titleEditor, "eg: OPEN-AI").Layout(gtx) }),
				)
			})
		})
	}
	// 网址栏
	urlEditorFunc := func(gtx C) D {
		return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx C) D {
			return BorderBox(
				self.style,
				self.style.darkmode.currentColor.Fg,
				self.style.darkmode.currentColor.IdleBg,
				self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Https,
			).Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {
						return material.Editor(self.style.theme, self.urlEditor, "eg: https://api.openai.com/v1/chat/completions").Layout(gtx)
					}),
				)
			})
		})
	}
	// 模型栏
	modelEditorFunc := func(gtx C) D {
		return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx C) D {
			return BorderBox(
				self.style,
				self.style.darkmode.currentColor.Fg,
				self.style.darkmode.currentColor.IdleBg,
				self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Model,
			).Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D { return material.Editor(self.style.theme, self.modelEditor, "eg: gpt-5").Layout(gtx) }),
				)
			})
		})
	}
	// 密钥栏
	self.editor_key.Mask = '*'
	editor_keyFunc := func(gtx C) D {
		return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx C) D {
			return BorderBox(
				self.style,
				self.style.darkmode.currentColor.Fg,
				self.style.darkmode.currentColor.IdleBg,
				self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Key,
			).Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {
						return material.Editor(self.style.theme, self.editor_key, "eg: sk-xxxxxxxxxxxxxxxxxxxxxxxx").Layout(gtx)
					}),
				)
			})
		})
	}

	// 展开栏
	panelFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D { return btnFunc(gtx) }),
			layout.Rigid(func(gtx C) D { return titleEditorFunc(gtx) }),
			layout.Rigid(func(gtx C) D { return urlEditorFunc(gtx) }),
			layout.Rigid(func(gtx C) D { return modelEditorFunc(gtx) }),
			layout.Rigid(func(gtx C) D { return editor_keyFunc(gtx) }),
		)
	}

	d, _ := self.newModel_expandPanel.Update(gtx, headFunc, panelFunc)
	return d
}

// 创建高级控件 ---------------------------------------------------------------------
// 创建按钮
func (self *ModelSetting) build_Button(
	stackedWidget layout.Widget,
	clickable *widget.Clickable,
	inset layout.Inset,
	c color.NRGBA,
	width int, height int,
) layout.Widget {
	return func(gtx C) D {
		btn := material.Button(self.style.theme, clickable, "")
		btn.Background = c
		return inset.Layout(gtx, func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					gtx.Constraints.Min.X = width
					gtx.Constraints.Min.Y = height
					gtx.Constraints.Max.Y = height
					return btn.Layout(gtx)
				}),
				layout.Stacked(func(gtx C) D {
					gtx.Constraints.Min.X = width
					gtx.Constraints.Min.Y = height
					gtx.Constraints.Max.Y = height
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.End,
					}.Layout(gtx,
						layout.Rigid(stackedWidget),
					)
				}),
			)
		})
	}
}

// 创建输入栏
func (self *ModelSetting) build_Editor(gtx C, editorWidget *widget.Editor, desc string) D {
	editorWidget.LineHeightScale = float32(14)
	editor := material.Editor(self.style.theme, editorWidget, desc)
	return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(Spacer(gtx.Dp(15), 0)),
				layout.Rigid(Spacer(gtx.Dp(15), 0)),
				layout.Flexed(1, func(gtx C) D {
					return widget.Border{Width: unit.Dp(2), CornerRadius: unit.Dp(4), Color: self.style.theme.Fg}.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(4).Layout(gtx, func(gtx C) D {
							if gtx.Constraints.Max.Y >= gtx.Dp(50) {
								gtx.Constraints.Max.Y = gtx.Dp(50)
							}
							return editor.Layout(gtx)
						})
					})
				}),
				layout.Rigid(Spacer(gtx.Dp(15), 0)),
			)
		})
	})
}

// 格式化字符串
func (self *ModelSetting) tool_FormatString(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\r", "")
	return text
}

// 检测字符串
func (self *ModelSetting) tool_CheckText(text string) bool {
	if text == "" {
		return false
	} else {
		return true
	}
}

// 重载modelList
func (self *ModelSetting) tool_ReloadModelDetail() {
	self.modelList = []ModelPresetPanel{}
	mf := config.Config_Model{}
	list, _ := mf.Get_Preset()

	self.modelList = []ModelPresetPanel{}
	for _, p := range list {
		self.modelList = append(self.modelList, *New_ModelPresetPanel(self.style, mf.ConvertTo_Preset(p)))
	}
}

// 删除model
func (self *ModelSetting) tool_DeleteModelDetail(name string) {
	mf := config.Config_Model{}
	mf.Delete_Preset(name)

	self.tool_ReloadModelDetail()

	if name == self.modelEnum.Value {
		var name string = ""
		if len(self.modelList) > 0 {
			name = self.modelList[0].API_GetName()
		}
		self.modelEnum.Value = name
		config.API_MODEL_SetCurrentModelTitle(name)
	}
}

// 得到预设参数
// func (self *ModelSetting) api_GetPreset() config.Model_Preset {
// 	return config.Model_Preset{
// 		Name: self.titleEditor.Text(),
// 	}
// }

// 初始化

// 模型详情控件 -------------------------------------------------------------------
type ModelPresetPanel struct {
	style     Style
	isDefault bool

	name  string
	url   string
	model string
	key   string

	editor_url   widget.Editor
	editor_model widget.Editor
	editor_key   widget.Editor

	checkBtn  widget.Clickable
	checkIcon string

	deleteBtn  widget.Clickable
	confirmBtn widget.Clickable

	parent *ModelSetting

	preset_expandPanel ExpandPanel
}

func New_ModelPresetPanel(style Style, modelInfo config.Model_Preset) *ModelPresetPanel {
	input := widget.Editor{SingleLine: true, Mask: '*'}
	input.SetText(modelInfo.Key)
	self := ModelPresetPanel{
		style: style,

		name:      modelInfo.Name,
		isDefault: modelInfo.Default,
		url:       modelInfo.Url,
		model:     modelInfo.Model,
		key:       modelInfo.Key,

		editor_url:   widget.Editor{},
		editor_model: widget.Editor{},
		editor_key:   input,

		preset_expandPanel: *New_ExpandPanel(style, 190),

		checkBtn:  widget.Clickable{},
		checkIcon: "\uE739",

		deleteBtn:  widget.Clickable{},
		confirmBtn: widget.Clickable{},
	}

	self.editor_key.SingleLine = true

	self.editor_model.SetText(modelInfo.Model)
	self.editor_model.SingleLine = true
	self.editor_model.ReadOnly = true

	self.editor_url.SetText(modelInfo.Url)
	self.editor_url.SingleLine = true
	self.editor_url.ReadOnly = true

	return &self
}

func (self *ModelPresetPanel) Default() {
	self.editor_key.SetText("")
}

// 更新
func (self *ModelPresetPanel) Update(gtx C, parent *ModelSetting, enum *widget.Enum) D {
	self.parent = parent

	// 确认按钮
	if self.confirmBtn.Clicked(gtx) {
		self.key = self.editor_key.Text()
		detail := config.Model_Preset{
			Name:    self.name,
			Default: self.isDefault,
			Url:     self.url,
			Model:   self.model,
			Key:     self.key,
		}

		mf := config.Config_Model{}
		mf.Change_Preset(detail)
	}

	// 删除按钮
	if self.deleteBtn.Clicked(gtx) {
		if self.isDefault == false {
			parent.tool_DeleteModelDetail(self.name)
		}
	}

	// 点选按钮
	if self.checkBtn.Clicked(gtx) {
		self.parent.modelEnum.Value = self.name
		config.API_MODEL_SetCurrentModelTitle(self.name)
		self.parent.mainSettingPage.MainPage.trunk.Handler_AI.API_SetAiModelApi(config.Model_Preset{
			Name:    self.name,
			Default: self.isDefault,
			Url:     self.url,
			Model:   self.model,
			Key:     self.key,
		})

	}

	if self.name == enum.Value {
		self.checkIcon = "\ue73d"
	} else {
		self.checkIcon = "\uE739"
	}

	checkBtnStyle := func(gtx C) D {
		gtx.Constraints.Max = image.Pt(gtx.Dp(20), gtx.Dp(20))
		gtx.Constraints.Min = image.Pt(gtx.Dp(20), gtx.Dp(20))

		icon_label := material.Label(self.style.theme, unit.Sp(20), self.checkIcon)
		btn := material.Button(self.style.theme, &self.checkBtn, "")
		btn.Background = color.NRGBA{A: 0}
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return btn.Layout(gtx)
			}),
			layout.Stacked(func(gtx C) D {
				return icon_label.Layout(gtx)
			}),
		)
	}

	headFunc := func(gtx C) D {
		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(Spacer(gtx.Dp(15), 0)),
				layout.Rigid(func(gtx C) D { return checkBtnStyle(gtx) }),
				layout.Rigid(Spacer(gtx.Dp(5), 0)),
				layout.Rigid(material.Label(self.style.theme, unit.Sp(20), self.name).Layout),
				layout.Flexed(1, Flexer()),
			)
		})
	}

	panelFunc := func(gtx C) D {
		return self.layout_Panel(gtx)
	}

	d, _ := self.preset_expandPanel.Update(gtx, headFunc, panelFunc)
	return d
}

// 详细栏
func (self *ModelPresetPanel) layout_Panel(gtx C) D {
	confirmBtnFunc := func(gtx C) D {
		return layout.Inset{Right: unit.Dp(2), Left: unit.Dp(2), Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx C) D {
			return self.build_Button(
				func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(Spacer(gtx.Dp(5), 0)),
							layout.Rigid(func(gtx C) D { return material.Label(self.style.theme, unit.Sp(18), "\uE74E ").Layout(gtx) }),
							layout.Rigid(func(gtx C) D {
								return material.Label(self.style.theme, unit.Sp(18), self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_ChoiceAPI_Confirm).Layout(gtx)
							}),
							layout.Rigid(Spacer(gtx.Dp(5), 0)),
						)
					})
				},
				&self.confirmBtn,
				layout.Inset{},
				color.NRGBA{A: 0},
			)(gtx)
		})
	}
	deleteBtnFunc := func(gtx C) D {
		return layout.Inset{Right: unit.Dp(2), Left: unit.Dp(2), Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx C) D {
			return self.build_Button(
				func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(Spacer(gtx.Dp(5), 0)),
							layout.Rigid(func(gtx C) D { return material.Label(self.style.theme, unit.Sp(18), "\uE74D ").Layout(gtx) }),
							layout.Rigid(func(gtx C) D {
								return material.Label(self.style.theme, unit.Sp(18), self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_ChoiceAPI_Delete).Layout(gtx)
							}),
							layout.Rigid(Spacer(gtx.Dp(5), 0)),
						)
					})
				},
				&self.deleteBtn,
				layout.Inset{},
				color.NRGBA{A: 0},
			)(gtx)
		})
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
						return BorderBox(
							self.style,
							self.style.darkmode.currentColor.Fg,
							self.style.darkmode.currentColor.IdleBg,
							self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Https,
						).Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Flexed(1, func(gtx C) D { return material.Editor(self.style.theme, &self.editor_url, "").Layout(gtx) }),
							)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
						return BorderBox(
							self.style,
							self.style.darkmode.currentColor.Fg,
							self.style.darkmode.currentColor.IdleBg,
							self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Model,
						).Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Flexed(1, func(gtx C) D { return material.Editor(self.style.theme, &self.editor_model, "").Layout(gtx) }),
							)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
						return BorderBox(
							self.style,
							self.style.darkmode.currentColor.Fg,
							self.style.darkmode.currentColor.IdleBg,
							self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Model_NewAPI_Key,
						).Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Flexed(1, func(gtx C) D { return material.Editor(self.style.theme, &self.editor_key, "").Layout(gtx) }),
							)
						})
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D { return confirmBtnFunc(gtx) }),
					layout.Rigid(func(gtx C) D { return deleteBtnFunc(gtx) }),
				)
			})
		}),
	)
}

// 信息栏
func (self *ModelPresetPanel) build_infoLabel(gtx C, title string, content string) D {
	titleLabel := material.Label(self.style.theme, unit.Sp(16), title)
	contLabel := material.Label(self.style.theme, unit.Sp(16), content)
	return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40), Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if gtx.Constraints.Max.X >= gtx.Dp(50) {
					gtx.Constraints.Max.X = gtx.Dp(50)
				}
				return titleLabel.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D { return contLabel.Layout(gtx) }),
			layout.Flexed(1, Flexer()),
		)
	})
}

// 编辑栏
func (self *ModelPresetPanel) build_Editor(gtx C, title string, editorWidget *widget.Editor) D {
	titleLabel := material.Label(self.style.theme, unit.Sp(16), title)
	editorLine := material.Editor(self.style.theme, editorWidget, "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	return layout.Inset{Left: unit.Dp(40), Right: unit.Dp(40), Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if gtx.Constraints.Max.X >= gtx.Dp(50) {
					gtx.Constraints.Max.X = gtx.Dp(50)
				}
				return titleLabel.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D { return editorLine.Layout(gtx) }),
			layout.Flexed(1, Flexer()),
		)
	})
}

// 按钮控件
func (self *ModelPresetPanel) build_Button(stackedWidget layout.Widget,
	clickable *widget.Clickable,
	inset layout.Inset,
	c color.NRGBA,
) layout.Widget {
	return func(gtx C) D {
		btn := material.Button(self.style.theme, clickable, "")
		btn.Background = c
		return inset.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Dp(40)
			gtx.Constraints.Max.Y = gtx.Dp(40)
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return btn.Layout(gtx)
				}),
				layout.Stacked(func(gtx C) D {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Rigid(stackedWidget),
					)
				}),
			)
		})
	}
}

// 得到当前详情名称
func (self *ModelPresetPanel) API_GetName() string {
	return self.name
}
