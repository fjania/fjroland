package sequencer

import (
	"fmt"
	tm "github.com/buger/goterm"
	p "github.com/fjania/fjroland/pkg/pattern"
	"strings"
)

func RenderTimeline(track p.Track, pulse int, available bool) string {
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
	if !available {
		sb.WriteString(tm.Color(" // not available on any output", tm.YELLOW))
	}

	return sb.String()
}

func RenderPattern(s *Sequencer, pulse int) {
	pattern := s.Pattern
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
		isAvailable := s.IsInstrumentAvailable(t.Instrument)
		trackLabel := fmt.Sprintf(labelFormatter, t.Instrument)
		if !isAvailable {
			trackLabel = tm.Color(trackLabel, tm.YELLOW)
		} else {
			trackLabel = tm.Bold(trackLabel)
		}
		tm.Println(trackLabel, RenderTimeline(t, pulse, isAvailable))
	}

	tm.Println()
	if len(s.AudioOutputs) < 1 {
		tm.Println(tm.Color(tm.Bold("(No audio outputs specified.)"), tm.RED))
	} else {
		tm.Println(tm.Bold("Playing on:"), s.AudioDeviceNameList())
	}

	tm.Flush()
}
