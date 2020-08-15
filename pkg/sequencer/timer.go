package sequencer

import (
    "time"
)

// Timer is a struct that defines the basic synchronization
// behavior of the step sequencer
type Timer struct {
    Pulses           chan int
    Done             chan bool
    Tempo            int
    DivisionsPerBeat int
}

// NewTimer creates and returns a new pointer to a Timer.
func NewTimer() *Timer {
    t := &Timer{
        Pulses:           make(chan int),
        Done:             make(chan bool),
        Tempo:            60,
        DivisionsPerBeat: 4,
    }

    return t
}

func (t *Timer) SetTempo(tempo int) {
    t.Tempo = tempo
}

func (t *Timer) SetDivisionsPerBeat(dpb int) {
    t.DivisionsPerBeat = dpb
}

func (t *Timer) Start() {
    for {
        select {
        case <-t.Done:
            break
        default:
            interval := t.NanosecondsPerPulse()
            time.Sleep(interval)
            t.Pulses <- 1
        }
    }
}

// MicrosecondsPerPulse returns how many microseconds
// A client would need to wait for a "Pulse" to take place.
func (t *Timer) NanosecondsPerPulse() time.Duration {
    return time.Duration(time.Minute / time.Duration(t.Tempo*t.DivisionsPerBeat))
}
