package ui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"myapp/src/config"

	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SoundSettingNewVoice interface {
	Update(gtx C, parent *SoundSetting) W
	Default(gtx C) bool
	API_GetCoefficient() any
	API_SetLanguage(language string)
}
type SoundSettingPresetBox interface {
	Update(gtx C, parent *SoundSetting) W
	API_SetLanguage(language string)
	API_GetName() string
}

type SoundSetting struct {
	style           Style
	mainSettingPage *SettingPage

	mainList widget.List

	// 新建声音选项
	editor_Title   widget.Editor
	editor_Preview widget.Editor

	// 顶部三键
	btn_SaveVoice widget.Clickable
	btn_PlayVoice widget.Clickable
	btn_Reset     widget.Clickable

	// 中文转片假音
	langConv_switch widget.Bool
	langConv_bool   bool

	// 启用声音
	enableSound_btn   widget.Clickable
	enableSound_state bool

	// 当前声音预设
	preset_Enum widget.Enum
	preset_list []SoundSettingPresetBox

	// 新建语音预设的下拉栏
	newVoice_ExpandPanel *ExpandPanel
	newVoice_Height      int

	// 不同的语音引擎选择
	switchEngine_DropDown *BorderDropDown
	switchEngine_Value    string

	// 参数调节界面
	coe_Aquestalk *CoefficientAquestalk
	coe_OpenJTalk *CoefficientOpenJTalk
	coe_Teto      *CoefficientTeto
	coe_IsDefault bool
	coe_CurPanel  SoundSettingNewVoice

	// 语言
	language      string
	languageTable LangPack
}

func New_SoundSetting(style Style) *SoundSetting {

	self := SoundSetting{
		style:    style,
		mainList: widget.List{},

		// 新建声音
		btn_SaveVoice: widget.Clickable{},
		btn_PlayVoice: widget.Clickable{},
		btn_Reset:     widget.Clickable{},

		// 设置语言转换
		langConv_switch: widget.Bool{Value: true},
		langConv_bool:   false,

		// 编辑栏
		editor_Preview: widget.Editor{SingleLine: true},
		editor_Title:   widget.Editor{SingleLine: true},

		// 当前声音预设
		preset_Enum: widget.Enum{Value: config.API_SOUND_GetCurrentPresetTitle()},

		// 启用声音
		enableSound_btn: widget.Clickable{},

		// 切换引擎
		switchEngine_DropDown: New_BorderDropDown(style, "VoiceEngine", []string{"AquesTalk", "OpenJTalk", "VoiceVox"}, "AquesTalk"),
		switchEngine_Value:    "AquesTalk",

		// 参数复位
		coe_IsDefault: false,

		// 下拉栏
		newVoice_ExpandPanel: New_ExpandPanel(style, 0),
		newVoice_Height:      0,

		// 语言
		language: "Chinese",
	}

	self.coe_CurPanel = New_CoefficientAquestalk(style)
	self.switchEngine_DropDown.execute(self.changeVoiceCoe)
	self.languageTable = self.Language()

	cv := config.Config_Voice{}
	_presetString, _ := cv.Get_CurString()
	self.preset_Enum.Value = _presetString[1]
	self.enableSound_state = cv.Get_EnableState()

	preset_StrList, _ := cv.Get_Preset()
	for _, p := range preset_StrList {
		switch p[0] {
		case "AquesTalk":
			ap, _ := cv.ConvertTo_Aquestalk(p)
			self.preset_list = append(self.preset_list, New_PresetBoxAquesTalk(style, ap))
		case "OpenJTalk":
			op, _ := cv.ConvertTo_OpenJTalk(p)
			self.preset_list = append(self.preset_list, New_PresetBoxOpenJTalk(style, op))
		case "VoiceVox":
			vp, _ := cv.ConvertTo_VoiceVox(p)
			self.preset_list = append(self.preset_list, New_PresetBoxVoiceVox(style, vp))
		default:
		}
	}

	return &self
}
func (self *SoundSetting) Default() {
	self.tool_ReloadPreset()
	self.editor_Preview.SetText("")
	// self.isNewVoiceExpand = false
}

func (self *SoundSetting) Title() string {
	return self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_Title
}

func (self *SoundSetting) Update(gtx C, parent *SettingPage) D {
	self.mainSettingPage = parent

	self.language = self.mainSettingPage.MainPage.trunk.Language
	if !self.langConv_bool {
		self.langConv_bool = true
		if parent.MainPage.trunk.Language == "Chinese" {
			self.langConv_switch.Value = true
		} else {
			self.langConv_switch.Value = false
		}
	}

	// 列表控件, 包含是否启用语音,新建语音, 以及当前预设
	listElements := []layout.Widget{
		// 启用语音标题
		func(gtx C) D {
			return self.build_TitleLabel(
				gtx,
				self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_EnableVoice_Title)
		},
		// 启用语音
		func(gtx C) D { return self.layout_EnableSound(gtx) },

		// 新建语音标题
		func(gtx C) D {
			return self.build_TitleLabel(
				gtx,
				self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CreateVoice_Title)
		},
		// 新建语音
		func(gtx C) D { return self.layout_NewVoice(gtx) },

		// 预设列表标题
		func(gtx C) D {
			return self.build_TitleLabel(
				gtx,
				self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_ChoiceVoice)
		},
	}

	for _, p := range self.preset_list {
		p.API_SetLanguage(self.language)
		listElements = append(listElements, func(gtx C) D { return p.Update(gtx, self)(gtx) })
	}

	// 预设语音
	self.mainList.Axis = layout.Vertical
	self.mainList.Alignment = layout.Middle
	list := material.List(self.style.theme, &self.mainList)
	list.ScrollbarStyle.Indicator.Color.A = 64
	list.ScrollbarStyle.Indicator.MinorWidth = unit.Dp(4)

	if self.preset_Enum.Update(gtx) {
	}

	return list.Layout(gtx, len(listElements), func(gtx C, index int) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx C) D { return listElements[index](gtx) })
			}),
			layout.Rigid(Spacer(0, gtx.Dp(8))),
		)
	})
}

func (self *SoundSetting) layout_SetCurrentVoice(gtx C) D {
	label := material.Label(self.style.theme, unit.Sp(24), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_ChoiceVoice)
	return label.Layout(gtx)
}

// 启用声音 -----------------------------------------------------------------------
func (self *SoundSetting) layout_EnableSound(gtx C) D {
	if self.enableSound_btn.Clicked(gtx) {
		self.enableSound_state = !self.enableSound_state

		trunk := self.mainSettingPage.MainPage.trunk
		trunk.VoiceEnable = self.enableSound_state

		cv := config.Config_Voice{}
		cv.Set_EnableState(self.enableSound_state)
	}

	state_icon := "\uECCA"
	if self.enableSound_state {
		state_icon = "\uECCB"
	}

	// 内存泄露是由于stack内部设置尺寸造成的
	gtx.Constraints.Min.Y = gtx.Dp(40)
	gtx.Constraints.Max.Y = gtx.Dp(40)
	// gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			btn := material.Button(self.style.theme, &self.enableSound_btn, "")
			btn.Background = self.style.darkmode.currentColor.IdleBg
			return btn.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, FlexerY()),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(Spacer(gtx.Dp(15), 0)),
						layout.Rigid(material.Label(self.style.theme, unit.Sp(20), "\uE8AB").Layout),
						layout.Rigid(Spacer(gtx.Dp(5), 0)),
						layout.Rigid(material.Label(self.style.theme, unit.Sp(18), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_EnableVoice_Title).Layout),
						layout.Flexed(1, Flexer()),
						layout.Rigid(material.Label(self.style.theme, unit.Sp(20), state_icon).Layout),
						layout.Rigid(Spacer(gtx.Dp(15), 0)),
					)
				}),
				layout.Flexed(1, FlexerY()),
			)
		}),
	)
}

// 更改当前的声音
func (self *SoundSetting) changeVoiceCoe(title string) {
	if title == self.switchEngine_Value {
		return
	}
	switch title {
	case "AquesTalk":
		self.coe_CurPanel = New_CoefficientAquestalk(self.style)
	case "OpenJTalk":
		self.coe_CurPanel = New_CoefficientOpenJTalk(self.style)
	// case "Teto":
	// 	self.coe_CurPanel = New_CoefficientTeto(self.style)
	case "VoiceVox":
		self.coe_CurPanel = New_CoefficientVoiceVox(self.style)
	}
}

