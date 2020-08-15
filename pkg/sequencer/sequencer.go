package sequencer

import (
    "io/ioutil"
    "log"
    "os"
    p "github.com/fjania/froland/pkg/pattern"
    a "github.com/fjania/froland/pkg/audio"
    //w "github.com/fjania/froland/pkg/audio/waveform"
    m "github.com/fjania/froland/pkg/audio/midi"
)

// c := time.Tick(time.Minute/time.Duration(divisionsPerMinute))
type Sequencer struct {
    Timer   *Timer
    Pattern *p.Pattern
    Output   a.Output
    PatternFilePath string
}

func NewSequencer(patternFile, kitName string) (*Sequencer, error) {
    //x, err := w.NewSamplePack()
    x, err := m.NewMidi()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    s := &Sequencer{
        Timer:   NewTimer(),
        Output:   x,
    }

    s.LoadPattern(patternFile)

    return s, nil
}

func (s *Sequencer) LoadPattern(patternFile string) error {
    sep := string(os.PathSeparator)

    patternFilePath := ".." + sep + ".." +
        sep + "assets" + sep +
        "patterns" + sep + patternFile

    s.PatternFilePath = patternFilePath

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
