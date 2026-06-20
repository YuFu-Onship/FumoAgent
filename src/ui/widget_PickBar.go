package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type PickBar struct {
	style   Style
	value   string
	index   int
	choices []string

	btn_left  widget.Clickable
	btn_right widget.Clickable
}

func New_PickBar(style Style) *PickBar {
	p := PickBar{
		style:   style,
		value:   "a",
		choices: []string{"a", "b", "c"},
		index:   0,
	}
	return &p
}
func (self *PickBar) Update(gtx C) D {
	// 按钮逻辑
	if self.btn_left.Clicked(gtx) {
		self.index -= 1
	}
	if self.btn_right.Clicked(gtx) {
		self.index += 1
	}

	if self.index <= 0 {
		self.index = 0
	} else if self.index >= (len(self.choices) - 1) {
		self.index = len(self.choices) - 1
	}

	// 元素样式
	labelFunc := func(gtx C) D {
		// gtx.Constraints.Min = image.Pt(gtx.Dp(150), gtx.Dp(20))
		// gtx.Constraints.Max = gtx.Constraints.Min
		gtx.Constraints.Max.X = gtx.Dp(200)
		gtx.Constraints.Min.X = gtx.Dp(200)
		l := material.Label(self.style.theme, unit.Sp(16), self.choices[self.index])
		return layout.Center.Layout(gtx, l.Layout)
	}

	leftFunc := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(20), gtx.Dp(20))
		gtx.Constraints.Max = gtx.Constraints.Min
		b := material.Button(self.style.theme, &self.btn_left, "\uE96F")
		b.Background.A = 0
		b.Color = self.style.theme.Fg
		b.Inset = layout.Inset{}
		b.TextSize = unit.Sp(20)
		return b.Layout(gtx)
	}

	rightFunc := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(20), gtx.Dp(20))
		gtx.Constraints.Max = gtx.Constraints.Min
		b := material.Button(self.style.theme, &self.btn_right, "\uE970")
		b.Background.A = 0
		b.Color = self.style.theme.Fg
		b.Inset = layout.Inset{}
		b.TextSize = unit.Sp(20)
		return b.Layout(gtx)
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(leftFunc),
		layout.Rigid(labelFunc),
		layout.Rigid(rightFunc),
	)
}
func (self *PickBar) API_GetValue() string {
	return self.choices[self.index]
}
func (self *PickBar) API_SetChoices(choices []string) {
	self.choices = choices
}
func (self *PickBar) API_SetValue(value string) {
	self.value = value
}