// 新建声音的布局界面 -----------------------------------------------------------------------
func (self *SoundSetting) layout_NewVoice(gtx C) D {
	// 设置参数界面的语言
	self.coe_CurPanel.API_SetLanguage(self.language)
	// 展开栏整体高度
	curTotalHeight := 240
	if self.switchEngine_DropDown.API_GetState() {
		curTotalHeight += 140
	}

	self.switchEngine_Value = self.switchEngine_DropDown.API_GetValue()

	// 不同语音引擎拥有不同的高度
	switch self.switchEngine_Value {
	case "AquesTalk":
		curTotalHeight += 240
	case "OpenJTalk":
		curTotalHeight += 140
	case "VoiceVox":
		curTotalHeight += 290
	case "Teto":
		curTotalHeight += 0
	default:
		curTotalHeight += 0
	}

	self.newVoice_ExpandPanel.API_SetPanelHeight(int(unit.Dp(curTotalHeight)))

	// 参数复位动画
	if self.coe_IsDefault {
		if self.coe_CurPanel.Default(gtx) {
			self.coe_IsDefault = false
		}
	}

	// 保存按钮
	if self.btn_SaveVoice.Clicked(gtx) {
		name := strings.TrimSpace(self.editor_Title.Text())
		name = strings.ReplaceAll(name, "\n", "")
		name = strings.ReplaceAll(name, "\r", "")
		if name != "" {
			presetData := self.coe_CurPanel.API_GetCoefficient()
			cv := config.Config_Voice{}
			s := cv.ToString_Preset(presetData)
			s[1] = name
			preset, _ := cv.ToPreset_String(s)

			go func() {
				if cv.Save_Preset(preset) {
					self.mainSettingPage.MainPage.notice.Add_Notice(
						1,
						self.languageTable["Notice_Add_Title"][self.language],
						self.languageTable["Notice_Add_Desc"][self.language],
					)
				} else {
					self.mainSettingPage.MainPage.notice.Add_Notice(
						2,
						self.languageTable["Notice_AddRepeat_Title"][self.language],
						self.languageTable["Notice_AddRepeat_Desc"][self.language],
					)
				}
				self.tool_ReloadPreset()
			}()
		}
	}

	// 试听按钮
	if self.btn_PlayVoice.Clicked(gtx) {
		content := self.editor_Preview.Text()

		// 决定是否中文转片假名
		var convCont string
		if self.langConv_switch.Value {
			convCont = self.mainSettingPage.MainPage.trunk.LangConv.ChineseToKatana(content)
		} else {
			convCont = content
		}

		// 不同语音引擎的播放逻辑
		if strings.TrimSpace(content) == "" {
		} else {
			presetData := self.coe_CurPanel.API_GetCoefficient()
			switch data := presetData.(type) {
			default:
			case config.Preset_Aquestalk:
				go func() {
					self.mainSettingPage.MainPage.trunk.Player.ReleaseFile()
					path := config.SOUND_GetVoicePath_Aquestalk(convCont, data)
					self.mainSettingPage.MainPage.trunk.Player.PlayFile(path)
				}()
			case config.Preset_OpenJTalk:
				go func() {
					self.mainSettingPage.MainPage.trunk.Player.ReleaseFile()
					path := config.SOUND_GetVoicePath_OpenJTalk(convCont, data)
					self.mainSettingPage.MainPage.trunk.Player.PlayFile(path)
				}()
			case config.Preset_VoiceVox:
				go func() {
					self.mainSettingPage.MainPage.trunk.Player.ReleaseFile()
					if self.mainSettingPage.MainPage.trunk.VoiceVox.API_GetRunState() {
					} else {
						self.mainSettingPage.MainPage.trunk.VoiceVox.StartServer()
					}

					index := self.mainSettingPage.MainPage.trunk.VoiceVox.API_GetSpeakIndex(data.Speaker)
					path, _ := self.mainSettingPage.MainPage.trunk.VoiceVox.GenerateWAV(
						convCont,
						strconv.Itoa(index),
						data.Speed,
						data.Pitch,
						data.Intonation,
						data.Volume,
					)
					self.mainSettingPage.MainPage.trunk.Player.PlayFile(path)
				}()
			}
		}
	}

	// 恢复按钮
	if self.btn_Reset.Clicked(gtx) {
		self.coe_IsDefault = true
	}

	// 播放/保存/复位按钮样式
	playBtnFunc := self.build_IconButton(gtx, &self.btn_PlayVoice, "\uf5b0", self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CreateVoice_Audition)
	saveBtnFunc := self.build_IconButton(gtx, &self.btn_SaveVoice, "\ue74e", self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CreateVoice_Save)
	restoreBtnFunc := self.build_IconButton(gtx, &self.btn_Reset, "\ue777", self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CreateVoice_Restore)
	btnFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, playBtnFunc),
			layout.Flexed(1, saveBtnFunc),
			layout.Flexed(1, restoreBtnFunc),
		)
	}

	// 语言转换控件
	convFunc := func(gtx C) D {
		gtx.Constraints.Max.Y = gtx.Dp(20)
		gtx.Constraints.Min.Y = gtx.Dp(20)
		s := material.Switch(self.style.theme, &self.langConv_switch, "")
		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(material.Label(self.style.theme, unit.Sp(16), self.languageTable["LangConv"][self.language]).Layout),
					layout.Flexed(1, Flexer()),
					layout.Rigid(s.Layout),
				)
			})
		})
	}
	// 预览编辑框与标题编辑框
	previewEditorFunc := self.build_EditorBox(gtx, &self.editor_Preview, self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CreateVoice_Preview)
	titleEditorFunc := self.build_EditorBox(gtx, &self.editor_Title, self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CreateVoice_VoiceName)
	editorFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(previewEditorFunc),
			layout.Rigid(titleEditorFunc),
		)
	}

	// 引擎选择控件
	self.switchEngine_DropDown.API_SetTitle(self.languageTable["VoiceEngine"][self.language])
	engineFunc := func(gtx C) D {
		return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
			return self.switchEngine_DropDown.Update(gtx)
		})
	}

	// 展开栏顶部样式
	headFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
			layout.Rigid(material.Label(self.style.theme, unit.Sp(20), "\uf8aa").Layout),
			layout.Rigid(Spacer(gtx.Dp(5), 0)),
			layout.Rigid(material.Label(self.style.theme, unit.Sp(18), self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CreateVoice_Btn).Layout),
			layout.Flexed(1, Flexer()),
		)
	}

	// 展开栏内容样式
	panelFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(btnFunc),
			layout.Rigid(convFunc),
			layout.Rigid(editorFunc),
			layout.Rigid(engineFunc),
			layout.Rigid(self.coe_CurPanel.Update(gtx, self)),
		)
	}

	// 具体选项参数
	d, _ := self.newVoice_ExpandPanel.Update(gtx, headFunc, panelFunc)
	return d
}

// 创建一些自定义的元素 ----------------------------------------------------
// 创建自定义按钮
func (self *SoundSetting) build_Button(
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

func (self *SoundSetting) build_IconButton(gtx C, btn *widget.Clickable, icon string, title string) W {
	BtnStyle := func(gtx C) D {
		icon_label := material.Label(self.style.theme, unit.Sp(20), icon)
		text_label := material.Label(self.style.theme, unit.Sp(18), title)
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
		BtnStyle,
		btn,
		layout.Inset{},
		color.NRGBA{A: 0},
		gtx.Constraints.Max.X, gtx.Dp(40),
	)
}

func (self *SoundSetting) build_EditorBox(gtx C, editorWidget *widget.Editor, title string) W {
	return func(gtx C) D {
		return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Dp(50)
			gtx.Constraints.Max.Y = gtx.Dp(50)
			return BorderBox(
				self.style,
				self.style.darkmode.currentColor.Fg,
				self.style.darkmode.currentColor.IdleBg,
				title,
			).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return material.Editor(self.style.theme, editorWidget, "").Layout(gtx)
					}),
				)

			})
		})
	}
}

// 创建文本标题
func (self *SoundSetting) build_TitleLabel(gtx C, title string) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D { return material.Label(self.style.theme, unit.Sp(16), title).Layout(gtx) }),
		layout.Flexed(1, func(gtx C) D { size := gtx.Constraints.Min; return D{Size: size} }),
	)
}

// 内部小工具 ---------------------------------------------------------
// 计算整数
func (self *SoundSetting) tool_ComputeValue(cur float32, min int, max int) int {
	return min + int(cur*float32(max-min))
}

// 重设预设
func (self *SoundSetting) tool_ReloadPreset() {
	self.preset_list = []SoundSettingPresetBox{}
	sp := config.Config_Voice{}
	preset_StrList, _ := sp.Get_Preset()
	for _, p := range preset_StrList {
		switch p[0] {
		case "AquesTalk":
			ap, _ := sp.ConvertTo_Aquestalk(p)
			self.preset_list = append(self.preset_list, New_PresetBoxAquesTalk(self.style, ap))
		case "OpenJTalk":
			op, _ := sp.ConvertTo_OpenJTalk(p)
			self.preset_list = append(self.preset_list, New_PresetBoxOpenJTalk(self.style, op))
		case "VoiceVox":
			vp, _ := sp.ConvertTo_VoiceVox(p)
			self.preset_list = append(self.preset_list, New_PresetBoxVoiceVox(self.style, vp))
		default:
		}
	}
}

// 删除功能
func (self *SoundSetting) tool_DeletePreset(name string) {
	go func() {
		sp := config.Config_Voice{}
		sp.Delete_Preset(name)
		self.tool_ReloadPreset()
	}()
	if name == self.preset_Enum.Value {
		if len(self.preset_list) > 0 {
			self.preset_Enum.Value = self.preset_list[0].API_GetName()
			// config.API_SOUND_SetCurrentPresetTitle(self.preset_list[0].API_GetName())
		} else {
			self.preset_Enum.Value = ""
			// config.API_SOUND_SetCurrentPresetTitle("")
		}
	}
	self.mainSettingPage.MainPage.notice.Add_Notice(1, self.languageTable["Notice_Delete_Title"][self.language], self.languageTable["Notice_Delete_Desc"][self.language])
}

