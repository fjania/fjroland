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
		// Bail if we try to create a sequencer but we can't find the pattern
		// file. If we assume we're only creating a sequencer from the CLI, then
		// this makes sense. This could changed if we needed to load a sequencer
		// before we had any content to play on it.
		log.Fatalf("Error: Can't start up the seqencer without a pattern file.")
	}

	// Watch for changes to the pattern file. This will let us update the pattern
	// as the user edits it. (We assume we're not loading new pattern files in this
	// version of the sequencer, so we don't worry about being able to change the
	// file we're watching)
	s.WatchPatternFile()

	return s, nil
}

// We'll allow multiple sample outputs for no good reason other than
// maybe people want to experiment with layering different intstruments?
func (s *Sequencer) ConfigureSamplesOutput(samplePackPath string) error {

	o, err := w.NewSamplePack(samplePackPath)
	if err != nil {
		log.Fatalf("Failed to load sample pack: '%s'", samplePackPath)
	}
	s.AudioOutputs = append(s.AudioOutputs, o)

	return nil
}

// We'll allow multiple midi outputs for no good reason other than
// maybe people want to experiment with layering different intstruments?
func (s *Sequencer) ConfigureMidiOutput(deviceName string) error {

	o, err := m.NewMidi(deviceName)
	if err != nil {
		log.Fatalf("Failed to configure midi device: '%s'", deviceName)
	}
	s.AudioOutputs = append(s.AudioOutputs, o)

	return nil
}

func (s *Sequencer) LoadPattern(patternFilePath string) error {

	jsonFile, err := os.Open(patternFilePath)
	if err != nil {
		// Note - We don't fatal here in case the sqeucer is running
		// and we want to let the user try to load a new patter again.
		log.Printf("Failed to load pattern file: '%s'", patternFilePath)
		return err
	}
	jsonBlob, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()

	pattern, err := p.ParsePattern(jsonBlob)
	if err != nil {
		// Note - We don't fatal here in case the sqeucer is running
		// and we want to let the user try to load a new patter again.
		log.Printf("Failed to parse pattern file: '%s'", patternFilePath)
		log.Printf("Error: %s", err)
		return err
	}

	s.Pattern = pattern
	s.Timer.SetStepInterval(s.Pattern.BPM, s.Pattern.DivisionsPerBeat)
	return nil
}

// All we really need to do to start is to tell the sequencer to start
// listening for pusles, and then start the pulses.
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

// Play a note any time the note's level is > 0, for each audio output
// as long as that instrument exists in that audio output
func (s *Sequencer) playAtPulse(pulse int) {
	for _, t := range s.Pattern.Tracks {
		noteLevel := t.Steps[pulse]
		if noteLevel > 0 {
			for _, o := range s.AudioOutputs {
				if o.HasInstrument(t.Instrument) {
					o.Play(t.Instrument, float32(noteLevel))
				}
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
