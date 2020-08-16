package sequencer

import (
    "io/ioutil"
    "log"
    "os"
    p "github.com/fjania/fjroland/pkg/pattern"
    a "github.com/fjania/fjroland/pkg/audio"
    w "github.com/fjania/fjroland/pkg/audio/waveform"
    m "github.com/fjania/fjroland/pkg/audio/midi"
)

// c := time.Tick(time.Minute/time.Duration(divisionsPerMinute))
type Sequencer struct {
    Timer   *Timer
    Pattern *p.Pattern
    Output   a.Output
    PatternFilePath string
}

func NewSequencer(patternFilePath, output, samplePackPath string) (*Sequencer, error) {
    s := &Sequencer{
        Timer:   NewTimer(),
    }

    s.PatternFilePath = patternFilePath
    s.LoadPattern(patternFilePath)

    if output == "samples" {
        o, err := w.NewSamplePack(samplePackPath)
        if err != nil {
            log.Fatal(err)
            return nil, err
        }
        s.Output = o
    } else if output == "midi" {
        o, err := m.NewMidi()
        if err != nil {
            log.Fatal(err)
            return nil, err
        }
        s.Output = o
    }

    return s, nil
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

    return nil
}

func (s *Sequencer) Start() {
    go func() {
        pulseCount := 0

        for {
            select {
            case <-s.Timer.Pulses:
                pulse := pulseCount % s.Pattern.Divisions
                RenderPattern(s, pulse)
                pulseCount++
                for _, t := range s.Pattern.Tracks {
                    hit := t.Steps[pulse]
                    if hit > 0  && s.Output.HasInstrument(t.Instrument){
                        s.Output.Play(t.Instrument, float32(hit))
                    }
                }
            }
        }
    }()

    go s.Timer.Start()
}
