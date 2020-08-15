package sequencer

import (
	"fmt"
	tm "github.com/buger/goterm"
	p "github.com/fjania/froland/pkg/pattern"
	"strings"
)

func RenderTimeline(track p.Track, pulse int) string {
	// The downbeat we're currently on
	beat := pulse / track.DivisionsPerBeat
	beatCount := 0

	var sb strings.Builder
	for i, s := range track.Steps {
		if i%track.DivisionsPerBeat == 0 {
			if beat == beatCount {
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

func RenderPattern(pattern *p.Pattern, pulse int) {
	tm.Clear()
	tm.MoveCursor(1, 1)

	tm.Println(tm.Bold("Song:"), pattern.Title)
	tm.Println(tm.Bold("BPM :"), pattern.BPM)
	tm.Println()

	maxLabel := 0
	for _, t := range pattern.Tracks {
		if len(t.Instrument) > maxLabel {
			maxLabel = len(t.Instrument)
		}
	}
	labelFormatter := fmt.Sprintf("%%-%ds", maxLabel+2)

	for _, t := range pattern.Tracks {
		tm.Println(fmt.Sprintf(labelFormatter, t.Instrument), RenderTimeline(t, pulse))
	}

	tm.Flush()
}