// 点击点选按钮后的功能
func (self *SoundSetting) tool_SetCurPreset(preset any) {
	trunk := self.mainSettingPage.MainPage.trunk
	trunk.VoicePreset = preset

	cv := config.Config_Voice{}
	cv.Set_CurPreset(preset)
	self.preset_Enum.Value = cv.ToString_Preset(preset)[1]
}

// 设置全局预设
func (self *SoundSetting) tool_SetPreset(presetName string) {
	self.preset_Enum.Value = presetName
	// config.API_SOUND_SetCurrentPresetTitle(presetName)
}

// 语言
func (self *SoundSetting) Language() LangPack {
	return LangPack{
		"LangConv": LangZone{
			"Chinese":  "中文转片假名",
			"Japanese": "中国語→カタカナ",
			"English":  "Chinese To Katakana",
		},
		"VoiceEngine": LangZone{
			"Chinese":  "语音引擎",
			"English":  "Voice Engine",
			"Japanese": "音声エンジン",
		},
		"Notice_Add_Title": LangZone{
			"Chinese":  "语音预设已保存",
			"English":  "Voice Preset Saved",
			"Japanese": "プリセットを保存",
		},
		"Notice_Add_Desc": LangZone{
			"Chinese":  "...",
			"English":  "...",
			"Japanese": "...",
		},
		"Notice_AddRepeat_Title": LangZone{
			"Chinese":  "语音预设已存在",
			"English":  "Voice Preset Already Exists",
			"Japanese": "重複しています",
		},
		"Notice_AddRepeat_Desc": LangZone{
			"Chinese":  "...",
			"English":  "...",
			"Japanese": "...",
		},
		"Notice_Delete_Title": LangZone{
			"Chinese":  "语音预设已删除",
			"English":  "Voice Preset Deleted",
			"Japanese": "プリセットを削除",
		},
		"Notice_Delete_Desc": LangZone{
			"Chinese":  "...",
			"English":  "...",
			"Japanese": "...",
		},
	}
}

// 参数设置界面 --------------------------------------------------------------------------------------------
// aquestalk参数界面 ---------------------------------------------------------------------------------------
type CoefficientAquestalk struct {
	slider_speed      CoefficientSlider
	slider_volume     CoefficientSlider
	slider_pitch      CoefficientSlider
	slider_quality    CoefficientSlider
	slider_intonation CoefficientSlider
	slider_accent     CoefficientSlider
	language          string
	languageTable     LangPack
}

func New_CoefficientAquestalk(style Style) *CoefficientAquestalk {
	ca := CoefficientAquestalk{
		slider_speed:      *New_CoefficientSlider(style, "语速", 300, 50, 100),
		slider_volume:     *New_CoefficientSlider(style, "音量", 300, 0, 100),
		slider_pitch:      *New_CoefficientSlider(style, "音高", 300, 20, 100),
		slider_accent:     *New_CoefficientSlider(style, "顿音", 200, 0, 100),
		slider_quality:    *New_CoefficientSlider(style, "音质", 200, 0, 100),
		slider_intonation: *New_CoefficientSlider(style, "夹音", 200, 50, 100),
	}
	ca.language = "Chinese"
	ca.languageTable = ca.Language()
	return &ca
}

func (self *CoefficientAquestalk) Default(gtx C) bool {
	done := true
	done = self.smoothReturnToPosition(gtx, &self.slider_speed, 100) && done
	done = self.smoothReturnToPosition(gtx, &self.slider_volume, 100) && done
	done = self.smoothReturnToPosition(gtx, &self.slider_pitch, 100) && done
	done = self.smoothReturnToPosition(gtx, &self.slider_accent, 100) && done
	done = self.smoothReturnToPosition(gtx, &self.slider_quality, 100) && done
	done = self.smoothReturnToPosition(gtx, &self.slider_intonation, 100) && done
	return done
}

func (self *CoefficientAquestalk) Update(gtx C, parent *SoundSetting) W {
	self.slider_speed.API_SetTitle(self.languageTable["Speed"][self.language])
	self.slider_volume.API_SetTitle(self.languageTable["Volume"][self.language])
	self.slider_pitch.API_SetTitle(self.languageTable["Pitch"][self.language])
	self.slider_accent.API_SetTitle(self.languageTable["Accent"][self.language])
	self.slider_quality.API_SetTitle(self.languageTable["Quality"][self.language])
	self.slider_intonation.API_SetTitle(self.languageTable["Intonation"][self.language])

	return func(gtx C) D {
		return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(self.slider_speed.Update),
				layout.Rigid(self.slider_volume.Update),
				layout.Rigid(self.slider_pitch.Update),
				layout.Rigid(self.slider_accent.Update),
				layout.Rigid(self.slider_quality.Update),
				layout.Rigid(self.slider_intonation.Update),
			)
		})
	}
}

func (self *CoefficientAquestalk) API_SetLanguage(language string) { self.language = language }

func (self *CoefficientAquestalk) API_GetCoefficient() any {
	isDefault := false
	isMonoTone := false

	speed := self.slider_speed.API_GetIntValue()
	volume := self.slider_volume.API_GetIntValue()
	pitch := self.slider_pitch.API_GetIntValue()
	accent := self.slider_accent.API_GetIntValue()
	quality := self.slider_quality.API_GetIntValue()
	intonation := self.slider_intonation.API_GetIntValue()

	pa := config.Preset_Aquestalk{
		Name:       "",
		Default:    isDefault,
		MonoTone:   isMonoTone,
		Type:       "AquesTalk10",
		Engine:     "F1E",
		Speed:      speed,
		Volume:     volume,
		Pitch:      pitch,
		Accent:     accent,
		Quality:    quality,
		Intonation: intonation,
	}
	return pa
}

// 滑块条顺滑归位
func (self *CoefficientAquestalk) smoothReturnToPosition(gtx C, slider *CoefficientSlider, targetValue int) bool {
	current := float64(slider.API_GetIntValue())
	target := float64(targetValue)
	if (current-target)*(current-target) <= 25 {
		slider.API_SetIntValue(targetValue)
		return true
	}

	newVal := easing(gtx, current, target, 0.3)
	slider.API_SetIntValue(int(newVal))
	return int(newVal) == targetValue
}
func (self *CoefficientAquestalk) Language() LangPack {
	return LangPack{
		"Speed": LangZone{
			"Chinese":  "语速",
			"English":  "Speed",
			"Japanese": "話速",
		},
		"Volume": LangZone{
			"Chinese":  "音量",
			"English":  "Volume",
			"Japanese": "音量",
		},
		"Pitch": LangZone{
			"Chinese":  "音高",
			"English":  "Pitch",
			"Japanese": "ピッチ",
		},
		"Accent": LangZone{
			"Chinese":  "声调",
			"English":  "Accent",
			"Japanese": "アクセント",
		},
		"Intonation": LangZone{
			"Chinese":  "语调",
			"English":  "Intonation",
			"Japanese": "抑揚",
		},
		"Quality": LangZone{
			"Chinese":  "音质",
			"English":  "Quality",
			"Japanese": "音質",
		},
	}
}

// openjtalk语音参数 -----------------------------------------------------------------------------------------
type CoefficientOpenJTalk struct {
	style          Style
	slider_speed   CoefficientSlider
	slider_volume  CoefficientSlider
	slider_pitch   CoefficientSlider
	picker_speaker PickBar
	language       string
	languagePack   LangPack
}

