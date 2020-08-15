package sequencer

import (
    "io/ioutil"
    "log"
    "os"
    p "github.com/fjania/froland/pkg/pattern"
    a "github.com/fjania/froland/pkg/audio"
)

// c := time.Tick(time.Minute/time.Duration(divisionsPerMinute))
type Sequencer struct {
    Timer   *Timer
    Pattern *p.Pattern
    Synth   *a.Kit
    Instruments map[string]bool
}

func NewSequencer(patternFile, kitName string) (*Sequencer, error) {
    x, err := a.Synth()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    instruments := make(map[string]bool)
    for i, _ := range x.Samples {
        instruments[i] = true
    }

    s := &Sequencer{
        Timer:   NewTimer(),
        Synth:   x,
        Instruments:   instruments,
    }

    s.LoadPattern(patternFile)

    return s, nil
}

func (s *Sequencer) LoadPattern(patternFile string) error {
    sep := string(os.PathSeparator)
    patternFilePath := ".." + sep + ".." + sep + "assets" + sep + "patterns" + sep + patternFile
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

    s.Timer.SetTempo(s.Pattern.BPM)
    s.Timer.SetDivisionsPerBeat(s.Pattern.DivisionsPerBeat)

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
                    if hit > 0  && s.Instruments[t.Instrument]{
                        s.Synth.Play(t.Instrument, float32(hit))
                    }
                }
            }
        }
    }()

    go s.Timer.Start()
}
