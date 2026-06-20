package ui

import (
	"image"
	"image/color"
	"myapp/src/config"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// CharacterInfo 模拟配置数据结构
type CharacterInfo struct {
	Name   string
	Prompt string
}

// 人设设置 ----------------------------------------------------------------------------------------
type CharacterSetting struct {
	style    Style
	mainList widget.List

	// 新建人设相关
	newCharBtn          widget.Clickable
	newCharSaveBtn      widget.Clickable
	newCharRestoreBtn   widget.Clickable
	newChar_expandPanel ExpandPanel

	nameEditor   *widget.Editor
	promptEditor *widget.Editor

	charEnum widget.Enum
	charList []CharacterDetailPanel

	mainSettingPage *SettingPage
}

func New_CharacterSetting(style Style) *CharacterSetting {
	self := CharacterSetting{
		style:    style,
		mainList: widget.List{},

		newCharBtn:          widget.Clickable{},
		newChar_expandPanel: *New_ExpandPanel(style, 350),

		nameEditor:   &widget.Editor{SingleLine: true},
		promptEditor: &widget.Editor{},

		newCharSaveBtn:    widget.Clickable{},
		newCharRestoreBtn: widget.Clickable{},

		charEnum: widget.Enum{},
	}
	cc := config.Config_Character{}
	self.charEnum.Value = cc.Get_CurName()
	self.ReloadInfo()
	return &self
}

func (self *CharacterSetting) Default() {
	self.newChar_expandPanel.API_SetPanelState(false)
}
func (self *CharacterSetting) Title() string {
	return self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Character_Title
}
func (self *CharacterSetting) Update(gtx C, mainSettingPage *SettingPage) D {
	self.mainSettingPage = mainSettingPage

	listElements := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Label(self.style.theme, unit.Sp(16), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Character_NewCharacter_Title).Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D { size := gtx.Constraints.Min; return D{Size: size} }),
			)
		},
		func(gtx C) D { return self.layout_NewCharacter(gtx) },
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Label(self.style.theme, unit.Sp(16), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Character_NewCharacter_ChoiceCharacter).Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D { size := gtx.Constraints.Min; return D{Size: size} }),
			)
		},
	}

	// 动态添加角色列表
	for i := range self.charList {
		charPtr := &self.charList[i]
		listElements = append(listElements, func(gtx C) D {
			return charPtr.Update(gtx, self, &self.charEnum)
		})
	}

	self.mainList.Axis = layout.Vertical

	list := material.List(self.style.theme, &self.mainList)
	list.ScrollbarStyle.Indicator.Color.A = 64
	list.ScrollbarStyle.Indicator.MinorWidth = unit.Dp(4)

	if self.charEnum.Update(gtx) {
	}

	return list.Layout(gtx, len(listElements), func(gtx C, index int) D {
		return layout.Inset{Right: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx C) D { return listElements[index](gtx) })
	})
}

// 布局：新建人设 ---------------------------------------------------------------------
func (self *CharacterSetting) layout_NewCharacter(gtx C) D {
	// 恢复按钮
	if self.newCharRestoreBtn.Clicked(gtx) {
		self.nameEditor.SetText("")
		self.promptEditor.SetText("")
	}
	if self.newCharSaveBtn.Clicked(gtx) {
		self.CreateNewChar()
	}

	// 保存按钮样式
	saveBtnFunc := func(gtx C) layout.Widget {
		saveBtnStyle := func(gtx C) D {
			icon_label := material.Label(self.style.theme, unit.Sp(20), "\ue74e")
			text_label := material.Label(self.style.theme, unit.Sp(18), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Character_NewCharacter_Save)
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
			&self.newCharSaveBtn,
			layout.Inset{},
			color.NRGBA{A: 0},
			gtx.Constraints.Max.X, gtx.Dp(40),
		)
	}

	// 恢复按钮样式
	restoreBtnFunc := func(gtx C) layout.Widget {
		restoreBtnStyle := func(gtx C) D {
			icon_label := material.Label(self.style.theme, unit.Sp(20), "\ue777")
			text_label := material.Label(self.style.theme, unit.Sp(18), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Character_NewCharacter_Restore)
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
			restoreBtnStyle,
			&self.newCharRestoreBtn,
			layout.Inset{},
			color.NRGBA{A: 0},
			gtx.Constraints.Max.X, gtx.Dp(40),
		)
	}

	// 按钮栏
	// btnFunc := func(gtx C) D {
	// 	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
	// 		layout.Flexed(1, func(gtx C) D { return saveBtnFunc(gtx)(gtx) }),
	// 		layout.Flexed(1, func(gtx C) D { return restoreBtnFunc(gtx)(gtx) }),
	// 	)
	// }
	// 标题编辑框样式
	nameEditorFunc := func(gtx C) D {
		return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
			return BorderBox(
				self.style,
				self.style.darkmode.currentColor.Fg,
				self.style.darkmode.currentColor.IdleBg,
				self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Character_NewCharacter_Name,
			).Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {
						return material.Editor(self.style.theme, self.nameEditor, "").Layout(gtx)
					}),
				)
			})
		})
	}
	// 内容编辑框样式
	contEditorFunc := func(gtx C) D {
		return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
			return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
				return BorderBox(
					self.style,
					self.style.darkmode.currentColor.Fg,
					self.style.darkmode.currentColor.IdleBg,
					self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Character_NewCharacter_Prompt,
				).Layout(gtx, func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							gtx.Constraints.Min.Y = 220
							gtx.Constraints.Max.Y = 220
							return material.Editor(self.style.theme, self.promptEditor, "").Layout(gtx)
						}),
					)
				})
			})
		})
	}

	headFunc := func(gtx C) D {
		icon_label := material.Label(self.style.theme, unit.Sp(20), "\uf8aa")
		text_label := material.Label(self.style.theme, unit.Sp(18), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Character_NewCharacter_Btn)
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
	panelFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, saveBtnFunc(gtx)),
					layout.Flexed(1, restoreBtnFunc(gtx)),
				)
			}),
			layout.Rigid(nameEditorFunc),
			layout.Rigid(contEditorFunc),
		)
	}

	d, _ := self.newChar_expandPanel.Update(gtx, headFunc, panelFunc)
	return d
}

