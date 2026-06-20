package ui

import (
	"image/color"
	"math"
	"myapp/src/config"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget/material"
)

// 暗色与浅色模式的颜色配色
type ColorStruct struct {
	Bg             color.NRGBA
	Fg             color.NRGBA
	Fg_disable     color.NRGBA
	ContrastBg     color.NRGBA
	ContrastFg     color.NRGBA
	IdleBg         color.NRGBA
	IdleFg_disable color.NRGBA
	IdleFg         color.NRGBA
	SelfMsgBg      color.NRGBA
	SelfMsgFg      color.NRGBA
	OtherMsgBg     color.NRGBA
	OtherMsgFg     color.NRGBA

	Notice_Fg_1 color.NRGBA // 普通提示 - 前景
	Notice_Bg_1 color.NRGBA // 普通提示 - 背景
	Notice_Fg_2 color.NRGBA // 中度警告 - 前景
	Notice_Bg_2 color.NRGBA // 中度警告 - 背景
	Notice_Fg_3 color.NRGBA // 严重危机 - 前景
	Notice_Bg_3 color.NRGBA // 严重危机 - 背景
}

// ==================== 琪露诺 (Cirno) ====================
var cirno_lightColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255},
	Fg:             color.NRGBA{R: 0, G: 60, B: 120, A: 255},
	Fg_disable:     color.NRGBA{R: 120, G: 140, B: 160, A: 255},
	ContrastBg:     color.NRGBA{R: 65, G: 185, B: 255, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 210, G: 235, B: 250, A: 255},
	IdleFg_disable: color.NRGBA{R: 140, G: 170, B: 190, A: 255},
	IdleFg:         color.NRGBA{R: 90, G: 140, B: 180, A: 255},
	SelfMsgBg:      color.NRGBA{R: 30, G: 144, B: 255, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgFg:     color.NRGBA{R: 0, G: 100, B: 180, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 0, G: 130, B: 110, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 215, G: 245, B: 240, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 180, G: 90, B: 0, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 255, G: 235, B: 210, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 190, G: 0, B: 40, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 255, G: 220, B: 225, A: 255},
}

var cirno_darkColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},
	Fg:             color.NRGBA{R: 200, G: 230, B: 255, A: 255},
	Fg_disable:     color.NRGBA{R: 100, G: 120, B: 140, A: 255},
	ContrastBg:     color.NRGBA{R: 0, G: 162, B: 255, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 45, G: 55, B: 70, A: 255},
	IdleFg_disable: color.NRGBA{R: 70, G: 95, B: 120, A: 255},
	IdleFg:         color.NRGBA{R: 100, G: 130, B: 160, A: 255},
	SelfMsgBg:      color.NRGBA{R: 0, G: 120, B: 215, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 60, G: 80, B: 100, A: 255},
	OtherMsgFg:     color.NRGBA{R: 180, G: 220, B: 240, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 0, G: 240, B: 210, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 25, G: 50, B: 48, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 255, G: 170, B: 0, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 55, G: 45, B: 25, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 255, G: 70, B: 110, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 60, G: 30, B: 38, A: 255},
}

// ==================== 古明地恋 (Koishi) ====================
var koishi_lightColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255},
	Fg:             color.NRGBA{R: 40, G: 70, B: 60, A: 255},
	Fg_disable:     color.NRGBA{R: 150, G: 165, B: 160, A: 255},
	ContrastBg:     color.NRGBA{R: 70, G: 190, B: 150, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 250, G: 240, B: 190, A: 255},
	IdleFg_disable: color.NRGBA{R: 180, G: 180, B: 160, A: 255},
	IdleFg:         color.NRGBA{R: 140, G: 130, B: 80, A: 255},
	SelfMsgBg:      color.NRGBA{R: 255, G: 130, B: 170, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 220, G: 240, B: 230, A: 255},
	OtherMsgFg:     color.NRGBA{R: 30, G: 100, B: 80, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 150, G: 110, B: 0, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 255, G: 248, B: 215, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 190, G: 50, B: 100, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 255, G: 225, B: 235, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 160, G: 20, B: 30, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 255, G: 215, B: 215, A: 255},
}

