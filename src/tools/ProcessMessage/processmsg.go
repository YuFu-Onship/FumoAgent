package processMsg

import (
	"myapp/src/model"
	"regexp"
)

type ProcessMessage struct {
	trunk model.Trunk
}

func New_ProcessMessage(trunk model.Trunk) *ProcessMessage {
	return &ProcessMessage{
		trunk: trunk,
	}
}

// 处理ai消息, 包括删除关键词, 调用相关功能
// <name:func,args:test>
func (self *ProcessMessage) Process_AIMessage(text string) {
	re := regexp.MustCompile(`<name:(?P<func>[^,]+),args:(?P<args>[^>]+)>`)
	match := re.FindStringSubmatch(text)

	if len(match) > 0 {
		funcName := match[1]
		// argsValue := match[2]
		switch funcName {
		case "Emotion":
			// self.trunk.Handler_Server.SendCommand("func", "Emotion", []string{argsValue})
		}
	}
}
