package ui

import (
	_ "embed"
	"image"
	_ "image/png"
	"log"
	"myapp/src/config"
	"myapp/src/model"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// 变量设置 --------------------------------------------

// 字体相关

func loadFonts() []font.FontFace {
	sansFontData, err := os.ReadFile(filepath.Join(config.RootPath, "assets/font/MiSansVF.ttf"))
	if err != nil {
		config.Error_Save(err, "读取MiSans字体失败")
	}

	emojiFontData, err := os.ReadFile(filepath.Join(config.RootPath, "assets/font/NotoEmoji-Bold.ttf"))
	if err != nil {
		config.Error_Save(err, "读取Emoji字体失败")
	}

	iconFontData, err := os.ReadFile(filepath.Join(config.RootPath, "assets/font/SegoeFluentIcons.ttf"))
	if err != nil {
		config.Error_Save(err, "读取图标字体失败")
	}

	// 解析字体文件
	sans_face, err := opentype.Parse(sansFontData)
	emoji_face, err := opentype.Parse(emojiFontData)
	icon_face, err := opentype.Parse(iconFontData)

	// 注意：建议对上面三个 Parse 的 err 分别进行校验，
	// 像你之前那样只在最后检查 err，如果前两个报错了会被第三个覆盖掉。
	config.Error_Save(err, "字体解析")

	return []font.FontFace{
		{
			Face: sans_face,
			Font: font.Font{Typeface: "sans"},
		},
		{
			Face: emoji_face,
			Font: font.Font{Typeface: "noto"},
		},
		{
			Face: icon_face,
			Font: font.Font{Typeface: "icons"},
		},
	}
}

// 定义:侧边栏
var (
	btnHome      *BtnWithDesc
	btnChat      *BtnWithDesc
	btnDarkmode  *BtnWithDesc
	btnSetting   *BtnWithDesc
	sideBarRigid []layout.FlexChild
)

// 定义:页面接口
type Component interface {
	Update(gtx C) D
	Default()
	API() map[string]any
}

// ui界面 ----------------------------------------------------
// 定义:总页面
type Page struct {
	Window   *app.Window
	darkmode *DarkMode
	title    string
	style    Style
	gtx      C

	// 占位页高度, 用于实现页面自下而上的加载
	pageOffsetHeight float64

	// 主干
	trunk *model.Trunk

	// 子页面
	currentPageID string
	lastPageID    string
	pages         map[string]Component

	// 侧边栏
	sideBar *SideBar

	// 通知层
	notice *Notice

	// 刷新
	IsRefresh bool
}

// 类创建: 总页面
func New_Page(trunk *model.Trunk) *Page {
	self := &Page{}

	self.trunk = trunk
	self.Window = new(app.Window)
	self.Window.Option(
		app.Title(trunk.WindowTitle),
		app.Size(unit.Dp(450), unit.Dp(600)),
		app.Decorated(true),
	)

	// 窗口标题
	self.title = trunk.WindowTitle

	// 页面偏移高度
	self.pageOffsetHeight = 0

	// 初始化:主题
	self.style.theme = material.NewTheme()
	self.style.darkmode = NewDarkMode(self.title)
	self.style.darkmode.Apply(self.style.theme)

	// 初始化:加载文字
	fonts := loadFonts()
	shaper := text.NewShaper(text.WithCollection(fonts))
	self.style.theme.Shaper = shaper
	self.style.theme.TextSize = unit.Sp(18)

	// 初始化:通知层
	self.notice = New_Notice(self.style)

	// 初始化:不同页面
	self.currentPageID = "Home"
	self.lastPageID = self.currentPageID
	self.pages = make(map[string]Component)
	self.pages["Home"] = New_HomePage(self.style, self)
	self.pages["Chat"] = New_ChatPage(self.style, self)
	self.pages["Plugin"] = New_PluginPage(self.style, self)
	self.pages["Setting"] = New_SettingPage(self.style, self)

	self.IsRefresh = false

	return self
}

// 开始运行
func (self *Page) Run() {
	go func() {
		if err := self.draw(self.Window); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (self *Page) draw(w *app.Window) error {
	// 初始化
	var ops op.Ops

	// 暗色模式设置窗口颜色
	self.style.darkmode.RetrySetWindow(w, 10)

	// 初始化:侧边栏
	self.sideBar = New_SideBar(self.gtx, []SideBarElement{
		{
			ID:       "Home",
			Btn:      New_BtnWidthDesc(self.gtx, self.style, "\ue80f", self.trunk.LanguageTable.SideBar_Home),
			IsBottom: false,
			Callback: self.method_BtnHome,
		},
		{
			ID:       "Chat",
			Btn:      New_BtnWidthDesc(self.gtx, self.style, "\ued0d", self.trunk.LanguageTable.SideBar_Chat),
			IsBottom: false,
			Callback: self.method_BtnChat,
		},
		{
			ID:       "Plugin",
			Btn:      New_BtnWidthDesc(self.gtx, self.style, "\uE71D", self.trunk.LanguageTable.SideBar_Plugin),
			IsBottom: false,
			Callback: self.method_BtnPlugin,
		},
		{
			ID:       "Darkmode",
			Btn:      New_BtnWidthDesc(self.gtx, self.style, "\uf08c", self.trunk.LanguageTable.SideBar_DarkMode),
			IsBottom: true,
			Callback: self.method_BtnDarkMode,
		},
		{
			ID:       "Setting",
			Btn:      New_BtnWidthDesc(self.gtx, self.style, "\ue713", self.trunk.LanguageTable.SideBar_Setting),
			IsBottom: true,
			Callback: self.method_BtnSetting,
		},
	})
	sideBarRigid = self.sideBar.makeFlexChild()
	self.sideBar.API_SetCurrentID(self.currentPageID)

	// 事件驱动
	for {
		switch typ := w.Event().(type) {
		case app.DestroyEvent:
			return typ.Err

		case app.FrameEvent:
			ops.Reset()
			gtx := app.NewContext(&ops, typ)

			self.gtx = gtx
			self.style.darkmode.Update(gtx)
			paint.FillShape(self.gtx.Ops, self.style.darkmode.currentColor.Bg, clip.Rect{Max: self.gtx.Constraints.Max}.Op())
			self.Update(self.gtx)
			typ.Frame(self.gtx.Ops)
		}
	}
}

// 主要布局
func (self *Page) Update(gtx C) D {
	if self.IsRefresh {
		gtx.Execute(op.InvalidateCmd{At: time.Now().Add(time.Second / 60)})
		self.IsRefresh = false
	}

	for _, s := range self.sideBar.elements {
		id := s.ID
		switch id {
		case "Home":
			s.Btn.API_SetDesc(self.trunk.LanguageTable.SideBar_Home)
		case "Chat":
			s.Btn.API_SetDesc(self.trunk.LanguageTable.SideBar_Chat)
		case "Plugin":
			s.Btn.API_SetDesc(self.trunk.LanguageTable.SideBar_Plugin)
		case "Darkmode":
			s.Btn.API_SetDesc(self.trunk.LanguageTable.SideBar_DarkMode)
		case "Setting":
			s.Btn.API_SetDesc(self.trunk.LanguageTable.SideBar_Setting)
		}
	}

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// 左侧占位
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Dp(70)
					return D{Size: gtx.Constraints.Min}
				}),
				// 右侧主界面
				layout.Flexed(1, self.mainPageOffset),
			)
		}),

		// 侧边栏
		layout.Expanded(func(gtx C) D {
			return layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(15),
			}.Layout(gtx, func(gtx C) D {
				return self.sideBarFunc(gtx)
			})
		}),

		// 通知层
		layout.Expanded(func(gtx C) D {
			return self.notice.Update(gtx, self)
		}),
	)
}

