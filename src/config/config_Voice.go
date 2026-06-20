package config

import (
	"encoding/json"
	"fmt"
	"io"
	fumovoice "myapp/src/tools/FumoVoice"
	"sync"

	"strings"

	"os"
	"path/filepath"
	"strconv"
)

// 语音方面	------------------------------------------------
// 直接从终端播放音频 --------------------------------------

func SOUND_PlayVoiceTemp_Aquestalk(content string, preset Preset_Aquestalk) {
	presetPath := filepath.Join(RootPath, "apps/aquestalkplayer/AquestalkPlayer.preset")

	name := "TEMP"

	csvData := []string{
		name,
		strconv.FormatBool(preset.MonoTone),
		preset.Type,
		preset.Engine,
		strconv.Itoa(preset.Speed),
		strconv.Itoa(preset.Volume),
		strconv.Itoa(preset.Pitch),
		strconv.Itoa(preset.Accent),
		strconv.Itoa(preset.Quality),
		strconv.Itoa(preset.Intonation),
		"",
	}
	fo := FileOperation{}
	fo.SaveCSV(presetPath, csvData)

	go func() {
		// aquestalk := fumovoice.NewAquestalk(RootPath)
		// aquestalk.TalkWithCoefficient("TEMP", content)
	}()
}

func SOUND_GetVoicePath_Aquestalk(content string, preset Preset_Aquestalk) string {
	presetPath := filepath.Join(RootPath, "apps/aquestalkplayer/AquestalkPlayer.preset")

	name := "TEMP"

	csvData := []string{
		name,
		strconv.FormatBool(preset.MonoTone),
		preset.Type,
		preset.Engine,
		strconv.Itoa(preset.Speed),
		strconv.Itoa(preset.Volume),
		strconv.Itoa(preset.Pitch),
		strconv.Itoa(preset.Accent),
		strconv.Itoa(preset.Quality),
		strconv.Itoa(preset.Intonation),
		"",
	}
	fo := FileOperation{}
	fo.SaveCSV(presetPath, csvData)

	aquestalk := fumovoice.NewAquestalk(RootPath)
	path := aquestalk.ExportWAV(content, "TEMP", filepath.Join(RootPath, "/data/temp/temp_aquestalk.wav"))
	return path
}

// 使用openjtalk播放,只接受平假音文本
func SOUND_PlayVoiceTemp_OpenJTalk(content string, preset Preset_OpenJTalk) {
	openJTalk := fumovoice.New_OpenJTalk()
	exePath := filepath.Join(RootPath, "/apps/openjtalk-windows-x64/open_jtalk.exe")
	dicPath := filepath.Join(RootPath, "/apps/openjtalk-windows-x64/dic")

	voiceFile := fmt.Sprintf("%s.htsvoice", preset.Speaker)
	speakerPath := filepath.Join(RootPath, "/apps/openjtalk-windows-x64/voice", voiceFile)

	outPath := filepath.Join(RootPath, "/data/temp/openJTalk.wav")

	speed := strconv.FormatFloat(preset.Speed, 'f', 2, 64)
	volume := strconv.FormatFloat(preset.Volume, 'f', 2, 64)
	pitch := strconv.FormatFloat(preset.Pitch, 'f', 2, 64)

	args := []string{exePath, dicPath, outPath, speakerPath, speed, volume, pitch}
	go func() {
		ct := fumovoice.NewCharToKata(filepath.Join(RootPath, Path_PinyinTable))
		openJTalk.Play(ct.GenKata(content), args)
	}()
}