func New_CoefficientOpenJTalk(style Style) *CoefficientOpenJTalk {
	co := CoefficientOpenJTalk{
		style:          style,
		slider_speed:   *New_CoefficientSlider(style, "Speed", 200, 50, 100),
		slider_volume:  *New_CoefficientSlider(style, "Volume", 30, -30, 1),
		slider_pitch:   *New_CoefficientSlider(style, "Pitch", 30, -30, 0),
		picker_speaker: *New_PickBar(style),
	}
	choices := config.API_SOUND_Get_AllHtsVoiceFiles()
	co.picker_speaker.API_SetChoices(choices)
	co.picker_speaker.API_SetValue(choices[0])
	co.language = "Chinese"
	co.languagePack = co.Language()
	return &co
}
func (self *CoefficientOpenJTalk) Default(gtx C) bool {
	done := true
	done = self.smoothReturnToPosition(gtx, &self.slider_speed, 100) && done
	done = self.smoothReturnToPosition(gtx, &self.slider_volume, 1) && done
	done = self.smoothReturnToPosition(gtx, &self.slider_pitch, 0) && done
	return done
}
func (self *CoefficientOpenJTalk) Update(gtx C, parent *SoundSetting) W {
	self.slider_speed.API_SetTitle(self.languagePack["Speed"][self.language])
	self.slider_pitch.API_SetTitle(self.languagePack["Pitch"][self.language])
	self.slider_volume.API_SetTitle(self.languagePack["Volume"][self.language])
	return func(gtx C) D {
		return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(material.Label(self.style.theme, unit.Sp(16), self.languagePack["Speaker"][self.language]).Layout),
						layout.Flexed(1, Flexer()),
						layout.Rigid(self.picker_speaker.Update),
					)
				}),
				layout.Rigid(self.slider_pitch.Update),
				layout.Rigid(self.slider_speed.Update),
				layout.Rigid(self.slider_volume.Update),
			)
		})
	}
}
func (self *CoefficientOpenJTalk) API_SetLanguage(lang string) { self.language = lang }
func (self *CoefficientOpenJTalk) API_GetCoefficient() any {
	fmt.Println()
	po := config.Preset_OpenJTalk{
		Name:    "",
		Default: false,
		Speaker: self.picker_speaker.API_GetValue(),
		Speed:   float64(self.slider_speed.API_GetIntValue()) * 0.01,
		Volume:  float64(self.slider_volume.API_GetIntValue()),
		Pitch:   float64(self.slider_pitch.API_GetIntValue()),
	}
	return po
}
func (self *CoefficientOpenJTalk) smoothReturnToPosition(gtx C, slider *CoefficientSlider, targetValue int) bool {
	current := float64(slider.API_GetIntValue())
	target := float64(targetValue)
	if (current-target)*(current-target) <= 25 {
		slider.API_SetIntValue(targetValue)
		return true
	}

	newVal := easing(gtx, current, target, 0.3)
	slider.API_SetIntValue(int(newVal))
	return int(newVal) == targetValue
}
func (self *CoefficientOpenJTalk) Language() LangPack {
	return LangPack{
		"Speaker": LangZone{
			"Chinese":  "讲者",
			"English":  "Speaker",
			"Japanese": "話者",
		},
		"Speed": LangZone{
			"Chinese":  "语速",
			"English":  "Speed",
			"Japanese": "話速",
		},
		"Pitch": LangZone{
			"Chinese":  "音高",
			"English":  "Pitch",
			"Japanese": "ピッチ",
		},
		"Volume": LangZone{
			"Chinese":  "音量",
			"English":  "Volume",
			"Japanese": "音量",
		},
	}
}

// Teto语音参数 ------------------------------------------------------------------------------------------
type CoefficientTeto struct {
	title string
}

func New_CoefficientTeto(style Style) *CoefficientTeto {
	ct := CoefficientTeto{}
	return &ct
}
func (self CoefficientTeto) Update(gtx C, parent *SoundSetting) W {
	return func(gtx C) D { return D{} }
}
func (self CoefficientTeto) Default(gtx C) bool { return true }
func (self CoefficientTeto) API_GetCoefficient() any {
	return config.Preset_Teto{
		Name:    self.title,
		Default: false,
	}
}

// VoiceVox语音参数 ------------------------------------------------------------------------------------------
type CoefficientVoiceVox struct {
	style       Style
	title       string
	Speaker     string
	SpeakerList []string

	box_Speaker       MultiParmmeterBox
	slider_Speed      CoefficientSlider
	slider_Pitch      CoefficientSlider
	slider_Intonation CoefficientSlider
	slider_Volume     CoefficientSlider

	isSpeakers       bool
	isSpeakerLoading bool

	language      string
	languageTable LangPack
}

func New_CoefficientVoiceVox(style Style) *CoefficientVoiceVox {
	v := CoefficientVoiceVox{
		style:             style,
		box_Speaker:       *New_MultiParmmeterBox(style, []string{"...", "...", "...", "..."}, "..."),
		slider_Speed:      *New_CoefficientSlider(style, "speed", 200, 50, 100),
		slider_Pitch:      *New_CoefficientSlider(style, "pitch", 15, -15, 0),
		slider_Intonation: *New_CoefficientSlider(style, "intonation", 200, 0, 100),
		slider_Volume:     *New_CoefficientSlider(style, "volume", 200, 0, 100),
		isSpeakers:        false,
		isSpeakerLoading:  false,
	}

	v.language = "Chinese"
	v.languageTable = v.Language()
	return &v
}
func (self *CoefficientVoiceVox) Default(gtx C) bool {
	result := true
	result = self.smoothReturnToPosition(gtx, &self.slider_Speed, 100) && result
	result = self.smoothReturnToPosition(gtx, &self.slider_Pitch, 0) && result
	result = self.smoothReturnToPosition(gtx, &self.slider_Intonation, 100) && result
	result = self.smoothReturnToPosition(gtx, &self.slider_Volume, 100) && result
	self.box_Speaker.API_SetValueIndex(0)

	return result
}
func (self *CoefficientVoiceVox) Update(gtx C, parent *SoundSetting) W {
	self.box_Speaker.API_SetTitle(self.languageTable["Speaker"][self.language])
	self.slider_Speed.API_SetTitle(self.languageTable["Speed"][self.language])
	self.slider_Pitch.API_SetTitle(self.languageTable["Pitch"][self.language])
	self.slider_Intonation.API_SetTitle(self.languageTable["Intonation"][self.language])
	self.slider_Volume.API_SetTitle(self.languageTable["Volume"][self.language])

	if !self.isSpeakers && !self.isSpeakerLoading {
		self.isSpeakerLoading = true
		voiceVox := parent.mainSettingPage.MainPage.trunk.VoiceVox
		go func() {
			speakers := voiceVox.API_GetSpeakers()
			if len(speakers) > 0 {
				self.box_Speaker.API_SetChoices(speakers)
				self.box_Speaker.API_SetValueIndex(0)
				self.isSpeakers = true
			}
			self.isSpeakerLoading = false
		}()
	}
	return func(gtx C) D {
		return layout.Inset{Left: unit.Dp(15), Right: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(self.box_Speaker.Update),
				layout.Rigid(self.slider_Speed.Update),
				layout.Rigid(self.slider_Pitch.Update),
				layout.Rigid(self.slider_Intonation.Update),
				layout.Rigid(self.slider_Volume.Update),
			)
		})
	}
}
func (self *CoefficientVoiceVox) API_GetCoefficient() any {
	speedRaw := float64(self.slider_Speed.API_GetFloatValue()*1.5 + 0.5)
	pitchRaw := float64(self.slider_Pitch.API_GetFloatValue()*0.3 - 0.15)
	intonationRaw := float64(self.slider_Intonation.API_GetFloatValue() * 2)
	volumeRaw := float64(self.slider_Volume.API_GetFloatValue() * 2)

	speed := math.Round(speedRaw*100) / 100
	pitch := math.Round(pitchRaw*100) / 100
	intonation := math.Round(intonationRaw*100) / 100
	volume := math.Round(volumeRaw*100) / 100

	return config.Preset_VoiceVox{
		Name:       "",
		Default:    false,
		Speaker:    self.box_Speaker.API_GetValue(),
		Speed:      speed,
		Pitch:      pitch,
		Intonation: intonation,
		Volume:     volume,
	}
}
func (self *CoefficientVoiceVox) smoothReturnToPosition(gtx C, slider *CoefficientSlider, targetValue int) bool {
	current := float64(slider.API_GetIntValue())
	target := float64(targetValue)
	if (current-target)*(current-target) <= 25 {
		slider.API_SetIntValue(targetValue)
		return true
	}

	newVal := easing(gtx, current, target, 0.3)
	slider.API_SetIntValue(int(newVal))
	return int(newVal) == targetValue
}
func (self *CoefficientVoiceVox) API_SetLanguage(language string) { self.language = language }
func (self *CoefficientVoiceVox) Language() LangPack {
	return LangPack{
		"Speaker": LangZone{
			"Chinese":  "讲者",
			"English":  "Speaker",
			"Japanese": "話者",
		},
		"Speed": LangZone{
			"Chinese":  "语速",
			"English":  "Speed",
			"Japanese": "話速",
		},
		"Pitch": LangZone{
			"Chinese":  "音高",
			"English":  "Pitch",
			"Japanese": "ピッチ",
		},
		"Intonation": LangZone{
			"Chinese":  "语调",
			"English":  "Intonation",
			"Japanese": "抑揚",
		},
		"Volume": LangZone{
			"Chinese":  "音量",
			"English":  "Volume",
			"Japanese": "音量",
		},
	}
}

// 多参数单选控件
type MultiParmmeterBox struct {
	style   Style
	choices []string
	value   string

	arrowLeft  widget.Clickable
	arrowRight widget.Clickable

	pageCount int
	pageNum   int

	clickables []widget.Clickable
	enum       widget.Enum
	title      string
}

