package ui

import (
	"fmt"
	"image"
	"math"
	"myapp/src/config"
	"os"
	"path/filepath"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ncruces/zenity"
)

type Live2DSetting struct {
	style           Style
	mainSettingPage *SettingPage

	btn_addModel  widget.Clickable
	explorerCount int
	selectedPath  string

	isLoading   bool
	btn_refresh widget.Clickable
	btn_reset   widget.Clickable
	btn_save    widget.Clickable

	isReset     bool
	sliderScale widget.Float
	sliderX     widget.Float
	sliderY     widget.Float
	btnScale    widget.Clickable
	btnX        widget.Clickable
	btnY        widget.Clickable
	tarS        float32
	tarX        float32
	tarY        float32
	lastS       float32
	lastX       float32
	lastY       float32
	stableTimer int

	list       widget.List
	enum       widget.Enum
	choices    []string
	choiceBtns []widget.Clickable
	deleteBtns []widget.Clickable
	deleteName string
	isAdd      bool

	language      string
	languageTable LangPack
}

func New_Live2DSetting(style Style) *Live2DSetting {
	self := Live2DSetting{
		style:         style,
		btn_addModel:  widget.Clickable{},
		explorerCount: 0,

		list:       widget.List{},
		enum:       widget.Enum{},
		choices:    []string{},
		choiceBtns: []widget.Clickable{},

		isReset: false,
		isAdd:   false,

		btn_reset:   widget.Clickable{},
		btn_save:    widget.Clickable{},
		btn_refresh: widget.Clickable{},
		isLoading:   false,
	}
	self.list.Axis = layout.Vertical
	self.language = "Chinese"
	self.languageTable = self.Language()

	// 预设信息
	self.enum.Value = (&config.Config_Live2D{}).Get_CurModelName()
	self.choices = (&config.Config_Live2D{}).Get_ModelList()
	self.choiceBtns = make([]widget.Clickable, len(self.choices))
	self.deleteBtns = make([]widget.Clickable, len(self.choices))

	// 创建滑块条
	item := (&config.Config_Live2D{}).Get_TarCoe(self.enum.Value)
	s, _ := item["scale"].(float64)
	x, _ := item["x"].(float64)
	y, _ := item["y"].(float64)

	// 这里的参数不是归一化的
	self.tarS = float32(s)
	self.tarX = float32(x)
	self.tarY = float32(y)

	self.sliderScale = widget.Float{Value: 0}
	self.sliderX = widget.Float{Value: 0}
	self.sliderY = widget.Float{Value: 0}
	self.btnScale = widget.Clickable{}

	self.btnX = widget.Clickable{}
	self.btnY = widget.Clickable{}

	self.initSlider(&self.sliderScale, 3, 0.2, self.tarS)
	self.initSlider(&self.sliderX, 500, -500, self.tarY)
	self.initSlider(&self.sliderY, 500, -500, self.tarX)

	// 上一次参数,记录是否发生变化,防止一直发送信息
	self.lastS = self.sliderScale.Value
	self.lastX = self.sliderX.Value
	self.lastY = self.sliderY.Value
	self.stableTimer = 0

	return &self
}

func (self *Live2DSetting) Default() {}
func (self *Live2DSetting) Title() string {
	return "Live2D"
}

