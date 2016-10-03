package bme280

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

type fakeBus struct {
	data map[byte][]byte
	wErr map[byte]error
	rErr map[byte]error

	sync.Mutex
}

var goodData = map[byte][]byte{
	IDAddr:        []byte{IDVal},
	TempCompAddr:  make([]byte, 6),
	PressCompAddr: make([]byte, 18),
	H1CompAddr:    []byte{0},
	H2CompAddr:    make([]byte, 7),
	DataAddr:      make([]byte, 8),
}

func (fb *fakeBus) ReadReg(a byte, data []byte) error {
	fb.Lock()
	defer fb.Unlock()
	if err, ok := fb.rErr[a]; ok {
		return err
	}
	ldata, ok := fb.data[a]
	if !ok {
		return fmt.Errorf("no data for registry: %X", a)
	}
	for i, d := range ldata {
		data[i] = d
	}
	return nil
}

func (fb *fakeBus) WriteReg(a byte, d []byte) error {
	fb.Lock()
	defer fb.Unlock()
	if err, ok := fb.wErr[a]; ok {
		return err
	}
	return nil
}

func (fb *fakeBus) set(data map[byte][]byte, rErr, wErr map[byte]error) {
	fb.data = data
	fb.rErr = rErr
	fb.wErr = wErr
}

func TestInit(t *testing.T) {
	testData := []struct {
		desc    string
		data    map[byte][]byte
		rErr    map[byte]error
		wErr    map[byte]error
		wantErr bool
	}{
		{
			desc: "Failed to read ID data",
			rErr: map[byte]error{
				IDAddr: errors.New("error reading ID"),
			},
			wantErr: true,
		},
		{
			desc: "Bad ID data",
			data: map[byte][]byte{
				IDAddr: []byte{0},
			},
			wantErr: true,
		},
		{
			desc: "Failed to write config data.",
			data: goodData,
			wErr: map[byte]error{
				ConfigAddr: errors.New("error"),
			},
			wantErr: true,
		},
		{
			desc: "Failed to write control humidity data",
			data: goodData,
			wErr: map[byte]error{
				CtrlHumAddr: errors.New("error"),
			},
			wantErr: true,
		},
		{
			desc: "Failed to write control measurment data",
			data: goodData,
			wErr: map[byte]error{
				CtrlMeasAddr: errors.New("error"),
			},
			wantErr: true,
		},
		{
			desc: "Failed to read temperature compensation data",
			data: goodData,
			rErr: map[byte]error{
				TempCompAddr: errors.New("error"),
			},
			wantErr: true,
		},
		{
			desc: "Failed to read pressure compensation data",
			data: goodData,
			rErr: map[byte]error{
				PressCompAddr: errors.New("error"),
			},
			wantErr: true,
		},
		{
			desc: "Failed to read first humidity compensation data",
			data: goodData,
			rErr: map[byte]error{
				H1CompAddr: errors.New("error"),
			},
			wantErr: true,
		},
		{
			desc: "Failed to read second humidity compensation data",
			data: goodData,
			rErr: map[byte]error{
				H2CompAddr: errors.New("error"),
			},
			wantErr: true,
		},
		{
			desc: "Success",
			data: goodData,
		},
	}

	fb := &fakeBus{}
	b := &BME280{
		dev: fb,
	}

	for i, tc := range testData {
		fb.set(tc.data, tc.rErr, tc.wErr)
		err := b.Init()

		if (err == nil) != !tc.wantErr {
			t.Errorf("Test(%d) Desc(%s): Init => unexpected error got(%v) want(%v)", i, tc.desc, err, tc.wantErr)
		}
	}
}

func TestEnvData(t *testing.T) {
	testData := []struct {
		desc    string
		mode    byte
		data    map[byte][]byte
		rErr    map[byte]error
		wErr    map[byte]error
		wantErr bool
	}{
		{
			desc: "Failed to read sensor data",
			data: goodData,
			rErr: map[byte]error{
				DataAddr: errors.New("error reading data"),
			},
			wantErr: true,
		},
		{
			desc: "Forced mode, failed to write control humidity data.",
			mode: ForcedMode,
			data: goodData,
			wErr: map[byte]error{
				CtrlHumAddr: errors.New("error reading data"),
			},
			wantErr: true,
		},
		{
			desc: "Forced mode, failed to write control measurment data",
			mode: ForcedMode,
			data: goodData,
			wErr: map[byte]error{
				CtrlMeasAddr: errors.New("error"),
			},
			wantErr: true,
		},
		{
			desc: "Success",
			data: goodData,
		},
	}

	fb := &fakeBus{}
	b := &BME280{
		dev: fb,
	}

	for i, tc := range testData {
		fb.set(tc.data, tc.rErr, tc.wErr)
		b.mode = tc.mode
		_, _, _, err := b.EnvData()

		if (err == nil) != !tc.wantErr {
			t.Errorf("Test(%d) Desc(%s): Init => unexpected error got(%v) want(%v)", i, tc.desc, err, tc.wantErr)
		}
	}
}