// 重新加载角色列表
func (self *CharacterSetting) ReloadInfo() {
	cc := config.Config_Character{}
	preset := cc.Get_Preset()

	self.charList = []CharacterDetailPanel{}
	for _, p := range preset {
		c := New_CharacterDetailPanel(self.style, []string{p[0], "false", p[1]})
		self.charList = append(self.charList, *c)
	}
}

// 设定当前角色
func (self *CharacterSetting) SetCurCharacter(name string, desc string) {
	self.charEnum.Value = name
	trunk := self.mainSettingPage.MainPage.trunk
	trunk.CharacterTitle = name
	trunk.CharacterDesc = desc

	cc := config.Config_Character{}
	cc.Set_CurName(name)
}

// 删除角色
func (self *CharacterSetting) DeleteCharacter(name string) {
	cc := config.Config_Character{}
	cc.Delete_Character(name)
	self.ReloadInfo()
}

// 保存角色设定
func (self *CharacterSetting) CreateNewChar() {
	name := strings.TrimSpace(self.nameEditor.Text())
	desc := strings.TrimSpace(self.promptEditor.Text())
	if name != "" && desc != "" {
		cc := config.Config_Character{}
		cc.Save_Character(name, desc)
		self.ReloadInfo()
	}
}

// 人设详情控件 -------------------------------------------------------------------
type CharacterDetailPanel struct {
	style Style
	info  []string

	title     string
	isDefault string

	isExpand           bool
	panelHeight_tar    float32
	panelHeight_cur    float32
	preset_expandPanel ExpandPanel

	expandBtn  widget.Clickable
	deleteBtn  widget.Clickable
	checkBtn   widget.Clickable
	promptEdit *widget.Editor

	parent      *CharacterSetting
	confirmBtn  widget.Clickable
	confirmIcon string
}

func New_CharacterDetailPanel(style Style, info []string) *CharacterDetailPanel {
	editor := &widget.Editor{}
	editor.SetText(info[2])
	editor.ReadOnly = true

	c := CharacterDetailPanel{
		style:              style,
		info:               info,
		promptEdit:         editor,
		confirmBtn:         widget.Clickable{},
		title:              info[0],
		isDefault:          info[1],
		preset_expandPanel: *New_ExpandPanel(style, 240),
	}
	return &c
}

