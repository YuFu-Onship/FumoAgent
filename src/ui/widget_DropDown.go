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

type DropDown struct {
	style     Style
	ease      Ani_Easing
	easeValue float32

	options []string

	isOpen     bool
	btn_Expand widget.Clickable
	expandIcon string

	choiceEnum    widget.Enum
	choiceBtnList []*DropDownChoiceBtn
	expandList    widget.List
}

func New_DropDown(style Style, options []string, defaultValue string) *DropDown {
	list := widget.List{}
	list.Axis = layout.Vertical

	dd := DropDown{
		style: style,
		ease:  *New_Ani_Esaing(0, 0, 0.3, true),

		options: options,

		isOpen:     false,
		btn_Expand: widget.Clickable{},
		expandIcon: "\uf090",

		choiceEnum: widget.Enum{Value: defaultValue},
		expandList: list,
	}

	cbl := []*DropDownChoiceBtn{}
	for _, o := range options {
		cbl = append(cbl, New_DropDownChoiceBtn(style, o))
	}
	dd.choiceBtnList = cbl
	dd.ease.API_SetMaxValue(float32(len(options) * 50))

	return &dd
}
func (self *DropDown) Default() {
}

func (self *DropDown) Update(gtx C) D {

	// 操作逻辑
	if self.isOpen {
		self.expandIcon = "\uf090"
		self.ease.API_SetDirection(false)
	} else {
		self.expandIcon = "\uf08e"
		self.ease.API_SetDirection(true)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(self.layout_Top),
		layout.Rigid(Spacer(0, gtx.Dp(4))),
		layout.Rigid(self.layout_Expand),
	)
}

// 顶部栏
func (self *DropDown) layout_Top(gtx C) D {
	if self.btn_Expand.Clicked(gtx) {
		self.isOpen = !self.isOpen
	}

	top_layout := func(gtx C) D {
		text_label := material.Label(self.style.theme, unit.Sp(18), self.choiceEnum.Value)
		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(Spacer(gtx.Dp(15), 0)),
				layout.Rigid(text_label.Layout),
				layout.Flexed(1, Flexer()),
				layout.Rigid(material.Label(self.style.theme, unit.Sp(18), self.expandIcon).Layout),
				layout.Rigid(Spacer(gtx.Dp(15), 0)),
			)
		})
		// })
	}
	btn := self.build_Button(top_layout, &self.btn_Expand, layout.Inset{}, color.NRGBA{A: 0}, gtx.Constraints.Max.X, gtx.Dp(40))
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			pt := gtx.Constraints.Min
			paint.FillShape(gtx.Ops, self.style.darkmode.currentColor.IdleBg, clip.RRect{
				Rect: image.Rectangle{Max: pt},
				SE:   gtx.Dp(4),
				SW:   gtx.Dp(4),
				NE:   gtx.Dp(4),
				NW:   gtx.Dp(4),
			}.Op(gtx.Ops))
			return D{}
		}),
		layout.Stacked(btn),
	)
}

// 下拉栏
func (self *DropDown) layout_Expand(gtx C) D {
	// 展开栏高度
	self.easeValue = self.ease.Update(gtx)
	value := gtx.Dp(unit.Dp(self.easeValue))
	// 控制下拉栏高度
	cl := clip.Rect{
		Max: image.Point{X: gtx.Constraints.Max.X, Y: value},
	}.Push(gtx.Ops)

	//下拉栏内容
	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(self.build_RoundRect(gtx.Dp(8), self.style.darkmode.currentColor.IdleBg)),
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(Spacer(0, gtx.Dp(2))),
				layout.Rigid(func(gtx C) D {
					return self.expandList.Layout(gtx, len(self.choiceBtnList), func(gtx C, index int) D {
						return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4), Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx C) D {
							return self.choiceBtnList[index].Update(gtx, self)
						})
					})
				}),
				layout.Rigid(Spacer(0, gtx.Dp(2))),
			)
		}),
	)

	cl.Pop()
	return D{
		Size: image.Point{
			X: dims.Size.X,
			Y: value,
		},
		Baseline: dims.Baseline,
	}
}

// 创建按钮
func (self *DropDown) build_Button(
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

// 绘制圆角矩形
func (self *DropDown) build_RoundRect(radius int, c color.NRGBA) layout.Widget {
	return func(gtx C) D {
		pt := gtx.Constraints.Min
		paint.FillShape(gtx.Ops, c, clip.RRect{
			Rect: image.Rectangle{Max: pt},
			SE:   radius,
			SW:   radius,
			NE:   radius,
			NW:   radius,
		}.Op(gtx.Ops))
		return D{}
	}
}

// api
func (self *DropDown) API_GetValue() string {
	return self.choiceEnum.Value
}

// 选项按钮
type DropDownChoiceBtn struct {
	style Style
	btn   widget.Clickable
	title string
	value bool
}

func New_DropDownChoiceBtn(style Style, title string) *DropDownChoiceBtn {
	ddcb := DropDownChoiceBtn{
		style: style,
		btn:   widget.Clickable{},
		title: title,
		value: true,
	}
	return &ddcb
}

func (self *DropDownChoiceBtn) Default() {}
func (self *DropDownChoiceBtn) Update(gtx C, parent *DropDown) D {
	// 逻辑
	if self.btn.Clicked(gtx) {
		parent.choiceEnum.Value = self.title
		parent.isOpen = false
	}
	// 布局
	btnFunc := func(gtx C) D {
		btn := material.Button(self.style.theme, &self.btn, "")
		btn.Background = color.NRGBA{A: 0}
		return btn.Layout(gtx)
	}
	titleLabel := material.Label(self.style.theme, unit.Sp(18), self.title)
	iconLabel := material.Label(self.style.theme, unit.Sp(20), "\uE73E")
	if self.title != parent.choiceEnum.Value {
		iconLabel.Color = color.NRGBA{A: 0}
	}
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(btnFunc),
			layout.Stacked(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(Spacer(gtx.Dp(15), 0)),
					layout.Rigid(iconLabel.Layout),
					layout.Rigid(Spacer(gtx.Dp(5), 0)),
					layout.Rigid(titleLabel.Layout),
					layout.Flexed(1, Flexer()),
				)
			}),
		)
	})
}

func (self *DropDownChoiceBtn) API_SetValue(value bool) {
	self.value = value
}
