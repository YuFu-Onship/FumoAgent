package fumovoice

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-audio/wav"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// 音素结构
type PhonemeStructure struct {
	Start   []byte
	Content []byte
	End     []byte
}

type PhonemeEntry struct {
	Data         []byte
	Alias        string // 别名，如 "あ"
	Offset       int    // 左边界 (ms)
	Consonant    int    // 固定范围 (ms)
	Cutoff       int    // 右边界 (ms)
	Preutterance int    // 先行发声 (ms)
	Overlap      int    // 叠加 (ms)
}

// 音素处理
type PhonemeProcesser struct {
	wavData   []PhonemeEntry
	CurData   []byte
	OutBuffer []byte
	// wavData     [][]byte
	LastPhoneme PhonemeEntry
}

// Teto 语音部分
type OtoEntry struct {
	Alias        string // 别名，如 "あ"
	Offset       int    // 左边界 (ms)
	Consonant    int    // 固定范围 (ms)
	Cutoff       int    // 右边界 (ms)
	Preutterance int    // 先行发声 (ms)
	Overlap      int    // 叠加 (ms)
}

// KataToHira 片假转平假
type ContrastTable map[string]string
type ContrastJson struct {
	FullWidth ContrastTable `json:"fullWidth"`
	HalfWidth ContrastTable `json:"halfWidth"`
}

// 从终端播放teto音频
//
// coe[0] pinyin路径
//
// coe[1] teto音索引文件
//
// coe[2] teto单独音路径
//
// coe[3] wav路径
func PlayTotoVoice(content string, coe []string) {
	voiceDir := coe[2]
	outPath := coe[3]

	char2katana := NewCharToKata(coe[0])
	text := char2katana.GenKata(content)
	hiraText := KatakanaToHiragana(text)

	otoDict := ReturnOtoData(coe[1])

	// var clipdata [][]byte
	var wavData []PhonemeEntry

	for _, char := range []rune(hiraText) {
		key := string(char)
		fileName := fmt.Sprintf("_%s.wav", string(char))
		fullPath := filepath.Join(voiceDir, fileName)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// fmt.Printf("警告: 找不到音频文件 %s，已跳过\n", fullPath)
			// if char == 'ー' {
			// 	key_ := string([]rune(hiraText)[index-1])
			// 	fmt.Println(key_)
			// 	fileName_ := fmt.Sprintf("_%s.wav", key_)
			// 	fullPath_ := filepath.Join(voiceDir, fileName_)
			// 	wavData = append(wavData, voice.PhonemeEntry{
			// 		Data:         voice.ReturnWavData(fullPath_),
			// 		Alias:        otoDict[key_].Alias,
			// 		Offset:       otoDict[key_].Offset,
			// 		Consonant:    otoDict[key_].Consonant,
			// 		Cutoff:       otoDict[key_].Cutoff,
			// 		Preutterance: otoDict[key_].Preutterance,
			// 		Overlap:      otoDict[key_].Overlap,
			// 	})
			// }
			continue
		}

		wavData = append(wavData, PhonemeEntry{
			Data:         ReturnWavData(fullPath),
			Alias:        otoDict[key].Alias,
			Offset:       otoDict[key].Offset,
			Consonant:    otoDict[key].Consonant,
			Cutoff:       otoDict[key].Cutoff,
			Preutterance: otoDict[key].Preutterance,
			Overlap:      otoDict[key].Overlap,
		})
	}

	PP := New_PhonemeSplicing(wavData)
	fullData := PP.PhonemeSplicing()
	WriteWAV(outPath, fullData)
	PlayWavSimple(outPath)
}

// KatakanaToHiragana 将字符串中的片假名转为平假名
func KatakanaToHiragana(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		// 片假名范围是 0x30A1 (ア) 到 0x30F6 (ヶ)
		if r >= 0x30A1 && r <= 0x30F6 {
			runes[i] = r - 0x60
		}
	}
	return string(runes)
}

