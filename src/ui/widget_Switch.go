// SPDX-License-Identifier: Unlicense OR MIT

package ui

import (
	"image"
	"image/color"

	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

type MD2SwitchStyle struct {
	Description string
	Color       struct {
		// 激活状态下的颜色
		Accent color.NRGBA
		// 未激活状态下的滑轨填充色
		TrackOff color.NRGBA
		// 未激活状态下的滑块填充色
		ThumbOff color.NRGBA
	}
	Switch *widget.Bool
}

// NewMD2Switch 创建一个 MD2 结构、小圆角矩形特征的开关
func NewMD2Switch(swtch *widget.Bool, description string) MD2SwitchStyle {
	return MD2SwitchStyle{
		Switch:      swtch,
		Description: description,
	}
}

// Layout 绘制自定义矩形 MD2 风格开关
func (s MD2SwitchStyle) Layout(gtx layout.Context) layout.Dimensions {
	s.Switch.Update(gtx)

	// 保持与原版一致的外部尺寸 (36dp x 20dp)
	totalWidth := gtx.Dp(36)
	totalHeight := gtx.Dp(20)

	// 统一采用 4px 的圆角半径
	cornerRadius := gtx.Dp(4)

	// 1. 轨道尺寸与位置
	// MD2 轨道的经典高度比滑块小（原版高 16dp，这里也采用 16dp 以腾出视觉呼吸空间）
	trackWidth := totalWidth
	trackHeight := gtx.Dp(14)
	trackOffY := (totalHeight - trackHeight) / 2

	// 2. 滑块尺寸
	// 采用 20dp 边长的正方形（利用 6px 圆角后呈现为圆角矩形）
	thumbSize := totalHeight

	// 提取并处理禁用颜色
	accentColor := s.Color.Accent
	trackOffColor := TintColor(s.Color.TrackOff, 0x80)
	thumbOffColor := s.Color.ThumbOff

	if !gtx.Enabled() {
		accentColor = ToDisabled(accentColor)
		trackOffColor = ToDisabled(trackOffColor)
		thumbOffColor = ToDisabled(thumbOffColor)
	}

	// ---------------- 绘制滑轨 (Track) ----------------
	trackRect := image.Rectangle{Max: image.Point{X: trackWidth, Y: trackHeight}}
	var finalTrackColor color.NRGBA

	if s.Switch.Value {
		// 激活态：轨道通常是 Accent 颜色的浅色半透明版（MD2 经典规范）
		finalTrackColor = TintColor(accentColor, 0x80) // 约 50% 透明度
	} else {
		// 未激活态：使用传入的轨道原色
		finalTrackColor = trackOffColor
	}

	tStack := op.Offset(image.Point{Y: trackOffY}).Push(gtx.Ops)
	clTrack := clip.UniformRRect(trackRect, cornerRadius).Push(gtx.Ops)
	paint.ColorOp{Color: finalTrackColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	clTrack.Pop()
	tStack.Pop()

	// ---------------- 绘制滑块 (Thumb) ----------------
	var finalThumbColor color.NRGBA
	var xoff int

	if s.Switch.Value {
		finalThumbColor = accentColor
		xoff = trackWidth - thumbSize
	} else {
		finalThumbColor = thumbOffColor
		xoff = 0
	}

	thumbRect := image.Rectangle{Max: image.Point{X: thumbSize, Y: thumbSize}}

	pStack := op.Offset(image.Point{X: xoff, Y: 0}).Push(gtx.Ops)
	clThumb := clip.UniformRRect(thumbRect, cornerRadius).Push(gtx.Ops)
	paint.ColorOp{Color: finalThumbColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	clThumb.Pop()
	pStack.Pop()

	// ---------------- 热区与语义化 ----------------
	clickSize := gtx.Dp(40)
	clickOff := image.Point{
		X: (totalWidth - clickSize) / 2,
		Y: (totalHeight - clickSize) / 2,
	}

	stack := op.Offset(clickOff).Push(gtx.Ops)
	sz := image.Pt(clickSize, clickSize)
	clClick := clip.UniformRRect(image.Rectangle{Max: sz}, cornerRadius).Push(gtx.Ops)

	s.Switch.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if d := s.Description; d != "" {
			semantic.DescriptionOp(d).Add(gtx.Ops)
		}
		semantic.Switch.Add(gtx.Ops)
		return layout.Dimensions{Size: sz}
	})
	clClick.Pop()
	stack.Pop()

	return layout.Dimensions{Size: image.Point{X: totalWidth, Y: totalHeight}}
}
