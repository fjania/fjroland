package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "time"
    "github.com/fjania/froland/pkg/sequencer"
    "github.com/fsnotify/fsnotify"
)

func main() {
    var patternFile string
    var output string
    var samplePack string

    flag.StringVar(
        &patternFile,
        "pattern",
        "turn-down-for-what.json",
        "--pattern=pattern-file.json",
    )

    flag.StringVar(
        &output,
        "output",
        "samples",
        "--output=samples | midi",
    )

    flag.StringVar(
        &samplePack,
        "samplepack",
        "acoustic",
        "--samplepack=pack-name",
    )

    flag.Parse()

    s, _ := sequencer.NewSequencer(patternFile, output, samplePack)
    s.Start()

    if (output != "samples" &&  output != "midi") {
        Usage()
        os.Exit(1)
    }

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Printf("Error starting watcher %s", s.PatternFilePath)
    }
    defer watcher.Close()

    go func() {
        for {
            select {
            case event := <-watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    s.LoadPattern(s.PatternFilePath)
                }

            case err := <-watcher.Errors:
                log.Printf("Error in watching %s: %v", s.PatternFilePath, err)
            }
        }
    }()

    if err := watcher.Add(s.PatternFilePath); err != nil {
        log.Printf("Error watching %s", s.PatternFilePath)
    }

    for {
        time.Sleep(time.Second * 3)
    }
}

var Usage = func() {
    fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
    flag.PrintDefaults()
}