// 导入并解析对照表 --------------------------------------------------------------------------------------
func LoadOtoMap(path string) (map[string]string, error) {
	otoMap := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 使用 Shift-JIS 解码器转换流
	reader := transform.NewReader(f, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "=") {
			continue
		}

		// 分解行：_あ.wav=あ,24,56...
		parts := strings.Split(line, "=")
		wavFile := parts[0] // "_あ.wav"

		// 获取别名部分（逗号前的内容）
		meta := strings.Split(parts[1], ",")
		alias := meta[0] // "あ"

		// 存入 map：otoMap["あ"] = "_あ.wav"
		otoMap[alias] = wavFile
	}
	return otoMap, scanner.Err()
}

// 音素处理相关 ------------------------------------------------------------------------------------------

func New_PhonemeSplicing(wavData []PhonemeEntry) *PhonemeProcesser {
	pp := PhonemeProcesser{
		wavData: wavData,
	}
	return &pp
}

// 音频拼接
func (self *PhonemeProcesser) PhonemeSplicing() []byte {
	totalSize := self.count_DataLength()
	// fmt.Println(totalSize)
	fullFile := make([]byte, 0, 44+totalSize)

	// 创建并添加文件头
	header := self.build_Header(int(totalSize))
	fullFile = append(fullFile, header...)

	// 创建结构列表,用于后续拼接
	var phonemeStructures []PhonemeStructure = make([]PhonemeStructure, 0, len(self.wavData))
	for _, d := range self.wavData {
		structure := self.split_WavData(d)

		// 语速处理
		structure.Content = self.SpeedUpPCM(structure.Content, 1.1)
		phonemeStructures = append(phonemeStructures, structure)
	}

	// 开始拼接
	phonemeLength := len(phonemeStructures)
	if phonemeLength != 0 {

		for i := 0; i <= phonemeLength; i += 1 {
			switch i {
			case 0:
				fullFile = append(fullFile, self.clip_WavData(self.wavData[i].Data, 0, self.wavData[i].Offset)...)
				fullFile = append(fullFile, phonemeStructures[i].Content...)
			case phonemeLength:
				fullFile = append(fullFile, self.clip_WavData(
					self.wavData[i-1].Data,
					self.wavData[i-1].Offset+self.wavData[i-1].Consonant,
					self.wavData[i-1].Offset+self.wavData[i-1].Consonant+self.wavData[i-1].Cutoff)...)
			default:
				nextPart := self.splic_WavData(phonemeStructures[i].Start, phonemeStructures[i].Content)
				overlapPart := self.overlap_WavData(phonemeStructures[i-1].End, nextPart)
				fullFile = append(fullFile, overlapPart...)
			}
		}
		return fullFile
	}
	return fullFile
}

// 创建文件头
func (self *PhonemeProcesser) build_Header(totalSize int) []byte {
	// 准备 Header (44 字节)
	// 参数假设：单声道(1), 采样率 44100, 位深 16bit
	header := make([]byte, 44)

	// RIFF 标识
	copy(header[0:4], "RIFF")
	// ChunkSize: 36 + totalAudioSize
	chunkSize := 36 + totalSize
	header[4] = byte(chunkSize)
	header[5] = byte(chunkSize >> 8)
	header[6] = byte(chunkSize >> 16)
	header[7] = byte(chunkSize >> 24)

	copy(header[8:12], "WAVE")

	// fmt chunk
	copy(header[12:16], "fmt ")
	header[16] = 16 // Subchunk1Size (PCM 为 16)
	header[17] = 0
	header[18] = 0
	header[19] = 0

	header[20] = 1 // AudioFormat (1 为 PCM)
	header[21] = 0
	header[22] = 1 // NumChannels (单声道)
	header[23] = 0

	sampleRate := uint32(44100)
	header[24] = byte(sampleRate)
	header[25] = byte(sampleRate >> 8)
	header[26] = byte(sampleRate >> 16)
	header[27] = byte(sampleRate >> 24)

	byteRate := sampleRate * 1 * 2 // SampleRate * Channels * BitsPerSample/8
	header[28] = byte(byteRate)
	header[29] = byte(byteRate >> 8)
	header[30] = byte(byteRate >> 16)
	header[31] = byte(byteRate >> 24)

	header[32] = 2 // BlockAlign (Channels * BitsPerSample/8)
	header[33] = 0
	header[34] = 16 // BitsPerSample
	header[35] = 0

	// data chunk
	copy(header[36:40], "data")
	header[40] = byte(totalSize)
	header[41] = byte(totalSize >> 8)
	header[42] = byte(totalSize >> 16)
	header[43] = byte(totalSize >> 24)
	return header
}

