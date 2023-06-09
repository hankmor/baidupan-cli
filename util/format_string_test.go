package util

import (
	"fmt"
	"testing"
)

func TestConvReadableSize(t *testing.T) {
	type args struct {
		size int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "< 0", args: args{size: -1}, want: ""},
		{name: "= 0", args: args{size: 0}, want: ""},
		{name: "< 1KB", args: args{size: 1000}, want: "1000B"},
		{name: "1MB > x > 1KB", args: args{size: 1000 * KB}, want: fmt.Sprintf("%dKB", 1000)},
		{name: "1GB > x > 1MB", args: args{size: 1000 * MB}, want: fmt.Sprintf("%dMB", 1000)},
		{name: "1TB > x > 1GB", args: args{size: 1000 * GB}, want: fmt.Sprintf("%dGB", 1000)},
		{name: "x > 1TB", args: args{size: 1000 * TB}, want: fmt.Sprintf("%dTB", 1000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvReadableSize(tt.args.size); got != tt.want {
				t.Errorf("ConvReadableSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
