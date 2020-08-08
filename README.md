# RaspiGPS 📡
## RaspberryPi GPS library written in Go for interfacing with GPS serial hardware
This is a library that was written to work with any serial enabled gps module. I had only an option to test the Neo 8 gps, but it should work with any NMEA compliant device.
Please leave a star if you end up using it in a project. It means a tone to me ⭐ 
It was written for a series of serial enabled extension boards available from [Aliexpress](https://www.aliexpress.com/item/32325428866.html?spm=a2g0s.9042311.0.0.27424c4dH6kF1l)   

### Hardware list:   
__Tested hardware:__  
* [Ublox Neo M8N](https://www.u-blox.com/en/product/neo-m8-series). I tested [this](https://www.aliexpress.com/item/32325428866.html?spm=a2g0s.9042311.0.0.27424c4dH6kF1l) one  

__Probably working hardware:__  
All serial and Nmea compliant GPS modules
* Probably all Ublox series gps
* maybe others?

### Pinout:  
Pins on the raspberry are labeled according to the diagram below, referenced with numbers running from 1 to 40  

<img src="https://raw.githubusercontent.com/LanVukusic/RaspiGPS/master/pinout.png" width="500">

* __PIN 8__ serial TX (gpio15)    -  gps RX
* __PIN 10__ serial RX (gpio14)   -  gps TX



### Library is dependent on an external library:  
* [go serial library](https://github.com/jacobsa/go-serial) - github.com/jacobsa/go-serial  


Here is a sample code which connects and recieves data from a gps module.  


```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"./LORA"
	"github.com/warthog618/gpiod/device/rpi"
)

var a LORA.Comms

func main() {
	// https://pi4j.com/1.2/pins/model-zero-rev1.html
	a.Init("gpiochip0", rpi.J8p12, 1) //gpio_chip, rpi pin for settings lora mode, minimal read size
	a.GetInfo()  // prints our gpio chip and pin number

	// 1 indicates that we call back evey "1" recieved byte of data.
	//if you wish to wait for a bigger chunk, set number higher
	a.SetLoraListener(1, loraCallback)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		a.WriteString(text)
	}

}

func loraCallback(b []byte) {
	fmt.Println(string(b))
}
```

Following functions are available:  
| Name           | DataType | Description                                                   |
|----------------|----------|---------------------------------------------------------------|
| Lat            | float64  | Latitude                                                      |
| NS             | string   | North / South part of latitude                                |
| Lng            | float64  | Longitude                                                     |
| WE             | string   | West / East part if Longitude                                 |
| Fix            | uint64   | GPS fix type; described below                                 |
| SatTracking    | uint64   | Number of satellites used in tracking                         |
| SatInView      | uint64   | Number of satellites in view                                  |
| Hdop           | float64  | Horizontal dilution of precision aka accuracy on land.[meters]|
| Alti           | float64  | Altitude above sea level [meters]                            |
| TrueTrack      | float64  | Heading track [degrees azimuth]                    |
| GroundSpeedKmh | float64  | Speed above ground in [Km/h]                     |
| SNR            | float64  | Signal to noise ratio                                         |


__GPS fix types:__
  0 = invalid
  1 = GPS fix (SPS)
  2 = DGPS fix
  3 = PPS fix
  4 = Real Time Kinematic
  5 = Float RTK
  6 = estimated (dead reckoning) (2.3 feature)
  7 = Manual input mode
  8 = Simulation mode

