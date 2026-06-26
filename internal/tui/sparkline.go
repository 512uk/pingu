package tui

import "time"

// Unicode block elements, 8 levels from lowest to highest.
var blocks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// Render produces a sparkline string from a slice of durations.
// width controls how many of the most recent values are shown.
// Zero-duration entries (timeouts) are rendered as the lowest block.
func Render(values []time.Duration, width int) string {
	if len(values) == 0 || width <= 0 {
		return ""
	}

	// Take only the tail that fits in width.
	if len(values) > width {
		values = values[len(values)-width:]
	}

	// Find min/max of the visible window (ignoring zero/timeout values for scaling).
	var min, max time.Duration
	hasValid := false
	for _, v := range values {
		if v <= 0 {
			continue
		}
		if !hasValid || v < min {
			min = v
		}
		if !hasValid || v > max {
			max = v
		}
		hasValid = true
	}

	out := make([]rune, len(values))
	for i, v := range values {
		if v <= 0 || !hasValid {
			out[i] = blocks[0]
			continue
		}
		out[i] = blockFor(v, min, max)
	}

	return string(out)
}

func blockFor(v, min, max time.Duration) rune {
	if max == min {
		return blocks[len(blocks)/2] // all same value -> middle block
	}

	// Scale to 0..7
	ratio := float64(v-min) / float64(max-min)
	idx := int(ratio * float64(len(blocks)-1))
	if idx >= len(blocks) {
		idx = len(blocks) - 1
	}
	return blocks[idx]
}
