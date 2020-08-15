package pattern

import (
    "encoding/json"
    "fmt"
    "log"
    "strings"
)

type Pattern struct {
    Title   string
    BPM    int
    Tracks []Track
    Beats            int
    Divisions        int
    DivisionsPerBeat int
}

type Track struct {
    Instrument       string
    Steps            []int
    Beats            int
    Divisions        int
    DivisionsPerBeat int
}

const (
    SILENT = 0
    SILENTMARKER = "-"

    GHOST  = 1
    GHOSTMARKER  = "o"

    STRIKE = 2
    STRIKEMARKER = "X"

    ACCENT = 3
    ACCENTMARKER = ">"

    BEATMARKER = "|"
)

var IndicatorsAsLevels = map[string]int{
    SILENTMARKER: SILENT,
    GHOSTMARKER: GHOST,
    STRIKEMARKER: STRIKE,
    ACCENTMARKER: ACCENT,
}

var LevelsAsIndicators = map[int]string{
    SILENT: SILENTMARKER,
    GHOST:  GHOSTMARKER,
    STRIKE: STRIKEMARKER,
    ACCENT: ACCENTMARKER,
}

/*
* We'll assume the format of a pattern is a json file** of
* this format:
*
* Name: SongName
* BPM: bpm
* InstrumentName: |----|----|----|----|
* InstrumentName: |----|----|----|----|
* ...
* InstrumentName: |----|----|----|----|
*
*  **Why a json file? I don't want to write a parser for
*  a custom format, and json is portable. I'm also trying
*  to keep the format human readable, so one can modify the
*  pattern in a simple text editor.
*/

func (t *Track) UnmarshalJSON(b []byte) error {
    var s string
    if err := json.Unmarshal(b, &s); err != nil {
        return err
    }
    err := t.ParseTrack(s)
    return err
}

func ParsePattern(patternJson []byte) (*Pattern, error) {
    var p Pattern
    err := json.Unmarshal(patternJson, &p)
    if err != nil {
        log.Fatal("Error parsing JSON file. ", err)
        return nil, err
    }

    // Todo - fail on patterns with no tracks
    if len(p.Tracks) == 0 {
        return nil, fmt.Errorf(
            "Malformed track entry. There must be at least one track in the file.",
        )
    }

    p.Beats = p.Tracks[0].Beats
    p.Divisions = p.Tracks[0].Divisions
    p.DivisionsPerBeat = p.Tracks[0].DivisionsPerBeat

    return &p, nil
}

/*
* We'll assume the track as the following form:
*
* InstrumentName: |----|----|----|----|
*
* InstrumentName can be anything, but our sequencer might not know about it
* The patter can be any number of divisons, but parsing will fail if beats have
* not been divided uniformly.
*
 */
func (track *Track) ParseTrack(trackAsString string) error {

    parts := strings.Split(trackAsString, ":")

    // If we separated on colon and there are more than two parts, there's
    // a problem.
    if len(parts) != 2 {
        return fmt.Errorf(
            "Malformed track entry. Instrument and track parts could not be identified: %q",
            trackAsString,
        )
    }

    instrument := strings.TrimSpace(parts[0])
    trackSpec := strings.Trim(parts[1], " ")

    // The track layout must start and end with a | character
    start := string(trackSpec[0])
    end := string(trackSpec[len(trackSpec)-1])
    if start != BEATMARKER && end != BEATMARKER {
        return fmt.Errorf(
            "Malformed track entry. No leading/trailing '%s' characters: %q",
            BEATMARKER,
            trackAsString,
        )
    }

    divisions := 0
    divisionsPerBeat := 0
    divisionsInThisBeat := 0
    beats := 0

    for _, e := range trackSpec {
        char := string(e)
        if char == BEATMARKER {
            // If we reading this first beat, we'll be defining the beats per
            // division we expect from now on.
            if beats == 1 {
                divisionsPerBeat = divisionsInThisBeat

                // if it's not the first beat, we need to check that we're always
                // using that number of divisions
            } else {
                if divisionsInThisBeat != divisionsPerBeat {
                    return fmt.Errorf(
                        "Malformed track entry. Non-uniform beat divisions: %q",
                        trackAsString,
                    )
                }
            }

            divisionsInThisBeat = 0
            beats++

        } else {
            divisions++
            divisionsInThisBeat++

            if char != SILENTMARKER && IndicatorsAsLevels[char] == 0 {
                return fmt.Errorf(
                    "Malformed track entry. Invalid indicator %s: %s",
                    char,
                    trackAsString,
                )
            }
        }
    }

    // Remove the extra beat we counted
    beats--

    // Now that we knoe the input is valid, we make another pass to capture the
    // data we'll need. Build the array for steps, reset the divisions count and
    // start capturing
    steps := make([]int, divisions)
    division := 0
    for _, e := range trackSpec {
        char := string(e)
        if char != BEATMARKER {
            steps[division] = IndicatorsAsLevels[char]
            division++
        }
    }


    track.Instrument = instrument
    track.Steps = steps
    track.Beats = beats
    track.Divisions = divisions
    track.DivisionsPerBeat = divisionsPerBeat
    return nil
}
