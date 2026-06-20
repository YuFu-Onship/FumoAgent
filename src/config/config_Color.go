package config

import (
	"image/color"
)

type ColorPattle struct {
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
}

type ColorStyle struct {
	ID    string
	Name  LanguageStyle
	Dark  ColorPattle
	Light ColorPattle
}

func InitColorList() map[string]ColorStyle {
	colorList := map[string]ColorStyle{
		"cirno": {
			ID:    "cirno",
			Name:  LanguageStyle{Chinese: "琪露诺", English: "Cirno", Japanese: "チルノ"},
			Dark:  cirno_darkColor,
			Light: cirno_lightColor,
		},
		"koishi": {
			ID:    "koishi",
			Name:  LanguageStyle{Chinese: "古明地恋", English: "Koishi Komeiji", Japanese: "古明地こいし"},
			Dark:  koishi_darkColor,
			Light: koishi_lightColor,
		},
		"satori": {
			ID:    "satori",
			Name:  LanguageStyle{Chinese: "古明地觉", English: "Satori Komeiji", Japanese: "古明地さとり"},
			Dark:  satori_darkColor,
			Light: satori_lightColor,
		},
		"teto": {
			ID:    "teto",
			Name:  LanguageStyle{Chinese: "重音Teto", English: "Kasane Teto", Japanese: "重音テト"},
			Dark:  teto_darkColor,
			Light: teto_lightColor,
		},
		"windows11": {
			ID:    "windows11",
			Name:  LanguageStyle{Chinese: "Windows 11", English: "Windows 11", Japanese: "Windows 11"},
			Dark:  win11_darkColor,
			Light: win11_lightColor,
		},
	}
	return colorList
}

var cirno_lightColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255}, // 保持不变
	Fg:             color.NRGBA{R: 0, G: 60, B: 120, A: 255},    // 深蓝文字
	Fg_disable:     color.NRGBA{R: 120, G: 140, B: 160, A: 255},
	ContrastBg:     color.NRGBA{R: 65, G: 185, B: 255, A: 255},  // 冰翼亮蓝
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 纯白
	IdleBg:         color.NRGBA{R: 210, G: 235, B: 250, A: 255}, // 极浅冰蓝
	IdleFg_disable: color.NRGBA{R: 140, G: 170, B: 190, A: 255},
	IdleFg:         color.NRGBA{R: 90, G: 140, B: 180, A: 255}, // 天空蓝

	SelfMsgBg:  color.NRGBA{R: 30, G: 144, B: 255, A: 255}, // 闪烁蓝
	SelfMsgFg:  color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg: color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 白色衬衫感
	OtherMsgFg: color.NRGBA{R: 0, G: 100, B: 180, A: 255},
}
var cirno_darkColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},    // 保持不变
	Fg:             color.NRGBA{R: 200, G: 230, B: 255, A: 255}, // 浅冰蓝 (近白)
	Fg_disable:     color.NRGBA{R: 100, G: 120, B: 140, A: 255},
	ContrastBg:     color.NRGBA{R: 0, G: 162, B: 255, A: 255},   // 琪露诺蓝 (核心色)
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 纯白
	IdleBg:         color.NRGBA{R: 45, G: 55, B: 70, A: 255},    // 冷灰蓝
	IdleFg_disable: color.NRGBA{R: 70, G: 95, B: 120, A: 255},
	IdleFg:         color.NRGBA{R: 100, G: 130, B: 160, A: 255}, // 磨砂冰蓝
	SelfMsgBg:      color.NRGBA{R: 0, G: 120, B: 215, A: 255},   // 深蓝色 (裙子色)
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 60, G: 80, B: 100, A: 255}, // 冰砖色
	OtherMsgFg:     color.NRGBA{R: 180, G: 220, B: 240, A: 255},
}

