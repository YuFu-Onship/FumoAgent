package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type SideBarElement struct {
	ID       string
	Btn      *BtnWithDesc
	IsBottom bool
	Callback CallBackFunc
}

type SideBar struct {
	elements []SideBarElement
	gtx      layout.Context

	sideLineY_top_cur    float64
	sideLineY_top_tar    float64
	sideLineY_bottom_cur float64
	sideLineY_bottom_tar float64
	currentID            string
}

func New_SideBar(gtx layout.Context, elements []SideBarElement) *SideBar {
	return &SideBar{
		elements:  elements,
		gtx:       gtx,
		currentID: "",
	}
}

func (self *SideBar) Layout(gtx layout.Context, currentID string) layout.Dimensions {

	return layout.Dimensions{}
}

// 构造元素
func (self *SideBar) makeFlexChild() []layout.FlexChild {
	var children []layout.FlexChild
	// 定义间隙大小
	gap := layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout)

	for _, e := range self.elements {
		element := e
		if !element.IsBottom {
			// 如果不是第一个元素，在添加按钮前先加个间隙
			if len(children) > 0 {
				children = append(children, gap)
			}

			child := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return element.Btn.Update(gtx, self.currentID == element.ID, element.Callback)
			})
			children = append(children, child)
		}
	}

	// 中间的伸缩空间
	children = append(children, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		// 这里的减法操作可能引起约束越界，建议保持简单
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}))

	for _, e := range self.elements {
		element := e
		if element.IsBottom {
			// 在底部按钮上方添加间隙
			children = append(children, gap)

			child := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return element.Btn.Update(gtx, self.currentID == element.ID, element.Callback)
			})
			children = append(children, child)
		}
	}

	return children
}

// api ------------------------------------------------------------------------
func (self *SideBar) API_SetCurrentID(id string) {
	self.currentID = id
}

// 缓动动画
// func (self *SideBar) easing(cur float64, tar float64, coefficient float64) float64 {
// 	cur += (tar - cur) * coefficient
// 	if math.Abs(cur-tar) < 1 {
// 		cur = tar
// 	} else {
// 		self.gtx.Execute(op.InvalidateCmd{At: self.gtx.Now.Add(time.Second / 60)})
// 	}
// 	return cur
// }
