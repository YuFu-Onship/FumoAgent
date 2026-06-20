package netserver

import "encoding/json"

type Message struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 向ws客户端发送的命令结构
// type Command struct {
// 	Category  string   `json:"Category"`
// 	Name      string   `json:"Name"`
// 	Arguement []string `json:"Arguement"`
// }

type Command struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

// 将命令结构转换为json字符串
func CommandToJson(command Command) string {
	jsonByte, err := json.Marshal(command)
	if err != nil {
		panic(err)
	}
	return string(jsonByte)
}

func ToJson(msg Message) string {
	jsonByte, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(jsonByte)
}
