package sequencer

import (
	p "github.com/fjania/froland/pkg/pattern"
	"testing"
)

func TestRenderPattern(t *testing.T) {
	var jsonBlob = []byte(`
    {"title": "Turn Down for What",
    "bpm": 100,
    "tracks": [
        "Snare: |>-X-X->-|X-X->-X-|X->-X-X-|>->->>>>|",
        "Bass:  |X-------|--------|X-------|--------|"
    ]}`)

	pattern, _ := p.ParsePattern(jsonBlob)
	RenderPattern(pattern, 0)
	RenderPattern(pattern, 15)
	RenderPattern(pattern, 20)
	RenderPattern(pattern, 31)
}
