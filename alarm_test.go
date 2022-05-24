package ingnovus

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func TestParseAlarm(t *testing.T) {

	// raw := "01 41 41 41 30 30 30 30 48 30 34 53 30 37 00 00 00 82 00 33 00 16 00 00 00 00 00 25 00 00 00 00 40 0b bd 00 3d af d3 c0 0e 0a 17 0f 0a 1c c5"
	// raw := "02 41 41 41 30 30 30 30 48 30 34 53 30 37 00 00 00 81 00 32 00 16 00 00 00 00 00 23 00 00 00 00 40 0b bd 00 3d af d3 c0 0e 0a 17 0f 0a 1c 11 01 0e 0a 0f 0a 05"
	// raw := "02 41 41 41 30 30 30 30 48 30 34 53 30 37 00 00 00 80 00 31 00 16 00 00 00 00 00 21 00 00 00 00 40 0b bd 00 3d af d3 c0 0e 0a 17 0f 0a 1c 11 01 0e 0a 0f 0d 04"
	raw := "02 41 41 41 30 30 30 30 48 30 34 53 30 37 00 00 00 74 00 25 00 14 00 00 00 00 00 16 00 00 00 00 40 0b bd 00 3d af d3 c0 0e 0a 15 0f 0a 1c 05 05 0e 0a 0f 0a d2"
	data, err := hex.DecodeString(strings.ReplaceAll(raw, " ", ""))
	if err != nil {
		log.Fatalln(err)
	}

	evt, err := NewEvent(data)
	if err != nil {
		log.Fatalln(err)
	}

	type args struct {
		e Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				e: evt,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAlarm(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAlarm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("latitude: %v, longitude: %v", got.Latitude(), got.Longtitude())
			t.Logf("date: %v", got.Time())
			t.Logf("alarm date: %v", got.AlarmTime())
			t.Logf("alarm code: %v, sum: %v", got.Code(), got.AlarmTotal())

			res, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("alarm: %s", res)
		})
	}
}
