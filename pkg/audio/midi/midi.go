package audio

import (
    "log"
    "time"
    "github.com/rakyll/portmidi"
)

type Instrument struct {
    name     string
    MidiNote int
}

type Midi struct {
    Instruments map[string]bool
    MidiNotes map[string]int
    OutputStream *portmidi.Stream
}

var instruments = [15]Instrument{
    Instrument{"Bass", 35},
    Instrument{"Cross-stick", 37},
    Instrument{"Snare", 38},
    Instrument{"Clap", 39},
    Instrument{"High Tom", 50},
    Instrument{"Mid Tom", 48},
    Instrument{"Low Tom", 41},
    Instrument{"Closed Hi-Hat", 42},
    Instrument{"Pedal Hi-Hat", 44},
    Instrument{"Open Hi-Hat", 46},
    Instrument{"Crash Cymbal", 49},
    Instrument{"Ride Cymbal", 51},
    Instrument{"Ride Bell", 53},
    Instrument{"Tambourine", 54},
    Instrument{"Cowbell", 56},
}

func (m *Midi) Play(instrument string, level float32) {
    // input levels are 1, 2, 3 well map them linearly
    // as 40, 80 12
    midiNote := m.MidiNotes[instrument]
    m.OutputStream.WriteShort(0x9A, int64(midiNote), int64(level*40))
}

func (m *Midi) ListInstruments() map[string]bool {
    return m.Instruments
}

func (m *Midi) HasInstrument(instrument string) bool {
    return m.Instruments[instrument]
}

func NewMidi() (*Midi, error){
    m := &Midi{
        Instruments: make(map[string]bool),
        MidiNotes: make(map[string]int),
    }

    for _, ins := range instruments {
        m.Instruments[ins.name] = true
        m.MidiNotes[ins.name] = ins.MidiNote

    }

    portmidi.Initialize()

    log.Println("Midi Device Count>", portmidi.CountDevices())
    log.Printf("Default Device Info> %+v\n", portmidi.Info(portmidi.DefaultOutputDeviceID()))
    log.Println("Default Input Device>", portmidi.DefaultInputDeviceID())
    log.Println("Default Output Device>", portmidi.DefaultOutputDeviceID())

    var streamErr error
    m.OutputStream, streamErr = portmidi.NewOutputStream(
        portmidi.DefaultOutputDeviceID(), 1024, 0,
    )
    if streamErr != nil {
        log.Fatal(streamErr)
    }

    log.Printf("Default Device Info> %+v\n", portmidi.Info(portmidi.DefaultOutputDeviceID()))

    return m, nil
}

func (m *Midi) demo(){
    // Send "note on" events to play C major chord.
    for i := 0; i< 120; i++ {
        m.Play("Closed Hi-Hat", 2)
        time.Sleep(time.Second/8)
        m.Play("Clap", 2)
        time.Sleep(time.Second/8)
        m.Play("Snare", 2)
        time.Sleep(time.Second/8)
        m.Play("Bass", 2)
        time.Sleep(time.Second/8)
    }
}
