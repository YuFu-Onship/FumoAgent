package ui

import (
	"image"
	"image/color"
	"myapp/src/config"
	"strings"
	"time"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// 定义消息结构体 --------------------------------------------------------
type Message struct {
	Text string
	IsMe bool
}
type MsgBub struct {
	time    string
	message []rune
	label   string
	tar_len int
	cur_len int
	step    int
}

type ApiMessage struct {
	Content  string `json:"content"`
	Role     string `json:"role"`
	CreateAt int64  `json:"create_at"`
}

// 类 --------------------------------------------------------------------
type ChatPage struct {
	style Style

	editor     *widget.Editor
	sendBtn    *widget.Clickable
	photoBtn   *widget.Clickable
	historyBtn *widget.Clickable
	clearBtn   *widget.Clickable
	expandBtn  *widget.Clickable

	list        widget.List
	msgbub_list []MsgBub

	MainPage *Page
	mode     string

	placeholderHeight_cur float64
	placeholderHeight_tar float64
	editorBoxHeight       float64
}

func New_ChatPage(style Style, mainPage *Page) *ChatPage {
	self := ChatPage{
		MainPage: mainPage,

		style: style,
		editor: &widget.Editor{
			SingleLine: false,
			Submit:     true,
			InputHint:  key.HintAny,
		},
		sendBtn:    &widget.Clickable{},
		photoBtn:   &widget.Clickable{},
		historyBtn: &widget.Clickable{},
		clearBtn:   &widget.Clickable{},
		expandBtn:  &widget.Clickable{},

		mode: "default",

		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		msgbub_list: []MsgBub{},

		// 页面动画
		placeholderHeight_cur: 0,
		placeholderHeight_tar: 0,
		editorBoxHeight:       60,
	}

	// 获取到历史消息
	cm := config.Config_Message{}
	messages := cm.Get_Msg()

	// 加载消息列表
	for _, m := range messages {
		if m[1] == "user" || m[1] == "ai" || m[1] == "system" {
			mb := self.BubInit(m[2], m[1])
			self.msgbub_list = append(self.msgbub_list, mb)
		}
	}
	self.list.ScrollToEnd = true
	return &self
}

// 不同的界面 ------------------------------------------------------------------------

// 默认状态
func (self *ChatPage) Default() {}

// 更新
func (self *ChatPage) Update(gtx C) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return self.layout_chat(gtx)
		}),
		layout.Rigid(self.layout_Buttom),
	)
}

// 聊天界面布局, 包含两种状态
func (self *ChatPage) layout_chat(gtx C) D {
	if self.mode == "default" {
		self.placeholderHeight_tar = float64(gtx.Constraints.Min.Y - gtx.Dp(100))
	} else {
		self.placeholderHeight_tar = 0
	}
	self.placeholderHeight_cur = easing(gtx, self.placeholderHeight_cur, self.placeholderHeight_tar, 0.3)
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, self.layout_MsgList),
				layout.Rigid(Spacer(0, gtx.Dp(100))),
			)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(Spacer(0, int(unit.Dp(self.placeholderHeight_cur)))),
				layout.Flexed(1, func(gtx C) D {
					return layout.Stack{}.Layout(gtx,
						layout.Expanded(func(gtx C) D { DrawRRect(gtx, gtx.Constraints.Max, self.style.theme.Bg, 4); return D{} }),
						layout.Stacked(func(gtx C) D {
							gtx.Constraints.Min = gtx.Constraints.Max
							return self.layout_Editor(gtx)
						}),
					)
				}),
			)
		}),
	)
}

// 布局 ------------------------------------------------------------------------
// 编辑栏布局
func (self *ChatPage) layout_Editor(gtx C) D {
	editorBox := material.Editor(self.style.theme, self.editor, "Enter")
	editorBox.TextSize = unit.Sp(18)
	editorBox.Font.Weight = 900
	editorBorder := widget.Border{
		Color:        self.style.darkmode.currentColor.IdleFg,
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(2),
	}
	self.handleEditorEvents(gtx)
	return layout.Inset{Bottom: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
		return editorBorder.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, editorBox.Layout)
		})
	})
}