// 古明地恋
var koishi_lightColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255}, // 保持不变
	Fg:             color.NRGBA{R: 40, G: 70, B: 60, A: 255},    // 深墨绿 (裙色加深，保证可读性)
	Fg_disable:     color.NRGBA{R: 150, G: 165, B: 160, A: 255},
	ContrastBg:     color.NRGBA{R: 70, G: 190, B: 150, A: 255}, // 恋恋发色 (柔和青绿)
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 250, G: 240, B: 190, A: 255}, // 浅黄色 (衬衫色)
	IdleFg_disable: color.NRGBA{R: 180, G: 180, B: 160, A: 255},
	IdleFg:         color.NRGBA{R: 140, G: 130, B: 80, A: 255},  // 橄榄灰
	SelfMsgBg:      color.NRGBA{R: 255, G: 130, B: 170, A: 255}, // 桃粉色 (第三只眼/爱心)
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 220, G: 240, B: 230, A: 255}, // 极淡青绿
	OtherMsgFg:     color.NRGBA{R: 30, G: 100, B: 80, A: 255},
}
var koishi_darkColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},    // 保持不变
	Fg:             color.NRGBA{R: 210, G: 240, B: 220, A: 255}, // 浅淡青绿文字
	Fg_disable:     color.NRGBA{R: 110, G: 130, B: 120, A: 255},
	ContrastBg:     color.NRGBA{R: 50, G: 160, B: 130, A: 255}, // 裙绿色 (核心色)
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 50, G: 60, B: 55, A: 255}, // 暗森林绿
	IdleFg_disable: color.NRGBA{R: 80, G: 100, B: 90, A: 255},
	IdleFg:         color.NRGBA{R: 120, G: 155, B: 140, A: 255}, // 灰绿色
	SelfMsgBg:      color.NRGBA{R: 180, G: 80, B: 200, A: 255},  // 深紫色 (瞳色/血管线)
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 65, G: 85, B: 80, A: 255}, // 灰绿底色
	OtherMsgFg:     color.NRGBA{R: 180, G: 220, B: 200, A: 255},
}

// 古明地觉
var satori_lightColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255},
	Fg:             color.NRGBA{R: 80, G: 40, B: 80, A: 255}, // 深紫 (裙色加深)
	Fg_disable:     color.NRGBA{R: 160, G: 140, B: 160, A: 255},
	ContrastBg:     color.NRGBA{R: 200, G: 140, B: 220, A: 255}, // 觉的发色 (柔和紫)
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 255, G: 230, B: 240, A: 255}, // 极浅粉 (衬衫色)
	IdleFg_disable: color.NRGBA{R: 180, G: 160, B: 170, A: 255},
	IdleFg:         color.NRGBA{R: 170, G: 110, B: 150, A: 255}, // 藕粉色
	SelfMsgBg:      color.NRGBA{R: 180, G: 20, B: 60, A: 255},   // 第三只眼 (深红)
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 240, G: 225, B: 245, A: 255}, // 淡紫色
	OtherMsgFg:     color.NRGBA{R: 120, G: 60, B: 140, A: 255},
}
var satori_darkColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},
	Fg:             color.NRGBA{R: 230, G: 210, B: 240, A: 255}, // 浅紫粉文字
	Fg_disable:     color.NRGBA{R: 120, G: 100, B: 120, A: 255},
	ContrastBg:     color.NRGBA{R: 150, G: 80, B: 180, A: 255}, // 核心紫色
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 60, G: 50, B: 65, A: 255}, // 深暗紫灰
	IdleFg_disable: color.NRGBA{R: 90, G: 80, B: 100, A: 255},
	IdleFg:         color.NRGBA{R: 140, G: 120, B: 160, A: 255}, // 磨砂紫
	SelfMsgBg:      color.NRGBA{R: 120, G: 40, B: 60, A: 255},   // 暗红 (第三只眼)
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 80, G: 70, B: 85, A: 255}, // 灰紫底
	OtherMsgFg:     color.NRGBA{R: 210, G: 180, B: 220, A: 255},
}

