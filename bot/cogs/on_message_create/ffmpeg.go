package on_message_create

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
)

type Ffmpeg struct{
	ctx context.Context
}

func NewFfmpeg(ctx context.Context) *Ffmpeg {
	return &Ffmpeg{
		ctx: ctx,
	}
}

func (f Ffmpeg) ConversionAudioFile(tmpFile, tmpFileNotExt string) error {
	cmd := exec.CommandContext(
		f.ctx,
		"ffmpeg",
		"-i",
		tmpFile,
		tmpFileNotExt+".m4a",
	)
	return cmd.Run()
}

func (f Ffmpeg) GetAudioFileSecond(tmpFile, tmpFileNotExt string) (float64, error) {
	cmd := exec.CommandContext(
		f.ctx,
		"ffprobe",
		"-hide_banner",
		tmpFileNotExt+".m4a",
		"-show_entries",
		"format=duration",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0.0, err
	}
	re, err := regexp.Compile(`(\d+\.\d+)`)
	if err != nil {
		return 0.0, err
	}
	match := re.FindStringSubmatch(string(out))
	if len(match) == 0 {
		return 0.0, errors.New("not found duration")
	}
	sec, err := strconv.ParseFloat(match[0], 64)
	if err != nil {
		return 0.0, err
	}

	return sec, nil
}

func (f *Ffmpeg) Start(output string) error {
	cmd := exec.CommandContext(
		f.ctx,
		"ffmpeg",
		"-i",
		output,
		"-f",
		"s16le",
		"-ar",
		"48000",
		"-ac",
		"2",
		"pipe:1",
	)
	return cmd.Run()
}

func (f *Ffmpeg) Play(buf *bufio.Reader, send chan []int16) error {
	for {
		audiobuf := make([]int16, 960*2)
		if err := binary.Read(buf, binary.LittleEndian, &audiobuf); err != nil {
			return err
		}
		select {
		case send <- audiobuf:
			continue
		case <-f.ctx.Done():
			return nil
		}
	}
}

type FfmpegMock struct {
	ConversionAudioFileFunc func(tmpFile, tmpFileNotExt string) error
	GetAudioFileSecondFunc  func(tmpFile, tmpFileNotExt string) (float64, error)
}

func (f FfmpegMock) ConversionAudioFile(tmpFile, tmpFileNotExt string) error {
	return f.ConversionAudioFileFunc(tmpFile, tmpFileNotExt)
}

func (f FfmpegMock) GetAudioFileSecond(tmpFile, tmpFileNotExt string) (float64, error) {
	return f.GetAudioFileSecondFunc(tmpFile, tmpFileNotExt)
}

type FfmpegInterface interface {
	ConversionAudioFile(tmpFile, tmpFileNotExt string) error
	GetAudioFileSecond(tmpFile, tmpFileNotExt string) (float64, error)
}

var (
	_ FfmpegInterface = (*Ffmpeg)(nil)
	_ FfmpegInterface = (*FfmpegMock)(nil)
)
