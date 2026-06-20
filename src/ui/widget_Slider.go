package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type SimpleSliderStyle struct {
	Axis  layout.Axis
	Float *widget.Float

	// 颜色配置
	ActiveColor   color.NRGBA
	InactiveColor color.NRGBA
	ThumbColor    color.NRGBA

	// 尺寸配置
	FingerSize   unit.Dp
	TrackHeight  unit.Dp
	ThumbWidth   unit.Dp
	ThumbHeight  unit.Dp
	ThumbRadius  unit.Dp
	GapSize      unit.Dp // 调节竖线与矩形两侧的间距
	CornerRadius unit.Dp
}

func NewSimpleSlider(float *widget.Float, active, inactive, thumb color.NRGBA) SimpleSliderStyle {
	return SimpleSliderStyle{
		Axis:  layout.Horizontal,
		Float: float,

		ActiveColor:   active,
		InactiveColor: active,
		ThumbColor:    active,

		FingerSize:   44,
		TrackHeight:  16,
		ThumbWidth:   4,
		ThumbHeight:  44,
		GapSize:      4,
		CornerRadius: 8,
	}
}

func (s SimpleSliderStyle) Layout(gtx layout.Context) layout.Dimensions {
	radius := gtx.Dp(s.CornerRadius)
	line_radius := gtx.Dp(s.ThumbRadius)
	th := gtx.Dp(s.TrackHeight)
	tw := gtx.Dp(s.ThumbWidth)
	tH := gtx.Dp(s.ThumbHeight)
	gap := gtx.Dp(s.GapSize)

	thumbRadiusDp := unit.Dp(float32(s.ThumbWidth) / 2)
	tr := gtx.Dp(thumbRadiusDp)

	axis := s.Axis
	minLength := tr * 4

	touchSizePx := max(gtx.Dp(s.FingerSize), axis.Convert(gtx.Constraints.Min).Y)
	sizeMain := max(axis.Convert(gtx.Constraints.Min).X, minLength)
	sizeCross := max(tH, touchSizePx)
	size := axis.Convert(image.Pt(sizeMain, sizeCross))

	// 手势拦截
	o := axis.Convert(image.Pt(tr, 0))
	trans := op.Offset(o).Push(gtx.Ops)
	gtx.Constraints.Min = axis.Convert(image.Pt(sizeMain-2*tr, sizeCross))
	dims := s.Float.Layout(gtx, axis, thumbRadiusDp)
	gtx.Constraints.Min = gtx.Constraints.Min.Add(axis.Convert(image.Pt(0, sizeCross)))

	// 滑块中心绝对像素位置
	thumbPos := tr + int(s.Float.Value*float32(axis.Convert(dims.Size).X))
	trans.Pop()

	// 状态颜色
	activeColor := s.ActiveColor
	inactiveColor := s.ActiveColor
	inactiveColor.A = activeColor.A * 5 / 4
	thumbColor := s.ActiveColor
	if !gtx.Enabled() {
		activeColor = rgbToDisabled(activeColor)
		inactiveColor = rgbToDisabled(inactiveColor)
		thumbColor = rgbToDisabled(thumbColor)
	}

	// 坐标系转换函数
	rect := func(minx, miny, maxx, maxy int) image.Rectangle {
		r := image.Rect(minx, miny, maxx, maxy)
		if axis == layout.Vertical {
			r.Max.X, r.Min.X = sizeMain-r.Min.X, sizeMain-r.Max.X
		}
		r.Min = axis.Convert(r.Min)
		r.Max = axis.Convert(r.Max)
		return r
	}

	centerY := sizeCross / 2

	// ==================== 核心裁剪绘制逻辑 ====================

	// 1. 绘制一条完整长度的背景圆角矩形（未激活状态颜色）
	// 这条轨道永远是完美的 sizeMain 长度和 radius 圆角，不会因为滑块移动而变短
	fullTrackRect := rect(0, centerY-th/2, sizeMain, centerY+th/2)

	// 为了使两侧轨道都能应用间距，我们创建一个“挖空/裁剪”的总逻辑
	// 实际上，最优雅的方法是用全局背景，然后用 clip 限制住哪些地方能显示颜色

	// 【步骤 A】画未激活背景 (右侧部分)
	// 我们用 clip.Rect 限制只在滑块右侧 + 间距（inactiveStart）到终点之间可见
	inactiveStart := thumbPos + tw/2 + gap
	bgClipStack := op.Offset(image.Point{}).Push(gtx.Ops) // 占位
	if inactiveStart < sizeMain {
		// 转换为对应轴向的裁剪矩形
		clipMin := axis.Convert(image.Pt(inactiveStart, 0))
		clipMax := axis.Convert(image.Pt(sizeMain, sizeCross))
		// 修正垂直轴方向的裁剪边界
		if axis == layout.Vertical {
			clipMin.X, clipMax.X = 0, sizeMain
		}

		bgClip := clip.Rect{Min: clipMin, Max: clipMax}.Push(gtx.Ops)
		paint.FillShape(gtx.Ops, inactiveColor, clip.RRect{
			Rect: fullTrackRect, SE: radius, NE: radius, SW: radius, NW: radius,
		}.Op(gtx.Ops))
		bgClip.Pop()
	}
	bgClipStack.Pop()

	// 【步骤 B】画激活背景 (左侧部分)
	// 我们用 clip.Rect 限制只在 0 到滑块左侧 - 间距（activeEnd）之间可见
	activeEnd := thumbPos - tw/2 - gap
	fgClipStack := op.Offset(image.Point{}).Push(gtx.Ops)
	if activeEnd > 0 {
		clipMin := axis.Convert(image.Pt(0, 0))
		clipMax := axis.Convert(image.Pt(activeEnd, sizeCross))
		if axis == layout.Vertical {
			clipMin.X, clipMax.X = 0, sizeMain
		}

		fgClip := clip.Rect{Min: clipMin, Max: clipMax}.Push(gtx.Ops)
		paint.FillShape(gtx.Ops, activeColor, clip.RRect{
			Rect: fullTrackRect, SE: radius, NE: radius, SW: radius, NW: radius,
		}.Op(gtx.Ops))
		fgClip.Pop()
	}
	fgClipStack.Pop()

	// 3. 绘制中间：M3 竖条滑块 (不受裁剪影响，独立绘制)
	thumbRect := rect(thumbPos-tw/2, centerY-tH/2, thumbPos+tw/2, centerY+tH/2)
	paint.FillShape(gtx.Ops, thumbColor, clip.RRect{
		Rect: thumbRect,
		SE:   line_radius, NE: line_radius, SW: line_radius, NW: line_radius,
	}.Op(gtx.Ops))

	return layout.Dimensions{Size: size}
}

func rgbToDisabled(c color.NRGBA) color.NRGBA {
	col := color.NRGBAModel.Convert(c).(color.NRGBA)
	gray := uint8(float32(col.R)*0.299 + float32(col.G)*0.587 + float32(col.B)*0.114)
	return color.NRGBA{
		R: gray, G: gray, B: gray,
		A: uint8(float32(col.A) * 0.38),
	}
}
