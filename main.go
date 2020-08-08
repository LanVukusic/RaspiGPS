package main

import (
	"fmt"
	"time"

	"./gps"
)

func hello() {
	fmt.Println("Hello world goroutine")
}

var a gps.Neo8

func main() {
	a.Init("/dev/serial0")
	
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("Position: %g %s, %g %s with %d satelites\n", a.Lat, a.NS, a.Lng, a.WE, a.SatTracking)
	}

}