// 重叠数据
func (self *PhonemeProcesser) overlap_WavData(a []byte, b []byte) []byte {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	outData := make([]byte, maxLen)

	for i := 0; i < maxLen; i += 2 {
		var sampleA, sampleB int16
		if i+1 < len(a) {
			sampleA = int16(binary.LittleEndian.Uint16(a[i : i+2]))
		}
		if i+1 < len(b) {
			sampleB = int16(binary.LittleEndian.Uint16(b[i : i+2]))
		}
		mixed := int32(sampleA) + int32(sampleB)
		if mixed > 32767 {
			mixed = 32767
		} else if mixed < -32768 {
			mixed = -32768
		}
		binary.LittleEndian.PutUint16(outData[i:i+2], uint16(mixed))
	}
	return outData
}

// 拼接数据
func (self *PhonemeProcesser) splic_WavData(a []byte, b []byte) []byte {
	outData := make([]byte, 0, len(a)+len(b))
	outData = append(outData, a...)
	outData = append(outData, b...)
	return outData
}

// 切割原始数据,返回:预发音, 正文, 重叠 三部分
func (self *PhonemeProcesser) split_WavData(entry PhonemeEntry) PhonemeStructure {
	wavData := entry.Data
	start := self.clip_WavData(wavData, entry.Offset-entry.Preutterance, entry.Offset)
	end := self.clip_WavData(wavData, entry.Offset+entry.Consonant, entry.Offset+entry.Consonant+entry.Overlap)

	contentData := self.clip_WavData(wavData, entry.Offset, entry.Offset+entry.Consonant)
	startData := self.strengthen_phoneme(start)
	endData := self.weaken_phoneme(end)
	return PhonemeStructure{
		Start:   startData,
		Content: contentData,
		End:     endData,
	}
}

// 计算数据的总长度
func (self *PhonemeProcesser) count_DataLength() int {
	var totalSize uint32
	for _, entry := range self.wavData {
		totalSize += uint32(len(entry.Data))
	}
	return int(totalSize)
}

// 计算指定时长(ms)内的字节长度
func (self *PhonemeProcesser) count_msLength(duration int) int {
	return int(float64(duration) * 88.2)
}

// 切割数据
func (self *PhonemeProcesser) clip_WavData(wavdata []byte, startMs int, endMs int) []byte {
	const headerSize = 44

	// 1. 提取纯 PCM 数据部分
	if len(wavdata) <= headerSize {
		return nil
	}
	pcmData := wavdata[headerSize:]

	// 2. 计算字节位置
	startByte := int(float64(startMs) * 88.2) // bytesPerMs := 44100 * 2 * 1 / 1000 // = 88.2 字节/毫秒
	endByte := int(float64(endMs) * 88.2)

	// 3. 强制偶数对齐 (针对 16bit)
	startByte = (startByte >> 1) << 1
	endByte = (endByte >> 1) << 1

	// 4. 边界安全检查
	if startByte < 0 {
		startByte = 0
	}
	if endByte > len(pcmData) {
		endByte = len(pcmData)
	}
	if startByte >= endByte {
		return nil
	}

	// 5. 复制一份数据，避免底层数组共享导致的意外
	result := make([]byte, endByte-startByte)
	copy(result, pcmData[startByte:endByte])
	return result
}

