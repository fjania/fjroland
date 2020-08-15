package main

import (
    "flag"
    "fmt"
    "time"
    p "github.com/fjania/froland/pkg/pattern"
    s "github.com/fjania/froland/pkg/sequencer"
)

func main() {
    var patternPath string
    var kitPath string

    flag.StringVar(
        &patternPath,
        "pattern",
        "patterns/pattern_1.splice",
        "-pattern=path/to/pattern.splice",
    )

    flag.StringVar(
        &kitPath,
        "kit",
        "kits",
        "-kit=path/to/kits",
    )

    flag.Parse()
    fmt.Println(patternPath, kitPath)

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
        s.RenderPattern(pattern, i%32)
        i++
        time.Sleep(time.Second/8)
    }
}
