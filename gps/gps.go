package gps

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/jacobsa/go-serial/serial"
)

//Neo8 is a struct representing a gps object
type Neo8 struct {
	neoSerialPort  *io.ReadWriteCloser
	serialData     chan []byte
	serialReadData uint

	// data
	Lat            float64
	NS             string
	Lng            float64
	WE             string
	Fix            uint64
	SatTracking    uint64
	SatInVIew      uint64
	Hdop           float64
	Alti           float64
	TrueTrack      float64
	GroundSpeedKmh float64
	SNR            float64
}

//Init initializes the lora communication channel
func (n *Neo8) Init(serialDev string) {
	// OPEN SERIAL PORT
	options := serial.OpenOptions{
		PortName:        serialDev,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}
	neoSerialPort, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	} else {
		fmt.Println("serial opened successfuly")
	}

	n.neoSerialPort = &neoSerialPort

	var wg sync.WaitGroup
	wg.Add(1)

	n.setListener()
	fmt.Println("Init done")

}

func (n *Neo8) setListener() {
	serial0 := *n.neoSerialPort
	// variables to parse NMEA sentances
	sentance := string("")

	go func() {
		fmt.Println("listening")
		var data = make([]byte, 128)
		for {
			num, err := serial0.Read(data)
			if err != nil {
				// Handle error
				//return
			}
			if num != 0 {
				sentance += string(data[:num])
				// clear possibly wrong preamble
				del := 0
				for i := 0; i < len(sentance); i++ {
					if string(sentance[i]) == "$" {
						break
					} else {
						del = i
					}
				}
				sentance = sentance[del:]

				if len(sentance) != 0 {
					if strings.Contains(sentance, "\n") {
						for i := 0; i < len(sentance); i++ {
							if string(sentance[i]) == "\n" {
								if len(sentance[:i]) != 0 {
									n.updateState(sentance[:i])
									sentance = sentance[i:]
								}
							}
						}
					}
				}

			}
		}
	}()
}

func (n *Neo8) updateState(in string) {
	//fmt.Println(in[4:7])
	//fmt.Println(in)
	var err error
	var fl float64
	var ui uint64

	if in[4:7] == "GGA" {
		data := strings.Split(in, ",")
		//parse latitude
		fl, err = strconv.ParseFloat(data[2], 10)
		if err == nil {
			n.Lat = fl
		}
		//direction S / N needs no parsing
		n.NS = data[3]
		//parse longitude
		fl, err = strconv.ParseFloat(data[4], 10)
		if err == nil {
			n.Lng = fl
		}
		//direction W / E needs no parsing
		n.WE = data[5]
		//parse gps fix type
		ui, err = strconv.ParseUint(data[6], 10, 32)
		if err == nil {
			n.Fix = ui
		}
		//parse number of tracked satelites
		ui, err = strconv.ParseUint(data[7], 10, 32)
		if err == nil {
			n.SatTracking = ui
		}
		//parse horizontal accuracy
		fl, err = strconv.ParseFloat(data[8], 10)
		if err == nil {
			n.Hdop = fl
		}
		//parse altitude above sea
		fl, err = strconv.ParseFloat(data[9], 10)
		if err == nil {
			n.Alti = fl
		}
	}
	if in[4:7] == "VTG" {
		data := strings.Split(in, ",")
		//parse track direction
		fl, err = strconv.ParseFloat(data[1], 10)
		if err == nil {
			n.TrueTrack = fl
		}
		//parse ground speed
		fl, err = strconv.ParseFloat(data[7], 10)
		if err == nil {
			n.GroundSpeedKmh = fl
		}
	}
	if in[4:7] == "GSV" {
		data := strings.Split(in, ",")
		//parse all visible satelites
		ui, err = strconv.ParseUint(data[3], 10, 32)
		if err == nil {
			n.SatInVIew = ui
		}
		//parse signal to noise ratio
		fl, err = strconv.ParseFloat(data[6], 10)
		if err == nil {
			n.SNR = fl
		}
	}
}
