package audio

import (
	"github.com/gordonklaus/portaudio"
	"github.com/mkb218/gosndfile/sndfile"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Sample struct {
	Buffer   []float32
	Info     sndfile.Info
	Playhead int
	Volume   float32
}

type SamplePack struct {
	Samples        map[string]*Sample
	Instruments    map[string]bool
	SamplePackPath string
}

func NewSamplePack(samplePackPath string) (*SamplePack, error) {
	samplePack, err := LoadSamplePack(samplePackPath)
	if err != nil {
		return nil, err
	}

	err = portaudio.Initialize()
	if err != nil {
		return nil, err
	}

	// We'll use the default audio device and assume stereo wav files
	// sampled at 44.1k. This could be changed, but all the samples
	// provided conform, and it's a reasonable default.
	defaultOutputDevice, err := portaudio.DefaultOutputDevice()
	if err != nil {
		log.Fatal("Failed to access the default audio device on this computer")
		return nil, err
	}
	p := portaudio.LowLatencyParameters(nil, defaultOutputDevice)
	p.Input.Channels = 0
	p.Output.Channels = 2
	p.SampleRate = 44100.0
	p.FramesPerBuffer = portaudio.FramesPerBufferUnspecified
	stream, err := portaudio.OpenStream(p, samplePack.ProcessAudio)

	if err != nil {
		log.Fatal("Failed to open a stream to the default audio device on this computer")
		return nil, err
	}

	stream.Start()
	return samplePack, nil
}

func LoadSamplePack(samplePackPath string) (*SamplePack, error) {
	k := &SamplePack{
		Samples:        make(map[string]*Sample),
		SamplePackPath: samplePackPath,
	}

	files, err := ioutil.ReadDir(samplePackPath)
	if err != nil {
		log.Fatalf("Could not load sample pack: '%s'\n", samplePackPath)
	}

	for _, f := range files {
		sampleFilePath := samplePackPath + string(os.PathSeparator) + f.Name()
		if strings.HasSuffix(f.Name(), ".wav") {
			instrument := strings.TrimRight(f.Name(), ".wav")

			sample, err := LoadSample(sampleFilePath)
			if err != nil {
				log.Fatalf("Could not load sample: %s\n", sampleFilePath)
			}
			k.Samples[instrument] = sample
		}
	}

	instruments := make(map[string]bool)
	for i, _ := range k.Samples {
		instruments[i] = true
	}
	k.Instruments = instruments

	return k, nil
}

func LoadSample(filepath string) (*Sample, error) {
	var info sndfile.Info
	soundFile, err := sndfile.Open(filepath, sndfile.Read, &info)

	s := &Sample{
		Buffer:   make([]float32, info.Samplerate*info.Channels),
		Info:     info,
		Playhead: 0,
		Volume:   1,
	}

	if err != nil {
		log.Printf("Could not open file: %s\n", filepath)
		return nil, err
	}

	_, err = soundFile.ReadItems(s.Buffer)
	if err != nil {
		log.Printf("Error reading data from file: %s\n", filepath)
		return nil, err
	}

	s.Playhead = len(s.Buffer)

	defer soundFile.Close()

	return s, nil
}

func (s *Sample) Play(level float32) {
	// Input levels are 1, 2, 3. To make them sound different to the human
	// ear we'll square the level to make the volume the sample should be
	// played at.
	s.Volume = level * level
	s.Playhead = 0
}

// Use the "playhead" method of playing back the samples. Setting the
// playhead to 0 means that we'll loop back trough the sample's data and
// add that to the []float32 which portaudio uses to write to the audio
// device on the computer. We do this for every sample that should be
// played, and clamp the output at 1.0 to avoid clipping.
//
// Found this method of playing back the audio in this repo
// https://github.com/kellydunn/go-step-sequencer
func (k *SamplePack) ProcessAudio(out []float32) {
	for i := range out {
		var data float32
		for _, s := range k.Samples {
			if s.Playhead < len(s.Buffer) {
				data += s.Volume * s.Buffer[s.Playhead]
				s.Playhead++
			}
		}

		if data > 1.0 {
			data = 1.0
		}

		out[i] = data
	}
}

// Implement the AudioOutput interface
func (k *SamplePack) Name() string {
	return k.SamplePackPath
}

func (k *SamplePack) Play(instrument string, level float32) {
	k.Samples[instrument].Play(level)
}

func (k *SamplePack) ListInstruments() map[string]bool {
	return k.Instruments
}

func (k *SamplePack) HasInstrument(instrument string) bool {
	return k.Instruments[instrument]
}