// 返回openjtalk语音文件路径
func SOUND_GetVoicePath_OpenJTalk(content string, preset Preset_OpenJTalk) string {
	// outChan := make(chan string, 1)
	openJTalk := fumovoice.New_OpenJTalk()
	exePath := filepath.Join(RootPath, "/apps/openjtalk-windows-x64/open_jtalk.exe")
	dicPath := filepath.Join(RootPath, "/apps/openjtalk-windows-x64/dic")

	voiceFile := fmt.Sprintf("%s.htsvoice", preset.Speaker)
	speakerPath := filepath.Join(RootPath, "/apps/openjtalk-windows-x64/voice", voiceFile)

	outPath := filepath.Join(RootPath, "/data/temp/openJTalk.wav")

	speed := strconv.FormatFloat(preset.Speed, 'f', 2, 64)
	volume := strconv.FormatFloat(preset.Volume, 'f', 2, 64)
	pitch := strconv.FormatFloat(preset.Pitch, 'f', 2, 64)

	args := []string{exePath, dicPath, outPath, speakerPath, speed, volume, pitch}
	ct := fumovoice.NewCharToKata(filepath.Join(RootPath, Path_PinyinTable))
	path := openJTalk.ExportWav(ct.GenKata(content), args)
	return path
}

// 使用Teto语音播放,只接受日语,其他语言不能正常工作
func SOUND_PlayVoiceTemp_Teto(content string, preset Preset_Teto) {
	c := []string{
		filepath.Join(RootPath, Path_PinyinTable),
		filepath.Join(RootPath, "apps/重音テト音声ライブラリー/oto.ini"),
		filepath.Join(RootPath, "apps/重音テト音声ライブラリー/重音テト単独音"),
		filepath.Join(RootPath, "data/temp/teto_temp.wav"),
	}
	go func() {
		fumovoice.PlayTotoVoice(content, c)
	}()
}

// 使用voice vox播放语音
func SOUND_PlayVoiceInTemp_VoiceVox() {

}

// 读取htsvoices文件夹, 返回相关htsvoice后缀文件名
func API_SOUND_Get_AllHtsVoiceFiles() []string {
	folderPath := filepath.Join(RootPath, "apps/openjtalk-windows-x64/voice")
	var fileName []string

	files, err := os.ReadDir(folderPath)
	if err != nil {
		Error_Save(err, "读取HtsVoice文件失败")
		return fileName
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			// 判断文件名是否以 .htsvoice 结尾
			if strings.HasSuffix(name, ".htsvoice") {
				cleanName := strings.TrimSuffix(name, ".htsvoice")
				fileName = append(fileName, cleanName)
			}
		}
	}
	return fileName
}

// 预设文件读写操作 ---------------------------------------------------------------
type Config_Voice struct{ mu sync.RWMutex }

// 检查文件是否存在
func (self *Config_Voice) CheckInstance() bool {
	path := filepath.Join(RootPath, "/data/preset_Voice.json")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.WriteFile(path, []byte("[]"), 0644)
		if err != nil {
			return false
		}
	}
	return true
}

// 追加预设
func (self *Config_Voice) AppendContent(cont []string) bool {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.CheckInstance()
	path := filepath.Join(RootPath, "/data/preset_Voice.json")
	file, err := os.Open(path)
	if err != nil {
		return false
	}

	var data [][]string
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)

	if self.Check_repeat(data, cont[1]) {
		file.Close()
		return false
	}
	file.Close()

	if err != nil {
		return false
	}

	data = append(data, cont)
	wf, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return false
	}
	defer wf.Close()
	encoder := json.NewEncoder(wf)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return false
	}
	return true
}