// 更新
func (self *Live2DSetting) Update(gtx C, mainSettingPage *SettingPage) D {
	self.language = mainSettingPage.MainPage.trunk.Language
	self.mainSettingPage = mainSettingPage

	// 检测参数是否发生变化
	self.stableTimer = 0
	func() {
		isStable := true
		isStable = (self.lastS == self.sliderScale.Value) && isStable
		isStable = (self.lastX == self.sliderX.Value) && isStable
		isStable = (self.lastY == self.sliderY.Value) && isStable
		defer func() {
			self.lastS = self.sliderScale.Value
			self.lastX = self.sliderX.Value
			self.lastY = self.sliderY.Value
		}()
		if !isStable {
			self.ws_live2d_trans()
		}
	}()

	// 功能 ----------------------------------------------------
	if self.btn_addModel.Clicked(gtx) {
		self.selectFile()
	}

	if self.btn_reset.Clicked(gtx) {
		self.tarS = 1
		self.tarX = 0
		self.tarY = 0
		self.isReset = true
	}
	if self.btn_save.Clicked(gtx) {
		id := self.enum.Value
		scale := self.transCoe(self.sliderScale.Value, 3, 0.2, 2)
		x := self.transCoe(self.sliderX.Value, 500, -500, 0)
		y := self.transCoe(self.sliderY.Value, 500, -500, 0)
		(&config.Config_Live2D{}).Save_Live2DCoe(id, scale, int(x), int(y))
	}
	if self.btn_refresh.Clicked(gtx) {
		if !self.isLoading {
			self.reload()
		}
	}

	// 复位 ----------------------------------------------------
	if self.isReset {
		result := true
		result = self.resetSlider(gtx, &self.sliderScale, 3, 0.2, self.tarS) && result
		result = self.resetSlider(gtx, &self.sliderX, 500, -500, self.tarX) && result
		result = self.resetSlider(gtx, &self.sliderY, 500, -500, self.tarY) && result
		if result {
			self.isReset = false
		}
	}

	// 样式 -----------------------------------------------------

	// 加载按钮
	btnFunc := func(gtx C) D {
		gtx.Constraints.Max.Y = gtx.Dp(40)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		btn := material.Button(self.style.theme, &self.btn_addModel, "\uF8AA")
		btn.Color = self.style.theme.Fg
		btn.Background = self.style.darkmode.currentColor.IdleBg
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layout.Flexed(1, btn.Layout))
	}

	// 复位按钮与保存按钮
	resetBtnFunc := func(gtx C) D {
		gtx.Constraints.Max.Y = gtx.Dp(40)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X

		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return self.buildBtn(gtx, &self.btn_save, self.languageTable["save"][self.language], "\uE74E")
			}),
			layout.Rigid(Spacer(gtx.Dp(4), 0)),
			layout.Flexed(1, func(gtx C) D {
				return self.buildBtn(gtx, &self.btn_refresh, self.languageTable["refresh"][self.language], "\uE72C")
			}),
			layout.Rigid(Spacer(gtx.Dp(4), 0)),
			layout.Flexed(1, func(gtx C) D {
				return self.buildBtn(gtx, &self.btn_reset, self.languageTable["reset"][self.language], "\uE7A7")
			}),
		)
	}

	// 控件列表
	listElements := []layout.Widget{
		LittleTitle(gtx, self.style, self.languageTable["load_title"][self.language], unit.Sp(16)),
		btnFunc,
		resetBtnFunc,
		func(gtx C) D {
			return self.buildSlider(gtx, &self.sliderScale, &self.btnScale, self.languageTable["scale"][self.language], 3, 0.2, 2)
		},
		func(gtx C) D {
			return self.buildSlider(gtx, &self.sliderX, &self.btnX, self.languageTable["transX"][self.language], 500, -500, 0)
		},
		func(gtx C) D {
			return self.buildSlider(gtx, &self.sliderY, &self.btnY, self.languageTable["transY"][self.language], 500, -500, 0)
		},
		LittleTitle(gtx, self.style, self.languageTable["select_title"][self.language], unit.Sp(16)),
	}

	// 添加可选项
	for i := range self.choices {
		listElements = append(listElements, func(gtx C) D {
			return self.buildChoiceBtn(
				gtx,
				self.choices[i],
				&self.choiceBtns[i],
				&self.deleteBtns[i],
			)
		})
	}

	list := material.List(self.style.theme, &self.list)
	list.ScrollbarStyle.Indicator.Color.A = 64
	list.ScrollbarStyle.Indicator.MinorWidth = unit.Dp(4)
	out := list.Layout(gtx, len(listElements), func(gtx C, index int) D {
		return layout.Inset{Right: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx C) D { return listElements[index](gtx) })
	})

	if self.deleteName != "" {
		self.delete(self.deleteName)
		self.deleteName = ""
	}

	if self.isAdd {
		self.reload()
		self.isAdd = false
	}

	return out
}

