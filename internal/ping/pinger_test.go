package ping

import (
	"errors"
	"testing"
	"time"
)

func TestResultFields(t *testing.T) {
	tests := []struct {
		name    string
		result  Result
		wantErr bool
		wantRTT bool
	}{
		{
			name:    "successful ping",
			result:  Result{Seq: 1, RTT: 12 * time.Millisecond, At: time.Now()},
			wantErr: false,
			wantRTT: true,
		},
		{
			name:    "timeout (zero RTT with error)",
			result:  Result{Seq: 2, RTT: 0, Err: errors.New("i/o timeout"), At: time.Now()},
			wantErr: true,
			wantRTT: false,
		},
		{
			name:    "error with no RTT",
			result:  Result{Seq: 3, RTT: 0, Err: errors.New("network unreachable"), At: time.Now()},
			wantErr: true,
			wantRTT: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.result.Err != nil) != tt.wantErr {
				t.Errorf("Err = %v, wantErr = %v", tt.result.Err, tt.wantErr)
			}
			if (tt.result.RTT > 0) != tt.wantRTT {
				t.Errorf("RTT = %v, wantRTT = %v", tt.result.RTT, tt.wantRTT)
			}
		})
	}
}

func TestNewPinger(t *testing.T) {
	p := New("8.8.8.8", 500*time.Millisecond)
	if p.target != "8.8.8.8" {
		t.Errorf("target = %q, want %q", p.target, "8.8.8.8")
	}
	if p.interval != 500*time.Millisecond {
		t.Errorf("interval = %v, want %v", p.interval, 500*time.Millisecond)
	}
}
