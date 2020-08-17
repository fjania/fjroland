package audio

type AudioOutput interface {
	Name() string
	Play(instrument string, level float32)
	ListInstruments() map[string]bool
	HasInstrument(instrument string) bool
}
