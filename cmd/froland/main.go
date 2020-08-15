package main

import (
    "flag"
    "time"
    "github.com/fjania/froland/pkg/sequencer"
)

func main() {
    var patternFile string
    var kitName string

    flag.StringVar(
        &patternFile,
        "pattern",
        "turn-down-for-what.json",
        "-pattern=pattern-file.json",
    )

    flag.StringVar(
        &kitName,
        "kit",
        "acoustic",
        "-kit=kit-name",
    )

    flag.Parse()

    s, _ := sequencer.NewSequencer(patternFile, kitName)
    s.Start()

    for {
        time.Sleep(time.Second)
    }
}
