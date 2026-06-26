package ping

import (
	"context"
	"errors"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// Result holds the outcome of a single ping.
type Result struct {
	Seq int
	RTT time.Duration // zero on timeout/error
	Err error
	At  time.Time
}

// Pinger sends ICMP echo requests in a loop and reports results on a channel.
type Pinger struct {
	target   string
	interval time.Duration
}

// New creates a Pinger for the given target at the given interval.
func New(target string, interval time.Duration) *Pinger {
	return &Pinger{target: target, interval: interval}
}

// Start begins pinging in a goroutine and returns a receive-only channel of Results.
// The channel is closed when ctx is cancelled.
func (p *Pinger) Start(ctx context.Context) <-chan Result {
	ch := make(chan Result, 1)

	go func() {
		defer close(ch)

		pinger, err := probing.NewPinger(p.target)
		if err != nil {
			ch <- Result{Seq: 0, Err: err, At: time.Now()}
			return
		}

		pinger.Interval = p.interval
		pinger.Count = -1           // run forever until ctx cancels
		pinger.SetPrivileged(false) // SOCK_DGRAM — no sudo needed on macOS

		// OnSend and OnRecv fire from different goroutines inside pro-bing,
		// so we need a mutex to guard the pending map.
		var mu sync.Mutex
		pending := make(map[int]time.Time) // seq -> send time

		pinger.OnSend = func(pkt *probing.Packet) {
			mu.Lock()
			defer mu.Unlock()
			// Flush any pending pings older than the interval — they timed out.
			for seq, t := range pending {
				if time.Since(t) > p.interval {
					ch <- Result{Seq: seq, Err: errors.New("timeout"), At: t}
					delete(pending, seq)
				}
			}
			pending[pkt.Seq] = time.Now()
		}

		pinger.OnRecv = func(pkt *probing.Packet) {
			mu.Lock()
			delete(pending, pkt.Seq)
			mu.Unlock()
			ch <- Result{Seq: pkt.Seq, RTT: pkt.Rtt, At: time.Now()}
		}

		pinger.RunWithContext(ctx)
	}()

	return ch
}
