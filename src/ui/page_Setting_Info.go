package ui

import (
	"os/exec"
	"runtime"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type InfoSetting struct {
	style           Style
	mainSettingPage *SettingPage

	// 为 GitHub 链接添加点击状态
	githubClick widget.Clickable
}

func New_InfoSetting(style Style) *InfoSetting {
	i_s := InfoSetting{
		style: style,
	}
	return &i_s
}

func (self *InfoSetting) Default() {

}

func (self *InfoSetting) Title() string {
	return self.mainSettingPage.MainPage.trunk.LanguageTable.SettingPage_Info_Title
}

// 辅助函数：调用系统默认浏览器打开网址
func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}

func (self *InfoSetting) Update(gtx C, parent *SettingPage) D {
	self.mainSettingPage = parent

	// 获取当前全局语言
	lang := self.mainSettingPage.MainPage.trunk.Language
	langData := self.InitLangauge()

	// 处理 GitHub 链接的点击事件
	if self.githubClick.Clicked(gtx) {
		openURL(langData["githubDesc"].GetByLang(lang))
	}

	// 辅助函数：渲染一行信息（支持普通文本和可点击链接）
	renderInfoRow := func(gtx C, titleKey, descKey string, isLink bool, click *widget.Clickable) D {
		return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				// 标题渲染
				layout.Rigid(func(gtx C) D {
					title := material.Label(self.style.theme, unit.Sp(13), langData[titleKey].GetByLang(lang))
					title.Color = self.style.theme.Palette.ContrastBg
					return title.Layout(gtx)
				}),
				// 紧凑间距
				layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
				// 内容渲染
				layout.Rigid(func(gtx C) D {
					textStr := langData[descKey].GetByLang(lang)

					// 如果是链接，特殊渲染并捕获点击
					if isLink && click != nil {
						return click.Layout(gtx, func(gtx C) D {
							// 鼠标悬停时变成手型光标
							pointer.CursorPointer.Add(gtx.Ops)

							label := material.Label(self.style.theme, unit.Sp(16), textStr)
							// 使用主题的主题色（通常是蓝色或亮色）来表现它是超链接
							label.Color = self.style.theme.Palette.Fg

							// 极简风格通常不加粗重的下划线，如果您想加，可以使用纯文本并在多语言里加下划线，
							// 或者保持颜色区分，这是最符合现代极简 UI 的做法。
							return label.Layout(gtx)
						})
					}

					// 普通文本渲染
					desc := material.Label(self.style.theme, unit.Sp(16), textStr)
					return desc.Layout(gtx)
				}),
			)
		})
	}

	// 整体内边距布局
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top:  unit.Dp(16),
				Left: unit.Dp(12),
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{
					Axis: layout.Vertical,
					// Alignment: layout.Start,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return renderInfoRow(gtx, "versionTitle", "versionDesc", false, nil)
					}),
					layout.Rigid(func(gtx C) D {
						return renderInfoRow(gtx, "developerTitle", "developerDesc", false, nil)
					}),
					layout.Rigid(func(gtx C) D {
						return renderInfoRow(gtx, "githubTitle", "githubDesc", true, &self.githubClick)
					}),
					layout.Rigid(func(gtx C) D {
						return renderInfoRow(gtx, "creditsVoiceTitle", "creditsVoiceDesc", false, nil)
					}),
					layout.Rigid(func(gtx C) D {
						return renderInfoRow(gtx, "creditsRenderTitle", "creditsRenderDesc", false, nil)
					}),
					layout.Rigid(func(gtx C) D {
						return renderInfoRow(gtx, "creditsUiTitle", "creditsUiDesc", false, nil)
					}),
				)
			})
		}),
		layout.Flexed(1, Flexer()),
	)
}

func (self *InfoSetting) InitLangauge() map[string]LanguageTable {
	return map[string]LanguageTable{
		"versionTitle": {
			Chinese:  "软件版本",
			English:  "Software Version",
			Japanese: "ソフトウェアバージョン",
		},
		"versionDesc": {
			Chinese:  "v2026-6",
			English:  "v2026-6",
			Japanese: "v2026-6",
		},
		"developerTitle": {
			Chinese:  "开发者",
			English:  "Developer",
			Japanese: "開発者",
		},
		"developerDesc": {
			Chinese:  "YUFU",
			English:  "YUFU",
			Japanese: "YUFU",
		},
		"githubTitle": {
			Chinese:  "GitHub 主页",
			English:  "GitHub Homepage",
			Japanese: "GitHub ホームページ",
		},
		"githubDesc": {
			Chinese:  "https://github.com/YuFu-Onship",
			English:  "https://github.com/YuFu-Onship",
			Japanese: "https://github.com/YuFu-Onship",
		},
		"creditsVoiceTitle": {
			Chinese:  "鸣谢：语音与声码器",
			English:  "Credits: TTS & Voice Library",
			Japanese: "谢辞：音声合成＆音源ライブラリ",
		},
		"creditsVoiceDesc": {
			Chinese:  "AquestaTalk, Open JTalk 开发者，以及 VOICEVOX 开发者与全体声优角色库提供方",
			English:  "AquestaTalk, Open JTalk Developers, VOICEVOX Teams & All Voice Library Providers",
			Japanese: "AquestaTalk, Open JTalk 開発者、VOICEVOX 開発者及びすべてのキャラクター・音源制作者の方々",
		},
		"creditsRenderTitle": {
			Chinese:  "鸣谢：Live2D 渲染",
			English:  "Credits: Live2D Rendering",
			Japanese: "谢辞：Live2D レンダリング",
		},
		"creditsRenderDesc": {
			Chinese:  "pixi-live2d-display 框架开源社区",
			English:  "pixi-live2d-display Framework Community",
			Japanese: "pixi-live2d-display フレームワークコミュニティ",
		},
		"creditsUiTitle": {
			Chinese:  "鸣谢：界面驱动",
			English:  "Credits: UI Framework",
			Japanese: "谢辞：UI 框架",
		},
		"creditsUiDesc": {
			Chinese:  "Gio UI (gioui.org) 开源项目",
			English:  "Gio UI (gioui.org) Open Source Project",
			Japanese: "Gio UI (gioui.org) オープンソースプロジェクト",
		},
	}
}
