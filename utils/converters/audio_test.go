package converters

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/princesslunaequestrian/vinilify/utils"
)

func TestMixFull(t *testing.T) {
	effect := filepath.Join(utils.GetAssets(), "sounds", "vinyl.mp3")
	music := filepath.Join(utils.GetAssets(), "sounds", "music.mp3")
	outfile := "./mix.mp3"
	err := Mix(effect, music, outfile)
	if err != nil {
		t.Fatal("error: ", err)
	}
}

func TestConvertAudio(t *testing.T) {
	inputFile := "../../users/286946560/mix.mp3"
	outputFile := "../../users/286946560/mix.wav"
	cmd := exec.Command("ffmpeg", "-i", inputFile, outputFile)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func TestAudioEffect(t *testing.T) {
	inputFile := "../../users/286946560/mix.wav"
	outputFile := "../../users/286946560/mix_filtered.wav"
	cmd := exec.Command("sox", inputFile, outputFile, "sinc", "220-12k")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
}
