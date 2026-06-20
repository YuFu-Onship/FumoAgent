package ui

import (
	"image"
	"image/color"
	"math"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// 动画 ---------------------------------------------------------------------------
// 缓动动画
func easing(gtx C, cur float64, tar float64, coefficient float64) float64 {
	cur += (tar - cur) * coefficient
	if math.Abs(cur-tar) < 1 {
		cur = tar
	} else {
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
	}
	return cur
}

func easing32(gtx C, cur float32, tar float32, coefficient float32) float32 {
	cur += (tar - cur) * coefficient
	if math.Abs(float64(cur-tar)) < 1 {
		cur = tar
	} else {
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
	}
	return cur
}

// 占位符 ---------------------------------------------------------------------------

// 创建占位符
func Spacer(width, height int) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := image.Point{X: width, Y: height}
		return layout.Dimensions{Size: size}
	}
}

func Flexer() layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Max.Y = 0
		return D{Size: gtx.Constraints.Max}
	}
}
func FlexerY() layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Max.X = 0
		return D{Size: gtx.Constraints.Max}
	}
}

// 一些简单预设 --------------------------------------------------------------------

// 小标题
func LittleTitle(gtx C, style Style, text string, size unit.Sp) W {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return material.Label(style.theme, size, text).Layout(gtx)
			}),
			layout.Flexed(1, func(gtx C) D { size := gtx.Constraints.Min; return D{Size: size} }),
		)
	}
}

// 绘制
// 绘制圆角矩形 --------------------------------------------------------------------
func DrawRRect(gtx C, pt image.Point, color color.NRGBA, radius int) {
	paint.FillShape(gtx.Ops, color, clip.RRect{
		Rect: image.Rectangle{Max: pt},
		SE:   radius,
		SW:   radius,
		NE:   radius,
		NW:   radius,
	}.Op(gtx.Ops))
}

// 绘制线条
func DrawLine(gtx C, x1 float32, y1 float32, x2 float32, y2 float32, width float32, c color.NRGBA) {
	var path clip.Path

	path.Begin(gtx.Ops)

	path.MoveTo(f32.Pt(x1, y1))
	path.LineTo(f32.Pt(x2, y2)) // 第一段线

	shape := clip.Stroke{
		Path:  path.End(),
		Width: float32(unit.Dp(width)),
	}.Op()

	// 5. 激活这个形状的裁剪区域（告诉 Gio 接下来的填充只在这个线条范围内生效）
	defer shape.Push(gtx.Ops).Pop()

	// 6. 用蓝色填充这个线条形状
	paint.ColorOp{Color: c}.Add(gtx.Ops)

	// 7. 执行绘制
	paint.PaintOp{}.Add(gtx.Ops)
}

// 颜色处理 -------------------------------------------------------------------------
// Argb 将 0xAARRGGBB 格式的 uint32 数字转换为 color.NRGBA
func Argb(c uint32) color.NRGBA {
	return color.NRGBA{
		A: uint8(c >> 24),
		R: uint8(c >> 16),
		G: uint8(c >> 8),
		B: uint8(c),
	}
}

// ToDisabled 将颜色转换为不透明度较低的禁用态表现
func ToDisabled(c color.NRGBA) color.NRGBA {
	return TintColor(c, 0x99)
}

// TintColor 改变颜色的不透明度
func TintColor(c color.NRGBA, alpha uint8) color.NRGBA {
	a := uint32(c.A) * uint32(alpha) / 255
	return color.NRGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: uint8(a),
	}
}

// 缓动动画对象 ---------------------------------------------------------------------
type Ani_Easing struct {
	curValue   float32
	tarValue   float32
	minValue   float32
	maxValue   float32
	Coeffcient float32
	isMin      bool
}

func New_Ani_Esaing(min float32, max float32, Coefficient float32, isMin bool) *Ani_Easing {
	var tar float32
	if isMin {
		tar = min
	} else {
		tar = max
	}
	return &Ani_Easing{
		curValue:   tar,
		tarValue:   tar,
		minValue:   min,
		maxValue:   max,
		Coeffcient: Coefficient,
		isMin:      isMin,
	}
}
func (self *Ani_Easing) Update(gtx C) float32 {
	if self.curValue == self.tarValue {
		return self.curValue
	}
	if math.Abs(float64(self.curValue-self.tarValue)) < 0.1 {
		self.curValue = self.tarValue
	}
	self.curValue += (self.tarValue - self.curValue) * self.Coeffcient
	gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
	return self.curValue
}

func (self *Ani_Easing) API_Toggle() {
	self.isMin = !self.isMin
	if self.isMin {
		self.tarValue = self.minValue
	} else {
		self.tarValue = self.maxValue
	}
}

func (self *Ani_Easing) API_SetDirection(isMin bool) {
	self.isMin = isMin
	if self.isMin {
		self.tarValue = self.minValue
	} else {
		self.tarValue = self.maxValue
	}
}

func (self *Ani_Easing) API_SetMaxValue(value float32) {
	self.maxValue = value
}

// 颜色处理: 微调输入的颜色
func ColorMultiper(c color.NRGBA, r float32, g float32, b float32, a float32) color.NRGBA {
	nr := c.R * uint8(r)
	ng := c.G * uint8(g)
	nb := c.B * uint8(b)
	na := c.A * uint8(a)
	return color.NRGBA{R: nr, G: ng, B: nb, A: na}
}
