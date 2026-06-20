package model

import "myapp/src/config"

type LangZone map[string]string
type LangPack map[string]LangZone

// ai api
type ModelAI interface {
	API_GetAIResponse(input string, prompt string) string
	API_SetAiModelApi(detail config.Model_Preset)
}

// 历史消息相关
type ModelMsg interface {
	API_GetContent() [][]string
	API_Clear()
	API_Add(msg []string) error
}

// 与ws客户端进行通信
type ModelComm interface {
	SendCommand(category string, name string, arguement []string)
}

// 处理消息, 调用api相关
type ModelProcessMsg any

// 个性化相关
type ModelCustom interface {
	API_SetLanguage(string)
	API_SetColor(string)
}

// 语音
type ModelVoice interface {
	Say(content string, engine string)
}

// 插件结构
type PluginMeta struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}

// Model Plugin
type ModelPlugin interface {
	Execute(pm PluginMeta) (string, bool)
	Name() string
	Desc() string
	Params() string
	API_SetLanguage(language string)
	API_GetID() string
	API_SetValueString(map[string]string)
}

// GUI部分接口
type GuiApi interface {
	API() map[string]any
}