func New_MultiParmmeterBox(style Style, choices []string, value string) *MultiParmmeterBox {
	m := MultiParmmeterBox{
		style:   style,
		choices: choices,
		value:   value,
		enum:    widget.Enum{},

		arrowLeft:  widget.Clickable{},
		arrowRight: widget.Clickable{},

		pageCount: 1,
		pageNum:   1,

		clickables: make([]widget.Clickable, 9),
		title:      "MultiParmmeterBox",
	}
	m.pageNum = (len(m.choices) + 9 - 1) / 9
	return &m
}
func (self *MultiParmmeterBox) Update(gtx C) D {
	// 控件功能
	if self.arrowLeft.Clicked(gtx) {
		self.pageCount -= 1
	}
	if self.arrowRight.Clicked(gtx) {
		self.pageCount += 1
	}
	totalPages := 0
	if len(self.choices) > 0 {
		totalPages = (len(self.choices) + 9 - 1) / 9
	}
	if self.pageCount <= 1 {
		self.pageCount = 1
	}
	if self.pageCount >= totalPages {
		self.pageCount = totalPages
	}

	btns := make([]W, 9)

	for i := range 9 {
		index := i + (self.pageCount-1)*9
		var desc string
		if index < len(self.choices) {
			desc = self.choices[index]
		} else {
			desc = ""
		}

		b := material.Button(self.style.theme, &self.clickables[i], desc)
		b.Inset = layout.Inset{}
		b.Background.A = 0
		b.Color = self.style.theme.Fg
		if desc == "" {
			btns[i] = func(gtx C) D {
				return D{}
			}
		} else {
			btns[i] = func(gtx C) D {
				return layout.UniformInset(unit.Dp(2)).Layout(gtx, b.Layout)
			}
			if self.clickables[i].Clicked(gtx) {
				self.value = desc
			}
		}
	}

	// 左右箭头
	leftFunc := func(gtx C) D {
		gtx.Constraints.Max = image.Pt(gtx.Dp(20), gtx.Dp(20))
		gtx.Constraints.Min = gtx.Constraints.Max
		b := material.Button(self.style.theme, &self.arrowLeft, "\uF08D")
		b.Inset = layout.Inset{}
		b.Background.A = 0
		b.Color = self.style.theme.Fg
		return b.Layout(gtx)
	}
	rightFunc := func(gtx C) D {
		gtx.Constraints.Max = image.Pt(gtx.Dp(20), gtx.Dp(20))
		gtx.Constraints.Min = gtx.Constraints.Max
		b := material.Button(self.style.theme, &self.arrowRight, "\uF08F")
		b.Inset = layout.Inset{}
		b.Background.A = 0
		b.Color = self.style.theme.Fg
		return b.Layout(gtx)
	}

	// 左右箭头中间的页码
	pageCountFunc := func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Dp(40)
		gtx.Constraints.Max.X = gtx.Dp(40)
		text := fmt.Sprintf("%d/%d", self.pageCount, self.pageNum)
		l := material.Label(self.style.theme, unit.Sp(16), text)
		return layout.Center.Layout(gtx, l.Layout)
	}

	headFunc := func(gtx C) D {
		label := material.Label(self.style.theme, unit.Sp(16), self.value)
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(label.Layout),
			layout.Flexed(1, Flexer()),
			layout.Rigid(leftFunc),
			layout.Rigid(pageCountFunc),
			layout.Rigid(rightFunc),
		)
	}

	panelFunc := func(gtx C) D {
		gtx.Constraints.Min.Y = gtx.Dp(100)
		gtx.Constraints.Max.Y = gtx.Dp(100)
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, btns[0]),
					layout.Flexed(1, btns[1]),
					layout.Flexed(1, btns[2]),
				)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, btns[3]),
					layout.Flexed(1, btns[4]),
					layout.Flexed(1, btns[5]),
				)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, btns[6]),
					layout.Flexed(1, btns[7]),
					layout.Flexed(1, btns[8]),
				)
			}),
		)
	}
	return BorderBox(
		self.style,
		self.style.darkmode.currentColor.IdleBg,
		self.style.theme.Fg,
		self.title,
	).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(headFunc),
			layout.Flexed(1, func(gtx C) D {
				DrawLine(gtx, float32(gtx.Dp(1)), 0, float32(gtx.Constraints.Max.X-2), 0, 1, self.style.theme.Fg)
				return D{}
			}),
			layout.Rigid(panelFunc),
		)
	})

}
func (self *MultiParmmeterBox) API_GetValue() string        { return self.value }
func (self *MultiParmmeterBox) API_SetValue(value string)   { self.value = value }
func (self *MultiParmmeterBox) API_SetValueIndex(index int) { self.value = self.choices[index] }
func (self *MultiParmmeterBox) API_SetTitle(title string)   { self.title = title }
func (self *MultiParmmeterBox) API_SetChoices(choices []string) {
	self.choices = choices
	self.pageNum = (len(self.choices) + 9 - 1) / 9
}

// 类: 参数条 ------------------------------------------------
type CoefficientSlider struct {
	style  Style
	slider widget.Float

	title    string
	maxValue float32
	minValue float32

	value      float32
	floatValue float32
}

func New_CoefficientSlider(style Style, title string, max float32, min float32, defaultValue float32) *CoefficientSlider {
	value := (defaultValue - min) / (max - min)

	slider := widget.Float{}
	slider.Value = value

	return &CoefficientSlider{
		style:  style,
		slider: slider,

		title:      title,
		maxValue:   max,
		minValue:   min,
		floatValue: value,
		value:      defaultValue,
	}
}
func (self *CoefficientSlider) Update(gtx C) D {
	// s := material.Slider(self.style.theme, &self.slider)
	s := SimpleSliderStyle{
		Axis:         layout.Horizontal,
		Float:        &self.slider,
		ActiveColor:  self.style.theme.Palette.ContrastBg,
		FingerSize:   20, // 触摸热区 44dp
		TrackHeight:  10, // 轨道高度 8dp
		ThumbWidth:   4,  // 竖线宽度 6dp
		ThumbHeight:  20, // 竖线高度 20dp
		ThumbRadius:  2,
		CornerRadius: 4, // 圆角半径 4dp
		GapSize:      4,
	}
	// s.Color = self.style.darkmode.currentColor.ContrastBg
	s.FingerSize = unit.Dp(12)

	// 更新数值逻辑
	self.floatValue = self.slider.Value
	self.value = self.minValue + ((self.maxValue - self.minValue) * self.floatValue)

	// 定义组件
	valueLabel := material.Label(self.style.theme, unit.Sp(16), strconv.FormatFloat(float64(self.value), 'f', 0, 64))
	titleLabel := material.Label(self.style.theme, unit.Sp(16), self.title)

	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(titleLabel.Layout),
				layout.Flexed(1, Flexer()),
				layout.Rigid(valueLabel.Layout),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			gtx.Constraints.Max.Y = gtx.Dp(16)
			gtx.Constraints.Min.Y = gtx.Dp(16)
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, s.Layout),
			)
		}),
	)
}

func (self *CoefficientSlider) API_GetFloatValue() float32 {
	return self.floatValue
}
func (self *CoefficientSlider) API_SetFloatValue(value float32) {
	self.slider.Value = value
}

// 返回参数
func (self *CoefficientSlider) API_GetIntValue() int {
	return int(self.value)
}

func (self *CoefficientSlider) API_SetIntValue(value int) {
	self.value = float32(value)
	self.floatValue = (self.value - self.minValue) / (self.maxValue - self.minValue)
	self.slider.Value = self.floatValue
}

func (self *CoefficientSlider) API_SetTitle(title string) {
	self.title = title
}

// 类: 表示参数选择的滑块条
type CoefficientDiscreteSlider struct {
	style Style

	title string
	coes  []string
	value string
	index int

	slider widget.Float
}

func New_CoefficientDiscreteSlider(style Style, title string, coes []string) *CoefficientDiscreteSlider {
	cds := CoefficientDiscreteSlider{
		style:  style,
		title:  title,
		coes:   coes,
		slider: widget.Float{Value: 0.3},
	}

	if len(coes) == 0 {
		cds.value = "nil"
		cds.index = 0
	} else {
		cds.value = coes[0]
		cds.index = 0
	}
	cds.slider.Value = 0
	return &cds
}

func (self *CoefficientDiscreteSlider) Update(gtx C) D {
	length := len(self.coes)
	if length > 1 {
		step := 1.0 / float32(length-1)
		self.index = int(self.slider.Value/step + 0.5)

		if self.index <= 0 {
			self.index = 0
		}
		if self.index >= length {
			self.index = length - 1
		}

		self.slider.Value = float32(self.index) * step
	} else {
		self.index = 0
		self.slider.Value = 0
	}

	// 上述标题与参数信息
	titleFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(material.Label(self.style.theme, unit.Sp(16), self.title).Layout),
			layout.Flexed(1, Flexer()),
			layout.Rigid(material.Label(self.style.theme, unit.Sp(16), self.coes[self.index]).Layout),
		)
	}

	// 下部滑块条
	sliderFunc := func(gtx C) D {

		gtx.Constraints.Max.Y = gtx.Dp(16)
		gtx.Constraints.Min.Y = gtx.Dp(16)
		slider := material.Slider(self.style.theme, &self.slider)
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layout.Flexed(1, slider.Layout))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(titleFunc),
		layout.Rigid(Spacer(0, gtx.Dp(4))),
		layout.Rigid(sliderFunc),
	)
}

func (self *CoefficientDiscreteSlider) API_SetTitle(title string) {
	self.title = title
}

func (self *CoefficientDiscreteSlider) API_GetValue() string {
	return self.coes[self.index]
}
func (self *CoefficientDiscreteSlider) API_GetValueFloat() float32 {
	return self.slider.Value
}
func (self *CoefficientDiscreteSlider) API_SetValueFloat(value float32) {
	self.slider.Value = value
}