// 消息列表布局
func (self *ChatPage) layout_MsgList(gtx C) D {
	result := true
	for i := range self.msgbub_list {
		result = self.BubUpdate(&self.msgbub_list[i]) && result
	}

	if !result {
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
		self.list.ScrollToEnd = true
	}

	list := material.List(self.style.theme, &self.list)
	list.ScrollbarStyle.Indicator.Color.A = 64
	list.ScrollbarStyle.Indicator.MinorWidth = unit.Dp(2)

	return list.Layout(gtx, len(self.msgbub_list), func(gtx C, index int) D {
		switch self.msgbub_list[index].label {
		case "user":
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(Spacer(gtx.Dp(30), 0)),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return self.BubBuild(gtx, &self.msgbub_list[index])
					})
				}),
			)
		case "ai":
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return self.BubBuild(gtx, &self.msgbub_list[index])
					})
				}),
				layout.Rigid(Spacer(gtx.Dp(30), 0)),
			)
		case "system":
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return self.BubBuild(gtx, &self.msgbub_list[index])
					})
				}),
				layout.Rigid(Spacer(gtx.Dp(30), 0)),
			)
		default:
			return D{}
		}
	})
}

// 底部功能键布局
func (self *ChatPage) layout_Buttom(gtx C) D {
	// 功能 -----------------------------------
	if self.sendBtn.Clicked(gtx) {
		self.performSend(gtx)
	}

	var expandIcon string = "\ue70e"
	if self.mode == "input" {
		expandIcon = "\ue70d"
	}
	if self.expandBtn.Clicked(gtx) {
		switch self.mode {
		case "default":
			self.mode = "input"
			self.placeholderHeight_cur = float64(gtx.Constraints.Max.Y)
			self.placeholderHeight_tar = 0
		case "input":
			self.mode = "default"
			self.placeholderHeight_cur = 0
		default:
		}
	}

	if self.clearBtn.Clicked(gtx) {
		self.msgbub_list = []MsgBub{}
		self.OUTAPI_ClearFunc()
	}

	// 样式 --------------------------------------
	sendBtnFunc := func(gtx C) D {
		size := image.Pt(gtx.Dp(80), gtx.Dp(40))
		gtx.Constraints.Max = size
		gtx.Constraints.Min = size

		btn := material.Button(self.style.theme, self.sendBtn, "\ue725")
		btn.Color = self.style.darkmode.currentColor.SelfMsgFg
		btn.Background = self.style.darkmode.currentColor.ContrastBg
		btn.Inset = layout.Inset{}
		btn.TextSize = unit.Sp(20)
		return btn.Layout(gtx)
	}

	photoBtnFunc := func(gtx C) D {
		size := image.Pt(gtx.Dp(40), gtx.Dp(40))
		gtx.Constraints.Max = size
		gtx.Constraints.Min = size

		btn := material.Button(self.style.theme, self.photoBtn, "\ue722")
		btn.Color = self.style.darkmode.currentColor.SelfMsgFg
		btn.Background = self.style.darkmode.currentColor.ContrastBg
		btn.Inset = layout.Inset{}
		btn.TextSize = unit.Sp(20)
		return btn.Layout(gtx)
	}

	clearBtnFunc := func(gtx C) D {
		size := image.Pt(gtx.Dp(40), gtx.Dp(40))
		gtx.Constraints.Max = size
		gtx.Constraints.Min = size

		btn := material.Button(self.style.theme, self.clearBtn, "\ue710")
		btn.Color = self.style.darkmode.currentColor.SelfMsgFg
		btn.Background = self.style.darkmode.currentColor.ContrastBg
		btn.Inset = layout.Inset{}
		btn.TextSize = unit.Sp(20)
		return btn.Layout(gtx)
	}

	expandBtnFunc := func(gtx C) D {
		size := image.Pt(gtx.Dp(40), gtx.Dp(40))
		gtx.Constraints.Max = size
		gtx.Constraints.Min = size

		btn := material.Button(self.style.theme, self.expandBtn, expandIcon)
		btn.Color = self.style.darkmode.currentColor.SelfMsgFg
		btn.Background = self.style.darkmode.currentColor.ContrastBg
		btn.Inset = layout.Inset{}
		btn.TextSize = unit.Sp(20)
		return btn.Layout(gtx)
	}

	return layout.Flex{
		Alignment: layout.Middle,
		Axis:      layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(clearBtnFunc),
				layout.Rigid(Spacer(gtx.Dp(5), 0)),
				layout.Rigid(photoBtnFunc),
				layout.Flexed(1, Flexer()),
				layout.Rigid(expandBtnFunc),
				layout.Rigid(Spacer(gtx.Dp(5), 0)),
				layout.Rigid(sendBtnFunc),
				layout.Rigid(Spacer(gtx.Dp(10), 0)),
			)
		}),
	)
}

