package on_message_create

import (
	"context"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
)

type Ffmpeg struct {
	tmpFile       string
}

func NewFfmpeg(tmpFile string) *Ffmpeg {
	return &Ffmpeg{
		tmpFile:       tmpFile,
	}
}

func (f Ffmpeg) ConversionAudioFile(ctx context.Context, tmpFile, tmpFileNotExt string) error {
	err := exec.CommandContext(ctx, "ffmpeg", "-i", tmpFile, tmpFileNotExt+".m4a").Run()
	return err
}

func (f Ffmpeg) GetAudioFileSecond(ctx context.Context, tmpFile, tmpFileNotExt string) (float64, error) {
	cmd := exec.CommandContext(
		ctx,
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

type FfmpegMock struct {
	ConversionAudioFileFunc func(ctx context.Context, tmpFile, tmpFileNotExt string) error
	GetAudioFileSecondFunc  func(ctx context.Context, tmpFile, tmpFileNotExt string) (float64, error)
}

func (f FfmpegMock) ConversionAudioFile(ctx context.Context, tmpFile, tmpFileNotExt string) error {
	return f.ConversionAudioFileFunc(ctx, tmpFile, tmpFileNotExt)
}

func (f FfmpegMock) GetAudioFileSecond(ctx context.Context, tmpFile, tmpFileNotExt string) (float64, error) {
	return f.GetAudioFileSecondFunc(ctx, tmpFile, tmpFileNotExt)
}

type FfmpegInterface interface {
	ConversionAudioFile(ctx context.Context, tmpFile, tmpFileNotExt string) error
	GetAudioFileSecond(ctx context.Context, tmpFile, tmpFileNotExt string) (float64, error)
}

var (
	_ FfmpegInterface = (*Ffmpeg)(nil)
	_ FfmpegInterface = (*FfmpegMock)(nil)
)
