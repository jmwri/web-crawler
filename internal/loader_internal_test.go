package internal

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewHTTPGetLoader(t *testing.T) {
	type args struct {
		client *http.Client
	}
	tests := []struct {
		name string
		args args
		want *HTTPGetLoader
	}{
		{
			name: "creates successfully",
			args: args{client: http.DefaultClient},
			want: &HTTPGetLoader{http.DefaultClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPGetLoader(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPGetLoader() = %v, want %v", got, tt.want)
			}
		})
	}
}
