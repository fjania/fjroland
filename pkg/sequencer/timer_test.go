package sequencer

import (
    "testing"
    "time"
)

func TestNanosecondsPerPulse(t *testing.T) {
    timer := NewTimer()

    timer.SetStepInterval(120, 8)

    output := timer.NanosecondsPerPulse()
    expect := time.Duration((1000000000 * 60) / (120 * 8))

    if output != expect {
        t.Errorf("NanosecondsPerPulse Epected: %d got %d", expect, output)
    }

    timer.SetStepInterval(80, 3)

    output = timer.NanosecondsPerPulse()
    expect = time.Duration((1000000000 * 60) / (80 * 3))

    if output != expect {
        t.Errorf("NanosecondsPerPulse Epected: %d got %d", expect, output)
    }
}
