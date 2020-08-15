package main

import (
    "fmt"
    "strings"
    "time"
    p "github.com/fjania/froland/pkg/pattern"
    tm "github.com/buger/goterm"
)


func RenderTimeline(track p.Track, pulse int) string {
    // The downbeat we're currently on
    beat := pulse/track.DivisionsPerBeat
    beatCount := 0

    var sb strings.Builder
    for i, s := range track.Steps {
        if (i % track.DivisionsPerBeat == 0){
            if (beat == beatCount) {
                sb.WriteString(tm.Bold(p.BEATMARKER))
            } else {
                sb.WriteString(p.BEATMARKER)
            }
            beatCount++
        }
        if i == pulse {
            sb.WriteString(tm.Background(p.LevelsAsIndicators[s], tm.RED))
        } else {
            sb.WriteString(p.LevelsAsIndicators[s])
        }
    }
    sb.WriteString(p.BEATMARKER)

    return sb.String()
}

func RenderPattern(pattern *p.Pattern, pulse int){
    tm.Clear()
    tm.MoveCursor(1,1)

    tm.Println(tm.Bold("Song:"),pattern.Title)
    tm.Println(tm.Bold("BPM :"),pattern.BPM)
    tm.Println()

    labelFormatter := fmt.Sprintf("%%-%ds", 16)

    for _, t := range pattern.Tracks {
        tm.Println(fmt.Sprintf(labelFormatter, t.Instrument), RenderTimeline(t, pulse))
    }

    tm.Flush()
}

func main() {
    var jsonBlob = []byte(`
    {"title": "Turn Down for What",
    "bpm": 100,
    "tracks": [
        "Snare: |>-X-X->-|X-X->-X-|X->-X-X-|>->->>>>|",
        "Bass:  |X-------|--------|X-------|--------|"
    ]}`)

    pattern, _ := p.ParsePattern(jsonBlob)
    i := 0
    for {
        RenderPattern(pattern, i%32)
        i++
        time.Sleep(time.Second/8)
    }
}