// aquestalk预设 ----------------------------------------------------------------------------
type PresetBoxAquesTalk struct {
	style  Style
	title  string
	info   config.Preset_Aquestalk
	coes   []string
	expand ExpandPanel

	checkBtn  widget.Clickable
	checkIcon string
	deleteBtn widget.Clickable

	language      string
	languageTable LangPack
}

func New_PresetBoxAquesTalk(style Style, info config.Preset_Aquestalk) *PresetBoxAquesTalk {
	p := PresetBoxAquesTalk{
		style:  style,
		info:   info,
		expand: *New_ExpandPanel(style, 160),
	}

	speed := strconv.Itoa(info.Speed)
	volume := strconv.Itoa(info.Volume)
	pitch := strconv.Itoa(info.Pitch)
	accent := strconv.Itoa(info.Accent)
	quality := strconv.Itoa(info.Quality)
	intonation := strconv.Itoa(info.Intonation)
	p.coes = []string{speed, volume, pitch, accent, quality, intonation}

	p.language = "Chinese"
	p.languageTable = p.Language()

	return &p
}

func (self *PresetBoxAquesTalk) Update(gtx C, parent *SoundSetting) W {

	// 点选按钮功能
	if self.checkBtn.Clicked(gtx) {
		// parent.preset_Enum.Value = self.info.Name
		// parent.tool_SetPreset(self.info.Name)
		parent.tool_SetCurPreset(self.info)
	}
	var checkIcon string
	if parent.preset_Enum.Value == self.info.Name {
		checkIcon = "\uE73D"
	} else {
		checkIcon = "\uE739"
	}

	// 删除按钮功能
	if self.deleteBtn.Clicked(gtx) {
		parent.tool_DeletePreset(self.info.Name)
	}

	// 点选按钮样式
	checkBtnStyle := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(22), gtx.Sp(24))
		btn := material.Button(self.style.theme, &self.checkBtn, checkIcon)
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
				layout.Rigid(checkBtnStyle),
				layout.Rigid(Spacer(gtx.Dp(5), 0)),
				layout.Rigid(material.Label(self.style.theme, unit.Sp(20), self.info.Name).Layout),
				layout.Flexed(1, Flexer()),
			)
		})
	}

	panelFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)

						btn := material.Button(self.style.theme, &self.deleteBtn, "\uE74D")
						btn.Background.A = 0
						btn.Color = self.style.darkmode.currentColor.Fg
						btn.TextSize = unit.Sp(20)

						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Engine"][self.language], self.info.Type)),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Type"][self.language], self.info.Engine)),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Rigid(btn.Layout),
						)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Speed"][self.language], self.coes[0])),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Volume"][self.language], self.coes[1])),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Pitch"][self.language], self.coes[2])),
						)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Accent"][self.language], self.coes[3])),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Quality"][self.language], self.coes[4])),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Intonation"][self.language], self.coes[5])),
						)
					}),
				)
			}),
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
		)
	}
	return func(gtx C) D {
		d, _ := self.expand.Update(gtx, headFunc, panelFunc)
		return d
	}
}
func (self *PresetBoxAquesTalk) build_LabelBox(gtx C, title string, content string) W {
	return func(gtx C) D {
		return BorderBox(
			self.style,
			self.style.darkmode.currentColor.Fg,
			// self.style.darkmode.currentColor.IdleBg,
			self.style.darkmode.currentColor.IdleBg,
			title,
		).Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return material.Label(self.style.theme, unit.Sp(18), content).Layout(gtx)
				}),
			)
		})
	}
}
func (self *PresetBoxAquesTalk) API_GetName() string              { return self.info.Name }
func (self *PresetBoxAquesTalk) API_SetLanguage(languagge string) { self.language = languagge }
func (self *PresetBoxAquesTalk) Language() LangPack {
	return LangPack{
		"Engine": LangZone{
			"Chinese":  "引擎",
			"English":  "Engine",
			"Japanese": "エンジン",
		},
		"Type": LangZone{
			"Chinese":  "类型",
			"English":  "Engine",
			"Japanese": "タイプ",
		},
		"Speed": LangZone{
			"Chinese":  "语速",
			"English":  "Speed",
			"Japanese": "話速",
		},
		"Volume": LangZone{
			"Chinese":  "音量",
			"English":  "Volume",
			"Japanese": "音量",
		},
		"Pitch": LangZone{
			"Chinese":  "音高",
			"English":  "Pitch",
			"Japanese": "ピッチ",
		},
		"Accent": LangZone{
			"Chinese":  "声调",
			"English":  "Accent",
			"Japanese": "アクセント",
		},
		"Intonation": LangZone{
			"Chinese":  "语调",
			"English":  "Intonation",
			"Japanese": "抑揚",
		},
		"Quality": LangZone{
			"Chinese":  "音质",
			"English":  "Quality",
			"Japanese": "音質",
		},
	}
}

// openJTalk预设 -------------------------------------------------------------------------------------------
type PresetBoxOpenJTalk struct {
	style     Style
	info      config.Preset_OpenJTalk
	expand    ExpandPanel
	coes      []string
	checkBtn  widget.Clickable
	deleteBtn widget.Clickable

	language      string
	languageTable LangPack
}

func New_PresetBoxOpenJTalk(style Style, info config.Preset_OpenJTalk) *PresetBoxOpenJTalk {
	p := PresetBoxOpenJTalk{
		style:     style,
		info:      info,
		expand:    *New_ExpandPanel(style, 110),
		checkBtn:  widget.Clickable{},
		deleteBtn: widget.Clickable{},
	}

	speak := info.Speaker
	speed := strconv.FormatFloat(info.Speed*100, 'f', 0, 64)
	volume := strconv.FormatFloat(info.Volume, 'f', 0, 64)
	pitch := strconv.FormatFloat(info.Pitch, 'f', 0, 64)
	p.coes = []string{speak, speed, volume, pitch}

	p.language = "Chinese"
	p.languageTable = p.Language()

	return &p
}
func (self *PresetBoxOpenJTalk) Update(gtx C, parent *SoundSetting) W {
	// 点选按钮功能
	if self.checkBtn.Clicked(gtx) {
		// parent.tool_SetPreset(self.info.Name)
		parent.tool_SetCurPreset(self.info)
	}
	var checkIcon string
	if parent.preset_Enum.Value == self.info.Name {
		checkIcon = "\uE73D"
	} else {
		checkIcon = "\uE739"
	}

	// 删除按钮功能
	if self.deleteBtn.Clicked(gtx) {
		parent.tool_DeletePreset(self.info.Name)
	}

	// 点选按钮样式
	checkBtnStyle := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(22), gtx.Sp(24))
		btn := material.Button(self.style.theme, &self.checkBtn, checkIcon)
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
				layout.Rigid(checkBtnStyle),
				layout.Rigid(Spacer(gtx.Dp(5), 0)),
				layout.Rigid(material.Label(self.style.theme, unit.Sp(20), self.info.Name).Layout),
				layout.Flexed(1, Flexer()),
			)
		})
	}

	panelFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)

						btn := material.Button(self.style.theme, &self.deleteBtn, "\uE74D")
						btn.Background.A = 0
						btn.Color = self.style.darkmode.currentColor.Fg
						btn.TextSize = unit.Sp(20)

						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Engine"][self.language], "OpenJTalk")),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Speaker"][self.language], self.info.Speaker)),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Rigid(btn.Layout),
						)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Speed"][self.language], self.coes[1])),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Volume"][self.language], self.coes[2])),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Pitch"][self.language], self.coes[3])),
						)
					}),
				)
			}),
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
		)
	}
	return func(gtx C) D {
		d, _ := self.expand.Update(gtx, headFunc, panelFunc)
		return d
	}
}
func (self *PresetBoxOpenJTalk) build_LabelBox(gtx C, title string, content string) W {
	return func(gtx C) D {
		return BorderBox(
			self.style,
			self.style.darkmode.currentColor.Fg,
			self.style.darkmode.currentColor.IdleBg,
			title,
		).Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return material.Label(self.style.theme, unit.Sp(18), content).Layout(gtx)
				}),
			)
		})
	}
}
func (self *PresetBoxOpenJTalk) API_GetName() string             { return self.info.Name }
func (self *PresetBoxOpenJTalk) API_SetLanguage(language string) { self.language = language }
func (self *PresetBoxOpenJTalk) Language() LangPack {
	return LangPack{
		"Engine": LangZone{
			"Chinese":  "引擎",
			"English":  "Engine",
			"Japanese": "エンジン",
		},
		"Speaker": LangZone{
			"Chinese":  "讲者",
			"English":  "Speaker",
			"Japanese": "話者",
		},
		"Speed": LangZone{
			"Chinese":  "语速",
			"English":  "Speed",
			"Japanese": "話速",
		},
		"Pitch": LangZone{
			"Chinese":  "音高",
			"English":  "Pitch",
			"Japanese": "ピッチ",
		},
		"Volume": LangZone{
			"Chinese":  "音量",
			"English":  "Volume",
			"Japanese": "音量",
		},
	}
}

