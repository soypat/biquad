package biquad

import "math"

type HighPass struct {
	blt
}

// NewHighPass creates a high pass filter from
//  Fs: sampling frequency
//  f0: working frequency
//  BW: bandwidth of the filter in octaves
func NewHighPass(Fs, f0, BW float64) (*HighPass, error) {
	switch {
	case f0 >= Fs:
		return nil, ErrBadWorkingFreq
	case BW <= 0:
		return nil, ErrNegBandwidth
	case f0 <= 0 || Fs <= 0:
		return nil, ErrBadFreq
	}
	w0 := 2 * math.Pi * (f0 / Fs)
	cos := math.Cos(w0)
	alpha := alphaCalc{}.bw(w0, BW)
	var (
		b0 = (1 + cos) / 2
		b1 = -(1 + cos)
		b2 = b0
		a0 = 1 + alpha
		a1 = -2 * cos
		a2 = 1 - alpha
	)
	return &HighPass{
		blt: newBLT(a0, a1, a2, b0, b1, b2),
	}, nil
}
