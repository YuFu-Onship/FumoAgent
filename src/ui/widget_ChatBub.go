package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// 消息气泡控件 --------------------------------------------------------
type MsgBox struct {
	style Style

	title string

	tar_msg []rune
	cur_msg string
	cur_len int
	tar_len int
	step    int

	bg color.NRGBA
	fg color.NRGBA

	time_limit int
	isfinish   bool
	// editor     widget.Editor
}

func New_MsgBox(style Style, msg string) *MsgBox {
	self := MsgBox{
		style: style,
		// editor:  widget.Editor{},
		tar_msg: []rune(msg),
		cur_msg: "",
		cur_len: 0,
		tar_len: len(msg),
		step:    10,

		fg: style.theme.Fg,
		bg: style.theme.Bg,

		isfinish: false,
	}

	desiredFrames := 120
	self.step = len(self.tar_msg) / desiredFrames
	self.step = max(len(self.tar_msg)/desiredFrames, 1)

	// self.editor.ReadOnly = true
	return &self
}

// 界面更新
func (self *MsgBox) Update(gtx C) D {
	self.tar_len = len(self.tar_msg)
	if self.cur_len != self.tar_len {
		// gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
		self.cur_len += self.step
		if self.cur_len >= self.tar_len {
			self.cur_len = self.tar_len
		}
		// self.editor.SetText(self.cur_msg)
		self.isfinish = false
	} else {
		self.isfinish = true
	}
	self.cur_msg = string(self.tar_msg[:self.cur_len])
	return self.buildBubble(gtx)
}

// 构建聊天气泡
func (self *MsgBox) buildBubble(gtx C) D {
	// 	editor := material.Editor(self.style.theme, &self.editor, "")
	// 	editor.Font.Weight = 900
	// 	editor.TextSize = unit.Sp(16)
	// 	editor.Color = self.fg
	// 	editor.Editor.Alignment = text.Alignment(layout.Start)

	label := material.Label(self.style.theme, unit.Sp(16), self.cur_msg)
	label.Color = self.fg

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			DrawRRect(gtx, gtx.Constraints.Min, self.bg, gtx.Dp(8))
			return D{}
		}),
		layout.Stacked(func(gtx C) D {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
				// return editor.Layout(gtx)
				return label.Layout(gtx)
			})
		}),
	)
}

// 获取是否绘制信息
func (self *MsgBox) API_GetState() bool {
	return self.isfinish
}

// 设置标签属性
func (self *MsgBox) API_SetLabel(label string) { self.title = label }

// 获取标签属性
func (self *MsgBox) API_GetLabel() string { return self.title }

// 设置是否有打字机效果
func (self *MsgBox) API_IsType(value bool) {
	if !value {
		self.cur_len = self.tar_len
		self.cur_msg = string(self.tar_msg)
		self.isfinish = true
	}
}

// 修改颜色
func (self *MsgBox) API_SetColor(fg color.NRGBA, bg color.NRGBA) { self.fg = fg; self.bg = bg }
func (self *MsgBox) API_SetBg(bg color.NRGBA)                    { self.bg = bg }
func (self *MsgBox) API_SetFg(fg color.NRGBA)                    { self.fg = fg }
