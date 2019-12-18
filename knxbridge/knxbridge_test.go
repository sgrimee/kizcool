package knxbridge

import (
	"reflect"
	"testing"

	"github.com/sgrimee/kizcool"
	"github.com/sgrimee/kizcool/config/knxcfg"
)

func TestNew(t *testing.T) {
	type args struct {
		k *kizcool.Kiz
		d []knxcfg.Device
	}
	tests := []struct {
		name string
		args args
		want *Bridge
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.k, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestBridge_Start(t *testing.T) {
// 	type fields struct {
// 		kiz     *kizcool.Kiz
// 		devices []knxcfg.Device
// 	}
// 	type args struct {
// 		finish chan struct{}
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			br := &Bridge{
// 				kiz:     tt.fields.kiz,
// 				devices: tt.fields.devices,
// 			}
// 			if err := br.Start(tt.args.finish); (err != nil) != tt.wantErr {
// 				t.Errorf("Bridge.Start() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestBridge_processKizEvent(t *testing.T) {
// 	type fields struct {
// 		kiz     *kizcool.Kiz
// 		devices []knxcfg.Device
// 	}
// 	type args struct {
// 		kizEvent kizcool.Event
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			br := &Bridge{
// 				kiz:     tt.fields.kiz,
// 				devices: tt.fields.devices,
// 			}
// 			if err := br.processKizEvent(tt.args.kizEvent); (err != nil) != tt.wantErr {
// 				t.Errorf("Bridge.processKizEvent() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestBridge_processKnxEvent(t *testing.T) {
// 	type fields struct {
// 		kiz     *kizcool.Kiz
// 		devices []knxcfg.Device
// 	}
// 	type args struct {
// 		knxEvent knx.GroupEvent
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			br := &Bridge{
// 				kiz:     tt.fields.kiz,
// 				devices: tt.fields.devices,
// 			}
// 			if err := br.processKnxEvent(tt.args.knxEvent); (err != nil) != tt.wantErr {
// 				t.Errorf("Bridge.processKnxEvent() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestBridge_processCommand(t *testing.T) {
// 	type fields struct {
// 		kiz     *kizcool.Kiz
// 		devices []knxcfg.Device
// 	}
// 	type args struct {
// 		gcmd knxcfg.Command
// 		d    kizcool.Device
// 		data []byte
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			br := &Bridge{
// 				kiz:     tt.fields.kiz,
// 				devices: tt.fields.devices,
// 			}
// 			if err := br.processCommand(tt.args.gcmd, tt.args.d, tt.args.data); (err != nil) != tt.wantErr {
// 				t.Errorf("Bridge.processCommand() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