// 返回所有的数据
func (self *Config_Voice) Get_Preset() ([][]string, error) {
	var data [][]string

	self.mu.Lock()
	defer self.mu.Unlock()

	path := filepath.Join(RootPath, "/data/preset_Voice.json")

	f, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// 返回指定id的预设,字符串
func (self *Config_Voice) Get_SpecifyPreset(name string) []string {
	data, _ := self.Get_Preset()
	for _, p := range data {
		if p[1] == name {
			return p
		}
	}
	if len(data) > 0 {
		return data[0]
	}
	return []string{
		"AquesTalk",
		"cirno",
		"false",
		"false",
		"AquesTalk10",
		"F1E",
		"154",
		"100",
		"100",
		"100",
		"100",
		"100",
		"",
	}
}

// 检查是否有重复
func (self *Config_Voice) Check_repeat(data [][]string, name string) bool {
	for _, p := range data {
		if p[1] == name {
			return true
		}
	}
	return false
}

// 删除指定预设
func (self *Config_Voice) Delete_Preset(name string) {
	new_data := [][]string{}
	data, _ := self.Get_Preset()
	for _, p := range data {
		if p[1] != name {
			new_data = append(new_data, p)
		}
	}
	path := filepath.Join(RootPath, "/data/preset_Voice.json")
	wf, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer wf.Close()
	encoder := json.NewEncoder(wf)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(new_data)
	if err != nil {
		return
	}
}

// 保存预设
func (self *Config_Voice) Save_Preset(preset any) bool {
	preset_str := self.ToString_Preset(preset)
	return self.AppendContent(preset_str)
}
func (self *Config_Voice) save_Aquestalk(preset Preset_Aquestalk) bool {
	engineType := "AquesTalk"
	preset_name := preset.Name
	preset_default := strconv.FormatBool(preset.Default)
	preset_monotype := strconv.FormatBool(preset.MonoTone)
	preset_type := preset.Type
	preset_engine := preset.Engine
	preset_speed := strconv.Itoa(preset.Speed)
	preset_volume := strconv.Itoa(preset.Volume)
	preset_pitch := strconv.Itoa(preset.Pitch)
	preset_accent := strconv.Itoa(preset.Accent)
	preset_quality := strconv.Itoa(preset.Quality)
	preset_intonation := strconv.Itoa(preset.Intonation)

	preset_str := []string{
		engineType,
		preset_name,
		preset_default,
		preset_monotype,
		preset_type,
		preset_engine,
		preset_speed,
		preset_volume,
		preset_pitch,
		preset_accent,
		preset_quality,
		preset_intonation,
		preset.Note,
	}

	return self.AppendContent(preset_str)
}
func (self *Config_Voice) save_OpenJTalk(preset Preset_OpenJTalk) bool {
	engineType := "OpenJTalk"
	preset_name := preset.Name
	preset_default := strconv.FormatBool(preset.Default)
	preset_speaker := preset.Speaker
	preset_speed := strconv.FormatFloat(preset.Speed, 'f', 1, 64)
	preset_volume := strconv.FormatFloat(preset.Volume, 'f', 0, 64)
	preset_pitch := strconv.FormatFloat(preset.Pitch, 'f', 0, 64)

	preset_set := []string{
		engineType,
		preset_name,
		preset_default,
		preset_speaker,
		preset_speed,
		preset_volume,
		preset_pitch,
	}
	return self.AppendContent(preset_set)
}
func (self *Config_Voice) save_VoiceVox(preset Preset_VoiceVox) bool {
	engineType := "VoiceVox"
	preset_name := preset.Name
	preset_default := strconv.FormatBool(preset.Default)
	preset_speaker := preset.Speaker
	preset_speed := strconv.FormatFloat(preset.Speed, 'f', 2, 64)
	preset_pitch := strconv.FormatFloat(preset.Pitch, 'f', 2, 64)
	preset_intonation := strconv.FormatFloat(preset.Intonation, 'f', 2, 64)
	preset_volume := strconv.FormatFloat(preset.Volume, 'f', 2, 64)

	preset_str := []string{
		engineType,
		preset_name,
		preset_default,
		preset_speaker,
		preset_speed,
		preset_pitch,
		preset_intonation,
		preset_volume,
	}
	return self.AppendContent(preset_str)
}

// 字符串->预设
func (self *Config_Voice) ToPreset_String(r []string) (any, error) {
	var preset any
	var err error
	switch r[0] {
	case "AquesTalk":
		preset, err = self.ConvertTo_Aquestalk(r)
	case "OpenJTalk":
		preset, err = self.ConvertTo_OpenJTalk(r)
	case "VoiceVox":
		preset, err = self.ConvertTo_VoiceVox(r)
	default:
		preset, err = nil, nil
	}
	return preset, err
}
func (self *Config_Voice) ConvertTo_Aquestalk(r []string) (Preset_Aquestalk, error) {
	if len(r) < 13 || r[0] != "AquesTalk" {
		return Preset_Aquestalk{}, io.ErrUnexpectedEOF
	}
	defaultBool, _ := strconv.ParseBool(r[2])
	monoTone, _ := strconv.ParseBool(r[3])
	speed, _ := strconv.Atoi(r[6])
	volume, _ := strconv.Atoi(r[7])
	pitch, _ := strconv.Atoi(r[8])
	accent, _ := strconv.Atoi(r[9])
	quality, _ := strconv.Atoi(r[10])
	intonation, _ := strconv.Atoi(r[11])

	return Preset_Aquestalk{
		Name:       r[1],
		Default:    defaultBool,
		MonoTone:   monoTone,
		Type:       r[4],
		Engine:     r[5],
		Speed:      speed,
		Volume:     volume,
		Pitch:      pitch,
		Accent:     accent,
		Quality:    quality,
		Intonation: intonation,
		Note:       r[12],
	}, nil
}
func (self *Config_Voice) ConvertTo_OpenJTalk(r []string) (Preset_OpenJTalk, error) {
	if len(r) < 7 || r[0] != "OpenJTalk" {
		return Preset_OpenJTalk{}, io.ErrUnexpectedEOF
	}
	defaultBool, _ := strconv.ParseBool(r[2])
	speed, _ := strconv.ParseFloat(r[4], 64)
	volume, _ := strconv.ParseFloat(r[5], 64)
	pitch, _ := strconv.ParseFloat(r[6], 64)

	return Preset_OpenJTalk{
		Name:    r[1],
		Default: defaultBool,
		Speaker: r[3],
		Speed:   speed,
		Volume:  volume,
		Pitch:   pitch,
	}, nil
}
func (self *Config_Voice) ConvertTo_VoiceVox(r []string) (Preset_VoiceVox, error) {
	if len(r) < 8 || r[0] != "VoiceVox" {
		return Preset_VoiceVox{}, io.ErrUnexpectedEOF
	}
	defaultBool, _ := strconv.ParseBool(r[2])
	speed, _ := strconv.ParseFloat(r[4], 64)
	pitch, _ := strconv.ParseFloat(r[5], 64)
	intonation, _ := strconv.ParseFloat(r[6], 64)
	volume, _ := strconv.ParseFloat(r[7], 64)

	return Preset_VoiceVox{
		Name:       r[1],
		Default:    defaultBool,
		Speaker:    r[3],
		Speed:      speed,
		Pitch:      pitch,
		Intonation: intonation,
		Volume:     volume,
	}, nil
}

// 预设->字符串
func (self *Config_Voice) ToString_Preset(preset any) []string {
	var s []string = []string{}
	switch preset := preset.(type) {
	default:
		return s
	case Preset_Aquestalk:
		s = self.toString_AquesTalk(preset)
	case Preset_OpenJTalk:
		s = self.toString_OpenJTalk(preset)
	case Preset_VoiceVox:
		s = self.toString_VoiceVox(preset)
	}
	return s
}
func (self *Config_Voice) toString_AquesTalk(preset Preset_Aquestalk) []string {
	engineType := "AquesTalk"
	preset_name := preset.Name
	preset_default := strconv.FormatBool(preset.Default)
	preset_monotype := strconv.FormatBool(preset.MonoTone)
	preset_type := preset.Type
	preset_engine := preset.Engine
	preset_speed := strconv.Itoa(preset.Speed)
	preset_volume := strconv.Itoa(preset.Volume)
	preset_pitch := strconv.Itoa(preset.Pitch)
	preset_accent := strconv.Itoa(preset.Accent)
	preset_quality := strconv.Itoa(preset.Quality)
	preset_intonation := strconv.Itoa(preset.Intonation)
	s := []string{
		engineType,
		preset_name,
		preset_default,
		preset_monotype,
		preset_type,
		preset_engine,
		preset_speed,
		preset_volume,
		preset_pitch,
		preset_accent,
		preset_quality,
		preset_intonation,
		preset.Note,
	}
	return s
}
func (self *Config_Voice) toString_OpenJTalk(preset Preset_OpenJTalk) []string {
	engineType := "OpenJTalk"
	preset_name := preset.Name
	preset_default := strconv.FormatBool(preset.Default)
	preset_speaker := preset.Speaker
	preset_speed := strconv.FormatFloat(preset.Speed, 'f', 1, 64)
	preset_volume := strconv.FormatFloat(preset.Volume, 'f', 0, 64)
	preset_pitch := strconv.FormatFloat(preset.Pitch, 'f', 0, 64)

	preset_set := []string{
		engineType,
		preset_name,
		preset_default,
		preset_speaker,
		preset_speed,
		preset_volume,
		preset_pitch,
	}
	return preset_set
}
func (self *Config_Voice) toString_VoiceVox(preset Preset_VoiceVox) []string {
	engineType := "VoiceVox"
	preset_name := preset.Name
	preset_default := strconv.FormatBool(preset.Default)
	preset_speaker := preset.Speaker
	preset_speed := strconv.FormatFloat(preset.Speed, 'f', 2, 64)
	preset_pitch := strconv.FormatFloat(preset.Pitch, 'f', 2, 64)
	preset_intonation := strconv.FormatFloat(preset.Intonation, 'f', 2, 64)
	preset_volume := strconv.FormatFloat(preset.Volume, 'f', 2, 64)

	preset_str := []string{
		engineType,
		preset_name,
		preset_default,
		preset_speaker,
		preset_speed,
		preset_pitch,
		preset_intonation,
		preset_volume,
	}
	return preset_str
}

// 获取 语音启用状态
func (self *Config_Voice) Get_EnableState() bool {
	data := ReadConfigDate()
	value := data.Voice_Enable
	switch value {
	case "true":
		return true
	case "false":
		return false
	default:
		return false
	}
}

// 设置 语音启用状态
func (self *Config_Voice) Set_EnableState(value bool) {
	data := ReadConfigDate()
	if value {
		data.Voice_Enable = "true"
	} else {
		data.Voice_Enable = "false"
	}
	SaveConfig(data)
}

// 获取 当前预设
func (self *Config_Voice) Get_CurPreset() any {
	data := ReadConfigDate()
	s := data.Voice_Preset
	if len(s) == 0 {
		return nil
	}
	p, _ := self.ToPreset_String(s)
	return p
}

// 获取 当前预设字符串
func (self *Config_Voice) Get_CurString() ([]string, bool) {
	data := ReadConfigDate()
	s := data.Voice_Preset
	if len(s) == 0 {
		return s, false
	}
	return s, true
}

// 设置 当前预设
func (self *Config_Voice) Set_CurPreset(preset any) {
	s := self.ToString_Preset(preset)
	data := ReadConfigDate()
	data.Voice_Preset = s
	SaveConfig(data)
}

// 老代码 -----------------------------------------------------------------------------------------
// 语音预设
type PresetDetail struct {
	Name       string // 标题
	Default    *bool  // 是否为默认
	MonoTone   *bool  // 1 	机械读
	VoiceType  string // 3 	声音类型
	Engine     string // 2 	引擎
	Speed      *int   // 4	语速	默认:100 	50-300
	Volume     *int   // 5 	音量	默认:100	0-300
	Pitch      *int   // 6	音高	默认:100	20-200
	Accent     *int   // 7	顿音	默认:100	0-200
	Quality    *int   // 8 	音质	默认:100	0-200
	Intonation *int   // 9	夹音	默认:100	50-200
	Note       string // 10	备注
}

// 得到当前的预设名
func API_SOUND_GetCurrentPresetTitle() string {
	date := ReadConfigDate()
	return date.CurrentVoice
}
func API_SOUND_GetCurrentPresetTitlle_NOREAD(data *Config) string {
	return data.CurrentVoice
}

// 获取语音是否开启
func API_SOUND_GetEnableState() bool {
	data := ReadConfigDate()
	return API_SOUND_GetEnableState_NOREAD(data)
}
func API_SOUND_GetEnableState_NOREAD(data *Config) bool {
	value := data.Voice_Enable
	switch value {
	case "true":
		return true
	case "false":
		return false
	default:
		return true
	}
}

// 启用或关闭语音
func API_SOUND_SetEnableState(state bool) {
	data := ReadConfigDate()
	if state {
		data.Voice_Enable = "true"
	} else {
		data.Voice_Enable = "false"
	}
	SaveConfig(data)
}