// 加强字节 (Fade In - 淡入)
func (self *PhonemeProcesser) strengthen_phoneme(data []byte) []byte {
	size := len(data)
	if size < 2 {
		return data
	}

	newData := make([]byte, size)
	// 以 2 字节为一个采样点处理
	for i := 0; i < size-1; i += 2 {
		// 1. 读取原始 16位 采样点
		sample := int16(binary.LittleEndian.Uint16(data[i : i+2]))

		// 2. 计算缩放比例 (0.0 到 1.0)
		// 注意：我们要按采样点的进度算，而不是字节进度
		ratio := float64(i) / float64(size)

		// 3. 缩放并转回 int16
		newSample := int16(float64(sample) * ratio)

		// 4. 写回新数组
		binary.LittleEndian.PutUint16(newData[i:i+2], uint16(newSample))
	}
	return newData
}

// 减弱字节 (Fade Out - 淡出)
func (self *PhonemeProcesser) weaken_phoneme(data []byte) []byte {
	size := len(data)
	if size < 2 {
		return data
	}

	newData := make([]byte, size)
	for i := 0; i < size-1; i += 2 {
		sample := int16(binary.LittleEndian.Uint16(data[i : i+2]))

		// 比例从 1.0 降到 0.0
		ratio := float64(size-i) / float64(size)

		newSample := int16(float64(sample) * ratio)
		binary.LittleEndian.PutUint16(newData[i:i+2], uint16(newSample))
	}
	return newData
}

// 增加语速
func (self *PhonemeProcesser) SpeedUpPCM(data []byte, speed float64) []byte {
	if speed <= 1.0 {
		return data // 暂不支持减速，或原速返回
	}

	const bitsPerSample = 2
	originalSamples := len(data) / bitsPerSample
	newSamples := int(float64(originalSamples) / speed)
	newData := make([]byte, newSamples*bitsPerSample)

	for i := 0; i < newSamples; i++ {
		// 计算在原数据中的位置
		pos := float64(i) * speed
		index := int(pos)

		if index*bitsPerSample+1 < len(data) {
			// 这里简单取值，进阶可以用线性插值 (Linear Interpolation)
			sample := data[index*bitsPerSample : index*bitsPerSample+2]
			newData[i*bitsPerSample] = sample[0]
			newData[i*bitsPerSample+1] = sample[1]
		}
	}
	return newData
}

func readHiraFile() ContrastJson {
	path := "src/voice/KataToHira.json"
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}
	var cj ContrastJson
	if err := json.Unmarshal(file, &cj); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}
	// fmt.Printf(cj.FullWidth["ア"])
	return cj
}

func TranslateKataToHira(text string) string {
	data := readHiraFile()

	var outText []rune
	runes := []rune(text)
	for _, char := range runes {
		strChar := string(char)
		if val, ok := data.FullWidth[strChar]; ok {
			outText = append(outText, []rune(val)...)
			continue
		}
		if val, ok := data.HalfWidth[strChar]; ok {
			outText = append(outText, []rune(val)...)
			continue
		}
		outText = append(outText, char)
	}
	return string(outText)
}

func LoadVoiceData() {
	file, err := os.Open("src/voice/oto.ini")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	sjisDecoder := japanese.ShiftJIS.NewDecoder()
	utf8Descoder := transform.NewReader(file, sjisDecoder)
	reader := csv.NewReader(utf8Descoder)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record[0])
	}
}

func GetWAVMetaData() {
	inputFile, _ := os.Open(`assets\TETO-tandoku-100619\重音テト音声ライブラリー\重音テト単独音\_あ.wav`)
	defer inputFile.Close()

	decoder := wav.NewDecoder(inputFile)
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		panic(err)
	}
	// fmt.Println(buf)
	//

	sampleRate := decoder.SampleRate
	numChannels := int(decoder.NumChans)

	startOffset := int(24*time.Millisecond/time.Second*time.Duration(sampleRate)) * numChannels
	endOffset := int(56*time.Millisecond/time.Second*time.Duration(sampleRate)) * numChannels

	// 边界检查
	if endOffset > len(buf.Data) {
		endOffset = len(buf.Data)
	}
	fmt.Println(len(buf.Data))

	// 3. 提取切片
	clippedData := buf.Data[startOffset:endOffset]

	// 4. 写入新文件
	outPath := "output_clipped.wav"
	outputFile, _ := os.Create(outPath)
	defer outputFile.Close()

	encoder := wav.NewEncoder(outputFile, int(sampleRate), int(decoder.BitDepth), numChannels, int(decoder.WavAudioFormat))

	// 创建一个新的 Buffer 存放截取后的数据
	newBuf := buf
	newBuf.Data = clippedData

	encoder.Write(newBuf)
	encoder.Close()

	// ShowWavData(outPath)
}

