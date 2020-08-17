# fjroland

Pronounced _/fērōland, fyrōland/_ - A multifunction drum sequencer, written in Go.

## Features
- Command line based rendering of sequencer timeline in real time
- Pattern files are easily editable in plain text
- Sequencer watches for pattern file changes to allow interactive editing of patterns
- Sequencer can output audio via .wav file sample packs
- Sequencer can output midi to play pattern a midi device
- Example patterns and sample packs are included to play with. (Sorry, you have to buy your own midi device)

<p align="center">
  <img src="https://user-images.githubusercontent.com/184340/90358502-36bcaa80-e024-11ea-8dd2-c7d7948871e7.gif" width=649 height=536>
</p>


## Building

fjrolad was developed and tested with go 1.15 on OS X 10.15.

Before building the binary make sure you have the devlopment libraries for `portmidi`, `portaudio` and `libsndfile` installed. On OS X with Homebrew configure you can do that with this command:

```
brew install portmidi portaudio libsndfile
```

After you have the development libraries installed, `make build` in the root directory will build the binary and store it as `bin/fjroland`.

## Running

```
Usage:
  fjroland [OPTIONS] PatternFile

Application Options:
  -m, --midi=        A midi device name to output to
  -s, --samples=     A directory of samples to use for waveform playback

Help Options:
  -h, --help         Show this help message
```

There are exmaple pattern files in `assets/patterns/` and a couple of sample packs in `assets/samplepacks`. Try this command to play "Vivrant Thing" on an accoustic drum set:

```
bin/fjroland assets/patterns/vivrant-thing.json -s assets/samplepacks/acoustic/
```

Or, you can try a TR-909:

```
bin/fjroland assets/patterns/vivrant-thing.json -s assets/samplepacks/909/
```

## Pattern File Format

The pattern file is a json document - it's easily parsed in any language, human readable, and easily editable. When you start the sequencer for a pattern, the sequencer watches for changes in the file and hot-reloads the pattern so you can play with patterns and hear the chnages interactively.

The pattern file has three fields:
- `title` - The title of the pattern being played
- `bpm` - The tempo of the pattern, in beats per minute


Each track takes the form of:
```
Instrument Name: timeline
```
The instrument names are standard names, but their spellings are important as the instrument name is how the samples or midi codes are generated. They are _case-insensitive_ though.

The timeline is a `|`-separated set of markers can be any length. Each `|` represents a beat. The other markers are as follows:
- `o` - Ghost note
- `X` - Standard note
- `>` - Accented note

While timelines can be any length, they must conform to the following rules or they will will not be valid:
- Each beat in a timeline must have the same number of divisions
- All timelines must have the same number of beats and divisions

The following is a simple _"Four on the Floor"_ pattern. It's 4 beats, divided into 4 divisions per beat.
```json
{
    "title": "Four on the Floor",
    "bpm": 100,
    "tracks": [
    "Snare:         |----|X---|----|X---|",
    "Closed Hi-Hat: |X-X-|X-X-|X-X-|X-X-|",
    "Bass:          |X---|X---|X---|X---|"
    ]
}
```
Here is a more sophisticated example - the rythm for _"Don't Say Nuthin'"_ by The Roots. This also has 4 divisisons per beat, but the pattern has 8 beats since it's a two bar phrase. You can also see it uses accent notes (`>`).
```json
{
    "title": "Don't Say Nothin' - The Roots",
    "bpm": 99,
    "tracks": [
    "Closed Hi-Hat: |X-X-|X-X-|X-X-|X-X-|X-X-|X-X-|X-X-|X-X-|",
    "Snare:         |----|>---|----|>---|----|>---|----|>---|",
    "Bass:          |X---|----|--X-|-X--|X---|----|----|-X--|"
    ]
}
```
Finally, here is the pattern for the intro to _"Turn Down for What"_ by DJ Snake and Lil Jon. It's only 4 beats long, but each beat is broken down into 8 divisions to accomodate the the 32nd notes in the final beat.
```json
{
    "title": "Turn Down for What (Intro)- DJ Snake, Lil Jon",
    "bpm": 100,
    "tracks": [
    "Snare: |>-X-X->-|X-X->-X-|X->-X-X-|>->->>>>|"
    ]
}
```

## Sample Packs
Sample packs are just directories of properly named wav files. This is the listing of the `assets/samples/acoustic` directory:

```
assets/samplepacks/acoustic/
├── Bass.wav
├── Closed Hi-Hat.wav
├── Crash Cymbal.wav
├── Cross-stick.wav
├── High Tom.wav
├── Low Tom.wav
├── Open Hi-Hat.wav
├── Pedal Hi-Hat.wav
├── Rimshot.wav
└── Snare.wav
```

Sample files must be wav files in stereo, sampled at 44.1k.

## Midi Devices

You can also output you pattern to an eligible midi device. This has been tested against the [Alesis SamplePad4](https://www.alesis.com/products/view2/samplepad-4).

To output to a midi device pass the name of the device, as it's identified to the host system, via the `-m` command line option. e.g.:


```
bin/fjroland assets/patterns/vivrant-thing.json -m SamplePad
```
