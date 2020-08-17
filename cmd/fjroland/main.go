package main

import (
    "log"
    "os"
    "time"
    "github.com/fjania/fjroland/pkg/sequencer"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	var opts struct {
		Positional struct {
			PatternFile string `required:"1"`
		} `positional-args:"yes" required:"yes"`
		MidiDevices []string `short:"m" long:"midi" description:"A midi device name to output to"`
		SamplePacks []string `short:"s" long:"samples" description:"A directory of samples to use for waveform playback"`
	}

	parser := flags.NewParser(&opts, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	log.Printf("Pattern> %v", opts.Positional.PatternFile)
	log.Printf("Devices> %v", opts.MidiDevices)
	log.Printf("Samples> %v", opts.SamplePacks)
	log.Printf("Args> %v", args)

    s := sequencer.NewSequencer()
	serr := s.LoadPattern(opts.Positional.PatternFile)
    if serr != nil {
        log.Fatalf("Could not load pattern file: '%s'", opts.Positional.PatternFile)
    }
    s.Start()

	for _, e := range opts.MidiDevices{
		s.ConfigureMidiOutput(e)
	}

	for _, e := range opts.SamplePacks{
		s.ConfigureSamplesOutput(e)
	}

    for {
        time.Sleep(time.Second * 3)
    }
}
