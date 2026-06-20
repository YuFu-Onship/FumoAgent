package ui

type CheckBox struct {
	style Style
}

func New_CheckBox(style Style) *CheckBox {
	return &CheckBox{
		style: style,
	}
}
func (self *CheckBox) Update(gtx C) D {
	return D{}
}
