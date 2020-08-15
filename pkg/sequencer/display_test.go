package sequencer

import (
    "testing"
    p "github.com/fjania/froland/pkg/pattern"
)

func TestRenderPattern(t *testing.T){
    var jsonBlob = []byte(`
    {"title": "Turn Down for What",
    "bpm": 100,
    "tracks": [
        "Snare: |>-X-X->-|X-X->-X-|X->-X-X-|>->->>>>|",
        "Bass:  |X-------|--------|X-------|--------|"
    ]}`)

    pattern, _ := p.ParsePattern(jsonBlob)
    RenderPattern(pattern)
}
