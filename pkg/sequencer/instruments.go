package sequencer

type Instrument struct {
    name string
    midiNote int
}

var instruments = [15]Instrument {
    Instrument{"Bass", 35},
    Instrument{"Cross-stick", 37},
    Instrument{"Snare", 38},
    Instrument{"Clap", 39},
    Instrument{"High Tom", 50},
    Instrument{"Mid Tom", 48},
    Instrument{"Low Tom", 41},
    Instrument{"Closed Hi-hat", 42},
    Instrument{"Pedal Hi-hat", 44},
    Instrument{"Open Hi-hat", 46},
    Instrument{"Crash Cymbal", 49},
    Instrument{"Ride Cymbal", 51},
    Instrument{"Ride Bell", 53},
    Instrument{"Tambourine", 54},
    Instrument{"Cowbell", 56},
}
