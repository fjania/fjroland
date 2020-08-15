package audio

type Output interface{
    Play(instrument string, level float32)
    ListInstruments() map[string]bool
    HasInstrument(instrument string) bool
}
