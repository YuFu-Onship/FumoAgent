package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// 内部控件
type DropDownChild struct {
	name   string
	widget layout.Widget
}

// 下拉选择框
type DropDownBox struct {
	index int
	list  widget.List
	style Style

	mainBtn     widget.Clickable
	mainBtnText string
	isOpen      bool

	// 选择框
	buttons       []widget.Clickable
	options       []string
	boxHeight     float64
	boxHeight_tar float64
	boxHeight_cur float64
}

func New_DropDownBox(style Style, options []string) *DropDownBox {
	lines := len(options)
	boxHeight := lines*34 + 8
	if boxHeight > 200 {
		boxHeight = 200
	}
	return &DropDownBox{
		index:       1,
		style:       style,
		mainBtn:     widget.Clickable{},
		mainBtnText: options[0],
		isOpen:      false,

		options:       options,
		boxHeight:     float64(boxHeight),
		boxHeight_tar: 0,
		boxHeight_cur: 0,
		buttons:       make([]widget.Clickable, len(options)),
		// list: widget.List{
		// 	List: layout.List{
		// 		Axis: layout.Vertical,
		// 	},
		// },
	}
}
func New_Label(gtx C) D {
	return D{}
}

// 调用主函数
func (self *DropDownBox) Update(gtx layout.Context) layout.Dimensions {
	// 使用 material.List 包装滚动状态
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.Spacing(4),
	}.Layout(gtx,
		layout.Rigid(self.mainButtonFunc),
		layout.Rigid(self.layout_choiceBox),
	)
}

// 选择框布局
func (self *DropDownBox) layout_choiceBox(gtx C) D {
	if self.isOpen {
		self.boxHeight_tar = self.boxHeight
	} else {
		self.boxHeight_tar = 0
	}
	self.boxHeight_cur = easing(
		gtx,
		self.boxHeight_cur,
		self.boxHeight_tar,
		0.2,
	)
	gtx.Constraints.Min.Y = int(self.boxHeight_cur)
	gtx.Constraints.Max.Y = int(self.boxHeight_cur)
	if self.boxHeight_cur <= 30 {
		return D{}
	} else {
		return self.choiceBox(gtx)
	}
}

// 主按钮
func (self *DropDownBox) mainButtonFunc(gtx C) D {
	btn := material.Button(self.style.theme, &self.mainBtn, "")
	btn.Background = color.NRGBA{A: 0}
	label := material.Label(self.style.theme, unit.Sp(20), self.mainBtnText)

	var expandIcon string
	if self.isOpen {
		expandIcon = "\uf090"
	} else {
		expandIcon = "\uf08e"
	}
	expandLabel := material.Label(self.style.theme, unit.Sp(20), expandIcon)

	if self.mainBtn.Clicked(gtx) {
		self.isOpen = !self.isOpen
	}

	if self.mainBtn.Hovered() {

	}
	return layout.Center.Layout(gtx, func(gtx C) D {

		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				gtx.Constraints.Min.Y = 40
				gtx.Constraints.Max.Y = 40
				return widget.Border{CornerRadius: unit.Dp(4), Width: unit.Dp(2), Color: self.style.theme.Fg}.Layout(gtx, btn.Layout)
			}),
			layout.Stacked(func(gtx C) D {
				gtx.Constraints.Min.Y = 40
				gtx.Constraints.Max.Y = 40
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D { gtx.Constraints.Max.X = 15; return D{Size: gtx.Constraints.Max} }),
					layout.Rigid(func(gtx C) D { return layout.Center.Layout(gtx, label.Layout) }),
					layout.Flexed(1, func(gtx C) D { return D{Size: gtx.Constraints.Max} }),
					layout.Rigid(func(gtx C) D { return layout.Center.Layout(gtx, expandLabel.Layout) }),
					layout.Rigid(func(gtx C) D { gtx.Constraints.Max.X = 15; return D{Size: gtx.Constraints.Max} }),
				)
			}),
		)
	})

}

// 选择框 --------------------------------------------------------------------
func (self *DropDownBox) choiceBox(gtx C) D {
	self.list.Axis = layout.Vertical
	self.list.Alignment = layout.Middle
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					// 选择框背景
					layout.Expanded(func(gtx C) D {
						size := gtx.Constraints.Min
						paint.FillShape(gtx.Ops, self.style.darkmode.currentColor.IdleBg, clip.RRect{Rect: image.Rectangle{Max: size}, SE: 8, SW: 8, NE: 8, NW: 8}.Op(gtx.Ops))
						return D{Size: gtx.Constraints.Min}
					}),
					// 选择框前景选项
					layout.Stacked(func(gtx C) D {
						return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
							return self.list.List.Layout(gtx, len(self.options), func(gtx layout.Context, index int) layout.Dimensions {
								btnState := &self.buttons[index]
								if btnState.Clicked(gtx) {
									self.mainBtnText = self.options[index]
									self.isOpen = false
								}
								return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx C) D {
									return layout.Center.Layout(gtx, func(gtx C) D {
										return self.buildButton(gtx, self.options[index], &self.buttons[index],
											layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0), Left: unit.Dp(0), Right: unit.Dp(0)},
											color.NRGBA{A: 0},
											gtx.Constraints.Max.X, 30,
										)
									})
								})
							})
						})
					}),
				)
			})
		}),
	)
}

// 构建按钮
func (self *DropDownBox) buildButton(
	gtx C,
	iconText string,
	clickable *widget.Clickable,
	inset layout.Inset,
	c color.NRGBA,
	width int,
	height int,
) D {
	btn := material.Button(self.style.theme, clickable, "")
	btn.Background = c
	label := material.Label(self.style.theme, unit.Sp(18), iconText)
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
					layout.Rigid(func(gtx C) D { gtx.Constraints.Min.X = gtx.Dp(unit.Dp(15)); return D{Size: gtx.Constraints.Min} }),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions { return layout.Center.Layout(gtx, label.Layout) }),
				)
			}),
		)
	})
}

// api ---------------------------------------------------------------
func (self *DropDownBox) API_SetValue(value string) {
	self.mainBtnText = value
}
func (self DropDownBox) API_GetValue() string {
	return self.mainBtnText
}
