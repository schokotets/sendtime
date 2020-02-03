package main

import (
	"go.bug.st/serial"
	"time"
	"fmt"
	"os"
)

func main() {
	// https://godoc.org/go.bug.st/serial
	fmt.Println("Alle seriellen Schnittstellen suchen...")
	ports, err := serial.GetPortsList()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if len(ports) == 0 {
		fmt.Println("Keine seriellen Schnittstellen gefunden!")
		os.Exit(1)
	}

	fmt.Println("Tool, um die aktuelle Uhrzeit via USB zu senden")
	for i, port := range ports {
		fmt.Printf("%v) An %v senden\n", i, port)
	}
	fmt.Printf("Wähle ein USB-Gerät aus [%v-%v]: ", 0, len(ports)-1)
	var sel int
	_, err = fmt.Scanf("%d", &sel)

	if err != nil || sel < 0 || sel >= len(ports) {
		os.Exit(0)
	}
	selection := ports[sel]

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
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	n, err := port.Write(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v Bytes gesendet\n", n)
}