// presetBox VoiceVox ------------------------------------------------------------------------------
type PresetBoxVoiceVox struct {
	style     Style
	info      config.Preset_VoiceVox
	expand    ExpandPanel
	coes      []string
	checkBtn  widget.Clickable
	deleteBtn widget.Clickable

	language      string
	languageTable LangPack
}

func New_PresetBoxVoiceVox(style Style, info config.Preset_VoiceVox) *PresetBoxVoiceVox {
	p := PresetBoxVoiceVox{
		style:     style,
		info:      info,
		expand:    *New_ExpandPanel(style, 110), // 因增加了一行，高度从 110 调整至 160
		checkBtn:  widget.Clickable{},
		deleteBtn: widget.Clickable{},
	}

	speak := info.Speaker
	speed := strconv.FormatFloat(info.Speed*100, 'f', 0, 64)
	volume := strconv.FormatFloat(info.Volume*100, 'f', 0, 64)
	pitch := strconv.FormatFloat(info.Pitch*100, 'f', 0, 64)
	intonation := strconv.FormatFloat(info.Intonation*100, 'f', 0, 64)

	// 索引对应关系：0:讲者, 1:语速, 2:音量, 3:音高, 4:顿音
	p.coes = []string{speak, speed, volume, pitch, intonation}
	p.language = "Chinese"
	p.languageTable = p.Language()

	return &p

}

func (self *PresetBoxVoiceVox) Update(gtx C, parent *SoundSetting) W {
	// 点选按钮功能
	if self.checkBtn.Clicked(gtx) {
		// parent.tool_SetPreset(self.info.Name)
		parent.tool_SetCurPreset(self.info)
	}
	var checkIcon string
	if parent.preset_Enum.Value == self.info.Name {
		checkIcon = "\uE73D"
	} else {
		checkIcon = "\uE739"
	}

	// 删除按钮功能
	if self.deleteBtn.Clicked(gtx) {
		parent.tool_DeletePreset(self.info.Name)
	}

	// 点选按钮样式
	checkBtnStyle := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(22), gtx.Sp(24))
		btn := material.Button(self.style.theme, &self.checkBtn, checkIcon)
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
				layout.Rigid(checkBtnStyle),
				layout.Rigid(Spacer(gtx.Dp(5), 0)),
				layout.Rigid(material.Label(self.style.theme, unit.Sp(20), self.info.Name).Layout),
				layout.Flexed(1, Flexer()),
			)
		})
	}

	panelFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// 第一行：引擎、讲者、删除按钮
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)

						btn := material.Button(self.style.theme, &self.deleteBtn, "\uE74D")
						btn.Background.A = 0
						btn.Color = self.style.darkmode.currentColor.Fg
						btn.TextSize = unit.Sp(20)

						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Engine"][self.language], "VoiceVox")),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Speaker"][self.language], self.info.Speaker)),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Rigid(btn.Layout),
						)
					}),
					// 第二行：语速、音量
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Speed"][self.language], self.coes[1])),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Pitch"][self.language], self.coes[3])),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Intonation"][self.language], self.coes[4])),
							layout.Rigid(Spacer(gtx.Dp(10), 0)),
							layout.Flexed(1, self.build_LabelBox(gtx, self.languageTable["Volume"][self.language], self.coes[2])),
						)
					}),
				)
			}),
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
		)
	}
	return func(gtx C) D {
		d, _ := self.expand.Update(gtx, headFunc, panelFunc)
		return d
	}
}
func (self *PresetBoxVoiceVox) build_LabelBox(gtx C, title string, content string) W {
	return func(gtx C) D {
		return BorderBox(
			self.style,
			self.style.darkmode.currentColor.Fg,
			self.style.darkmode.currentColor.IdleBg,
			title,
		).Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return material.Label(self.style.theme, unit.Sp(18), content).Layout(gtx)
				}),
			)
		})
	}
}
func (self *PresetBoxVoiceVox) API_GetName() string             { return self.info.Name }
func (self *PresetBoxVoiceVox) API_SetLanguage(language string) { self.language = language }
func (self *PresetBoxVoiceVox) Language() LangPack {
	return LangPack{
		"Engine": LangZone{
			"Chinese":  "引擎",
			"English":  "Engine",
			"Japanese": "エンジン",
		},
		"Speaker": LangZone{
			"Chinese":  "讲者",
			"English":  "Speaker",
			"Japanese": "話者",
		},
		"Speed": LangZone{
			"Chinese":  "语速",
			"English":  "Speed",
			"Japanese": "話速",
		},
		"Pitch": LangZone{
			"Chinese":  "音高",
			"English":  "Pitch",
			"Japanese": "ピッチ",
		},
		"Intonation": LangZone{
			"Chinese":  "语调",
			"English":  "Intonation",
			"Japanese": "抑揚",
		},
		"Volume": LangZone{
			"Chinese":  "音量",
			"English":  "Volume",
			"Japanese": "音量",
		},
	}
}

// PresetBoxTeto 预设 ------------------------------------------------------------------------------
type PresetBoxTeto struct {
	style     Style
	checkBtn  widget.Clickable
	deleteBtn widget.Clickable
	expand    ExpandPanel
	info      config.Preset_Teto
}

func New_PresetBoxTeto(style Style, info config.Preset_Teto) *PresetBoxTeto {
	p := PresetBoxTeto{
		style:     style,
		checkBtn:  widget.Clickable{},
		deleteBtn: widget.Clickable{},
		expand:    *New_ExpandPanel(style, 60),
		info:      info,
	}
	return &p
}
func (self *PresetBoxTeto) Update(gtx C, parent *SoundSetting) W {

	// 点选按钮功能
	if self.checkBtn.Clicked(gtx) {
		parent.preset_Enum.Value = self.info.Name
	}
	var checkIcon string
	if parent.preset_Enum.Value == self.info.Name {
		checkIcon = "\uE73D"
	} else {
		checkIcon = "\uE739"
	}

	// 删除按钮功能
	if self.deleteBtn.Clicked(gtx) {
		parent.tool_DeletePreset(self.info.Name)
	}

	// 点选按钮样式
	checkBtnStyle := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(22), gtx.Sp(24))
		btn := material.Button(self.style.theme, &self.checkBtn, checkIcon)
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
				layout.Rigid(checkBtnStyle),
				layout.Rigid(Spacer(gtx.Dp(5), 0)),
				layout.Rigid(material.Label(self.style.theme, unit.Sp(20), self.info.Name).Layout),
				layout.Flexed(1, Flexer()),
			)
		})
	}

	panelFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(Spacer(gtx.Dp(40), 0)),
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Dp(50)
						gtx.Constraints.Max.Y = gtx.Dp(50)

						btn := material.Button(self.style.theme, &self.deleteBtn, "\uE74D")
						btn.Background.A = 0
						btn.Color = self.style.darkmode.currentColor.Fg
						btn.TextSize = unit.Sp(20)

						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, Flexer()),
							layout.Rigid(btn.Layout),
						)
					}),
				)
			}),
			layout.Rigid(Spacer(gtx.Dp(15), 0)),
		)
	}
	return func(gtx C) D {
		d, _ := self.expand.Update(gtx, headFunc, panelFunc)
		return d
	}

}

type PresetBox struct {
	style          Style
	name           string
	isExpand       bool
	expandBtn      widget.Clickable
	checkBtn       widget.Clickable
	deleteBtn      widget.Clickable
	checkIcon      string
	presetList     widget.List
	presetDetail   config.PresetDetail
	listHeight_tar float64
	listHeight_cur float64
	parent         *SoundSetting
}

// 旧的默认预设(已弃用) ------------------------------------------------------------------------------------------
// func New_PresetBox(style Style, preset config.PresetDetail) *PresetBox {
// 	return &PresetBox{
// 		style:        style,
// 		expandBtn:    widget.Clickable{},
// 		deleteBtn:    widget.Clickable{},
// 		presetDetail: preset,
// 		presetList:   widget.List{},
// 		name:         preset.Name,
// 	}
// }
// func (self *PresetBox) Default() {
// 	self.listHeight_tar = 0
// 	self.listHeight_cur = 0
// 	self.isExpand = false
// }
// func (self *PresetBox) Update(gtx C, parent *SoundSetting, enum *widget.Enum) D {
// 	self.parent = parent
// 	// 按钮操作
// 	// 点选按钮
// 	if self.checkBtn.Clicked(gtx) {
// 		// parent.tool_SetCurPreset(self.API_GetPreset())
// 		// enum.Value = self.name
// 		// config.API_SOUND_SetCurrentPresetTitle(self.name)

// 		// preset := self.API_GetPreset()
// 		// preset.Name = "CURRENT"
// 		// csvRow := fumovoice.PresetToCSV(preset)
// 		// AppendToCSV(csvRow)
// 	}
// 	if enum.Value == self.name {
// 		self.checkIcon = "\uE73D"
// 	} else {
// 		self.checkIcon = "\uE739"
// 	}

// 	// 展开按钮
// 	if self.expandBtn.Hovered() {
// 	}
// 	if self.expandBtn.Clicked(gtx) {
// 		self.isExpand = !self.isExpand
// 	}

// 	if self.isExpand {
// 		self.listHeight_tar = float64(gtx.Dp(180))
// 	} else {
// 		self.listHeight_tar = 0
// 	}

