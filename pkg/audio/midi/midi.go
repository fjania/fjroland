package audio

import (
	"github.com/rakyll/portmidi"
	"log"
)

type Instrument struct {
	name     string
	MidiNote int
}

type Midi struct {
	DeviceName   string
	Instruments  map[string]bool
	MidiNotes    map[string]int
	OutputStream *portmidi.Stream
}

const (
	MidiPercussionChannel = 0x9A
)

// A subset of the percussive instruments available on midi channel 10.
// There's no reason we couldn't support all of them, but this is a fine
// place to start for this example.
var instruments = [15]Instrument{
	Instrument{"bass", 35},
	Instrument{"cross-stick", 37},
	Instrument{"snare", 38},
	Instrument{"clap", 39},
	Instrument{"high tom", 50},
	Instrument{"mid tom", 48},
	Instrument{"low tom", 41},
	Instrument{"closed hi-hat", 42},
	Instrument{"pedal hi-hat", 44},
	Instrument{"open hi-hat", 46},
	Instrument{"crash cymbal", 49},
	Instrument{"ride cymbal", 51},
	Instrument{"ride bell", 53},
	Instrument{"tambourine", 54},
	Instrument{"cowbell", 56},
}

func NewMidi(deviceName string) (*Midi, error) {
	m := &Midi{
		DeviceName:  deviceName,
		Instruments: make(map[string]bool),
		MidiNotes:   make(map[string]int),
	}

	for _, ins := range instruments {
		m.Instruments[ins.name] = true
		m.MidiNotes[ins.name] = ins.MidiNote
	}

	portmidi.Initialize()

	// There isn't a map of devices we can use, so we'll loop through
	// all the devices trying to find the id of the one we were told
	// to use.
	var deviceID portmidi.DeviceID

	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		if info.IsOutputAvailable && info.Name == deviceName {
			deviceID = portmidi.DeviceID(i)
		}
	}

	var streamErr error
	m.OutputStream, streamErr = portmidi.NewOutputStream(deviceID, 1024, 0)
	if streamErr != nil {
		log.Printf("Failed to open Midi Device: '%v'", deviceName)
		return nil, streamErr
	}

	return m, nil
}

// AudioOutput Interface Methods
func (m *Midi) Name() string {
	return m.DeviceName
}
func (m *Midi) Play(instrument string, level float32) {
	// Input levels are 1, 2, 3 well map them linearly as 40, 80 12
	midiNote := m.MidiNotes[instrument]
	m.OutputStream.WriteShort(MidiPercussionChannel, int64(midiNote), int64(level*40))
}

func (m *Midi) ListInstruments() map[string]bool {
	return m.Instruments
}

func (m *Midi) HasInstrument(instrument string) bool {
	return m.Instruments[instrument]
}
