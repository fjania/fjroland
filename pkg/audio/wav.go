package audio

import (
    "fmt"
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

type Kit struct {
    Name string
    Samples map[string]*Sample
}

func LoadKit(kitName string) (*Kit, error) {
    k := &Kit{
        Name: kitName,
        Samples: make(map[string]*Sample),
    }

    dir := fmt.Sprintf("../../assets/kits/%s/", kitName)
    files, err := ioutil.ReadDir(dir);
    if err != nil {
        return nil, err
    }

    for _, f := range files {
        filepath := fmt.Sprintf("../../assets/kits/%s/%s", kitName, f.Name())
        if strings.HasSuffix(f.Name(), ".wav") {
            fmt.Println("Loading",filepath)
            instrument := strings.TrimRight(f.Name(), ".wav")

            sample, err := LoadSample(filepath)
            if err != nil {
                fmt.Printf("Could not load sample: %s\n", filepath)
                return nil, err
            }
            k.Samples[instrument] = sample
        }
    }
    return k, nil
}

func (k *Kit) Play(instrument string, level float32) {
    k.Samples[instrument].Play(level);
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
        fmt.Printf("Could not open file: %s\n", filepath)
        return nil, err
    }

    _, err = soundFile.ReadItems(s.Buffer)
    if err != nil {
        fmt.Printf("Error reading data from file: %s\n", filepath)
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

func Synth() (*Kit, error) {
    kit, err := LoadKit("acoustic")
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
        kit.ProcessAudio,
    )

    if err != nil {
        return nil, err
    }

    stream.Start()
    return kit, nil
}

func (k *Kit) ProcessAudio(out []float32) {
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
