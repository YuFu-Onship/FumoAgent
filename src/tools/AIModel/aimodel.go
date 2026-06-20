package aichat

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"myapp/src/config"
	"myapp/src/model"

	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 定义一些结构 ----------------------------------

// 消息格式
type Message struct {
	Role    string
	Content string
}

// 聊天信息构建
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type ImageUrlStruct struct {
	Url string `json:"url"`
}
type ImageContent struct {
	Type     string         `json:"type"`
	ImageUrl ImageUrlStruct `json:"image_url"`
}
type AIContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type UserContent struct {
	Role    string        `json:"role"`
	Content []interface{} `json:"content"`
}

// api key ----------------------------------------

// 类 ----------------------------------------------------
type AIClient struct {
	Url   string
	Key   string
	Model string
	trunk *model.Trunk
}

func New_AIClient(trunk *model.Trunk) *AIClient {
	url := trunk.AiUrl
	model := trunk.AiModel
	key := trunk.AiKey
	return &AIClient{
		Url:   url,
		Key:   key,
		Model: model,
		trunk: trunk,
	}
}

// 图像处理
func (self *AIClient) ImageToBase64(image_data []byte) string {
	base64_str := base64.StdEncoding.EncodeToString(image_data)
	return base64_str
}

func (self *AIClient) ImagePathToContent(image_path string) []byte {
	image_data, err := os.ReadFile(image_path)
	if err != nil {
		log.Fatalf("无法读取图片文件: %v", err)
	}
	return image_data
}

// 构建图片
func (self *AIClient) BuildImageUrl(image_path string) string {
	image_data := self.ImagePathToContent(image_path)
	image_ext := filepath.Ext(image_path)[1:]
	base64_str := self.ImageToBase64(image_data)
	image_url := "data:image/" + image_ext + ";base64," + base64_str
	return image_url
}

// 解析文本
func (self *AIClient) ParseText(input string) []interface{} {
	var parts []interface{}
	re := regexp.MustCompile(`\[CQ:image,file:(.+?)\]`)
	match := re.FindStringSubmatch(input)
	if match != nil {
		image_path := match[1]
		image_url := self.BuildImageUrl(image_path)
		parts = append(parts, ImageContent{
			Type: "image_url",
			ImageUrl: ImageUrlStruct{
				Url: image_url,
			},
		})
		input = re.ReplaceAllString(input, "")
	}
	text := strings.TrimSpace(input)
	if text != "" {
		parts = append(parts, TextContent{
			Type: "text",
			Text: text,
		})
	}
	return parts
}

// 构建发送消息
func (self *AIClient) BuildMessageContent(patrs []interface{}, ai_prompt string) map[string]interface{} {
	ai_content := AIContent{
		Role:    "system",
		Content: ai_prompt,
	}
	user_content := UserContent{
		Role:    "user",
		Content: patrs,
	}
	req_body := map[string]interface{}{
		"model":    self.Model,
		"messages": []interface{}{ai_content, user_content},
	}
	return req_body
}

// ai聊天方面
func (self *AIClient) AIChat(req_body map[string]interface{}) (string, error) {
	// fmt.Printf("【请求】当前使用模型: %s, 实例指针: %p\n", self.Model, self)
	data, _ := json.Marshal(req_body)
	req, _ := http.NewRequest("POST", self.Url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+self.Key)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	type Resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	var result Resp
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	// 检查 choices 是否为空
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("AI 返回了空的 choices, body: %s", string(body))
	}
	return result.Choices[0].Message.Content, nil
}

// 通用函数：输入任意结构，输出 JSON 字符串
func (self *AIClient) ToJSONString(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ") // 格式化输出，方便查看
	if err != nil {
		log.Printf("JSON 编码出错: %v", err)
		return ""
	}
	return string(data)
}

// 终极的 ai 回复调用方案
func (self *AIClient) API_GetAIResponse(input string, ai_prompt string) string {
	parts := self.ParseText(input)
	req_body := self.BuildMessageContent(parts, ai_prompt)

	ai_response, err := self.AIChat(req_body)
	if err != nil {
		log.Printf("ai回复消息出错: %v", err)
		return fmt.Sprintf("[ERROR] %v", err)
	}
	return ai_response
}

// api
func (self *AIClient) API_SetClientCoeffcient(url string, model string, key string) {
	self.Url = url
	self.Model = url
	self.Key = key
}
func (self *AIClient) API_SetAiModelApi(detail config.Model_Preset) {
	// fmt.Printf("【设置】正在切换模型到: %s, 地址: %p\n", detail.Model, self)
	self.Url = detail.Url
	self.Model = detail.Model
	self.Key = detail.Key
}

// 调用示例
//  input := "描述一下图片内容[CQ:image,file:E:/Users/YUFU/Pictures/Screenshots/屏幕截图 2026-03-13 223748.png]"
// 	ai_prompt := "你是kimi, 由月之暗面公司开发的一款人工智能"
// 	ai_url := "https://api.moonshot.cn/v1/chat/completions"
// 	ai_model := "kimi-k2.5"
// 	api_key := "sk-XOLqhRnwV5DRiCRx2YwYMXtXdonB6MM6fBRCBH9NZJL0Gfii"
// 	airesponse := aichat.GetAIResponse(input, ai_prompt, ai_url, ai_model, api_key)
// 	fmt.Println(airesponse)
