package main

import (
	"github.com/paulrademacher/climenu"
	"go.bug.st/serial"
	"time"
	"fmt"
	"log"
	"os"
)

func main() {
	// https://godoc.org/go.bug.st/serial
	fmt.Println("Alle seriellen Schnittstellen suchen...")
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("Keine seriellen Schnittstellen gefunden!")
	}

	menu := climenu.NewButtonMenu("\nTool, um die aktuelle Uhrzeit via USB zu senden\n"+
		"Auswahl mit den Pfeiltasten, Bestätigen mit Enter, Abbrechen mit Esc",
		"Wähle ein USB-Gerät aus")
	for _, port := range ports {
		menu.AddMenuItem(fmt.Sprintf("An %v senden", port), port)
	}

	selection, escaped := menu.Run()
	if escaped {
		os.Exit(0)
	}

	mode := &serial.Mode{
		BaudRate: 1000000,
		Parity: serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	t := time.Now()
	zone, _ := t.Zone()
	dstbyte := byte(0) //daylight savings
	if zone == "CEST" {
		dstbyte = byte(1)
	}

	t = time.Now()
	fmt.Println("Warten auf nächste volle Sekunde...")
	time.Sleep(time.Duration(1000000000-t.Nanosecond()))

	t = time.Now()
	data := []byte{255, byte(t.Hour()), byte(t.Minute()), byte(t.Second()), dstbyte}
	fmt.Printf("%v:%v:%v + DST(%v) an %v senden...\n", data[1], data[2], data[3], data[4], selection)

	port, err := serial.Open(selection, mode)
	if err != nil {
		log.Fatal(err)
	}

	n, err := port.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v Bytes gesendet\n", n)
}
