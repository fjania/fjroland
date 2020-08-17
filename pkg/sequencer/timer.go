package sequencer

import (
	"time"
)

type Timer struct {
	Pulses           chan int
	Tempo            int
	DivisionsPerBeat int
}

func NewTimer() *Timer {
	t := &Timer{
		Pulses:           make(chan int),
		Tempo:            60,
		DivisionsPerBeat: 4,
	}

	return t
}

func (t *Timer) SetStepInterval(tempo int, divisionsPerBeat int) {
	t.Tempo = tempo
	t.DivisionsPerBeat = divisionsPerBeat
}

func (t *Timer) Start() {
	for {
		select {
		default:
			interval := t.NanosecondsPerPulse()
			time.Sleep(interval)
			t.Pulses <- 1
		}
	}
}

func (t *Timer) NanosecondsPerPulse() time.Duration {
	return time.Duration(
		time.Minute / time.Duration(t.Tempo*t.DivisionsPerBeat),
	)
}
