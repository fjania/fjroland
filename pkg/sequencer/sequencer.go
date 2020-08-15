package sequencer

import (
    p "github.com/fjania/froland/pkg/pattern"
)

// c := time.Tick(time.Minute/time.Duration(divisionsPerMinute))
type Sequencer struct {
    Timer   *Timer
    Pattern *p.Pattern
}

func NewSequencer() (*Sequencer, error) {
    var jsonBlob = []byte(`
    {"title": "Turn Down for What",
    "bpm": 100,
    "tracks": [
        "Snare: |>-X-X->-|X-X->-X-|X->-X-X-|>->->>>>|",
        "Bass:  |X-------|--------|X-------|--------|"
    ]}`)

    pattern, _ := p.ParsePattern(jsonBlob)

    s := &Sequencer{
        Timer:   NewTimer(),
        Pattern: pattern,
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
            }
        }
    }()

    s.Timer.SetTempo(s.Pattern.BPM)
    s.Timer.SetDivisionsPerBeat(s.Pattern.DivisionsPerBeat)
    go s.Timer.Start()
}
