package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
	W = layout.Widget
)
type CallBackFunc func()

type Style struct {
	theme    *material.Theme
	darkmode *DarkMode
}

// 语言
type LanguageTable struct {
	Chinese  string
	English  string
	Japanese string
}

func (t LanguageTable) GetByLang(lang string) string {
	switch lang {
	case "Chinese":
		return t.Chinese
	case "Japanese":
		return t.Japanese
	default:
		return t.English
	}
}

type LangZone map[string]string
type LangPack map[string]LangZone

// 圆角边框
var myGridBorder = widget.Border{
	Color:        color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	CornerRadius: unit.Dp(8),
	Width:        unit.Dp(2),
}
