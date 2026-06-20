package ui

import (

	// 1. 引入 time 包

	"strconv"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type HomePage struct {
	style    Style
	MainPage *Page

	language  string
	langTable map[string]map[string]string

	btn_reteyClient widget.Clickable
}

func New_HomePage(style Style, mainPage *Page) *HomePage {
	self := &HomePage{
		style:    style,
		MainPage: mainPage,
	}
	self.langTable = self.Language()
	self.language = "Chinese"

	self.btn_reteyClient = widget.Clickable{}
	return self
}

func (self *HomePage) Update(gtx C) D {
	trunk := self.MainPage.trunk // 获取 trunk 引用
	self.language = trunk.Language

	if self.btn_reteyClient.Clicked(gtx) {
		trunk.EnableRender()
	}

	return layout.Inset{
		Top: unit.Dp(20), Bottom: unit.Dp(20),
		Left: unit.Dp(10), Right: unit.Dp(20),
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			// 标题栏：Home + 时间
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Alignment: layout.Baseline, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						l := material.H4(self.style.theme, trunk.CharacterTitle)
						l.Color = self.style.theme.Palette.Fg
						return l.Layout(gtx)
					}),
				)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),

			// 信息卡片区域
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D { return self.build_Lable(gtx, self.langTable["cur_llm"][self.language], trunk.AiModel) }),
					layout.Rigid(func(gtx C) D {
						return self.build_Lable(gtx, self.langTable["cur_live2d_model"][self.language], trunk.Live2D_CurName)
					}),
					layout.Rigid(func(gtx C) D {
						return self.build_Lable(gtx, self.langTable["cur_ws_port"][self.language], trunk.ServerPort)
					}),
					layout.Rigid(func(gtx C) D { return self.build_Lable_retey(gtx, trunk.IsConnect) }),
				)
			}),
			layout.Flexed(1, Flexer()),
		)
	})
}
func (self *HomePage) Default() {}

func (self *HomePage) build_Lable(gtx C, title string, desc string) D {
	gtx.Constraints.Max.Y = gtx.Dp(60)
	gtx.Constraints.Min.Y = gtx.Dp(60)
	l1 := material.Label(self.style.theme, unit.Sp(16), title)
	l2 := material.Label(self.style.theme, unit.Sp(14), desc)
	return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				DrawRRect(gtx, gtx.Constraints.Max, self.style.darkmode.currentColor.IdleBg, gtx.Dp(4))
				return D{}
			}),
			layout.Stacked(func(gtx C) D {
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(l1.Layout),
								layout.Rigid(l2.Layout),
							)
						}),
						layout.Flexed(1, Flexer()),
					)
				})
			}),
		)
	})
}

func (self *HomePage) build_Lable_retey(gtx C, isConnect bool) D {
	gtx.Constraints.Max.Y = gtx.Dp(60)
	gtx.Constraints.Min.Y = gtx.Dp(60)

	l1 := material.Label(self.style.theme, unit.Sp(16), self.langTable["cur_live2d_client"][self.language])
	l2 := material.Label(self.style.theme, unit.Sp(14), strconv.FormatBool(isConnect))
	btn_retryFunc := func(gtx C) D {
		if !isConnect {
			gtx.Constraints.Min.X = gtx.Dp(40)
			gtx.Constraints.Max.X = gtx.Dp(40)
			btn := material.Button(self.style.theme, &self.btn_reteyClient, "\ue777")
			btn.Inset = layout.Inset{}
			btn.TextSize = unit.Sp(20)
			btn.Background = self.style.darkmode.currentColor.IdleBg
			return layout.Inset{Left: unit.Dp(4)}.Layout(gtx, btn.Layout)
		}
		return D{}
	}

	return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						DrawRRect(gtx, gtx.Constraints.Min, self.style.darkmode.currentColor.IdleBg, gtx.Dp(4))
						return D{}
					}),
					layout.Stacked(func(gtx C) D {
						return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(l1.Layout),
										layout.Rigid(l2.Layout),
									)
								}),
								layout.Flexed(1, Flexer()),
							)
						})
					}),
				)
			}),
			layout.Rigid(btn_retryFunc),
		)
	})
}

func (self *HomePage) API() map[string]any {
	return map[string]any{}
}

func (self *HomePage) Language() map[string]map[string]string {
	return map[string]map[string]string{
		"cur_llm": {
			"Chinese":  "当前AI:",
			"English":  "Current AI:",
			"Japanese": "現在のAI:",
		},
		"cur_ws_port": {
			"Chinese":  "当前WebSocket端口:",
			"English":  "Current WebSocket Port:",
			"Japanese": "現在のWebSocketポート:",
		},
		"cur_live2d_client": {
			"Chinese":  "WebView连接状态:",
			"English":  "WebView Connection Status:",
			"Japanese": "WebView接続ステータス:",
		},
		"cur_live2d_model": {
			"Chinese":  "当前Live2d模型:",
			"English":  "Current Live2D Model:",
			"Japanese": "現在のLive2Dモデル:",
		},
	}
}
