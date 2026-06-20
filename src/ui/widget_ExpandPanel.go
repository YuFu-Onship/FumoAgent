package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ExpandPanel struct {
	style             Style
	topLayout         D
	detailLayout      D
	panelHeight_cur   float64
	panelHeight_tar   float64
	panelHeight       float64
	isExpand          bool
	isArrowVisibility bool
	expandBtn         widget.Clickable

	onExpandClicked func()
}

func New_ExpandPanel(style Style, height float64) *ExpandPanel {
	return &ExpandPanel{
		style:             style,
		isExpand:          false,
		isArrowVisibility: true,
		expandBtn:         widget.Clickable{},

		panelHeight:     height,
		panelHeight_cur: 0,
		panelHeight_tar: 0,
		onExpandClicked: func() {},
	}
}

func (self *ExpandPanel) Default() {
	self.isExpand = false
	self.panelHeight_cur = 0
	self.panelHeight_tar = 0
}

func (self *ExpandPanel) Update(gtx C, headLayout W, panelLayout W) (D, int) {
	var expandIcon string
	var animHeight int
	var realDims D

	// 1. 处理点击事件：切换展开/折叠状态
	if self.expandBtn.Clicked(gtx) {
		self.isExpand = !self.isExpand
		if self.isExpand {
			self.panelHeight_tar = self.panelHeight // 展开目标高度
		} else {
			self.panelHeight_tar = 0 // 折叠目标高度
		}
	}

	if self.panelHeight_cur != self.panelHeight_tar {
		gtx.Execute(op.InvalidateCmd{})
	}

	if self.isExpand {
		expandIcon = "\uf090"
	} else {
		expandIcon = "\uf08e"
	}
	if !self.isArrowVisibility {
		expandIcon = ""
	}

	// 2. 计算动画高度
	self.panelHeight_cur = easing(gtx, self.panelHeight_cur, self.panelHeight_tar, 0.3)

	topFunc := func(gtx C) D {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				btn := material.Button(self.style.theme, &self.expandBtn, "")
				btn.Background = self.style.darkmode.currentColor.IdleBg
				return btn.Layout(gtx)
			}),
			layout.Stacked(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, Flexer()),
					layout.Rigid(material.Label(self.style.theme, unit.Sp(20), expandIcon).Layout),
					layout.Rigid(Spacer(gtx.Dp(15), 0)),
				)
			}),
			layout.Stacked(headLayout),
		)
	}

	panelFunc := func(gtx C) D {
		if self.panelHeight_cur == 0 {
			return D{}
		}

		// var realDims D
		animHeight = gtx.Dp(unit.Dp(self.panelHeight_cur))

		cl := clip.Rect{
			Max: image.Point{X: gtx.Constraints.Max.X, Y: animHeight},
		}.Push(gtx.Ops)

		dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(Spacer(0, gtx.Dp(4))),
			layout.Rigid(func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						DrawRRect(gtx, image.Pt(gtx.Constraints.Max.X, animHeight-gtx.Dp(4)), self.style.darkmode.currentColor.IdleBg, 8)
						return D{}
					}),
					layout.Stacked(func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								realDims = panelLayout(gtx)
								return realDims
							}),
						)
					}),
				)
			}),
		)

		cl.Pop()
		return D{
			Size: image.Point{
				X: dims.Size.X,
				Y: animHeight,
			},
			Baseline: dims.Baseline,
		}
	}

	d := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(topFunc),
		layout.Rigid(panelFunc),
	)
	h := realDims.Size.Y
	if h == 0 {
		h = 1
	}
	return d, h
}

// 设置展开栏的高度
func (self *ExpandPanel) API_SetPanelHeight(height int) {
	self.panelHeight = float64(height)
	if self.isExpand {
		self.panelHeight_tar = float64(self.panelHeight)
	}
}

// 设置展开栏是否打开
func (self *ExpandPanel) API_SetPanelState(state bool) { self.isExpand = state }

// 得到展开栏是否打开
func (self *ExpandPanel) API_GetPanelState() bool { return self.isExpand }

// 展开按钮点击时的回调函数
func (self *ExpandPanel) API_OnExpandClicked(fn func()) {
	self.onExpandClicked = fn
}

// 设置是否显示右侧箭头标识
func (self *ExpandPanel) API_SetArrowVisibility(state bool) {
	self.isArrowVisibility = state
}

// 创建内部按钮控件
func (self *ExpandPanel) build_Button(
	stackedWidget layout.Widget,
	clickable *widget.Clickable,
	inset layout.Inset,
	c color.NRGBA,
	width int, height int,
) layout.Widget {
	return func(gtx C) D {
		btn := material.Button(self.style.theme, clickable, "")
		btn.Background = c
		gtx.Constraints.Min.X = width
		gtx.Constraints.Min.Y = height
		gtx.Constraints.Max.Y = height
		return inset.Layout(gtx, func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return btn.Layout(gtx)
				}),
				layout.Stacked(func(gtx C) D {
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