func (self *CharacterDetailPanel) Update(gtx C, parent *CharacterSetting, enum *widget.Enum) D {
	self.parent = parent

	// 展开按钮
	if self.expandBtn.Clicked(gtx) {
		self.isExpand = !self.isExpand
		if self.isExpand {
			self.panelHeight_tar = float32(unit.Dp(240))
		} else {
			self.panelHeight_tar = 0
		}
	}
	self.panelHeight_cur = easing32(gtx, self.panelHeight_cur, self.panelHeight_tar, 0.3)
	// 删除按钮
	if self.deleteBtn.Clicked(gtx) {
		if self.isDefault != "true" {
			self.parent.DeleteCharacter(self.title)
		}
	}

	// 点选按钮
	if self.checkBtn.Clicked(gtx) {
		// enum.Value = self.title
		// self.parent.charEnum.Value = self.title
		// config.API_CHARACTER_SetCurrentCharacter(self.title)
		// self.parent.mainSettingPage.MainPage.trunk.CharacterTitle = self.title
		// self.parent.mainSettingPage.MainPage.trunk.CharacterDesc = config.API_CHARACTER_GetCharacterInfo(self.title)[2]
		parent.SetCurCharacter(self.info[0], self.info[2])
	}
	if self.title == enum.Value {
		self.confirmIcon = "\ue73d"
	} else {
		self.confirmIcon = "\uE739"
	}

	checkBtnStyle := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(22), gtx.Sp(24))
		btn := material.Button(self.style.theme, &self.checkBtn, self.confirmIcon)
		btn.TextSize = unit.Sp(22)
		btn.Background.A = 0
		btn.Color = self.style.darkmode.currentColor.Fg
		btn.Inset = layout.UniformInset(0)
		return btn.Layout(gtx)
	}

	headFunc := func(gtx C) D {
		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(Spacer(gtx.Dp(15), 0)),
				layout.Rigid(func(gtx C) D { return checkBtnStyle(gtx) }),
				layout.Rigid(Spacer(gtx.Dp(5), 0)),
				layout.Rigid(material.Label(self.style.theme, unit.Sp(20), self.title).Layout),
				layout.Flexed(1, Flexer()),
			)
		})
	}

	panelFunc := func(gtx C) D {
		return self.layout_CharCont(gtx)
	}

	d, _ := self.preset_expandPanel.Update(gtx, headFunc, panelFunc)
	return d
}

// 人物描述的显示与删除按钮
func (self *CharacterDetailPanel) layout_CharCont(gtx C) D {
	btnStyle := func(gtx C) D {
		icon_label := material.Label(self.style.theme, unit.Sp(20), "\uE74D")
		btn := material.Button(self.style.theme, &self.deleteBtn, "")
		btn.Background = color.NRGBA{A: 0}
		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					gtx.Constraints.Max = image.Pt(gtx.Dp(30), gtx.Dp(30))
					gtx.Constraints.Min = image.Pt(gtx.Dp(30), gtx.Dp(30))
					return btn.Layout(gtx)
				}),
				layout.Stacked(func(gtx C) D {
					return icon_label.Layout(gtx)
				}),
			)
		})
	}

	editor := func(gtx C) D {
		return BorderBox(self.style, self.style.theme.Fg, self.style.darkmode.currentColor.IdleBg, "Prompt").Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Dp(200)
			gtx.Constraints.Max.Y = gtx.Dp(200)
			return material.Editor(self.style.theme, self.promptEdit, "").Layout(gtx)
		})
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(Spacer(gtx.Dp(15), 0)),
		layout.Flexed(1, editor),
		layout.Rigid(btnStyle),
		layout.Rigid(Spacer(gtx.Dp(10), 0)),
	)
}

func (self *CharacterDetailPanel) build_Button(stackedWidget layout.Widget, clickable *widget.Clickable, inset layout.Inset, c color.NRGBA, width int, height int) layout.Widget {
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
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, layout.Rigid(stackedWidget))
				}),
			)
		})
	}
}

// 辅助工具方法 ---------------------------------------------------------------------
func (self *CharacterSetting) layout_BackgroundContainer(gtx C, w layout.Widget) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			paint.FillShape(gtx.Ops, self.style.darkmode.currentColor.IdleBg, clip.RRect{
				Rect: image.Rectangle{Max: gtx.Constraints.Min},
				SE:   gtx.Dp(4),
				SW:   gtx.Dp(4),
				NE:   gtx.Dp(4),
				NW:   gtx.Dp(4),
			}.Op(gtx.Ops))
			return D{}
		}),
		layout.Stacked(func(gtx C) D {
			return widget.Border{
				Width: unit.Dp(2), Color: color.NRGBA{A: 0}, CornerRadius: unit.Dp(4),
			}.Layout(gtx, w)
		}),
	)
}

// 复用你提供的 build_Button
func (self *CharacterSetting) build_Button(stackedWidget layout.Widget, clickable *widget.Clickable, inset layout.Inset, c color.NRGBA, width int, height int) layout.Widget {
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
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, layout.Rigid(stackedWidget))
				}),
			)
		})
	}
}

// 复用你提供的 build_Editor
func (self *CharacterSetting) build_Editor(gtx C, editorWidget *widget.Editor, title string, desc string) D {
	editor := material.Editor(self.style.theme, editorWidget, desc)
	label := material.Label(self.style.theme, unit.Sp(18), title)
	return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(50)
				return label.Layout(gtx)
			}),
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
			layout.Flexed(1, func(gtx C) D {
				return widget.Border{Width: unit.Dp(2), CornerRadius: unit.Dp(4), Color: self.style.theme.Fg}.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
						// 动态调整高度，如果是非单行编辑器则允许更高
						if editorWidget.SingleLine && gtx.Constraints.Max.Y >= gtx.Dp(50) {
							gtx.Constraints.Max.Y = gtx.Dp(50)
						}
						return editor.Layout(gtx)
					})
				})
			}),
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
		)
	})
}
