# sendtime

"sendtime" is a command line utility to send the current local time to an embedded device (guest) from a host pc. Because we use it to configure a binary clock with different LED strips, it can also send GRB color values with 21-bit depth for each LED strip. The user selects whether the current time or which LED strip's color is to be sent.

## Transmission
Communication happens through a serial port.
A baud rate of 1000000 is used.
8 data bits are followed by 1 stop bit.
To ensure compatibility with the Arduino Nano, the RTSCTSFlowControl property is disabled.

## Data

### Sending the current time
6 bytes are sent, communicating the host's current time in the context of its time zone.
* 0x82 to start the communication of time
* the current second
* the current minute
* the current hour
* 0x00 (non-DST) or 0x01 (DST) depending on if it's daylight savings time(1)
* 0x81 to terminate the communication of time

1) The DST parameter is useful for fixed-installation clocks, so they could have a button for DST/Non-DST switching.

### Byte labels of the different LED strips
* hours' LEDs: 0x83
* minutes' LEDs: 0x84
* seconds' LEDs: 0x85
* labels' LEDs: 0x86

### Sending an RGB color value
4 bytes are sent.
* the LED strip's byte label (see above)
* 1 zero bit + 7 bits encoding the green value
* 1 zero bit + 7 bits encoding the red value
* 1 zero bit + 7 bits encoding the blue value

## Setup
Make sure the proper serial drivers are installed on the host. For the Arduino Nano, this often means installing CH340 drivers.

### Running fom source
When running from source, install [Go](https://golang.org) 1.16 on your pc.
Next, clone this repository.
Then open the command line, navigate to the repository's directory and install the required go dependencies.

`go mod download`

## Use
* attach the device to a serial interface, like usb
* start the command line utility by executing the binary or, after downloading the dependencies as described above, using `go run sendtime.go`
* select the serial device by typing the corresponding number / confirm the suggested serial device with `y` or `j`
* select whether the current time or which LED color should be sent
* answer the posed questions
* wait for the host to send the data
* done!
