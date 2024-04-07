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
	ffmpeg := Ffmpeg{
		tmpFile:       cwd+"/yumi_dannasama.mp3",
		tmpFileNotExt: cwd+"/yumi_dannasama",
	}
	_, err = os.Stat(ffmpeg.tmpFileNotExt+".m4a")
	// 既に変換済みの場合は削除
	if err == nil {
		err = os.Remove(ffmpeg.tmpFileNotExt+".m4a")
		require.NoError(t, err)
	}
	t.Run("正常系", func(t *testing.T) {
		err = ffmpeg.ConversionAudioFile(ctx)
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
	tmpFile := cwd+"/yumi_dannasama.mp3"
	tmpFileNotExt := cwd+"/yumi_dannasama"
	ffmpeg := Ffmpeg{
		tmpFile:       tmpFile,
		tmpFileNotExt: tmpFileNotExt,
	}
	t.Run("正常系", func(t *testing.T) {
		_, err := ffmpeg.GetAudioFileSecond(ctx)
		assert.NoError(t, err)
	})
}