// 	// 删除按钮
// 	if self.deleteBtn.Clicked(gtx) {
// 		if *self.presetDetail.Default == false {
// 			self.parent.tool_DeletePreset(self.name)
// 		}
// 	}
// 	// 详情列表高度更新
// 	self.listHeight_cur = easing(gtx, self.listHeight_cur, self.listHeight_tar, 0.3)
// 	expandFunc := func(gtx C) D {
// 		if self.listHeight_cur == 0 {
// 			return D{}
// 		}
// 		cl := clip.Rect{Max: image.Point{X: gtx.Constraints.Max.X, Y: int(self.listHeight_cur)}}.Push(gtx.Ops)
// 		dism := self.label_Detail(gtx)
// 		cl.Pop()
// 		return D{Size: image.Point{X: dism.Size.X, Y: int(self.listHeight_cur)}, Baseline: dism.Baseline}
// 	}

// 	// 返回布局
// 	return layout.Stack{}.Layout(gtx,
// 		layout.Expanded(func(gtx C) D {
// 			pt := gtx.Constraints.Min
// 			paint.FillShape(gtx.Ops, self.style.darkmode.currentColor.IdleBg, clip.RRect{
// 				Rect: image.Rectangle{Max: pt},
// 				SE:   gtx.Dp(4),
// 				SW:   gtx.Dp(4),
// 				NE:   gtx.Dp(4),
// 				NW:   gtx.Dp(4),
// 			}.Op(gtx.Ops))
// 			return D{}
// 		}),
// 		layout.Stacked(func(gtx C) D {
// 			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
// 				layout.Rigid(func(gtx C) D { return self.build_presetBox(gtx, self.name) }),
// 				layout.Rigid(expandFunc),
// 			)
// 		}),
// 	)
// }

// // 创建预设按钮 包括点选按钮 与 展开按钮
// func (self *PresetBox) build_presetBox(gtx C, text string) D {
// 	// 展开按钮操作

// 	var expandIcon string
// 	if self.isExpand {
// 		expandIcon = "\uf090"
// 	} else {
// 		expandIcon = "\uf08e"
// 	}

// 	// 点选按钮样式
// 	checkBtnStyle := func(gtx C) D {
// 		icon_label := material.Label(self.style.theme, unit.Sp(20), self.checkIcon)
// 		btn := material.Button(self.style.theme, &self.checkBtn, "")
// 		btn.Background = color.NRGBA{A: 0}
// 		return layout.Stack{}.Layout(gtx,
// 			layout.Expanded(func(gtx C) D {
// 				gtx.Constraints.Max = image.Pt(gtx.Dp(20), gtx.Dp(20))
// 				gtx.Constraints.Min = image.Pt(gtx.Dp(20), gtx.Dp(20))
// 				return btn.Layout(gtx)
// 			}),
// 			layout.Stacked(func(gtx C) D {
// 				return icon_label.Layout(gtx)
// 			}),
// 		)
// 	}

// 	// 展开按钮样式
// 	stackBtnStyle := func(gtx C) D {
// 		text_label := material.Label(self.style.theme, unit.Sp(18), text)
// 		expandIcon_label := material.Label(self.style.theme, unit.Sp(20), expandIcon)
// 		return layout.Center.Layout(gtx, func(gtx C) D {
// 			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
// 				layout.Rigid(Spacer(gtx.Dp(40), 0)),
// 				layout.Rigid(text_label.Layout),
// 				layout.Flexed(1, Flexer()),
// 				layout.Rigid(expandIcon_label.Layout),
// 				layout.Rigid(Spacer(gtx.Dp(15), 0)),
// 			)
// 		})
// 	}
// 	stackBtnFunc := self.build_Button(
// 		stackBtnStyle,
// 		&self.expandBtn,
// 		layout.Inset{},
// 		color.NRGBA{A: 0},
// 		gtx.Constraints.Max.X, gtx.Dp(40),
// 	)
// 	return layout.Center.Layout(gtx, func(gtx C) D {
// 		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
// 			layout.Rigid(func(gtx C) D {
// 				return layout.Stack{}.Layout(gtx,
// 					layout.Expanded(func(gtx C) D { return stackBtnFunc(gtx) }),
// 					layout.Expanded(func(gtx C) D {
// 						return layout.Center.Layout(gtx, func(gtx C) D {
// 							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
// 								layout.Rigid(Spacer(gtx.Dp(15), 0)),
// 								layout.Rigid(func(gtx C) D { return checkBtnStyle(gtx) }),
// 								layout.Flexed(1, Flexer()),
// 							)
// 						})
// 					}),
// 				)
// 			}),
// 		)
// 	})
// }

// // 创建适用于预设的按钮
// func (self *PresetBox) build_Button(
// 	stackedWidget layout.Widget,
// 	clickable *widget.Clickable,
// 	inset layout.Inset,
// 	c color.NRGBA,
// 	width int, height int,
// ) layout.Widget {
// 	return func(gtx C) D {
// 		btn := material.Button(self.style.theme, clickable, "")
// 		btn.Background = c
// 		return inset.Layout(gtx, func(gtx C) D {
// 			return layout.Stack{}.Layout(gtx,
// 				layout.Expanded(func(gtx C) D {
// 					gtx.Constraints.Min.X = width
// 					gtx.Constraints.Min.Y = height
// 					gtx.Constraints.Max.Y = height
// 					return btn.Layout(gtx)
// 				}),
// 				layout.Stacked(func(gtx C) D {
// 					gtx.Constraints.Min.X = width
// 					gtx.Constraints.Min.Y = height
// 					gtx.Constraints.Max.Y = height
// 					return layout.Flex{
// 						Axis:      layout.Horizontal,
// 						Alignment: layout.End,
// 					}.Layout(gtx,
// 						layout.Rigid(stackedWidget),
// 					)
// 				}),
// 			)
// 		})
// 	}
// }

// // 详情页部分
// func (self *PresetBox) label_Detail(gtx C) D {
// 	makeLabel := func(txt string) D {
// 		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(0), Right: unit.Dp(15)}.Layout(gtx, material.Label(self.style.theme, unit.Sp(14), txt).Layout)
// 	}
// 	speed := strconv.Itoa(*self.presetDetail.Speed)
// 	volume := strconv.Itoa(*self.presetDetail.Volume)
// 	pitch := strconv.Itoa(*self.presetDetail.Pitch)
// 	accent := strconv.Itoa(*self.presetDetail.Accent)
// 	quality := strconv.Itoa(*self.presetDetail.Quality)
// 	intonation := strconv.Itoa(*self.presetDetail.Intonation)

// 	btnStyle := func(gtx C) D {
// 		icon_label := material.Label(self.style.theme, unit.Sp(20), "\uE74D")
// 		btn := material.Button(self.style.theme, &self.deleteBtn, "")
// 		btn.Background = color.NRGBA{A: 0}
// 		return layout.Center.Layout(gtx, func(gtx C) D {
// 			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
// 				layout.Expanded(func(gtx C) D {
// 					gtx.Constraints.Max = image.Pt(gtx.Dp(30), gtx.Dp(30))
// 					gtx.Constraints.Min = image.Pt(gtx.Dp(30), gtx.Dp(30))
// 					return btn.Layout(gtx)
// 				}),
// 				layout.Stacked(func(gtx C) D {
// 					return icon_label.Layout(gtx)
// 				}),
// 			)
// 		})
// 	}

// 	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 		layout.Rigid(Spacer(gtx.Dp(40), 0)),
// 		layout.Rigid(func(gtx C) D {
// 			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 				layout.Rigid(func(gtx C) D {
// 					return makeLabel(self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CoeEngine + ": " + self.presetDetail.Engine)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return makeLabel(self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CoeType + ": " + self.presetDetail.VoiceType)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return makeLabel(self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CoeSpeed + ": " + speed)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return makeLabel(self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CoeVolume + ": " + volume)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return makeLabel(self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CoePitch + ": " + pitch)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return makeLabel(self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CoeAccent + ": " + accent)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return makeLabel(self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CoeQuality + ": " + quality)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return makeLabel(self.parent.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Voice_CoeIntonation + ": " + intonation)
// 				}),
// 			)
// 		}),
// 		layout.Flexed(1, Flexer()),
// 		layout.Rigid(func(gtx C) D {
// 			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 				layout.Rigid(btnStyle),
// 				layout.Rigid(Spacer(gtx.Dp(10), 0)),
// 			)
// 		}),
// 	)
// }

// // 返回当前的预设参数
// func (self *PresetBox) API_GetPreset() *fumovoice.AquestalkPreset {
// 	data := config.API_SOUND_GetPreset(self.name)
// 	return &fumovoice.AquestalkPreset{
// 		Name:       data.Name,
// 		MonoTone:   data.MonoTone,
// 		Engine:     data.Engine,
// 		VoiceType:  data.VoiceType,
// 		Speed:      data.Speed,
// 		Volume:     data.Volume,
// 		Pitch:      data.Pitch,
// 		Accent:     data.Accent,
// 		Quality:    data.Quality,
// 		Intonation: data.Intonation,
// 		Note:       data.Note,
// 	}
// }
