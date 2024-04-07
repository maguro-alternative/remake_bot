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
	tmpFileNotExt string
}

func (f Ffmpeg) ConversionAudioFile(ctx context.Context) error {
	err := exec.CommandContext(ctx, "ffmpeg", "-i", f.tmpFile, f.tmpFileNotExt+".m4a").Run()
	return err
}

func (f Ffmpeg) GetAudioFileSecond(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(
		ctx,
		"ffprobe",
		"-hide_banner",
		f.tmpFileNotExt+".m4a",
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
