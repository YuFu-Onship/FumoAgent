package conv

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

type LangConv struct {
	pinyinDictPath string
	pinyinDict     map[string]string
}

func New_LangConv(rootPath string) *LangConv {
	l := LangConv{
		pinyinDictPath: filepath.Join(rootPath, "/data/pinyin_table.json"),
	}
	l.pinyinDict = l.ReadPinyinFile()
	return &l
}

func (self *LangConv) ChineseToKatana(text string) string {
	var sb strings.Builder
	// 匹配中文字符的正则
	re := regexp.MustCompile(`[\p{Han}]`)
	// 拼音配置
	args := pinyin.NewArgs()

	for _, char := range text {
		charStr := string(char)
		if re.MatchString(charStr) {
			// 获取拼音
			py := pinyin.Pinyin(charStr, args)
			if len(py) > 0 && len(py[0]) > 0 {
				pinyinStr := py[0][0]
				if val, ok := self.pinyinDict[pinyinStr]; ok {
					sb.WriteString(val)
					continue
				}
			}
		}
		// 非中文或未匹配到则直接保留原样
		sb.WriteString(charStr)
	}
	return sb.String()
}

// 读取拼音对照文件
func (self *LangConv) ReadPinyinFile() map[string]string {
	file, err := os.ReadFile(self.pinyinDictPath)
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	var dict map[string]string
	if err := json.Unmarshal(file, &dict); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}
	return dict
}

// 获得拼音对照表
func (self *LangConv) API_GetPinyinDict() map[string]string {
	if self.pinyinDict == nil {
		self.pinyinDict = self.ReadPinyinFile()
	}
	return self.pinyinDict
}
