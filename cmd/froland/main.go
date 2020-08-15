package main

import (
    "flag"
    "log"
    "time"
    "github.com/fjania/froland/pkg/sequencer"
    "github.com/fsnotify/fsnotify"
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
