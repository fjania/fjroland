package sequencer

import (
	a "github.com/fjania/fjroland/pkg/audio"
	m "github.com/fjania/fjroland/pkg/audio/midi"
	w "github.com/fjania/fjroland/pkg/audio/waveform"
	p "github.com/fjania/fjroland/pkg/pattern"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Sequencer struct {
	Timer           *Timer
	Pattern         *p.Pattern
	AudioOutputs    []a.AudioOutput
	PatternFilePath string
}

func NewSequencer(patternFilePath string) (*Sequencer, error) {
	s := &Sequencer{
		Timer:           NewTimer(),
		PatternFilePath: patternFilePath,
	}

	err := s.LoadPattern(s.PatternFilePath)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Sequencer) ConfigureSamplesOutput(samplePackPath string) error {

	o, err := w.NewSamplePack(samplePackPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.AudioOutputs = append(s.AudioOutputs, o)

	return nil
}

func (s *Sequencer) ConfigureMidiOutput(deviceName string) error {

	o, err := m.NewMidi(deviceName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.AudioOutputs = append(s.AudioOutputs, o)

	return nil
}

func (s *Sequencer) LoadPattern(patternFilePath string) error {

	jsonFile, err := os.Open(patternFilePath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	jsonBlob, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()

	pattern, err := p.ParsePattern(jsonBlob)
	if err != nil {
		log.Fatal(err)
		return err
	}

	s.Pattern = pattern
	s.Timer.SetStepInterval(s.Pattern.BPM, s.Pattern.DivisionsPerBeat)

	s.WatchPatternFile()

	return nil
}

func (s *Sequencer) Start() {
	go handlePulses(s)
	go s.Timer.Start()
}

func handlePulses(s *Sequencer) {
	pulseCount := 0

	for {
		select {
		case <-s.Timer.Pulses:
			if pulseCount >= s.Pattern.Divisions {
				pulseCount = 0
			}

			RenderPattern(s, pulseCount)
			s.playAtPulse(pulseCount)

			pulseCount++
		}
	}
}

func (s *Sequencer) playAtPulse(pulse int) {
	for _, o := range s.AudioOutputs {
		for _, t := range s.Pattern.Tracks {
			hit := t.Steps[pulse]
			if hit > 0 && o.HasInstrument(t.Instrument) {
				o.Play(t.Instrument, float32(hit))
			}
		}
	}
}

func (s *Sequencer) IsInstrumentAvailable(i string) bool {
	for _, o := range s.AudioOutputs {
		if o.HasInstrument(i) {
			return true
		}
	}

	return false
}

func (s *Sequencer) AudioDeviceNameList() string {
	var outputs = make([]string, len(s.AudioOutputs))
	for i, o := range s.AudioOutputs {
		outputs[i] = o.Name()
	}
	return strings.Join(outputs, ", ")
}

func (s *Sequencer) WatchPatternFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Error starting watcher %s", s.PatternFilePath)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					s.LoadPattern(s.PatternFilePath)
				}

			case err := <-watcher.Errors:
				log.Printf("Error in watching %s: %v", s.PatternFilePath, err)
			}
		}
	}()

	if err := watcher.Add(s.PatternFilePath); err != nil {
		log.Printf("Error watching %s", s.PatternFilePath)
	}
}
