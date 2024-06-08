package ffmpeg

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
	require.NoError(t, os.Chdir("../../"))

	testFilesPath, err := os.Getwd()
	require.NoError(t, err)

	tmpFile := testFilesPath + "/testutil/files/yumi_dannasama.mp3"
	tmpFileNotExt := testFilesPath + "/testutil/files/yumi_dannasama"
	ffmpeg := Ffmpeg{
		ctx: ctx,
	}
	_, err = os.Stat(tmpFileNotExt + ".m4a")
	// 既に変換済みの場合は削除
	if err == nil {
		err = os.Remove(tmpFileNotExt + ".m4a")
		require.NoError(t, err)
	}
	t.Run("正常系", func(t *testing.T) {
		err = ffmpeg.ConversionAudioFile(tmpFile, tmpFileNotExt)
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
	require.NoError(t, os.Chdir("../../"))

	testFilesPath, err := os.Getwd()
	require.NoError(t, err)

	tmpFile := testFilesPath + "/testutil/files/yumi_dannasama.mp3"
	tmpFileNotExt := testFilesPath + "/testutil/files/yumi_dannasama"
	ffmpeg := Ffmpeg{
		ctx: ctx,
	}
	t.Run("正常系", func(t *testing.T) {
		_, err := ffmpeg.GetAudioFileSecond(tmpFile, tmpFileNotExt)
		assert.NoError(t, err)
	})
}

func TestPlayFfmpeg_Start(t *testing.T) {
	ctx := context.Background()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../../"))

	testFilesPath, err := os.Getwd()
	require.NoError(t, err)

	tmpFile := testFilesPath + "/testutil/files/yumi_dannasama.mp3"
	ffmpeg := Ffmpeg{
		ctx: ctx,
	}
	playFfmpeg := ffmpeg.NewPlayFFmpeg(tmpFile)
	err = playFfmpeg.Start()
	assert.NoError(t, err)
}

func TestPlayFfmpeg_Kill(t *testing.T) {
	ctx := context.Background()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../../"))

	testFilesPath, err := os.Getwd()
	require.NoError(t, err)

	tmpFile := testFilesPath + "/testutil/files/yumi_dannasama.mp3"
	ffmpeg := Ffmpeg{
		ctx: ctx,
	}
	playFfmpeg := ffmpeg.NewPlayFFmpeg(tmpFile)
	err = playFfmpeg.Start()
	assert.NoError(t, err)
	err = playFfmpeg.Kill()
	assert.NoError(t, err)
}

func TestPlayFfmpeg_StdoutPipe(t *testing.T) {
	ctx := context.Background()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../../"))

	testFilesPath, err := os.Getwd()
	require.NoError(t, err)

	tmpFile := testFilesPath + "/testutil/files/yumi_dannasama.mp3"
	ffmpeg := Ffmpeg{
		ctx: ctx,
	}
	playFfmpeg := ffmpeg.NewPlayFFmpeg(tmpFile)
	_, err = playFfmpeg.StdoutPipe()
	assert.NoError(t, err)
	err = playFfmpeg.Start()
	assert.NoError(t, err)
}
