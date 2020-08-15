package sequencer

import (
    "fmt"
    "strings"
    p "github.com/fjania/froland/pkg/pattern"
)


func ChunkIntoBeats(track p.Track) string {
    var sb strings.Builder
    for i, s := range track.Steps {
        if (i % track.DivisionsPerBeat == 0){
            sb.WriteString(p.BEATMARKER)
        }
        sb.WriteString(p.LevelsAsIndicators[s])
    }
    sb.WriteString(p.BEATMARKER)

    return sb.String()
}

func RenderPattern(pattern *p.Pattern){
    fmt.Println("Song:",pattern.Title)
    fmt.Println("BPM :",pattern.BPM)

    for _, t := range pattern.Tracks {
        fmt.Println(t.Instrument, ChunkIntoBeats(t))
    }
}
