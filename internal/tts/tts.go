package tts

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/X3NOOO/suri/routes"
	"github.com/ebitengine/oto/v3"
)

const PIPER_OUT_FORMAT = "s16le"

var (
	ErrMissingConfig = errors.New("missing required fields")
	ErrNilAudio      = errors.New("the audio is nil")
)

type piperModelConfig struct {
	Audio struct {
		SampleRate int `json:"sample_rate"`
	} `json:"audio"`
	NumSpeakers int `json:"num_speakers"`
}

type PiperAudio struct {
	audio       []byte
	sampleRate  int
	numSpeakers int
}

type PiperTTS struct {
	BinPath         string
	ModelPath       string
	ModelConfigPath string
}

// This is merely an interface for the Piper binary
func (p PiperTTS) Generate(text string) (routes.Audio, error) {
	if p.BinPath == "" || p.ModelPath == "" || p.ModelConfigPath == "" {
		return nil, ErrMissingConfig
	}

	piperPath, err := exec.LookPath(p.BinPath)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(p.ModelPath)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(p.ModelConfigPath)
	if err != nil {
		return nil, err
	}

	modelConfigJson, err := os.ReadFile(p.ModelConfigPath)
	if err != nil {
		return nil, err
	}

	var piperModel piperModelConfig
	err = json.Unmarshal(modelConfigJson, &piperModel)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(piperPath, "--quiet", "--model", p.ModelPath, "--config", p.ModelConfigPath, "--output-raw")
	log.Println("Running piper:", cmd.String())

	cmd.Stdin = strings.NewReader(text)

	output := bytes.Buffer{}
	cmd.Stdout = bufio.NewWriter(&output)

	err = cmd.Run()
	if err != nil {
		log.Println("Error running piper:", err)
		return nil, err
	}

	return &PiperAudio{
		audio:       output.Bytes(),
		sampleRate:  piperModel.Audio.SampleRate,
		numSpeakers: piperModel.NumSpeakers,
	}, nil
}

func (a *PiperAudio) Play() error {
	if a.audio == nil {
		return ErrNilAudio
	}

	op := &oto.NewContextOptions{
		SampleRate:   a.sampleRate,
		ChannelCount: a.numSpeakers,
		Format:       oto.FormatSignedInt16LE,
	}

	log.Printf("Audio player config: %+v\n", *op)

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		log.Println("Error creating new audio context:", err)
		return err
	}
	<-readyChan

	player := otoCtx.NewPlayer(bytes.NewReader(a.audio))

	player.Play()

	return nil
}

func (a *PiperAudio) Wav() ([]byte, error) {
	if a.audio == nil {
		return nil, ErrNilAudio
	}

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(ffmpegPath, "-f", PIPER_OUT_FORMAT, "-ar", strconv.Itoa(a.sampleRate), "-ac", strconv.Itoa(a.numSpeakers), "-i", "pipe:0", "-f", "wav", "pipe:1")
	log.Println("Running ffmpeg:", cmd.String())

	cmd.Stdin = bytes.NewReader(a.audio)

	output := bytes.Buffer{}
	cmd.Stdout = bufio.NewWriter(&output)

	err = cmd.Run()
	if err != nil {
		log.Println("Error running ffmpeg:", err)
		return nil, err
	}

	return output.Bytes(), nil
}
