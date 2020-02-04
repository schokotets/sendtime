package main

import (
	bserial "go.bug.st/serial"
	"github.com/jacobsa/go-serial/serial"
	"time"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Alle seriellen Schnittstellen suchen...")
	ports, err := bserial.GetPortsList()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if len(ports) == 0 {
		fmt.Println("Keine seriellen Schnittstellen gefunden!")
		os.Exit(1)
	}

	fmt.Println("Tool, um die aktuelle Uhrzeit via USB zu senden")
	var selection string
	if len(ports) == 1 {
		fmt.Printf("Das USB-Gerät %v auswählen (j/n)? ", ports[0])
		var sel string
		_, err = fmt.Scanf("%s", &sel)

		if err != nil || sel != "j" && sel != "J" && sel != "y" && sel != "Y" {
			os.Exit(0)
		}
		selection = ports[0]
	} else {
		for i, port := range ports {
			fmt.Printf("%v) An %v senden\n", i, port)
		}
		fmt.Printf("Wähle ein USB-Gerät aus (%v-%v): ", 0, len(ports)-1)
		var sel int
		_, err = fmt.Scanf("%d", &sel)

		if err != nil || sel < 0 || sel >= len(ports) {
			os.Exit(0)
		}
		selection = ports[sel]
	}

	options := serial.OpenOptions{
		PortName: selection,
		BaudRate: 1000000,
		DataBits: 8,
		StopBits: 1,
		MinimumReadSize: 4,
		RTSCTSFlowControl: false,
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

	port, err := serial.Open(options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer port.Close()

	n, err := port.Write(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v Bytes gesendet\n", n)
}