// 键盘事件
func (self *ChatPage) handleEditorEvents(gtx C) {
	for {
		ev, ok := self.editor.Update(gtx)
		if !ok {
			break
		}
		switch ev.(type) {
		case widget.SubmitEvent:
			self.performSend(gtx)
		case widget.ChangeEvent:
		}
	}
}

// 发送消息事件,以及接收来自模型的回复
func (self *ChatPage) performSend(gtx C) {
	editorText := strings.TrimSpace(self.editor.Text())
	if editorText == "" {
		return
	}
	text := editorText

	// 用户消息
	self.editor.SetText("")
	self.mode = "default"
	self.list.ScrollToEnd = true
	self.msgbub_list = append(self.msgbub_list, self.BubInit(text, "user"))

	// ai消息
	go func() {
		appendGuiMessage := func(role string, content string) {
			self.msgbub_list = append(self.msgbub_list, self.BubInit(content, role))
			self.list.ScrollToEnd = true
			gtx.Execute(op.InvalidateCmd{At: gtx.Now})
		}
		_ = self.MainPage.trunk.Get_AIRes(text, "user", appendGuiMessage)
	}()
}

// api 相关 --------------------------------------------------------------------

// 清除消息记录
func (self *ChatPage) OUTAPI_ClearFunc() {
	self.MainPage.trunk.Handler_Message.API_Clear()
}

// 加载历史消息
func (self *ChatPage) loadHistoryMsg() []Message {
	content := self.MainPage.trunk.Handler_Message.API_GetContent()
	var messages []Message
	isMe := true
	for _, v := range content {
		switch v[1] {
		case "ai":
			isMe = false
			message := Message{IsMe: isMe, Text: v[2]}
			messages = append(messages, message)
		case "user":
			isMe = true
			message := Message{IsMe: isMe, Text: v[2]}
			messages = append(messages, message)
		case "system":
			isMe = false
			message := Message{IsMe: isMe, Text: v[2]}
			messages = append(messages, message)

		default:
		}
	}
	return messages
}

// 创建气泡
func (self *ChatPage) BubInit(text string, label string) MsgBub {
	mb := MsgBub{message: []rune(text), label: label}
	mb.tar_len = len(mb.message)
	mb.step = max(mb.tar_len/120, 1)
	switch label {
	case "ai":
		mb.cur_len = 0
	case "user":
		mb.cur_len = mb.tar_len
	case "system":
		mb.cur_len = mb.tar_len
	}
	return mb
}

// 更新气泡
func (self *ChatPage) BubUpdate(msg *MsgBub) bool {
	msg.cur_len += msg.step
	if msg.cur_len >= msg.tar_len {
		msg.cur_len = msg.tar_len
		return true
	}
	return false
}

// 绘制气泡
func (self *ChatPage) BubBuild(gtx C, msg *MsgBub) D {
	var fg color.NRGBA
	var bg color.NRGBA
	switch msg.label {
	case "user":
		fg = self.style.darkmode.currentColor.SelfMsgFg
		bg = self.style.darkmode.currentColor.SelfMsgBg
	case "ai":
		fg = self.style.darkmode.currentColor.OtherMsgFg
		bg = self.style.darkmode.currentColor.OtherMsgBg
	case "system":
	default:
	}
	if msg.label == "system" {
		label := material.Label(self.style.theme, unit.Sp(12), string(msg.message[:msg.cur_len]))
		label.Color = self.style.theme.Fg
		label.Color.A = 90
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, Flexer()),
			layout.Rigid(label.Layout),
			layout.Flexed(1, Flexer()),
		)
	} else {
		label := material.Label(self.style.theme, unit.Sp(16), string(msg.message[:msg.cur_len]))
		label.Color = fg
		return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					DrawRRect(gtx, gtx.Constraints.Min, bg, gtx.Dp(8))
					return D{}
				}),
				layout.Stacked(func(gtx C) D {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
						return label.Layout(gtx)
					})
				}),
			)
		})
	}
}

// API
func (self *ChatPage) API() map[string]any {
	return map[string]any{
		"add_msg": func(role string, msg string) {
			go func() {
				self.msgbub_list = append(self.msgbub_list, self.BubInit(msg, role))
				self.list.ScrollToEnd = true
				self.MainPage.Window.Invalidate()
			}()
		},
	}
}
