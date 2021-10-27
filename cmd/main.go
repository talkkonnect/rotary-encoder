package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/talkkonnect/max7219"
	"github.com/talkkonnect/rotary-encoder/device"
	"periph.io/x/conn/gpio/gpioreg"
	"periph.io/x/host"
)

func main() {
	var TKChannel int = 9
	pinAName := flag.String("pina", "", "pin name for a channel of rotary encoder")
	pinBName := flag.String("pinb", "", "pin name for a channel of rotary encoder")
	help := flag.Bool("h", false, "print help page")

	flag.Parse()

	if *help || *pinAName == "" || *pinBName == "" {
		flag.Usage()
		os.Exit(0)
	}

	if _, err := host.Init(); err != nil {
		panic(err)
	}

	aPin := gpioreg.ByName(*pinAName)
	bPin := gpioreg.ByName(*pinBName)

	re := device.NewRotaryEncoderWithCustomTimeout(aPin, bPin, 3*time.Second)

	mtx := max7219.NewMatrix(1)
	err := mtx.Open(0, 0, 7)
	if err != nil {
		log.Fatal(err)
	}
	defer mtx.Close()

	fmt.Println("reading...")
	mtx.Device.SevenSegmentDisplay(strconv.Itoa(TKChannel))

	for {
		//spew.Dump(re.Read())
		result := re.Read()
		if result == "counterClockwise" {
			TKChannel++
		}
		if result == "clockwise" {
			TKChannel--
		}
		if TKChannel < 1 {
			TKChannel = 40
		}
		if TKChannel > 40 {
			TKChannel = 1
		}
		log.Print("debug: Channel ", TKChannel)
		mtx.Device.SevenSegmentDisplay(strconv.Itoa(TKChannel))
	}
}
