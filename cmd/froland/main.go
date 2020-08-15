package main

import (
    "flag"
    "log"
    "os"
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

    sep := string(os.PathSeparator)

    patternFilePath := ".." + sep + ".." +
        sep + "assets" + sep +
        "patterns" + sep + patternFile

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Printf("Error starting watcher %s", patternFile)
    }
    defer watcher.Close()

    go func() {
        for {
            select {
            case event := <-watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    s.LoadPattern(patternFile)
                }

            case err := <-watcher.Errors:
                log.Printf("Error in watching %s: %v", patternFile, err)
            }
        }
    }()

    if err := watcher.Add(patternFilePath); err != nil {
        log.Printf("Error watching %s", patternFile)
    }

    for {
        time.Sleep(time.Second * 3)
    }
}