// 预设选择控件
func (self *Live2DSetting) buildChoiceBtn(gtx C, choice string, checkBtn *widget.Clickable, deleteBtn *widget.Clickable) D {
	// gtx.Constraints.Max.Y = gtx.Dp(40)

	// 逻辑 -------------------------------------------------

	// 删除按钮
	if deleteBtn.Clicked(gtx) {
		// self.delete(choice)
		self.deleteName = choice
	}

	// 点选按钮
	if checkBtn.Clicked(gtx) && (self.enum.Value != choice) {
		self.enum.Value = choice
		trunk := self.mainSettingPage.MainPage.trunk

		// 设置并保存当前的全局参数
		trunk.Live2D_CurName = choice
		(&config.Config_Live2D{}).Save_CurName(choice)

		// 客户端更新
		self.ws_live2d_model(choice)

		// 应用参数
		coe := (&config.Config_Live2D{}).Get_TarCoe(choice)
		s, _ := coe["scale"].(float64)
		x, _ := coe["x"].(float64)
		y, _ := coe["y"].(float64)
		self.tarS = float32(s)
		self.tarX = float32(x)
		self.tarY = float32(y)
		self.isReset = true
	}

	// 样式 --------------------------------------------------

	// 点选按钮
	confirmIcon := "\uE739"
	if choice == self.enum.Value {
		confirmIcon = "\ue73d"
	}
	icon_label := material.Label(self.style.theme, unit.Sp(18), confirmIcon)
	text_label := material.Label(self.style.theme, unit.Sp(16), choice)

	// 删除按钮
	deleteFunc := func(gtx C) D {
		gtx.Constraints.Min = image.Pt(gtx.Dp(40), gtx.Sp(40))
		gtx.Constraints.Max = image.Pt(gtx.Dp(40), gtx.Sp(40))

		btn := material.Button(self.style.theme, deleteBtn, "\uE74D")
		btn.TextSize = unit.Sp(22)
		btn.Color = self.style.darkmode.currentColor.Fg
		btn.Background = self.style.darkmode.currentColor.IdleBg
		btn.Inset = layout.UniformInset(0)
		return btn.Layout(gtx)
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					btn := material.Button(self.style.theme, checkBtn, "")
					btn.Background = self.style.darkmode.currentColor.IdleBg
					return btn.Layout(gtx)
				}),
				layout.Stacked(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(Spacer(gtx.Dp(15), 0)),
							layout.Rigid(icon_label.Layout),
							layout.Rigid(Spacer(gtx.Dp(5), 0)),
							layout.Rigid(text_label.Layout),
							layout.Flexed(1, Flexer()),
						)
					})
				}),
			)
		}),
		layout.Rigid(Spacer(gtx.Dp(4), 0)),
		layout.Rigid(deleteFunc),
	)
}

// 选择文件逻辑
func (self *Live2DSetting) selectFile() string {
	self.explorerCount += 1
	if self.explorerCount <= 1 {
		go func() {
			filePath, err := zenity.SelectFile(
				zenity.Title(self.languageTable["select_explore_title"][self.language]),
				zenity.FileFilter{
					Name:     self.languageTable["select_explore_name"][self.language],
					Patterns: []string{"*.zip", "*.7z", "*.rar"},
				},
			)

			if err != nil {
				self.explorerCount = 0
				return
			}

			self.selectedPath = filePath
			self.explorerCount = 0
			result := (&config.Config_Live2D{}).CheckZipForModelFile(filePath)
			if !result {
				self.mainSettingPage.MainPage.notice.Add_Notice(
					2,
					self.languageTable["zip_warning_title"][self.language],
					self.languageTable["zip_warning_desc"][self.language],
				)
			}

			// 将zip解压到对应路径
			(&config.Config_Live2D{}).Unzip(filePath, filepath.Join(config.RootPath, "/live2d"))
			self.isAdd = true
		}()
	}
	return self.selectedPath
}

// 删除操作
func (self *Live2DSetting) delete(name string) {
	if name == self.enum.Value {
		if len(self.choices) > 0 {
			self.enum.Value = self.choices[0]
		}
	}

	path := filepath.Join(config.RootPath, "/live2d", fmt.Sprintf("/%s", name))
	err := os.RemoveAll(path)
	if err != nil {
		self.mainSettingPage.MainPage.notice.Add_Notice(2, "Delete Error", "")
	}
	self.reload()
}

