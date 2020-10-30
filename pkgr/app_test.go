package pkgr

import (
	"reflect"
	"testing"
)

func Test_appinfo(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name  string
		args  args
		wantA App
	}{
        // TODO: Add test cases.
         
		{
			args: args{
				appid: "2890",
			},
			wantA: App{
				Appid: "Splunk_ML_Toolkit",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotA := appinfo(tt.args.id); !reflect.DeepEqual(gotA, tt.wantA) {
				t.Errorf("appinfo() = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}
