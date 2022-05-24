package ingnovus

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
)

func (d *device) ListennEvents() chan Event {

	ch := make(chan Event)

	go func() {

		defer close(ch)

		t1 := time.NewTicker(1000 * time.Millisecond)
		defer t1.Stop()

		for {

			select {
			case <-d.quit:
				return
			case <-t1.C:
				read := bufio.NewReader(d.port)
				buf1 := make([]byte, 128)
				tn1 := time.Now()
				n1, err := read.Read(buf1)
				if err != nil {
					if errors.Is(err, io.EOF) {
						// fmt.Printf("read request error: %s, timeout: %s\n", err, time.Since(tn1))
						if d.config.ReadTimeout/10 > time.Since(tn1) {
							fmt.Printf("read request error: %s, timeout: %s\n", err, time.Since(tn1))
							return
						}
						continue
					}
					fmt.Printf("read request error: %s, timeout: %s\n", err, time.Since(tn1))
					return
				}
				if n1 <= 0 {
					continue
				}
				data1 := buf1[:n1]
				if !bytes.Contains(data1, []byte("SttReq")) {
					fmt.Printf("wrong expected data: %s, %X\n", data1, data1)
					continue
				}
				raw := "STT;0970000212;BFFFFF;97;1.0.2;0;%s;00000000;0;0;0000;0;+3.467287;-76.525963;0.00;0.00;11;1;00000000;00000000;0;0;2802;00038003;0.0;13.79;37;0;0\n"
				frame := fmt.Sprintf(raw, time.Now().Format("20060102;15:04:05"))
				// fmt.Printf("frame: %s\n", frame)
				tn3 := time.Now()
				if _, err := d.port.Write([]byte(frame)); err != nil {
					fmt.Printf("write error: %s, timeout: %s\n", err, time.Since(tn3))
					return
				}
				buf2 := make([]byte, 128)
				tn2 := time.Now()
				n2, err := read.Read(buf2)
				if err != nil {
					if errors.Is(err, io.EOF) {
						// fmt.Printf("read respone error: %s, timeout: %s\n", err, time.Since(tn2))
						if d.config.ReadTimeout/10 > time.Since(tn2) {
							fmt.Printf("read respone error: %s, timeout: %s\n", err, time.Since(tn2))
							return
						}
						continue
					}
					fmt.Printf("read response error: %s, timeout: %s\n", err, time.Since(tn2))
					return
				}
				data2 := buf2[:n2]
				evt, err := NewEvent(data2)
				if err != nil {
					fmt.Println(err)
					continue
				}

				select {
				case ch <- evt:
				case <-time.After(300 * time.Millisecond):
				}

			}
		}
	}()

	return ch
}
