package on_message_create

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFfmpeg_ConversionAudioFile(t *testing.T) {
	ctx := context.Background()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../.."))
	tmpFile := cwd + "/yumi_dannasama.mp3"
	tmpFileNotExt := cwd + "/yumi_dannasama"
	ffmpeg := Ffmpeg{
		tmpFile: tmpFile,
	}
	_, err = os.Stat(tmpFileNotExt + ".m4a")
	// 既に変換済みの場合は削除
	if err == nil {
		err = os.Remove(tmpFileNotExt + ".m4a")
		require.NoError(t, err)
	}
	t.Run("正常系", func(t *testing.T) {
		err = ffmpeg.ConversionAudioFile(ctx, tmpFile, tmpFileNotExt)
		assert.NoError(t, err)
	})
}

func TestFfmpeg_GetAudioFileSecond(t *testing.T) {
	ctx := context.Background()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../.."))
	tmpFile := cwd + "/yumi_dannasama.mp3"
	tmpFileNotExt := cwd + "/yumi_dannasama"
	ffmpeg := Ffmpeg{
		tmpFile: tmpFile,
	}
	t.Run("正常系", func(t *testing.T) {
		_, err := ffmpeg.GetAudioFileSecond(ctx, tmpFile, tmpFileNotExt)
		assert.NoError(t, err)
	})
}