func ShowWavData(path string) {
	inputFile, _ := os.Open(path)
	defer inputFile.Close()

	decoder := wav.NewDecoder(inputFile)
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		panic(err)
	}
	fmt.Println(buf)
}

func ClipWav(wavPath string, start int, end int) {
	data, _ := os.ReadFile(wavPath)
	headerSize := 44
	sampleRate := 44100
	bitDepth := 16
	channels := 1
	bytesPerSample := (bitDepth / 8) * channels

	startByte := int(float64(start) * 0.001 * float64(sampleRate) * float64(bytesPerSample))
	endByte := int(float64(end) * 0.001 * float64(sampleRate) * float64(bytesPerSample))

	audioData := data[headerSize:]
	clippedData := audioData[startByte:endByte]

	newHeader := make([]byte, 44)
	copy(newHeader, data[:44]) // 复制原头信息

	// 更新 Header 中的 Subchunk2Size (Offset 40, 4字节) -> 即音频数据大小
	dataSize := uint32(len(clippedData))
	newHeader[40] = byte(dataSize)
	newHeader[41] = byte(dataSize >> 8)
	newHeader[42] = byte(dataSize >> 16)
	newHeader[43] = byte(dataSize >> 24)

	// 更新 Header 中的 ChunkSize (Offset 4, 4字节) -> 即数据大小 + 36
	chunkSize := dataSize + 36
	newHeader[4] = byte(chunkSize)
	newHeader[5] = byte(chunkSize >> 8)
	newHeader[6] = byte(chunkSize >> 16)
	newHeader[7] = byte(chunkSize >> 24)

	// 6. 写入文件
	output, _ := os.Create("clipped.wav")
	output.Write(newHeader)
	output.Write(clippedData)
	output.Close()

	// fmt.Printf("截取完成，截取数据大小: %d 字节\n", dataSize)

}

// 返回数据
func ReturnWavData(path string) []byte {
	data, _ := os.ReadFile(path)
	return data
}

// 切片数据文件
func ClipDataWav(wavdata []byte, startMs int, endMs int) []byte {
	const headerSize = 44
	// 假设 44100Hz, 16bit, Mono
	//

	// 1. 提取纯 PCM 数据部分
	if len(wavdata) <= headerSize {
		return nil
	}
	pcmData := wavdata[headerSize:]

	// 2. 计算字节位置
	startByte := int(float64(startMs) * 88.2)
	endByte := int(float64(endMs) * 88.2)

	// 3. 强制偶数对齐 (针对 16bit)
	startByte = (startByte >> 1) << 1
	endByte = (endByte >> 1) << 1

	// 4. 边界安全检查
	if startByte < 0 {
		startByte = 0
	}
	if endByte > len(pcmData) {
		endByte = len(pcmData)
	}
	if startByte >= endByte {
		return nil
	}

	// 5. 复制一份数据，避免底层数组共享导致的意外
	result := make([]byte, endByte-startByte)
	copy(result, pcmData[startByte:endByte])
	return result
}

