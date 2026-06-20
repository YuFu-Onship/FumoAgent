package ui

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// type CallBackFunc func()
type BtnWithDesc struct {
	iconText     string
	descText     string
	descShowText string
	iconBtn      widget.Clickable
	descBtn      widget.Clickable
	iconBtnWidth float64

	descBtbWidth_01  float64
	descBtnWidth_02  float64
	descBtnWidth_tar float64
	descBtnWidth_cur int

	squareHeight_tar float64
	squareHeight_cur float64

	descRadius float32
	descHeight float64
	iconRadius float32

	desc_ease Ani_Easing

	style     Style
	isActive  bool
	isHovered bool

	gtx C
}

func New_BtnWidthDesc(gtx C, style Style, Icon string, Desc string) *BtnWithDesc {
	var iconRadius float32 = float32(unit.Dp(4))
	var iconWidth float64 = float64(unit.Dp(40))
	return &BtnWithDesc{
		gtx:              gtx,
		iconText:         Icon,
		descText:         Desc,
		style:            style, // 包含theme与darkmode,其中darkmode包含一些全局配色
		descBtbWidth_01:  0,
		descBtnWidth_02:  float64(unit.Dp(150)),
		descBtnWidth_tar: 0,
		descBtnWidth_cur: 0,
		squareHeight_tar: 0,
		squareHeight_cur: 0,

		iconBtnWidth: iconWidth,
		descHeight:   iconWidth,
		iconRadius:   iconRadius,
		descRadius:   iconRadius,

		desc_ease: *New_Ani_Esaing(0, 150, 0.3, false),

		isHovered: false,
		isActive:  false,
	}
}

// 调用通用函数 --------------------------------------------------------------------------------------
func (self *BtnWithDesc) Update(gtx layout.Context, isActive bool, callback CallBackFunc) layout.Dimensions {
	self.isActive = isActive

	// 1. 交互判定：先确定目标值 (Target)
	if self.iconBtn.Clicked(gtx) {
		callback()
	}

	if self.iconBtn.Hovered() {
		self.isHovered = true
		self.desc_ease.API_SetDirection(false)
	} else {
		self.isHovered = false
		self.desc_ease.API_SetDirection(true)
	}

	// 2. 动画计算：获取当前帧宽度 (Current)
	// 确保这一步在 Hovered() 判定之后执行，以保证动画立刻响应
	currentWidth := int(self.desc_ease.Update(gtx))
	self.descBtnWidth_cur = gtx.Dp(unit.Dp(currentWidth))
	// fmt.Println(self.descBtnWidth_cur)

	// 3. 样式准备
	var FgColor color.NRGBA
	var BgColor color.NRGBA
	if self.isHovered || self.isActive {
		BgColor = self.style.theme.ContrastBg
		FgColor = self.style.darkmode.currentColor.SelfMsgFg
	} else {
		BgColor = self.style.darkmode.currentColor.IdleBg
		FgColor = self.style.darkmode.currentColor.OtherMsgFg
	}

	// 4. 侧边状态条动画
	if self.isActive {
		self.squareHeight_tar = self.iconBtnWidth * 0.3
	} else {
		self.squareHeight_tar = 0
		self.squareHeight_cur = 0
	}
	self.squareHeight_cur = easing(gtx, self.squareHeight_cur, self.squareHeight_tar, 0.3)

	// 5. 布局构建
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// 左侧装饰条
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Dp(14)
			if self.squareHeight_cur >= 1 {
				paint.FillShape(gtx.Ops, self.style.darkmode.currentColor.ContrastBg, clip.RRect{
					Rect: image.Rect(
						gtx.Dp(5),
						gtx.Dp(unit.Dp(-self.squareHeight_cur)),
						gtx.Dp(10),
						gtx.Dp(unit.Dp(self.squareHeight_cur)),
					),
					SE: gtx.Dp(2), SW: gtx.Dp(2), NE: gtx.Dp(2), NW: gtx.Dp(2),
				}.Op(gtx.Ops))
			}
			return D{Size: gtx.Constraints.Min}
		}),

		// 图标按钮
		layout.Rigid(func(gtx C) D {

			measuringGtx := gtx
			measuringGtx.Constraints.Min = image.Point{}
			measuringGtx.Constraints.Max = image.Point{X: gtx.Dp(1000), Y: gtx.Dp(40)}

			gtx.Constraints.Min = image.Pt(gtx.Dp(40), gtx.Dp(40))
			gtx.Constraints.Max = image.Pt(gtx.Dp(40), gtx.Dp(40))
			btn := material.Button(self.style.theme, &self.iconBtn, self.iconText)
			btn.CornerRadius = unit.Dp(self.iconRadius)
			btn.Background = BgColor
			btn.Color = FgColor
			btn.TextSize = unit.Sp(22)
			btn.Font.Weight = font.Weight(gtx.Dp(900))
			btn.Inset = layout.Inset{}

			btnfunc := btn.Layout(gtx)

			return btnfunc
		}),

		layout.Rigid(Spacer(gtx.Dp(5), 0)),
		layout.Rigid(func(gtx C) D {
			if self.descBtnWidth_cur < 1 {
				return D{}
			}

			btn := material.Button(self.style.theme, &self.descBtn, self.descText)
			btn.CornerRadius = unit.Dp(self.iconRadius)
			btn.Background = BgColor
			btn.Color = FgColor
			btn.TextSize = unit.Sp(20)
			btn.Font.Weight = font.Weight(gtx.Dp(900))
			btn.Inset.Bottom = 0
			btn.Inset.Left = unit.Dp(20)
			btn.Inset.Top = 0
			btn.Inset.Right = unit.Dp(20)

			cl := clip.Rect{Max: image.Point{X: int(self.descBtnWidth_cur), Y: gtx.Dp(40)}}.Push(gtx.Ops)
			dims := func(gtx C) D {
				gtx.Constraints.Min = image.Pt(gtx.Dp(150), gtx.Dp(40))
				gtx.Constraints.Max = image.Pt(gtx.Dp(150), gtx.Dp(40))
				return btn.Layout(gtx)
			}(gtx)

			cl.Pop()
			return D{Size: image.Point{X: dims.Size.X, Y: gtx.Dp(40)}, Baseline: dims.Size.X}
		}),
	)
}

// api ---------------------------------------------------------------------
func (self *BtnWithDesc) API_SetIcon(iconText string) {
	self.iconText = iconText
}
func (self *BtnWithDesc) API_SetDesc(descText string) {
	self.descText = descText
}
func (self *BtnWithDesc) API_SetColor() {

}
func (self *BtnWithDesc) API_SetState(isActive bool) {
	self.isActive = isActive
}

// 内部工具 ------------------------------------------------------------------
// func (self *BtnWithDesc) easing(cur float64, tar float64, coefficient float64) float64 {
// 	cur += (tar - cur) * coefficient
// 	if math.Abs(cur-tar) < 1 {
// 		cur = tar
// 	} else {
// 		self.gtx.Execute(op.InvalidateCmd{At: self.gtx.Now.Add(time.Second / 60)})
// 	}
// 	return cur
// }

func (self *BtnWithDesc) buildButton(
	gtx C,
	iconText string,
	btn material.ButtonStyle,
	inset layout.Inset,
) D {
	label := material.Label(self.style.theme, unit.Sp(20), iconText)
	return func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Dp(40)
		gtx.Constraints.Min.Y = gtx.Dp(40)
		return inset.Layout(gtx, func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return btn.Layout(gtx)
				}),
				layout.Stacked(func(gtx C) D {
					return layout.Center.Layout(gtx, label.Layout)
				}),
			)

		})
	}(gtx)
}
