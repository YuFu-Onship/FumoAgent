package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type PluginCtrl interface {
	Options() []map[string]any
	GUI_Language() map[string]map[string]string

	API_SetValue(name string, value any)
	API_SetValueString(map[string]string)
	API_GetValueString() map[string]string
	API_SetEnable(value bool)
	API_GetEnable() bool
	API_IsStable() bool
	Calibrate()
	API_GetID() string
}

type PluginGener struct {
	style         Style
	panel         ExpandPanel
	language      string
	langTable_gui map[string]map[string]string
	langTable_sys map[string]map[string]string

	Plugin PluginCtrl

	id         string
	titleText  string
	descText   string
	detailText string

	switch_enable widget.Bool
	list_option   widget.List
	list_Detail   widget.List
	widgetList    []map[string]any

	eleList []W
}

func New_PluginGener(style Style, plugin PluginCtrl) *PluginGener {
	self := PluginGener{
		style:  style,
		panel:  *New_ExpandPanel(style, 100),
		Plugin: plugin,

		switch_enable: widget.Bool{},
		eleList:       []W{},
	}

	// 启用按钮
	self.switch_enable.Value = self.Plugin.API_GetEnable()

	// 选项列表
	self.list_option = widget.List{}
	self.list_option.Axis = layout.Vertical
	self.list_Detail = widget.List{}
	self.list_Detail.Axis = layout.Vertical

	// 语言
	self.language = "Chinese"
	self.langTable_gui = self.Plugin.GUI_Language()
	self.langTable_sys = self.Language()

	// 获取插件内部的选项,并将其转化为控件
	options := self.Plugin.Options()
	self.widgetList = []map[string]any{}
	for i := range options {
		self.optionToWidget(options[i]["id"].(string), options[i]["value"])
	}

	return &self
}

func (self *PluginGener) Update(gtx C, parent *PluginPage) D {
	self.language = parent.MainPage.trunk.Language
	// 功能 ----------------------------------------------------------------------
	self.Plugin.API_SetEnable(self.switch_enable.Value)

	// 布尔检测
	if !self.Plugin.API_IsStable() {
		// self.Plugin.Calibrate()
	}

	// **重中之重** --------------------------------------------------------------------------------------------
	self.eleList = []W{}
	self.eleList = append(self.eleList, func(gtx C) D {
		return self.build_Switch(gtx, &self.switch_enable, self.langTable_sys["enable"][self.language])
	})
	for i := range self.widgetList {
		self.eleList = append(self.eleList, func(gtx C) D {
			switch v := self.widgetList[i]["value"].(type) {
			case *widget.Bool:
				self.Plugin.API_SetValue(self.widgetList[i]["id"].(string), v.Value)
				return self.build_Switch(gtx, v, self.langTable_gui[self.widgetList[i]["id"].(string)][self.language])
			default:
				return D{}
			}
		})
	}

	// 样式 ----------------------------------------------------------------------
	self.titleText = self.langTable_gui["Title"][self.language]
	self.descText = self.langTable_gui["Desc"][self.language]
	self.detailText = self.langTable_gui["Detail"][self.language]
	// 顶栏, 标题与简短描述
	headFunc := func(gtx C) D {
		l1 := material.Label(self.style.theme, unit.Sp(18), self.titleText)
		l2 := material.Label(self.style.theme, unit.Sp(14), self.descText)
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(12), Right: unit.Dp(50)}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(l1.Layout),
						layout.Rigid(l2.Layout),
					)
				})
			}),
			layout.Flexed(1, Flexer()),
		)
	}
	d, h := self.panel.Update(gtx, headFunc, self.PanelLayout)
	self.panel.API_SetPanelHeight(min(h, gtx.Dp(200)))
	return d
}

func (self *PluginGener) PanelLayout(gtx C) D {
	list := material.List(self.style.theme, &self.list_option)
	leftFunc := func(gtx C) D {
		return BorderBox(
			self.style,
			self.style.darkmode.currentColor.IdleBg,
			self.style.darkmode.currentColor.Fg,
			self.langTable_sys["option_title"][self.language],
		).Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Dp(150)
			gtx.Constraints.Max.X = gtx.Dp(150)
			return list.Layout(gtx, len(self.eleList), func(gtx C, index int) D {
				return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx C) D {
					return self.eleList[index](gtx)
				})
			})
		})
	}

	rightFunc := func(gtx C) D {
		l := BorderBox(
			self.style,
			self.style.darkmode.currentColor.IdleBg,
			self.style.darkmode.currentColor.Fg,
			self.langTable_sys["detail_title"][self.language],
		).Layout(gtx, func(gtx C) D {
			gtx.Constraints.Max.Y = gtx.Dp(140)
			label := material.Label(self.style.theme, unit.Sp(14), self.detailText)
			return material.List(self.style.theme, &self.list_Detail).Layout(gtx, 1, func(gtx C, index int) D {
				return layout.Flex{}.Layout(gtx, layout.Flexed(1, label.Layout))
			})
		})
		return l
	}

	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(leftFunc),
			layout.Rigid(Spacer(gtx.Sp(16), 0)),
			layout.Rigid(rightFunc),
		)
	})
}

func (self *PluginGener) optionToWidget(name string, o any) {
	switch v := o.(type) {
	case bool:
		b := &widget.Bool{Value: v}
		self.widgetList = append(self.widgetList, map[string]any{"id": name, "value": b})
	}
}

func (self *PluginGener) build_Switch(gtx C, widget_bool *widget.Bool, title string) D {
	sw := material.Switch(self.style.theme, widget_bool, "")

	l := material.Label(self.style.theme, unit.Sp(14), title)
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(l.Layout),
			layout.Flexed(1, Flexer()),
			layout.Rigid(sw.Layout),
		)
	})
}

func (self *PluginGener) Language() map[string]map[string]string {
	return map[string]map[string]string{
		"enable": {
			"Chinese":  "启用",
			"English":  "Enable",
			"Japanese": "有効",
		},
		"option_title": {
			"Chinese":  "选项",
			"English":  "Options",
			"Japanese": "オプション",
		},
		"detail_title": {
			"Chinese":  "详情",
			"English":  "Detail",
			"Japanese": "詳細",
		},
	}
}

// 参数是否稳定
func (self *PluginGener) API_IsStable() bool {
	return self.Plugin.API_IsStable()
}

// 对齐
func (self *PluginGener) Calibrate() {
	self.Plugin.Calibrate()
}

// 返回 ID
func (self *PluginGener) API_GetID() string {
	return self.Plugin.API_GetID()
}

// 返回参数
func (self *PluginGener) API_GetValueString() map[string]string {
	return self.Plugin.API_GetValueString()
}
