package ingnovus

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

type event struct {
	Ttype     EventType `json:"type"`
	Pplate    []byte    `json:"plate"`
	Ffirmware []byte    `json:"firmare"`
	Inputs    []int     `json:"inputs"`
	Outputs   []int     `json:"outputs"`
	Lat       float32   `json:"lat"`
	Long      float32   `json:"long"`
	Nnum      int64     `json:"numeration"`
	Daytotal  int       `json:"daytotal"`
	Ttime     time.Time `json:"time"`
}

type EventType int

func (e EventType) String() string {
	switch e {
	case CountingEvent:
		return "CountingEvent"
	case AlarmEvent:
		return "AlarmEvent"
	default:
		return ""
	}
}

const (
	CountingEvent  EventType = 1
	AlarmEvent     EventType = 2
	DoorStateEvent EventType = 4
)

type Event interface {
	Type() EventType
	Plate() []byte
	Firmware() []byte
	InputsP1() int
	InputsP2() int
	InputsP3() int
	OutputsP1() int
	OutputsP2() int
	OutputsP3() int
	Latitude() float32
	Longtitude() float32
	Numeration() int64
	DayTotal() int
	Time() time.Time
}

func (e *event) Type() EventType {
	return e.Ttype
}

func (e *event) Plate() []byte {
	return e.Pplate
}

func (e *event) Firmware() []byte {
	return e.Ffirmware
}

func (e *event) InputsP1() int {
	if len(e.Inputs) < 1 {
		return 0
	}
	return e.Inputs[0]
}

func (e *event) InputsP2() int {
	if len(e.Inputs) < 2 {
		return 0
	}
	return e.Inputs[1]
}

func (e *event) InputsP3() int {
	if len(e.Inputs) < 3 {
		return 0
	}
	return e.Inputs[2]
}

func (e *event) OutputsP1() int {
	if len(e.Outputs) < 1 {
		return 0
	}
	return e.Outputs[0]
}

func (e *event) OutputsP2() int {
	if len(e.Outputs) < 2 {
		return 0
	}
	return e.Outputs[1]
}

func (e *event) OutputsP3() int {
	if len(e.Outputs) < 3 {
		return 0
	}
	return e.Outputs[2]
}

func (e *event) Latitude() float32 {
	return e.Lat
}

func (e *event) Longtitude() float32 {
	return e.Long
}

func (e *event) Numeration() int64 {
	return e.Nnum
}

func (e *event) DayTotal() int {
	return e.Daytotal
}

func (e *event) Time() time.Time {
	return e.Ttime
}

type DoorState interface {
	Event
	DoorState1() int
	DoorState2() int
}

func NewEvent(data []byte) (Event, error) {

	if len(data) < 44 {
		return nil, fmt.Errorf("len data in less that 44 bytes, len: %d, data: [% X]", len(data), data)
	}
	sum := 0
	for _, x := range data[:len(data)-1] {
		sum += int(x) & 0xFF
	}

	if byte(sum&0xFF) != data[len(data)-1] {
		return nil, fmt.Errorf("checksum is wrong: %X != %X", sum, data[len(data)-1])
	}

	evt := new(event)

	switch data[0] {
	case 1, 2, 4:
	default:
		return nil, fmt.Errorf("event type is unknown")
	}

	evt.Ttype = EventType(int(data[0]) & 0xFF)

	switch evt.Ttype {
	case AlarmEvent:
		if len(data) < 53 {
			return nil, fmt.Errorf("event type AlarmEvent require len data = 52 bytes")
		}
	default:
	}

	evt.Pplate = make([]byte, 0)
	evt.Pplate = append(evt.Pplate, data[1:8]...)
	evt.Ffirmware = make([]byte, 0)
	evt.Ffirmware = append(evt.Pplate, data[8:14]...)
	evt.Nnum = int64(binary.BigEndian.Uint32(data[14:18]))
	evt.Daytotal = int(binary.BigEndian.Uint16(data[18:20]))
	evt.Inputs = make([]int, 0)
	evt.Inputs = append(evt.Inputs, int(binary.BigEndian.Uint16(data[20:22])))
	evt.Inputs = append(evt.Inputs, int(binary.BigEndian.Uint16(data[22:24])))
	evt.Inputs = append(evt.Inputs, int(binary.BigEndian.Uint16(data[24:26])))
	evt.Outputs = make([]int, 0)
	evt.Outputs = append(evt.Outputs, int(binary.BigEndian.Uint16(data[26:28])))
	evt.Outputs = append(evt.Outputs, int(binary.BigEndian.Uint16(data[28:30])))
	evt.Outputs = append(evt.Outputs, int(binary.BigEndian.Uint16(data[30:32])))

	lat := make([]byte, 4)
	for i, v := range data[32:36] {
		lat[i] = v
	}
	evt.Lat = float32(math.Float32frombits(binary.BigEndian.Uint32(lat)))

	long := make([]byte, 8)
	for i, v := range data[36:40] {
		long[i] = v
	}
	evt.Long = float32(math.Float32frombits(binary.BigEndian.Uint32(long)))

	loc := time.UTC
	month := time.Month(data[41])
	evt.Ttime = time.Date(2000+int(data[42]), month, int(data[40]), int(data[43]), int(data[44]), int(data[45]), 0, loc)

	switch evt.Ttype {
	case CountingEvent:
		return &counting{evt}, nil
	case AlarmEvent:
		alarm := &alarm{
			event: evt,
		}
		alarm.Ccode = AlarmCode(int(data[46]) & 0xFF)
		alarm.AAlarmTotal = int(data[47]) & 0xFF
		month := time.Month(data[49])
		alarm.Alarmtime = time.Date(evt.Ttime.Year(), month, int(data[48]), int(data[50]), int(data[51]), 0, 0, loc)
		return alarm, nil
	default:
		return evt, nil
	}
}
