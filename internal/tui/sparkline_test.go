package tui

import (
	"testing"
	"time"
)

func ms(n int) time.Duration {
	return time.Duration(n) * time.Millisecond
}

func TestRender(t *testing.T) {
	tests := []struct {
		name   string
		values []time.Duration
		width  int
		want   string
	}{
		{
			name:   "empty values",
			values: nil,
			width:  10,
			want:   "",
		},
		{
			name:   "zero width",
			values: []time.Duration{ms(10)},
			width:  0,
			want:   "",
		},
		{
			name:   "single value",
			values: []time.Duration{ms(10)},
			width:  10,
			want:   "▅", // single value -> middle block (index 4 of 8)
		},
		{
			name:   "all same values",
			values: []time.Duration{ms(50), ms(50), ms(50)},
			width:  10,
			want:   "▅▅▅",
		},
		{
			name:   "min and max only",
			values: []time.Duration{ms(10), ms(100)},
			width:  10,
			want:   "▁█",
		},
		{
			name:   "timeouts rendered as lowest block",
			values: []time.Duration{ms(10), 0, ms(100)},
			width:  10,
			want:   "▁▁█",
		},
		{
			name:   "width truncates to tail",
			values: []time.Duration{ms(100), ms(10), ms(50)},
			width:  2,
			want:   "▁█",
		},
		{
			name:   "all timeouts",
			values: []time.Duration{0, 0, 0},
			width:  5,
			want:   "▁▁▁",
		},
		{
			name:   "ascending values",
			values: []time.Duration{ms(10), ms(30), ms(50), ms(70), ms(90)},
			width:  5,
			want:   "▁▂▄▆█",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Render(tt.values, tt.width)
			if got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}
