package audio

import (
    "fmt"
    "strings"
    "time"
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

func (k *Kit) Play(instrument string, volume float32) {
    k.Samples[instrument].Play(volume);
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

func (s *Sample) Play(volume float32) {
    s.Volume = volume
    s.Playhead = 0
}

func Music() {
    kit, err := LoadKit("acoustic")
    fmt.Printf("%+v\n", kit)
    fmt.Printf("%+v\n", err)

    err = portaudio.Initialize()
    if err != nil {
        fmt.Println("Error>",err)
    }

    stream, err := portaudio.OpenDefaultStream(
        0,
        2,
        44100.0,
        portaudio.FramesPerBufferUnspecified,
        kit.ProcessAudio,
    )

    if err != nil {
        fmt.Println("Error>",err)
    }

    stream.Start()

    time.Sleep(time.Second*2)

    strike := float32(0.5)
    accent := float32(2)
    for i:= 0; i<4; i++ {
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", strike);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", strike);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", strike);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", strike);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", strike);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", strike);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", strike);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", strike);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    kit.Play("Snare", accent);
    time.Sleep( time.Duration(time.Minute / time.Duration(100*8)))
    }

    time.Sleep(time.Second*2)


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