// 重新加载
func (self *Live2DSetting) reload() {
	self.isLoading = true
	self.choices = (&config.Config_Live2D{}).Get_ModelList()
	self.choiceBtns = make([]widget.Clickable, len(self.choices))
	self.deleteBtns = make([]widget.Clickable, len(self.choices))
	self.isLoading = false
}

// 滑块
func (self *Live2DSetting) buildSlider(gtx C, widget_float *widget.Float, widget_clickable *widget.Clickable, title string, max float32, min float32, decimal int) D {
	// gtx.Constraints.Max.Y = gtx.Dp()
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	if widget_clickable.Clicked(gtx) {
		switch widget_float {
		case &self.sliderScale:
			self.tarS = 1
		case &self.sliderX:
			self.tarX = 0
		case &self.sliderY:
			self.tarY = 0
		default:
		}
		self.isReset = true
	}

	btnFunc := func(gtx C) D {
		gtx.Constraints.Max = image.Pt(gtx.Dp(20), gtx.Dp(20))
		gtx.Constraints.Min = image.Pt(gtx.Dp(20), gtx.Dp(20))
		btn := material.Button(self.style.theme, widget_clickable, "\uE7A7")
		btn.Background.A = 0
		btn.Inset = layout.Inset{}
		btn.TextSize = unit.Sp(14)
		return btn.Layout(gtx)
	}

	slider := SimpleSliderStyle{
		Axis:         layout.Horizontal,
		Float:        widget_float,
		ActiveColor:  self.style.theme.Palette.ContrastBg,
		FingerSize:   20, // 触摸热区 44dp
		TrackHeight:  12, // 轨道高度 8dp
		ThumbWidth:   4,  // 竖线宽度 6dp
		ThumbHeight:  20, // 竖线高度 20dp
		ThumbRadius:  2,
		CornerRadius: 4, // 圆角半径 4dp
		GapSize:      4,
	}

	value := fmt.Sprintf("%.*f", decimal, self.transCoe(widget_float.Value, max, min, decimal))

	sliderFunc := func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, slider.Layout),
			layout.Rigid(Spacer(gtx.Dp(4), 0)),
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Max = image.Pt(gtx.Dp(40), gtx.Dp(20))
				gtx.Constraints.Min = image.Pt(gtx.Dp(40), gtx.Dp(20))
				return layout.Center.Layout(gtx, func(gtx C) D {
					l := material.Label(self.style.theme, unit.Sp(16), value)
					l.Color = self.style.darkmode.currentColor.Fg
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, Flexer()),
						layout.Rigid(l.Layout),
						layout.Flexed(1, Flexer()),
					)
				})
			}),
			layout.Rigid(btnFunc),
		)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(LittleTitle(gtx, self.style, title, unit.Sp(16))),
		layout.Rigid(sliderFunc),
	)
}

// 滑块复位
func (self *Live2DSetting) resetSlider(
	gtx C,
	widget_float *widget.Float,
	max float32,
	min float32,
	tar float32,
) bool {
	tar_float := (tar - min) / (max - min)
	if widget_float.Value == tar_float {
		return true
	}
	gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
	widget_float.Value = widget_float.Value + (tar_float-widget_float.Value)*0.15
	if math.Abs(float64(widget_float.Value-tar_float)) <= 1e-4 {
		widget_float.Value = tar_float
	}
	return false
}

// 初始化滑块条参数
func (self *Live2DSetting) initSlider(
	widget_float *widget.Float,
	max float32,
	min float32,
	tar float32,
) {
	tar_float := (tar - min) / (max - min)
	widget_float.Value = tar_float
}

// 参数转换
func (self *Live2DSetting) transCoe(coe float32, max float32, min float32, decimal int) float32 {
	ori_v := coe*(max-min) + min

	f := float32(1)
	for range decimal {
		f *= 10
	}

	new_v := int(ori_v*f + 0.5)
	return float32(new_v) / f
}