var koishi_darkColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},
	Fg:             color.NRGBA{R: 210, G: 240, B: 220, A: 255},
	Fg_disable:     color.NRGBA{R: 110, G: 130, B: 120, A: 255},
	ContrastBg:     color.NRGBA{R: 50, G: 160, B: 130, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 50, G: 60, B: 55, A: 255},
	IdleFg_disable: color.NRGBA{R: 80, G: 100, B: 90, A: 255},
	IdleFg:         color.NRGBA{R: 120, G: 155, B: 140, A: 255},
	SelfMsgBg:      color.NRGBA{R: 180, G: 80, B: 200, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 65, G: 85, B: 80, A: 255},
	OtherMsgFg:     color.NRGBA{R: 180, G: 220, B: 200, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 245, G: 215, B: 70, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 50, G: 48, B: 30, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 255, G: 110, B: 185, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 55, G: 35, B: 45, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 255, G: 60, B: 60, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 65, G: 30, B: 30, A: 255},
}

// ==================== 古明地觉 (Satori) ====================
var satori_lightColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255},
	Fg:             color.NRGBA{R: 80, G: 40, B: 80, A: 255},
	Fg_disable:     color.NRGBA{R: 160, G: 140, B: 160, A: 255},
	ContrastBg:     color.NRGBA{R: 200, G: 140, B: 220, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 255, G: 230, B: 240, A: 255},
	IdleFg_disable: color.NRGBA{R: 180, G: 160, B: 170, A: 255},
	IdleFg:         color.NRGBA{R: 170, G: 110, B: 150, A: 255},
	SelfMsgBg:      color.NRGBA{R: 180, G: 20, B: 60, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 240, G: 225, B: 245, A: 255},
	OtherMsgFg:     color.NRGBA{R: 120, G: 60, B: 140, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 60, G: 100, B: 190, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 225, G: 235, B: 255, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 190, G: 80, B: 20, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 255, G: 232, B: 215, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 140, G: 0, B: 70, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 245, G: 210, B: 225, A: 255},
}

var satori_darkColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},
	Fg:             color.NRGBA{R: 230, G: 210, B: 240, A: 255},
	Fg_disable:     color.NRGBA{R: 120, G: 100, B: 120, A: 255},
	ContrastBg:     color.NRGBA{R: 150, G: 80, B: 180, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 60, G: 50, B: 65, A: 255},
	IdleFg_disable: color.NRGBA{R: 90, G: 80, B: 100, A: 255},
	IdleFg:         color.NRGBA{R: 140, G: 120, B: 160, A: 255},
	SelfMsgBg:      color.NRGBA{R: 120, G: 40, B: 60, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 80, G: 70, B: 85, A: 255},
	OtherMsgFg:     color.NRGBA{R: 210, G: 180, B: 220, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 120, G: 175, B: 255, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 35, G: 45, B: 60, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 255, G: 135, B: 50, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 55, G: 40, B: 30, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 245, G: 30, B: 110, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 60, G: 25, B: 38, A: 255},
}

// ==================== 重音 Teto (Kasane Teto) ====================
var teto_lightColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255},
	Fg:             color.NRGBA{R: 60, G: 60, B: 65, A: 255},
	Fg_disable:     color.NRGBA{R: 150, G: 150, B: 155, A: 255},
	ContrastBg:     color.NRGBA{R: 255, G: 45, B: 88, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 230, G: 230, B: 235, A: 255},
	IdleFg_disable: color.NRGBA{R: 180, G: 170, B: 170, A: 255},
	IdleFg:         color.NRGBA{R: 160, G: 80, B: 90, A: 255},
	SelfMsgBg:      color.NRGBA{R: 230, G: 0, B: 50, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 210, G: 210, B: 215, A: 255},
	OtherMsgFg:     color.NRGBA{R: 50, G: 50, B: 55, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 145, G: 95, B: 20, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 255, G: 242, B: 210, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 195, G: 70, B: 0, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 255, G: 230, B: 215, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 170, G: 0, B: 20, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 255, G: 215, B: 220, A: 255},
}

