package fumovoice

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type OpenJTalk struct {
}

func New_OpenJTalk() *OpenJTalk {
	ojt := OpenJTalk{}
	return &ojt
}

// coe
//
// [0] exePath
//
// [1] dicPath
//
// [2] outPath
//
// [3] speakerPath
//
// [4] speed
//
// [5] volume
//
// [6] pitch
func (self *OpenJTalk) Play(content string, coe []string) {
	binPath := coe[0]
	dicPath := coe[1]
	outputPath := coe[2]
	speakerPath := coe[3]

	inputText, _ := self.utf8ToShiftJIS(content)

	// 2. 构建命令参数
	args := []string{
		"-x", dicPath, // dicPath
		"-m", speakerPath, // speaker
		"-ow", outputPath, // outputPath
		"-p", "200",
		"-r", coe[4], // 语速 1 -- 2
		"-g", coe[5], // 音量 -30 -- 30
		"-fm", coe[6], // 音高 -30 -- 30
		"-z", "6000",
	}

	// 3. 创建命令对象
	cmd := exec.Command(binPath, args...)

	// 4. 处理编码问题 (关键！)
	// Open JTalk 在 Windows 上通常期望 Shift_JIS 输入。
	// 我们需要将文字写入命令的标准输入。
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close() // 这个 defer 属于这个匿名 goroutine！
		fmt.Fprint(stdin, inputText)
	}()

	// 5. 执行并等待结果
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("执行失败: %s\n错误信息: %s\n", err, string(output))
		return
	}
}

func (self *OpenJTalk) ExportWav(content string, coe []string) string {
	binPath := coe[0]
	dicPath := coe[1]
	outputPath := coe[2]
	speakerPath := coe[3]
	inputText, _ := self.utf8ToShiftJIS(content)
	args := []string{
		"-x", dicPath, // dicPath
		"-m", speakerPath, // speaker
		"-ow", outputPath, // outputPath
		"-p", "200",
		"-r", coe[4], // 语速 0.5 -- 2
		"-g", coe[5], // 音量 -30 -- 30
		"-fm", coe[6], // 音高 -30 -- 30
	}
	cmd := exec.Command(binPath, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close() // 这个 defer 属于这个匿名 goroutine！
		fmt.Fprint(stdin, inputText)
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("执行失败: %s\n错误信息: %s\n", err, string(output))
		return ""
	}
	return outputPath
}

// 编码转换器, utf-8转shiftjis
func (self *OpenJTalk) utf8ToShiftJIS(s string) (string, error) {
	I := strings.NewReader(s)
	O := transform.NewReader(I, japanese.ShiftJIS.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return "", e
	}
	return string(d), nil
}