// 重音Teto
var teto_lightColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255},
	Fg:             color.NRGBA{R: 60, G: 60, B: 65, A: 255}, // 深军灰
	Fg_disable:     color.NRGBA{R: 150, G: 150, B: 155, A: 255},
	ContrastBg:     color.NRGBA{R: 255, G: 45, B: 88, A: 255}, // 钻头红发色
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 230, G: 230, B: 235, A: 255}, // 浅服灰色
	IdleFg_disable: color.NRGBA{R: 180, G: 170, B: 170, A: 255},
	IdleFg:         color.NRGBA{R: 160, G: 80, B: 90, A: 255}, // 灰红色
	SelfMsgBg:      color.NRGBA{R: 230, G: 0, B: 50, A: 255},  // 标志性红色
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 210, G: 210, B: 215, A: 255}, // 服装灰
	OtherMsgFg:     color.NRGBA{R: 50, G: 50, B: 55, A: 255},
}
var teto_darkColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},
	Fg:             color.NRGBA{R: 255, G: 210, B: 220, A: 255}, // 浅粉红文字
	Fg_disable:     color.NRGBA{R: 120, G: 110, B: 110, A: 255},
	ContrastBg:     color.NRGBA{R: 190, G: 30, B: 60, A: 255}, // 核心红
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 60, G: 55, B: 55, A: 255}, // 暗暖灰
	IdleFg_disable: color.NRGBA{R: 100, G: 85, B: 85, A: 255},
	IdleFg:         color.NRGBA{R: 160, G: 120, B: 125, A: 255}, // 灰粉色
	SelfMsgBg:      color.NRGBA{R: 150, G: 20, B: 40, A: 255},   // 深红色
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 70, G: 70, B: 75, A: 255}, // 军灰色底
	OtherMsgFg:     color.NRGBA{R: 240, G: 200, B: 210, A: 255},
}

// Windows 11
var win11_lightColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 243, G: 243, B: 243, A: 255}, // 标准浅色背景
	Fg:             color.NRGBA{R: 27, G: 27, B: 27, A: 255},    // 近乎黑的深灰文字
	Fg_disable:     color.NRGBA{R: 150, G: 150, B: 150, A: 255},
	ContrastBg:     color.NRGBA{R: 0, G: 103, B: 192, A: 255}, // Windows 典型蓝
	ContrastFg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IdleBg:         color.NRGBA{R: 232, G: 232, B: 232, A: 255}, // 辅助控件灰
	IdleFg_disable: color.NRGBA{R: 160, G: 160, B: 160, A: 255},
	IdleFg:         color.NRGBA{R: 96, G: 96, B: 96, A: 255}, // 中性灰文字
	SelfMsgBg:      color.NRGBA{R: 0, G: 95, B: 184, A: 255}, // 气泡蓝
	SelfMsgFg:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	OtherMsgBg:     color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 白色气泡
	OtherMsgFg:     color.NRGBA{R: 27, G: 27, B: 27, A: 255},
}
var win11_darkColor ColorPattle = ColorPattle{
	Bg:             color.NRGBA{R: 32, G: 32, B: 32, A: 255},    // Mica 黑色背景
	Fg:             color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 纯白文字
	Fg_disable:     color.NRGBA{R: 120, G: 120, B: 120, A: 255},
	ContrastBg:     color.NRGBA{R: 76, G: 194, B: 255, A: 255}, // 黑暗模式下的高亮蓝
	ContrastFg:     color.NRGBA{R: 0, G: 0, B: 0, A: 255},      // 亮蓝背景下使用深色字
	IdleBg:         color.NRGBA{R: 45, G: 45, B: 45, A: 255},   // 深灰控件背景
	IdleFg_disable: color.NRGBA{R: 100, G: 100, B: 100, A: 255},
	IdleFg:         color.NRGBA{R: 180, G: 180, B: 180, A: 255}, // 浅灰说明字
	SelfMsgBg:      color.NRGBA{R: 76, G: 194, B: 255, A: 255},  // 亮蓝消息
	SelfMsgFg:      color.NRGBA{R: 0, G: 34, B: 64, A: 255},     // 深蓝字以提升对比度
	OtherMsgBg:     color.NRGBA{R: 60, G: 60, B: 60, A: 255},    // 深色气泡
	OtherMsgFg:     color.NRGBA{R: 240, G: 240, B: 240, A: 255},
}