var teto_darkColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},
	Fg:             color.NRGBA{R: 255, G: 210, B: 220, A: 255},
	Fg_disable:     color.NRGBA{R: 120, G: 110, B: 110, A: 255},
	ContrastBg:     color.NRGBA{R: 190, G: 30, B: 60, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 60, G: 55, B: 55, A: 255},
	IdleFg_disable: color.NRGBA{R: 100, G: 85, B: 85, A: 255},
	IdleFg:         color.NRGBA{R: 160, G: 120, B: 125, A: 255},
	SelfMsgBg:      color.NRGBA{R: 150, G: 20, B: 40, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 70, G: 70, B: 75, A: 255},
	OtherMsgFg:     color.NRGBA{R: 240, G: 200, B: 210, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 245, G: 190, B: 50, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 50, G: 44, B: 30, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 255, G: 120, B: 40, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 55, G: 38, B: 30, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 255, G: 20, B: 70, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 60, G: 25, B: 32, A: 255},
}

// ==================== Windows 11 ====================
var win11_lightColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255},
	Fg:             color.NRGBA{R: 27, G: 27, B: 27, A: 255},
	Fg_disable:     color.NRGBA{R: 150, G: 150, B: 150, A: 255},
	ContrastBg:     color.NRGBA{R: 0, G: 103, B: 192, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 232, G: 232, B: 232, A: 255},
	IdleFg_disable: color.NRGBA{R: 160, G: 160, B: 160, A: 255},
	IdleFg:         color.NRGBA{R: 96, G: 96, B: 96, A: 255},
	SelfMsgBg:      color.NRGBA{R: 0, G: 95, B: 184, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgFg:     color.NRGBA{R: 27, G: 27, B: 27, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 0, G: 95, B: 184, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 225, G: 243, B: 255, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 159, G: 107, B: 0, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 255, G: 244, B: 206, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 196, G: 43, B: 28, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 253, G: 231, B: 233, A: 255},
}

var win11_darkColor ColorStruct = ColorStruct{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},
	Fg:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	Fg_disable:     color.NRGBA{R: 120, G: 120, B: 120, A: 255},
	ContrastBg:     color.NRGBA{R: 0, G: 120, B: 212, A: 255},
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 45, G: 45, B: 45, A: 255},
	IdleFg_disable: color.NRGBA{R: 100, G: 100, B: 100, A: 255},
	IdleFg:         color.NRGBA{R: 180, G: 180, B: 180, A: 255},
	SelfMsgBg:      color.NRGBA{R: 0, G: 120, B: 212, A: 255},
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 60, G: 60, B: 60, A: 255},
	OtherMsgFg:     color.NRGBA{R: 240, G: 240, B: 240, A: 255},
	Notice_Fg_1:    color.NRGBA{R: 107, G: 187, B: 251, A: 255},
	Notice_Bg_1:    color.NRGBA{R: 36, G: 49, B: 61, A: 255},
	Notice_Fg_2:    color.NRGBA{R: 253, G: 212, B: 111, A: 255},
	Notice_Bg_2:    color.NRGBA{R: 67, G: 58, B: 37, A: 255},
	Notice_Fg_3:    color.NRGBA{R: 241, G: 112, B: 122, A: 255},
	Notice_Bg_3:    color.NRGBA{R: 68, G: 39, B: 44, A: 255},
}

type ColorStructFloat struct {
	Bg, Fg, ContrastBg, ContrastFg, IdleBg, IdleFg, SelfMsgBg, SelfMsgFg, OtherMsgBg, OtherMsgFg [4]float64
}

// 暗色模式类 ---------------------------------------------------------
type DarkMode struct {
	IsDarkMode bool
	hwnd       uintptr
	hwndSet    bool

	currentColor   ColorStruct
	currentColorID string
	targetColor    ColorStruct
	currentFloat   ColorStructFloat

	titleName string
	isSwitch  bool
}

func NewDarkMode(title string) *DarkMode {
	curid := config.API_CUSTOM_COLOR_GetColorID()
	isdark := config.API_CUSTOM_COLOR_GetDarkmode()
	dm := DarkMode{
		IsDarkMode:     isdark,
		hwndSet:        false,
		currentColorID: curid,
		titleName:      title,
		isSwitch:       false,
	}

	curColor := dm.API_ReturnColor(curid, isdark)
	dm.currentColor = curColor
	return &dm
}

func (self *DarkMode) Update(gtx C) {
	self.colorTransition(gtx)
}

func (self *DarkMode) Init() {
	id := config.API_CUSTOM_COLOR_GetColorID()
	isdark := config.API_CUSTOM_COLOR_GetDarkmode()
	color := self.API_ReturnColor(id, isdark)
	self.currentColor = color
}