// 重设按钮与保存按钮构建
func (self *Live2DSetting) buildBtn(gtx C, clickable *widget.Clickable, title string, icon string) D {
	gtx.Constraints.Max.Y = gtx.Dp(40)
	gtx.Constraints.Min.Y = gtx.Dp(40)

	btn := material.Button(self.style.theme, clickable, "")
	btn.Background = self.style.darkmode.currentColor.IdleBg

	icon_label := material.Label(self.style.theme, unit.Sp(20), icon)
	icon_label.Color = self.style.theme.Fg
	title_label := material.Label(self.style.theme, unit.Sp(20), title)
	title_label.Color = self.style.theme.Fg
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return btn.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, Flexer()),
					layout.Rigid(icon_label.Layout),
					layout.Rigid(Spacer(gtx.Dp(4), 0)),
					layout.Rigid(title_label.Layout),
					layout.Flexed(1, Flexer()),
				)
			})
		}),
	)
}

// 语言
func (self *Live2DSetting) Language() LangPack {
	return LangPack{
		"select_explore_title": LangZone{
			"Chinese":  "请选择 Live2D 模型压缩包",
			"English":  "Please select a Live2D model archive",
			"Japanese": "Live2Dモデルの圧縮ファイルを選択してください",
		},
		"select_explore_name": LangZone{
			"Chinese":  "Live2D 压缩包 (*.zip)",
			"English":  "Live2D Archive (*.zip)",
			"Japanese": "Live2D圧縮パック (*.zip)",
		},
		"load_title": LangZone{
			"Chinese":  "加载文件(.zip)",
			"English":  "Load File (.zip)",
			"Japanese": "ファイルを読み込む (.zip)",
		},
		"select_title": LangZone{
			"Chinese":  "选择文件",
			"English":  "Select File",
			"Japanese": "ファイルを選択",
		},
		"zip_warning_title": LangZone{
			"Chinese":  "不是合规的Zip",
			"English":  "Invalid ZIP File",
			"Japanese": "無効なZIPファイル",
		},
		"zip_warning_desc": LangZone{
			"Chinese":  "不是Live2D包或加密",
			"English":  "Not a valid Live2D package or it is encrypted",
			"Japanese": "Live2Dパッケージではないか、暗号化されています",
		},
		"scale": LangZone{
			"Chinese":  "缩放",
			"English":  "Scale",
			"Japanese": "拡大縮小", // 或者 "スケール"
		},
		"transX": LangZone{
			"Chinese":  "水平偏移",
			"English":  "Offset X", // 或者 "X Offset" / "Translation X"
			"Japanese": "水平オフセット",  // 或者 "位置(X)"
		},
		"transY": LangZone{
			"Chinese":  "竖直偏移",
			"English":  "Offset Y", // 或者 "Y Offset" / "Translation Y"
			"Japanese": "垂直オフセット",  // 或者 "位置(Y)"
		},
		"save": LangZone{
			"Chinese":  "保存",
			"English":  "Save",
			"Japanese": "保存",
		},
		"reset": LangZone{
			"Chinese":  "重设",
			"English":  "Reset",
			"Japanese": "リセット",
		},
		"refresh": LangZone{
			"Chinese":  "刷新",
			"English":  "Refresh",
			"Japanese": "刷新",
		},
	}
}

// ws通信 -------------------------------------------------------

// 更新大小,位移
func (self *Live2DSetting) ws_live2d_trans() {
	trunk := self.mainSettingPage.MainPage.trunk
	trunk.WsServer.SendCommand("server_set_live2d_trans", map[string]any{
		"scale":  self.transCoe(self.sliderScale.Value, 3.0, 0.2, 2),
		"transx": self.transCoe(self.sliderX.Value, 500, -500, 0),
		"transy": self.transCoe(self.sliderY.Value, 500, -500, 0),
	})
}

// 更新模型
func (self *Live2DSetting) ws_live2d_model(choice string) {
	trunk := self.mainSettingPage.MainPage.trunk
	folderPath := filepath.Join(config.RootPath, "/live2d", fmt.Sprintf("/%s", choice))
	model3name := (&config.Config_Live2D{}).Get_Model3Name(folderPath)
	trunk.Live2D_CurModel = []string{choice, model3name}
	fmt.Println("folder_path", trunk.Live2D_CurModel[0])
	trunk.WsServer.SendCommand("server_set_live2d_model", map[string]any{
		"FolderPath": trunk.Live2D_CurModel[0],
		"JsonName":   trunk.Live2D_CurModel[1],
	})
}
