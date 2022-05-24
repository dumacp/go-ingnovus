package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dumacp/go-ingnovus"
)

var socket string
var baudRate int

func init() {
	flag.StringVar(&socket, "port", "/dev/ttyS2", "serial port")
	flag.IntVar(&baudRate, "baud", 9600, "baudrate")
}

func main() {

	flag.Parse()

	dev, err := ingnovus.NewDevice(socket, baudRate, 1000*time.Millisecond)
	if err != nil {
		log.Fatalln(err)
	}
	defer dev.Close()

	ch := dev.ListennEvents()

	for v := range ch {
		data, err := json.Marshal(v)
		if err != nil {
			break
		}
		fmt.Printf("event: %s", data)
	}
}