// 内容主页面
func (self *Page) mainPageOffset(gtx C) D {
	// 内容偏移
	var offsetY int
	self.pageOffsetHeight = easing(gtx, self.pageOffsetHeight, 0, 0.3)
	offsetY = int(self.pageOffsetHeight)
	defer op.Offset(image.Pt(0, offsetY)).Push(gtx.Ops).Pop()
	dims := func(gtx C) D {
		return layout.Inset{
			Top:    unit.Dp(10),
			Bottom: unit.Dp(15),
			Left:   unit.Dp(0),
			Right:  unit.Dp(0),
		}.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(1200))
			return layout.Center.Layout(gtx, func(gtx C) D {
				return self.pages[self.currentPageID].Update(gtx)
			})
		})
	}(gtx)
	return dims
}

// 侧边栏
func (self *Page) sideBarFunc(gtx C) D {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(2, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, sideBarRigid...)
		}),
	)
}

// 页面 -------------------------------------------------------------------------------------------
// 侧边栏按钮功能
func (self *Page) method_BtnDarkMode() {
	self.style.darkmode.Toggle()
	self.style.darkmode.Apply(self.style.theme)
}
func (self *Page) method_BtnHome() {
	if self.lastPageID == "Home" {
		return
	}
	self.initPlaceholder()
	self.currentPageID = "Home"
	self.sideBar.API_SetCurrentID("Home")
	self.pages["Home"].Default()
	self.lastPageID = "Home"
	runtime.GC()
}
func (self *Page) method_BtnChat() {
	if self.lastPageID == "Chat" {
		return
	}
	self.initPlaceholder()
	self.currentPageID = "Chat"
	self.sideBar.API_SetCurrentID("Chat")
	self.pages["Chat"].Default()
	self.lastPageID = "Chat"
	runtime.GC()
}
func (self *Page) method_BtnPlugin() {
	if self.lastPageID == "Plugin" {
		return
	}
	self.initPlaceholder()
	self.currentPageID = "Plugin"
	self.sideBar.API_SetCurrentID("Plugin")
	self.pages["Plugin"].Default()
	self.lastPageID = "Plugin"
	runtime.GC()
}
func (self *Page) method_BtnSetting() {
	if self.lastPageID == "Setting" {
		return
	}
	self.initPlaceholder()
	self.currentPageID = "Setting"
	self.sideBar.API_SetCurrentID("Setting")
	self.pages["Setting"].Default()
	self.lastPageID = "Setting"
	runtime.GC()
}

// 当前偏移量设置为半页高
func (self *Page) initPlaceholder() {
	self.pageOffsetHeight = float64(self.gtx.Constraints.Max.Y) * 0.5
}

// API
func (self *Page) API() map[string]any {
	return map[string]any{
		"add_msg":    self.pages["Chat"].API()["add_msg"],
		"add_plugin": self.pages["Plugin"].API()["add_plugin"],
	}
}

// 将当前的页面跳转到聊天页面
func (self *Page) API_SetPageChat() {
	if self.currentPageID == "Chat" && self.sideBar.currentID == "Chat" {
		return
	}
	self.currentPageID = "Chat"
	self.sideBar.currentID = "Chat"
	if self.Window != nil {
		self.Window.Invalidate()
	}
}
