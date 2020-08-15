package sequencer

import (
    "log"
    p "github.com/fjania/froland/pkg/pattern"
    a "github.com/fjania/froland/pkg/audio"
)

// c := time.Tick(time.Minute/time.Duration(divisionsPerMinute))
type Sequencer struct {
    Timer   *Timer
    Pattern *p.Pattern
    Synth   *a.Kit
}

func NewSequencer() (*Sequencer, error) {
    /*
        "Snare:         |----|X---|----|X---|",
        "Closed Hi-Hat: |X-X-|X-X-|X-X-|X-X-|",
        "Bass:          |X---|X---|X---|X---|"
        "Snare:         |----|X---|----|X---|",
        "Snare:         |X---|----|X---|----|"
    */
    var jsonBlob = []byte(`
    {"title": "Turn Down for What",
    "bpm": 100,
    "tracks": [
        "Snare: |>-X-X->-|X-X->-X-|X->-X-X-|>->->>>>|",
        "Bass:  |X-------|--------|X-------|--------|"
    ]}`)

    pattern, err := p.ParsePattern(jsonBlob)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    x, err := a.Synth()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    s := &Sequencer{
        Timer:   NewTimer(),
        Pattern: pattern,
        Synth:   x,
    }

    return s, nil
}

func (s *Sequencer) Start() {
    go func() {
        pulseCount := 0

        for {
            select {
            case <-s.Timer.Pulses:
                pulse := pulseCount % s.Pattern.Divisions
                RenderPattern(s.Pattern, pulse)
                pulseCount++
                for _, t := range s.Pattern.Tracks {
                    hit := t.Steps[pulse]
                    if hit > 0 {
                        s.Synth.Play(t.Instrument, float32(hit))
                    }
                }
            }
        }
    }()

    s.Timer.SetTempo(s.Pattern.BPM)
    s.Timer.SetDivisionsPerBeat(s.Pattern.DivisionsPerBeat)
    go s.Timer.Start()
}
