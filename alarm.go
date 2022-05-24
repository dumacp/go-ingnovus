package ingnovus

import (
	"fmt"
	"time"
)

type AlarmCode int

const (
	CortoCircuito       AlarmCode = 1
	BateriaDesconectada AlarmCode = 2
	SensorObstruido     AlarmCode = 3
	TapaCajaAbierta     AlarmCode = 4
	SensorBloqueado1    AlarmCode = 5
	SensorBloqueado2    AlarmCode = 6
	SensorBloqueado3    AlarmCode = 7
	ResetSistem         AlarmCode = 12
	PuertaCerrada       AlarmCode = 20
	PuertaAbierta       AlarmCode = 21
)

func (a AlarmCode) String() string {
	switch a {
	case SensorObstruido:
		return "SensorObstruido"
	case SensorBloqueado1:
		return "SensorBloqueado1"
	case SensorBloqueado2:
		return "SensorBloqueado2"
	case SensorBloqueado3:
		return "SensorBloqueado3"
	case PuertaAbierta:
		return "PuertaAbierta"
	case PuertaCerrada:
		return "PuertaCerrada"
	case BateriaDesconectada:
		return "BateriaDesconectada"
	default:
		return ""
	}
}

type alarm struct {
	*event
	Alarmtime   time.Time `json:"alarmTime"`
	AAlarmTotal int       `json:"sumAlarm"`
	Ccode       AlarmCode `json:"alarmCode"`
}
type Alarm interface {
	Event
	Code() AlarmCode
	AlarmTotal() int
	AlarmTime() time.Time
}

func ParseAlarm(e Event) (Alarm, error) {
	if e.Type() == AlarmEvent {
		v, ok := e.(Alarm)
		if !ok {
			return nil, fmt.Errorf("event is not a AlarmEvent, type: %v", e.Type())
		}
		return v, nil
	}
	return nil, fmt.Errorf("event is not a AlarmEvent, type: %v", e.Type())
}

func (a *alarm) Code() AlarmCode {
	return a.Ccode
}

func (a *alarm) AlarmTotal() int {
	return a.AAlarmTotal
}

func (a *alarm) AlarmTime() time.Time {
	return a.Alarmtime
}
