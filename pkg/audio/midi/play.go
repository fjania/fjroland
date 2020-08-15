package audio

import (
    "fmt"
    "log"
    "time"
    "github.com/rakyll/portmidi"
)

func Tryit() {

    portmidi.Initialize()

    fmt.Println(portmidi.CountDevices())
    fmt.Printf("%+v\n", portmidi.Info(portmidi.DefaultOutputDeviceID()))
    fmt.Println(portmidi.DefaultInputDeviceID())
    fmt.Println(portmidi.DefaultOutputDeviceID())

    out, err := portmidi.NewOutputStream(portmidi.DefaultOutputDeviceID(), 1024, 0)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", portmidi.Info(portmidi.DefaultOutputDeviceID()))

    // Send "note on" events to play C major chord.
    for i := 0; i< 12; i++ {
        out.WriteShort(0x9A, 42, 127)
        time.Sleep(time.Second/8)
        out.WriteShort(0x9A, 39, 127)
        time.Sleep(time.Second/8)
        out.WriteShort(0x9A, 38, 127)
        time.Sleep(time.Second/8)
        out.WriteShort(0x9A, 35, 127)
        time.Sleep(time.Second/8)
    }

    // Notes will be sustained for 2 seconds.

    out.Close()


    portmidi.Terminate()
}
