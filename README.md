# sendtime

This tool is a utility to send the current time to an embedded device (guest) from a host pc.

Make sure the proper serial drivers are installed on the host. For the Arduino Nano, this often means installing CH340 drivers.

## Transmission
Communication happens through a serial port.
A baud rate of 1000000 is used.
8 data bits are followed by 1 stop bit.
To ensure compatibility with the Arduino Nano, the RTSCTSFlowControl property is disabled.

## Data
Five bytes are sent, communicating the host's current time in the context of its time zone.
* 0xFF to start the communication
* the current hour
* the current minute
* the current second
* 0x00 or 0x01 depending on if it's daylight savings time(1)

1) The DST parameter is useful for fixed-installation clocks, so they could have a button for DST/Non-DST switching.

## Procedure
* unless already done, install all go dependencies using `go get go.bug.st/serial github.com/jacobsa/go-serial/serial`
* attach the device to a serial interface, like usb
* start the command line utility using `go run sendtime.go`
* select the serial device by typing the corresponding number / confirm the selected serial device with `y`
* wait for the host to send the current time
* done!
