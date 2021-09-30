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
	fmt.Println("Programm für die Konfiguration der Binary-Clock via USB")

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

	actionCode, err := chooseAction()

	if err != nil {
		fmt.Printf("Fehler beim Wählen der Aktion: %v\n", err)
		os.Exit(1)
	}

	switch actionCode {
	case 1:
		err = sendTime(selection, &options)

		if err != nil {
			fmt.Printf("Fehler beim Senden der Zeit: %v\n", err)
			os.Exit(1)
		}

	case 2:
		err = sendColor(selection, &options, "Stunden", 0x83)

		if err != nil {
			fmt.Printf("Fehler beim Senden der Farbe der Stunden: %v\n", err)
			os.Exit(1)
		}

	case 3:
		err = sendColor(selection, &options, "Minuten", 0x84)

		if err != nil {
			fmt.Printf("Fehler beim Senden der Farbe der Minuten: %v\n", err)
			os.Exit(1)
		}

	case 4:
		err = sendColor(selection, &options, "Sekunden", 0x85)

		if err != nil {
			fmt.Printf("Fehler beim Senden der Farbe der Sekunden: %v\n", err)
			os.Exit(1)
		}

	case 5:
		err = sendColor(selection, &options, "Rand", 0x86)

		if err != nil {
			fmt.Printf("Fehler beim Senden der Farbe des Rands: %v\n", err)
			os.Exit(1)
		}

		// omitting default case
		// because only 1-5 is possible from chooseAction

	}

}

func chooseAction() (int, error) {
	fmt.Println()
	fmt.Println("1) Aktuelle Uhrzeit senden")
	fmt.Println("2) Farbe der Stunden-LEDs ändern")
	fmt.Println("3) Farbe der Minuten-LEDs ändern")
	fmt.Println("4) Farbe der Sekunden-LEDs ändern")
	fmt.Println("5) Farbe der Rand-LEDs ändern")
	fmt.Println()
	fmt.Print("Wähle eine Aktion aus (1-5): ")

	var sel int
	_, err := fmt.Scanf("%d", &sel)

	if err != nil {
		return 0, fmt.Errorf("kann Antwort auf Aktions-Frage nicht lesen: %v", err)
	}
	if sel < 1 || sel > 5 {
		return 0, fmt.Errorf("ungültige Aktions-Auswahl - 1-5 wären möglich")
	}
	return sel, nil
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
			return nil, fmt.Errorf("kann Antwort auf Port-Frage nicht lesen: %v", err)
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
			return nil, fmt.Errorf("kann Antwort auf Port-Frage nicht lesen: %v", err)
		}

		if sel < 1 || sel > len(ports) {
			return nil, fmt.Errorf("ungültige Port-Auswahl - %d-%d wären möglich", 1, len(ports))
		}
		selection = ports[sel-1]
	}
	return &selection, nil
}

func sendTime(selection *string, options *serial.OpenOptions) error {
	fmt.Print("Ist gerade Sommerzeit? (j/n): ")

	var sel string
	_, err := fmt.Scanf("%s", &sel)

	if err != nil {
		return fmt.Errorf("kann Antwort auf Sommerzeit-Frage nicht lesen: %v", err)
	}

	dstbyte := byte(0) //daylight savings

	if sel == "j" || sel == "J" || sel == "y" || sel == "Y" {
		dstbyte = byte(1)
	}

	fmt.Println()

	// wait for next full second
	// note: this won't be perfectly precise because logging and opening the port takes time
	t := time.Now()
	fmt.Println("Warten auf nächste volle Sekunde...")
	time.Sleep(time.Duration(1000000000 - t.Nanosecond()))

	// prepare data
	t = time.Now()
	data := []byte{0x82, byte(t.Second()), byte(t.Minute()), byte(t.Hour()), dstbyte, 0x81}
	fmt.Printf("0x82, %ds %dm %dh, DST(%x), 0x81 an %s senden...\n", data[1], data[2], data[3], data[4], *selection)

	return openAndSend(options, &data)
}

func sendColor(selection *string, options *serial.OpenOptions, area string, startbyte byte) error {
	fmt.Println("Die Farbe ist nach RGB-Modell in einen Rot-, Grün- und Blauwert aufgeteilt.")
	fmt.Printf("Konfiguration der Farbe der %s-LEDs\n", area)

	red, err := getColorValue("Rot")
	if err != nil {
		return err
	}

	green, err := getColorValue("Grün")
	if err != nil {
		return err
	}

	blue, err := getColorValue("Blau")
	if err != nil {
		return err
	}

	fmt.Println()

	data := []byte{startbyte, green, red, blue}
	fmt.Println("Die angegebenen Werte wurden durch 2 geteilt, da die Komponenten nur 0-127 sein können.")
	fmt.Printf("0x%x, Grün %d, Rot %d, Blau %d an %s senden...\n", data[0], data[1], data[2], data[3], *selection)

	return openAndSend(options, &data)
}

func getColorValue(colorname string) (byte, error) {
	fmt.Printf("%s-Wert (0-255)? ", colorname)

	var sel int
	_, err := fmt.Scanf("%d", &sel)

	if err != nil {
		return 0, fmt.Errorf("kann Antwort auf %s-Wert-Frage nicht lesen: %v", colorname, err)
	}
	if sel < 0 || sel > 255 {
		return 0, fmt.Errorf("ungültiger %s-Wert - 0-255 wären möglich", colorname)
	}

	return byte(sel) / 2, nil
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
