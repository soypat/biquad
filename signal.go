package biquad

// Signal represents a stored data representing a digital signal.
type Signal interface {
	// Should return the length of the data.
	Len() int
	// Returns the i-th data point (time and signal data).
	XY(int) (t, y float64)
}

// XYer wraps the Len and XY methods.
type filtered struct {
	// Original Signal.
	Signal
	// filtered values.
	fval []float64
}

func (f filtered) XY(i int) (t, y float64) {
	t, _ = f.Signal.XY(i)
	return t, f.fval[i]
}

// MakeSignal allocates a new signal interface from data and a sampling frequency.
func MakeSignal(Fs float64, data []float64) Signal {
	if len(data) == 0 || Fs <= 0 {
		return &signal{}
	}
	vals := make([]float64, len(data))
	copy(vals, data)
	return &signal{data: vals, ts: 1 / Fs}
}

type signal struct {
	// Sampling period
	ts   float64
	data []float64
}

func (s signal) XY(i int) (t, y float64) {
	if i > len(s.data) {
		return
	}
	return s.ts * float64(i), s.data[i]
}

func (s signal) Len() int {
	return len(s.data)
}
