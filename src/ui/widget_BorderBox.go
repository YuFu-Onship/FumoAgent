package ui

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type BorderBoxStyle struct {
	style Style
	title string
	fg    color.NRGBA
	bg    color.NRGBA
}

func BorderBox(style Style, fg color.NRGBA, bg color.NRGBA, title string) *BorderBoxStyle {
	bb := BorderBoxStyle{
		style: style,
		title: title,
		fg:    fg,
		bg:    bg,
	}
	return &bb
}

func (self *BorderBoxStyle) Layout(gtx C, w layout.Widget) D {
	labelSize := unit.Sp(14)

	return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return widget.Border{Width: unit.Dp(1), CornerRadius: unit.Dp(4), Color: self.style.theme.Fg}.Layout(gtx, func(gtx C) D {
					return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(4), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, w)
				})
			}),
			layout.Stacked(func(gtx C) D {
				trans := op.Offset(image.Pt(gtx.Dp(10), -gtx.Dp(8))).Push(gtx.Ops)
				defer trans.Pop()

				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						pt := gtx.Constraints.Min
						pt.Y = gtx.Dp(14)
						// paint.FillShape(gtx.Ops, self.bg, clip.Rect{Max: pt}.Op())
						DrawRRect(gtx, pt, self.style.darkmode.currentColor.IdleBg, 0)
						return D{}
					}),
					layout.Stacked(func(gtx C) D {
						return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
							return material.Label(self.style.theme, labelSize, self.title).Layout(gtx)
						})
					}),
				)
			}),
		)
	})
}

// func (self *BorderBoxStyle) Layout(gtx C, w layout.Widget) D {
// 	labelSize := unit.Sp(14)
// 	radius := gtx.Dp(4)      // 外部圆角
// 	strokeWidth := gtx.Dp(1) // 边框粗细，你之前代码里是 2px，可以根据需要调整

// 	return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
// 		return layout.Stack{}.Layout(gtx,
// 			// 【第一层：Expanded 自动同步尺寸，用来画边框背景】
// 			// 它会读取第二层组件撑开后的最终大小，并在底层完美对齐
// 			layout.Expanded(func(gtx C) D {
// 				size := gtx.Constraints.Min // 此时 Min 就是第二层子组件加上 Inset 后的精确大小

// 				// 1. 画外层边框大矩形
// 				paint.FillShape(gtx.Ops, self.style.theme.Fg, clip.RRect{
// 					Rect: image.Rectangle{Max: size},
// 					SE:   radius, SW: radius, NE: radius, NW: radius,
// 				}.Op(gtx.Ops))

// 				// 2. 挖空内层：向内缩进边框粗细，填充主题背景色
// 				innerSize := image.Rectangle{
// 					Min: image.Pt(strokeWidth, strokeWidth),
// 					Max: image.Pt(size.X-strokeWidth, size.Y-strokeWidth),
// 				}
// 				// 内部圆角稍微缩一点，视觉上更平滑
// 				innerRadius := radius - strokeWidth
// 				if innerRadius < 0 {
// 					innerRadius = 0
// 				}

// 				paint.FillShape(gtx.Ops, self.style.darkmode.currentColor.IdleBg, clip.RRect{
// 					Rect: innerSize,
// 					SE:   innerRadius, SW: innerRadius, NE: innerRadius, NW: innerRadius,
// 				}.Op(gtx.Ops))

// 				return D{}
// 			}),

// 			// 【第二层：Stacked 作为主导，负责用子组件把整个控件的尺寸撑开】
// 			layout.Stacked(func(gtx C) D {
// 				// 在这里正式、且只调用一次子组件 w
// 				return layout.Inset{
// 					Top:    unit.Dp(8),
// 					Bottom: unit.Dp(4),
// 					Left:   unit.Dp(16),
// 					Right:  unit.Dp(16),
// 				}.Layout(gtx, w)
// 			}),

// 			// 【第三层：Stacked 悬浮的标签文字】
// 			layout.Stacked(func(gtx C) D {
// 				trans := op.Offset(image.Pt(gtx.Dp(10), -gtx.Dp(8))).Push(gtx.Ops)
// 				defer trans.Pop()

// 				return layout.Stack{}.Layout(gtx,
// 					layout.Expanded(func(gtx C) D {
// 						pt := gtx.Constraints.Min
// 						pt.Y = gtx.Dp(14)
// 						// paint.FillShape(gtx.Ops, self.bg, clip.Rect{Max: pt}.Op())
// 						DrawRRect(gtx, pt, self.style.darkmode.currentColor.IdleBg, 0)
// 						return D{}
// 					}),
// 					layout.Stacked(func(gtx C) D {
// 						return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
// 							return material.Label(self.style.theme, labelSize, self.title).Layout(gtx)
// 						})
// 					}),
// 				)
// 			}),
// 		)
// 	})
// }

func (self *BorderBoxStyle) API_SetTitle(title string) {
	self.title = title
}

type OctagonBorder struct {
	Color color.NRGBA
	Width unit.Dp
	Cut   unit.Dp // 切角的大小（比如 4dp 或 8dp）
}

func (b OctagonBorder) Layout(gtx C, w layout.Widget) D {
	// 先布局内部子组件，获取实际大小
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	// 计算像素尺寸
	width := gtx.Dp(b.Width)
	cut := gtx.Dp(b.Cut)
	sz := f32.Point{X: float32(dims.Size.X), Y: float32(dims.Size.Y)}
	c := float32(cut)

	// 如果切角太大，做个保护
	if c > sz.X/2 {
		c = sz.X / 2
	}
	if c > sz.Y/2 {
		c = sz.Y / 2
	}

	// 使用 Path 绘制八边形外框
	var p clip.Path
	p.Begin(gtx.Ops)

	// 从左上角切角终点开始 (c, 0)
	p.MoveTo(f32.Pt(c, 0))
	p.LineTo(f32.Pt(sz.X-c, 0))    // 顶边
	p.LineTo(f32.Pt(sz.X, c))      // 右上切角
	p.LineTo(f32.Pt(sz.X, sz.Y-c)) // 右边
	p.LineTo(f32.Pt(sz.X-c, sz.Y)) // 右下切角
	p.LineTo(f32.Pt(c, sz.Y))      // 底边
	p.LineTo(f32.Pt(0, sz.Y-c))    // 左下切角
	p.LineTo(f32.Pt(0, c))         // 左边
	p.Close()                      // 闭合到 (c, 0)

	// 转换为描边 Stroke 模式（仅绘制边框，不填充）
	// 使用直线路径描边极其高效，不会引起内存暴增
	strokeSpec := clip.Stroke{
		Width: float32(width),
	}.Op()

	// 渲染边框颜色
	paint.FillShape(gtx.Ops, b.Color, strokeSpec)

	// 绘制刚才记录的子组件
	call.Add(gtx.Ops)

	return dims
}
