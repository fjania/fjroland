package main

import (
    "flag"
    "fmt"
    //"time"
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
        out, err := fmt.Scanln()
        fmt.Println(out, err)
        //time.Sleep(time.Second * 3)
        s.LoadPattern(patternFile)
    }
}
