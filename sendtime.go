package main

import (
	"errors"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	bserial "go.bug.st/serial"
	"os"
	"time"
)

func main() {
	fmt.Println("Programm für die Binary-Clock, um die aktuelle Uhrzeit via USB zu senden")

	selection, err := getPortSelection()
	if err != nil {
		fmt.Printf("Fehler bei der Port-Auswahl: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Port-Auswahl getroffen: %s\n", *selection)

	options := serial.OpenOptions{
		PortName:          *selection,
		BaudRate:          1000000,
		DataBits:          8,
		StopBits:          1,
		MinimumReadSize:   4,
		RTSCTSFlowControl: false,
	}

	sendTime(selection, &options)

	if err != nil {
		fmt.Printf("Fehler beim Senden der Zeit: %v\n", err)
		os.Exit(1)
	}
}

func getPortSelection() (*string, error) {
	fmt.Println("Alle seriellen Schnittstellen suchen...")

	ports, err := bserial.GetPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, errors.New("keine seriellen Schnittstellen gefunden")
	}

	var selection string
	if len(ports) == 1 {
		fmt.Printf("Das USB-Gerät %s auswählen (j/n)? ", ports[0])

		var sel string
		_, err = fmt.Scanf("%s", &sel)

		if err != nil {
			return fmt.Errorf("kann Antwort auf Port-Frage nicht lesen: %v", err)
		}

		if sel != "j" && sel != "J" && sel != "y" && sel != "Y" {
			return nil, errors.New("einzig auswählbarer Port nicht ausgewählt")
		}

		selection = ports[0]
	} else {
		for i, port := range ports {
			fmt.Printf("%d) An %s senden\n", i+1, port)
		}
		fmt.Printf("Wähle ein USB-Gerät aus (%d-%d): ", 1, len(ports))
		var sel int
		_, err = fmt.Scanf("%d", &sel)

		if err != nil {
			return fmt.Errorf("kann Antwort auf Port-Frage nicht lesen: %v", err)
		}

		if sel < 1 || sel > len(ports) {
			return nil, fmt.Errorf("ungültige Port-Auswahl - %d-%d wären möglich", 1, len(ports))
		}
		selection = ports[sel-1]
	}
	return &selection, nil
}

func sendTime(selection *string, options *serial.OpenOptions) error {
	fmt.Printf("Ist gerade Sommerzeit? (j/n): ")

	var sel string
	_, err := fmt.Scanf("%s", &sel)

	if err != nil {
		return fmt.Errorf("kann Antwort auf Sommerzeit-Frage nicht lesen: %v", err)
	}

	dstbyte := byte(0) //daylight savings

	if sel == "j" || sel == "J" || sel == "y" || sel == "Y" {
		dstbyte = byte(1)
	}

	// wait for next full second
	// note: this won't be perfectly precise because logging and opening the port takes time
	t := time.Now()
	fmt.Println("Warten auf nächste volle Sekunde...")
	time.Sleep(time.Duration(1000000000 - t.Nanosecond()))

	// prepare data
	t = time.Now()
	data := []byte{0x82, byte(t.Second()), byte(t.Minute()), byte(t.Hour()), dstbyte, 0x81}
	fmt.Printf("%v:%v:%v + DST(%v) an %v senden...\n", data[3], data[2], data[1], data[4], *selection)

	return openAndSend(options, &data)
}

func openAndSend(options *serial.OpenOptions, data *[]byte) error {
	// open port
	port, err := serial.Open(*options)
	if err != nil {
		return fmt.Errorf("beim Öffnen des Ports: %v", err)
	}
	defer port.Close()

	// write data
	n, err := port.Write(*data)
	if err != nil {
		return fmt.Errorf("beim Senden der Daten: %v", err)
	}
	fmt.Printf("%d Bytes gesendet\n", n)
	return nil
}
