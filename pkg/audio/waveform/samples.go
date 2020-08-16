package audio

import (
    "log"
    "os"
    "strings"
    "io/ioutil"
    "github.com/gordonklaus/portaudio"
    "github.com/mkb218/gosndfile/sndfile"
)

type Sample struct {
    Buffer []float32
    Info sndfile.Info
    Playhead int
    Volume float32
}

type SamplePack struct {
    Samples map[string]*Sample
    Instruments map[string]bool
}

func LoadSamplePack(samplePackPath string) (*SamplePack, error) {
    k := &SamplePack{
        Samples: make(map[string]*Sample),
    }

    files, err := ioutil.ReadDir(samplePackPath);
    if err != nil {
        return nil, err
    }

    for _, f := range files {
        sampleFilePath := samplePackPath + string(os.PathSeparator) + f.Name()
        if strings.HasSuffix(f.Name(), ".wav") {
            log.Println("Loading",sampleFilePath)
            instrument := strings.TrimRight(f.Name(), ".wav")

            sample, err := LoadSample(sampleFilePath)
            if err != nil {
                log.Printf("Could not load sample: %s\n", sampleFilePath)
                return nil, err
            }
            k.Samples[instrument] = sample
        }
    }

    instruments := make(map[string]bool)
    for i, _ := range k.Samples {
        instruments[i] = true
    }
    k.Instruments = instruments

    return k, nil
}

// Implement the audio Output interface
func (k *SamplePack) Play(instrument string, level float32) {
    k.Samples[instrument].Play(level);
}

func (k *SamplePack) ListInstruments() map[string]bool {
    return k.Instruments
}

func (k *SamplePack) HasInstrument(instrument string) bool {
    return k.Instruments[instrument]
}

func LoadSample(filepath string) (*Sample, error) {
    var info sndfile.Info
    soundFile, err := sndfile.Open(filepath, sndfile.Read, &info)

    s := &Sample {
        Buffer: make([]float32, info.Samplerate * info.Channels),
        Info: info,
        Playhead: 0,
        Volume: 1,
    }

    if err != nil {
        log.Printf("Could not open file: %s\n", filepath)
        return nil, err
    }

    _, err = soundFile.ReadItems(s.Buffer)
    if err != nil {
        log.Printf("Error reading data from file: %s\n", filepath)
        return nil, err
    }

    s.Playhead = len(s.Buffer)

    defer soundFile.Close()

    return s, nil
}

func (s *Sample) Play(level float32) {
    s.Volume = level*level
    s.Playhead = 0
}

func NewSamplePack(samplePackPath string) (*SamplePack, error) {
    samplePack, err := LoadSamplePack(samplePackPath)
    if err != nil {
        return nil, err
    }

    err = portaudio.Initialize()
    if err != nil {
        return nil, err
    }

    stream, err := portaudio.OpenDefaultStream(
        0,
        2,
        44100.0,
        portaudio.FramesPerBufferUnspecified,
        samplePack.ProcessAudio,
    )

    if err != nil {
        return nil, err
    }

    stream.Start()
    return samplePack, nil
}

func (k *SamplePack) ProcessAudio(out []float32) {
    for i := range out {
        var data float32
        for _, s := range k.Samples{
            if s.Playhead < len(s.Buffer) {
                data += s.Volume * s.Buffer[s.Playhead]
                s.Playhead++
            }
        }

        if data > 1.0 {
            data = 1.0
        }

        out[i] = data
    }
}
