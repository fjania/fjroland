package audio

import (
    "testing"
)

func TestMidi(t *testing.T) {
    m, err := NewMidi()
    if err != nil {
        t.Errorf("Error opening Midi: %s", err)
    }
    m.demo()
}