// 将配色改动保存在配置文件中
func (self *DarkMode) SaveColor() {
	id := self.currentColorID
	isdark := self.IsDarkMode
	config.API_CUSTOM_COLOR_SetColorID(id)
	config.API_CUSTOM_COLOR_SetDarkmode(isdark)
}

// 颜色过渡	----------------------------------------------------------------------------
func (self *DarkMode) colorTransition(gtx C) {
	if !self.isSwitch {
		return
	}

	// 更新所有颜色
	self.currentColor.Bg = self.lerpColor(gtx, &self.currentFloat.Bg, self.targetColor.Bg)
	self.currentColor.Fg = self.lerpColor(gtx, &self.currentFloat.Fg, self.targetColor.Fg)
}

// 颜色过渡
func (self *DarkMode) lerpColor(gtx C, cur *[4]float64, tar color.NRGBA) color.NRGBA {
	return color.NRGBA{
		R: self.lerpChannel(gtx, &cur[0], tar.R),
		G: self.lerpChannel(gtx, &cur[1], tar.G),
		B: self.lerpChannel(gtx, &cur[2], tar.B),
		A: self.lerpChannel(gtx, &cur[3], tar.A),
	}
}
func (self *DarkMode) lerpChannel(gtx C, cur *float64, tar uint8) uint8 {
	target := float64(tar)
	*cur += (target - *cur) * 0.2

	if math.Abs(*cur-target) < 1 {
		*cur = target
	} else {
		gtx.Execute(op.InvalidateCmd{})
	}
	return uint8(*cur)
}

// 切换颜色 -----------------------------------------------------------------------------
func (self *DarkMode) Toggle() {
	self.IsDarkMode = !self.IsDarkMode
	self.targetColor = self.API_ReturnColor(self.currentColorID, self.IsDarkMode)
	// self.isSwitch = true
	// if self.IsDarkMode {
	// 	self.targetColor = self.API_ReturnColor(self.currentColorID, true)
	// } else {
	// 	self.targetColor = self.API_ReturnColor(self.currentColorID, false)
	// }
}

// 设置当前颜色
func (self *DarkMode) setCurrentColor() {
	self.currentColor = self.API_ReturnColor(self.currentColorID, self.IsDarkMode)
}

// API 接口
func (self *DarkMode) API_SetColorMode(mode bool) {
	self.IsDarkMode = mode
	self.setCurrentColor()
	self.SetTitleBarDarkMode(self.IsDarkMode)
}

func (self *DarkMode) API_GetColorMode() bool {
	return self.IsDarkMode
}

//	func (self *DarkMode) API_GetColor() ColorStruct {
//		if self.IsDarkMode {
//			return darkColor
//		}
//		return lightColor
//	}
func (self *DarkMode) API_SetColor(name string) {
	self.currentColor = self.API_ReturnColor(name, self.IsDarkMode)
}
func (self *DarkMode) API_GetISDark() bool { return self.IsDarkMode }
func (self *DarkMode) API_ReturnColor(name string, isdark bool) ColorStruct {
	self.currentColorID = name
	switch name {
	default:
		if isdark {
			return cirno_darkColor
		} else {
			return cirno_lightColor
		}
	case "cirno":
		if isdark {
			return cirno_darkColor
		} else {
			return cirno_lightColor
		}
	case "satori":
		if isdark {
			return satori_darkColor
		} else {
			return satori_lightColor
		}
	case "koishi":
		if isdark {
			return koishi_darkColor
		} else {
			return koishi_lightColor
		}
	case "teto":
		if isdark {
			return teto_darkColor
		} else {
			return teto_lightColor
		}
	case "windows":
		if isdark {
			return win11_darkColor
		} else {
			return win11_lightColor
		}
	}
}

// 将某id的窗口标题栏设为暗色 ---------------------------------------------------------------
func (self *DarkMode) API_SetWindowTitle(title string) {
	self.titleName = title
	if self.SetWindowByTitle() {
		self.hwndSet = true
	}
}

