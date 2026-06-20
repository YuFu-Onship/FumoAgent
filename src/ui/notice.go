package ui

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Notice struct {
	style Style

	notice_list  []*NoticeBox
	notice_limit int

	testBox NoticeBox
}

func New_Notice(style Style) *Notice {
	self := Notice{
		style:        style,
		notice_list:  []*NoticeBox{},
		notice_limit: 3,
		testBox:      *New_NoticeBox(style, "Title", "Desc", 2),
	}
	return &self
}

func (self *Notice) Update(gtx C, parent *Page) D {
	var children []layout.FlexChild
	if len(self.notice_list) > self.notice_limit {
		self.notice_list = self.notice_list[len(self.notice_list)-self.notice_limit:]
	}

	new_list := []*NoticeBox{}
	for _, n := range self.notice_list {
		if !n.API_GetResult() {
			new_list = append(new_list, n)
		}
	}
	self.notice_list = new_list

	if len(self.notice_list) > 0 {
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
		for i := range self.notice_list {
			box := self.notice_list[i]
			children = append(children, layout.Rigid(func(gtx C) D {
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
					return box.Update(gtx)
				})
			}))
		}
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, Flexer()),
			layout.Rigid(func(gtx C) D { return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...) }),
		)
	}
	return D{}
}

// 新添加通知
func (self *Notice) Add_Notice(level int, title string, desc string) {
	self.notice_list = append(self.notice_list, New_NoticeBox(self.style, title, desc, level))
	// if len(self.notice_list) > 3 {
	// 	self.notice_list = self.notice_list[len(self.notice_list)-3:]
	// }
}

// 通知弹出栏
type NoticeBox struct {
	style    Style
	closeBtn widget.Clickable
	title    string
	desc     string
	level    int

	fg color.NRGBA
	bg color.NRGBA

	timer       float32
	timer_limit float32
	isFinished  bool
}

func New_NoticeBox(style Style, title string, desc string, level int) *NoticeBox {
	self := NoticeBox{
		style:       style,
		closeBtn:    widget.Clickable{},
		title:       title,
		desc:        desc,
		level:       level,
		timer:       0,
		timer_limit: 100,
		isFinished:  false,
	}
	self.update_color()
	return &self
}

func (self *NoticeBox) Update(gtx C) D {
	// 功能更新
	self.timer = self.timer + 1
	if self.timer >= self.timer_limit {
		self.timer = self.timer_limit
		self.isFinished = true
	} else {
		// gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
	}

	if self.closeBtn.Clicked(gtx) {
		self.isFinished = true
	}

	// 样式更新
	self.update_color()
	size := image.Pt(gtx.Dp(300), gtx.Dp(60))
	gtx.Constraints.Max = size
	gtx.Constraints.Min = size

	btnFunc := func(gtx C) D {
		s := image.Pt(gtx.Dp(20), gtx.Sp(20))
		gtx.Constraints.Max = s
		gtx.Constraints.Min = s
		btn := material.Button(self.style.theme, &self.closeBtn, "\uE711")
		btn.Background.A = 0
		btn.Color = self.fg
		btn.TextSize = unit.Sp(16)
		btn.Inset = layout.Inset{}
		return btn.Layout(gtx)
	}

	label_title := material.Label(self.style.theme, unit.Sp(18), self.title)
	label_title.Color = self.fg
	label_desc := material.Label(self.style.theme, unit.Sp(14), self.desc)
	label_desc.Color = self.fg

	showFunc := func(gtx C) D {
		cl := clip.RRect{
			Rect: image.Rectangle{Max: size},
			SE:   gtx.Dp(8),
			SW:   gtx.Dp(8),
			NE:   gtx.Dp(8),
			NW:   gtx.Dp(8),
		}.Push(gtx.Ops)
		defer cl.Pop()

		DrawLine(
			gtx,
			0,
			0,
			float32(gtx.Constraints.Max.X)*self.timer/self.timer_limit,
			0,
			float32(gtx.Dp(6)),
			self.fg,
		)
		dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(Spacer(gtx.Dp(12), 0)),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Flexed(1, FlexerY()),
					layout.Rigid(label_title.Layout),
					layout.Rigid(label_desc.Layout),
					layout.Flexed(1, FlexerY()),
				)
			}),
			layout.Flexed(1, Flexer()),
			layout.Rigid(Spacer(gtx.Dp(4), 0)),
			layout.Rigid(btnFunc),
			layout.Rigid(Spacer(gtx.Dp(12), 0)),
		)

		return D{
			Size:     image.Pt(dims.Size.X, dims.Size.Y),
			Baseline: dims.Baseline,
		}
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			DrawRRect(
				gtx,
				image.Pt(gtx.Dp(300), gtx.Dp(60)),
				self.bg,
				gtx.Dp(8),
			)

			return D{}
		}),
		layout.Stacked(showFunc),
	)
}

func (self *NoticeBox) API_GetResult() bool {
	return self.isFinished
}

func (self *NoticeBox) update_color() {
	switch self.level {
	case 1:
		self.bg = self.style.darkmode.currentColor.Notice_Bg_1
		self.fg = self.style.darkmode.currentColor.Notice_Fg_1
	case 2:
		self.bg = self.style.darkmode.currentColor.Notice_Bg_2
		self.fg = self.style.darkmode.currentColor.Notice_Fg_2
	case 3:
		self.bg = self.style.darkmode.currentColor.Notice_Bg_3
		self.fg = self.style.darkmode.currentColor.Notice_Fg_3
	default:
		self.bg = self.style.darkmode.currentColor.Notice_Bg_1
		self.fg = self.style.darkmode.currentColor.Notice_Fg_1
	}
}
