package main

import (
    "flag"
    "fmt"
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
}
