package pattern

import (
	"reflect"
	"testing"
)

func TestParseInvalidTrack_BadColon(t *testing.T) {
	var track *Track
	trackSpec := "Snare: : |----|X---|----|X---|"
	err := track.ParseTrack(trackSpec)
	if err == nil {
		t.Errorf("ParseTrack did not fail on multiple colons: '%s'", trackSpec)
	}
}

func TestParseInvalidTrack_NoInstrumentName(t *testing.T) {
	var track *Track
	trackSpec := "|----|X---|----|X---|"
	err := track.ParseTrack(trackSpec)
	if err == nil {
		t.Errorf("ParseTrack did not fail without an instrument name: '%s'", trackSpec)
	}
}

func TestParseInvalidTrack_BeatMarkers(t *testing.T) {
	var track *Track
	trackSpec := "Snare:   ----|X---|----|X---"
	err := track.ParseTrack(trackSpec)
	if err == nil {
		t.Errorf("ParseTrack did not fail without leading/trailing beat markers: '%s'", trackSpec)
	}
}

func TestParseInvalidTrack_BeatUniformity(t *testing.T) {
	var track *Track
	trackSpec := "Snare:   |----------|X---|----|X---|"
	err := track.ParseTrack(trackSpec)
	if err == nil {
		t.Errorf("ParseTrack did not fail with non-uniform beats: '%q'", trackSpec)
	}
}

func TestParseInvalidTrack_InvalidIndicators(t *testing.T) {
	var track *Track
	trackSpec := "Snare:   |----|q---|----|X---|"
	err := track.ParseTrack(trackSpec)
	if err == nil {
		t.Errorf("ParseTrack did not fail on invalid indicators: '%q'", trackSpec)
	}
}

func TestParseValidTrack(t *testing.T) {
	ValidateTrackSpec(
		t,
		"Snare:      |>-o-|X---|----|X---|",
		"Snare",
		4,
		16,
		4,
		[]int{3, 0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0},
	)
}

func ValidateTrackSpec(t *testing.T,
	trackSpec string,
	eInstrument string,
	eBeats int,
	eDivisions int,
	eDivisionsPerBeat int,
	eSteps []int,
) {

	var track *Track = new(Track)
	err := track.ParseTrack(trackSpec)

	if err != nil {
		t.Errorf("Valid track did not parse properly: %s", trackSpec)
	}

	if track.Instrument != eInstrument {
		t.Errorf(
			"ParseTrack got instrument %s but expected %s",
			track.Instrument,
			eInstrument,
		)
	}
	if track.Beats != eBeats {
		t.Errorf(
			"ParseTrack counted %d beats but expected %d",
			track.Beats,
			eBeats,
		)
	}
	if track.Divisions != eDivisions {
		t.Errorf(
			"ParseTrack counted %d divisions but expected %d",
			track.Divisions,
			eDivisions,
		)
	}
	if track.DivisionsPerBeat != eDivisionsPerBeat {
		t.Errorf(
			"ParseTrack counted %d divisionsPerBeat but expected %d",
			track.DivisionsPerBeat,
			eDivisionsPerBeat,
		)
	}
	if !reflect.DeepEqual(track.Steps, eSteps) {
		t.Errorf("ParseTrack got %d but expected %d", track.Steps, eSteps)
	}
}

func TestParsePattern(t *testing.T) {
	var p = []byte(`
    {"title": "Turn Down for What",
    "bpm": 100,
    "tracks": [
        "Snare: |>-X-X->-|X-X->-X-|X->-X-X-|>->->>>>|",
        "Bass:  |X-------|--------|X-------|--------|"
    ]}`)

	pattern, err := ParsePattern(p)

	if err != nil {
		t.Errorf("Valid pattern did not parse properly: %s", p)
	}

	eTitle := "Turn Down for What"
	if pattern.Title != eTitle {
		t.Errorf("Expected title of '%s' but got '%s'", eTitle, pattern.Title)
	}

	eBPM := 100
	if pattern.BPM != eBPM {
		t.Errorf("Expected BPM of %d but got %d", eBPM, pattern.BPM)
	}

	ValidateTrackSpec(
		t,
		"Snare: |>-X-X->-|X-X->-X-|X->-X-X-|>->->>>>|",
		"Snare",
		4,
		32,
		8,
		[]int{3, 0, 2, 0, 2, 0, 3, 0, 2, 0, 2, 0, 3, 0, 2, 0, 2, 0, 3, 0, 2, 0, 2, 0, 3, 0, 3, 0, 3, 3, 3, 3},
	)

	ValidateTrackSpec(
		t,
		"Bass:  |X-------|--------|X-------|--------|",
		"Bass",
		4,
		32,
		8,
		[]int{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	)

}
