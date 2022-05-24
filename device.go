package ingnovus

import (
	"time"

	"github.com/tarm/serial"
)

type device struct {
	config *serial.Config
	port   *serial.Port
	quit   chan int
}

type Device interface {
	ListennEvents() chan Event
	Close()
}

func NewDevice(port string, baud int, timeout time.Duration) (Device, error) {
	config := &serial.Config{
		Name:        port,
		Baud:        baud,
		ReadTimeout: timeout,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
	}

	s, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}
	dev := &device{
		config: config,
		port:   s,
	}
	dev.quit = make(chan int)
	return dev, nil
}

func (d *device) Close() {
	if d.quit != nil {
		select {
		case _, ok := <-d.quit:
			if ok {
				close(d.quit)
			}
		default:
			close(d.quit)
		}
	}
	d.port.Close()
}
