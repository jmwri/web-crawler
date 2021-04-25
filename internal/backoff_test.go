package internal_test

import (
	"github.com/jmwri/web-crawler/internal"
	"testing"
	"time"
)

func TestSimpleBackoff(t *testing.T) {
	type args struct {
		attempt int
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "attempt 1",
			args: args{attempt: 1},
			want: 0,
		},
		{
			name: "attempt 2",
			args: args{attempt: 2},
			want: time.Millisecond * 500,
		},
		{
			name: "attempt 11",
			args: args{attempt: 11},
			want: time.Second * 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := internal.SimpleBackoff(tt.args.attempt); got != tt.want {
				t.Errorf("SimpleBackoff() = %v, want %v", got, tt.want)
			}
		})
	}
}