// 应用颜色 ---------------------------------------------------------------------------------
func (self *DarkMode) Change(theme *material.Theme, color ColorStruct) {
	theme.Palette.Bg = color.Bg
	theme.Palette.Fg = color.Fg
	theme.Palette.ContrastBg = color.ContrastBg
	theme.Palette.ContrastFg = color.ContrastFg
}

func (self *DarkMode) Apply(theme *material.Theme) {
	self.SetTitleBarDarkMode(self.IsDarkMode)
	self.setCurrentColor()
	theme.Palette.Bg = self.currentColor.Bg
	theme.Palette.Fg = self.currentColor.Fg
	theme.Palette.ContrastBg = self.currentColor.ContrastBg
	theme.Palette.ContrastFg = self.currentColor.ContrastFg
	if !self.hwndSet {
		if self.SetWindowByTitle() {
			self.hwndSet = true
		}
	}
	self.SaveColor()
}

func (self *DarkMode) RetrySetWindow(w *app.Window, attempts int) {
	if self.hwndSet || attempts <= 0 {
		return
	}
	if self.SetWindowByTitle() {
		self.hwndSet = true
		w.Invalidate()
		return
	}
	// 如果没找到，200ms 后再试一次
	time.AfterFunc(200*time.Millisecond, func() {
		self.RetrySetWindow(w, attempts-1)
	})
}

// 获取窗口句柄 ----------------------------------------------------------------------------------
// SetWindowByTitle 通过窗口标题查找窗口并设置暗色标题栏, 当无法直接获取HWND时使用
func (self *DarkMode) SetWindowByTitle() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	self.hwnd = self.findWindow(self.titleName)
	if self.hwnd == 0 {
		return false
	}
	// 设置HWND并应用当前主题
	self.SetHWND(self.hwnd)
	return true
}

// findWindow 调用 Windows API FindWindow 查找窗口句柄
func (self *DarkMode) findWindow(title string) uintptr {
	user32 := syscall.NewLazyDLL("user32.dll")
	findWindowProc := user32.NewProc("FindWindowW")
	// 将标题转换为 UTF16 指针
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	// FindWindowW(NULL, title) - 通过窗口标题查找
	hwnd, _, _ := findWindowProc.Call(
		0, // NULL - 不限制窗口类名
		uintptr(unsafe.Pointer(titlePtr)),
	)
	return hwnd
}

// SetHWND 设置窗口句柄，用于Windows平台设置标题栏颜色
func (self *DarkMode) SetHWND(h uintptr) {
	self.hwnd = h
	self.SetTitleBarDarkMode(self.IsDarkMode)
}

// 使用 Windows DWM API 来设置标题栏颜色
func (self *DarkMode) SetTitleBarDarkMode(dark bool) {
	if runtime.GOOS != "windows" || self.hwnd == 0 {
		return
	}
	// DWMWA_USE_IMMERSIVE_DARK_MODE 属性值
	// Windows 10 20H1+ 和 Windows 11 使用值 20
	// 旧版本 Windows 10 使用值 19
	// 先尝试使用值 20，如果失败则尝试值 19
	const DWMWA_USE_IMMERSIVE_DARK_MODE_NEW = 20
	const DWMWA_USE_IMMERSIVE_DARK_MODE_OLD = 19

	var value int32
	if dark {
		value = 1
	} else {
		value = 0
	}

	// 先尝试新版本的属性值
	err := self.dwmSetWindowAttribute(self.hwnd, DWMWA_USE_IMMERSIVE_DARK_MODE_NEW, &value)
	if err != nil {
		// 如果失败，尝试旧版本的属性值
		self.dwmSetWindowAttribute(self.hwnd, DWMWA_USE_IMMERSIVE_DARK_MODE_OLD, &value)
	}
}

// dwmSetWindowAttribute 调用 Windows API DwmSetWindowAttribute
func (self *DarkMode) dwmSetWindowAttribute(hwnd uintptr, attr uint32, value *int32) error {
	dwmapi := syscall.NewLazyDLL("dwmapi.dll")
	dwmSetWindowAttributeProc := dwmapi.NewProc("DwmSetWindowAttribute")

	ret, _, _ := dwmSetWindowAttributeProc.Call(
		hwnd,
		uintptr(attr),
		uintptr(unsafe.Pointer(value)),
		unsafe.Sizeof(*value),
	)

	if ret != 0 {
		return syscall.Errno(ret)
	}
	return nil
}
