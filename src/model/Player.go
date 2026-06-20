package model

import (
	"math"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

type SimplePlayer struct {
	ctrl        *beep.Ctrl
	volume      *effects.Volume
	currentFile *os.File
	streamer    beep.StreamSeekCloser
	currentVol  float64
	initialized bool // 记录全局扬声器是否已经初始化
}

func NewSimplePlayer() *SimplePlayer {
	return &SimplePlayer{
		currentVol: 2.0,
	}
}

// 解析音频文件,返回字节串
func (self *SimplePlayer) ParamsWav(path string) []byte {
	return []byte{}
}

// 解析音频文件, 返回语音大小的归一化值
func (self *SimplePlayer) ParamsVolume(path string) ([]float64, error) {
	volumeList := []float64{}
	f, err := os.Open(path)
	if err != nil {
		return []float64{}, err
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return volumeList, err
	}
	defer streamer.Close()

	// 客户端的渲染帧率
	const targetFPS = 40
	samplesPerFrame := int(format.SampleRate) / targetFPS

	// 如果音频特别短，连一帧都凑不够，给个保底大小
	if samplesPerFrame <= 0 {
		samplesPerFrame = 512
	}
	var rawVolumes []float64
	maxVolume := 0.0

	buffer := make([][2]float64, samplesPerFrame)
	for {
		n, ok := streamer.Stream(buffer)
		if n == 0 || !ok {
			break
		}

		var sum float64 = 0
		for i := range n {
			sum += (math.Abs(buffer[i][0]) + math.Abs(buffer[i][1])) / 2.0
		}
		avgVolume := sum / float64(n)

		if avgVolume > maxVolume {
			maxVolume = avgVolume
		}
		rawVolumes = append(rawVolumes, avgVolume)
	}

	// 模糊处理
	volumeList = make([]float64, len(rawVolumes))
	radius := 2
	length := len(rawVolumes)

	valueAll := 0.0
	for i := range rawVolumes {
		valueAll += rawVolumes[i]
	}

	scale := 0.06 / (valueAll / float64(length))
	for i := range rawVolumes {
		rawVolumes[i] = rawVolumes[i] * scale
		if i < radius {
			volumeList[i] = rawVolumes[i]
		} else if i > (length - radius - 1) {
			volumeList[i] = rawVolumes[i]
		} else {
			var sum float64
			count := 0
			for j := i - radius; j <= i+radius; j++ {
				sum += rawVolumes[j]
				count++
			}
			volumeList[i] = math.Round(sum/float64(count)*100) / 100
		}
	}

	return volumeList, nil
}

// 播放传递进来的wav文件
func (self *SimplePlayer) PlayFile(path string) error {
	if self.streamer != nil {
		self.streamer.Close()
	}
	if self.currentFile != nil {
		self.currentFile.Close()
	}

	const targetSampleRate = beep.SampleRate(44100)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		f.Close()
		return err
	}
	self.currentFile = f
	self.streamer = streamer

	// init是全局的, 只能初始化一次
	if !self.initialized {
		err = speaker.Init(targetSampleRate, targetSampleRate.N(time.Second/10))
		if err != nil {
			self.streamer.Close()
			self.currentFile.Close()
			return err
		}
		self.initialized = true
	}

	ratio := float64(format.SampleRate) / float64(targetSampleRate)
	speedCorrectedStreamer := beep.ResampleRatio(4, ratio, streamer)
	self.ctrl = &beep.Ctrl{
		Streamer: speedCorrectedStreamer,
		Paused:   false,
	}

	self.volume = &effects.Volume{
		Streamer: self.ctrl,
		Base:     2,
		Volume:   self.currentVol,
		Silent:   false,
	}

	done := make(chan struct{})
	sequence := beep.Seq(self.volume, beep.Callback(func() {
		close(done)
	}))

	// speaker.Play(self.volume)
	speaker.Play(sequence)

	<-done

	return nil
}

// 设置音量
func (p *SimplePlayer) SetVolume(vol float64) {
	speaker.Lock()
	p.currentVol = vol
	if p.volume != nil {
		p.volume.Volume = vol
	}
	speaker.Unlock()
}

// 解除文件占用
func (self *SimplePlayer) ReleaseFile() {
	if self.currentFile == nil {
		return
	}
	self.currentFile.Close()
	self.streamer.Close()
	speaker.Clear()
}

// 停止播放
func (p *SimplePlayer) Close() {
	speaker.Lock()
	if p.ctrl != nil {
		p.ctrl.Paused = true
	}
	speaker.Unlock()

	if p.streamer != nil {
		p.streamer.Close()
	}
	if p.currentFile != nil {
		p.currentFile.Close()
	}
}