// 构造新文件
func BuildFullWav(clips [][]byte) []byte {
	// 1.计算所有音频数据的总长度
	var totalAudioSize uint32
	for _, clip := range clips {
		totalAudioSize += uint32(len(clip))
	}

	// 2. 准备 Header (44 字节)
	// 参数假设：单声道(1), 采样率 44100, 位深 16bit
	header := make([]byte, 44)

	// RIFF 标识
	copy(header[0:4], "RIFF")
	// ChunkSize: 36 + totalAudioSize
	chunkSize := 36 + totalAudioSize
	header[4] = byte(chunkSize)
	header[5] = byte(chunkSize >> 8)
	header[6] = byte(chunkSize >> 16)
	header[7] = byte(chunkSize >> 24)

	copy(header[8:12], "WAVE")

	// fmt chunk
	copy(header[12:16], "fmt ")
	header[16] = 16 // Subchunk1Size (PCM 为 16)
	header[17] = 0
	header[18] = 0
	header[19] = 0

	header[20] = 1 // AudioFormat (1 为 PCM)
	header[21] = 0
	header[22] = 1 // NumChannels (单声道)
	header[23] = 0

	sampleRate := uint32(44100)
	header[24] = byte(sampleRate)
	header[25] = byte(sampleRate >> 8)
	header[26] = byte(sampleRate >> 16)
	header[27] = byte(sampleRate >> 24)

	byteRate := sampleRate * 1 * 2 // SampleRate * Channels * BitsPerSample/8
	header[28] = byte(byteRate)
	header[29] = byte(byteRate >> 8)
	header[30] = byte(byteRate >> 16)
	header[31] = byte(byteRate >> 24)

	header[32] = 2 // BlockAlign (Channels * BitsPerSample/8)
	header[33] = 0
	header[34] = 16 // BitsPerSample
	header[35] = 0

	// data chunk
	copy(header[36:40], "data")
	header[40] = byte(totalAudioSize)
	header[41] = byte(totalAudioSize >> 8)
	header[42] = byte(totalAudioSize >> 16)
	header[43] = byte(totalAudioSize >> 24)

	// 3. 拼接所有数据
	fullFile := make([]byte, 0, 44+totalAudioSize)
	fullFile = append(fullFile, header...)
	for _, clip := range clips {
		fullFile = append(fullFile, clip...)
	}

	// fmt.Println(fullFile)
	return fullFile
}

func WriteWAV(outPath string, data []byte) {
	os.WriteFile(outPath, data, 0644)
}

// 返回索引数据
func ReturnOtoData(path string) map[string]OtoEntry {
	var otoDict map[string]OtoEntry = make(map[string]OtoEntry)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sjisDecoder := japanese.ShiftJIS.NewDecoder()
	utf8Descoder := transform.NewReader(file, sjisDecoder)
	reader := csv.NewReader(utf8Descoder)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var key string
		re := regexp.MustCompile(`_(.*?)\.wav`)
		match := re.FindStringSubmatch(record[0])
		if len(match) > 1 {
			key = match[1]
		}

		offset, _ := strconv.Atoi(record[1])
		consonant, _ := strconv.Atoi(record[2])
		cutoff, _ := strconv.Atoi(record[3])
		pre, _ := strconv.Atoi(record[4])
		overlap, _ := strconv.Atoi(record[5])

		otoDict[key] = OtoEntry{
			Alias:        record[0],
			Offset:       offset,
			Consonant:    consonant,
			Cutoff:       cutoff,
			Preutterance: pre,
			Overlap:      overlap,
		}
	}
	return otoDict
}

// play sound
func PlayWavSimple(filePath string) {
	powershellCmd := fmt.Sprintf(`(New-Object Media.SoundPlayer "%s").PlaySync()`, filePath)
	// 执行命令
	cmd := exec.Command("powershell", "-c", powershellCmd)
	cmd.Run()
}

// 日语处理 ---------------------------------------------------------------------------------------------
func JapaneseToKata(text string) string {
	t, _ := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	tokens := t.Tokenize(text)

	var sb strings.Builder
	// for _, token := range tokens {
	// 	if token.Class == tokenizer.DUMMY {
	// 		continue
	// 	}
	// 	reading, ok := token.Reading()
	// 	if ok {
	// 		sb.WriteString(reading)
	// 	} else {
	// 		sb.WriteString(token.Surface)
	// 	}
	// }

	for _, token := range tokens {
		if reading, ok := token.Reading(); ok {
			// 调用上面的手动转换函数
			sb.WriteString(KatakanaToHiragana(reading))
		} else {
			sb.WriteString(token.Surface)
		}
	}

	return sb.String()
}
