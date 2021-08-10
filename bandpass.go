package biquad

import "math"

type BandPass struct {
	blt
}

// NewBandPass creates a band pass filter from
//  Fs: sampling frequency
//  f0: working frequency
//  BW: bandwidth of the filter in octaves
// This filter has constant peak gain (0 dB).
func NewBandPass(Fs, f0, BW float64) (*BandPass, error) {
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
		b0 = alpha
		b1 = 0.
		b2 = -b0
		a0 = 1 + alpha
		a1 = -2 * cos
		a2 = 1 - alpha
	)
	return &BandPass{
		blt: newBLT(a0, a1, a2, b0, b1, b2),
	}, nil
}

// NewBandPass creates a band pass filter from
//  Fs: sampling frequency
//  Q: peak gain
//  BW: bandwidth of the filter in octaves
func NewBandPassFromQ(Fs, Q, BW float64) (*BandPass, error) {
	switch {
	case BW <= 0:
		return nil, ErrNegBandwidth
	case Q <= 0:
		return nil, ErrBadGain
	}
	// c = w0/sin(w0).
	c := 2 * math.Asinh(1/(2*Q)) / (math.Ln2 * BW)
	// apply two newton iterations to obtain w0 starting at sampling frequency.
	// f(x)  = c*sin(x) - x
	// f'(x) = c*cos(x) - 1
	w0 := Fs
	for i := 0; i < 2; i++ {
		w0 -= (c*math.Sin(w0) - w0) / (c*math.Cos(w0) - 1)
	}

	sin, cos := math.Sincos(w0)
	alpha := alphaCalc{}.bw(w0, BW)
	var (
		b0 = sin / 2
		b1 = 0.
		b2 = -b0
		a0 = 1 + alpha
		a1 = -2 * cos
		a2 = 1 - alpha
	)
	return &BandPass{
		blt: newBLT(a0, a1, a2, b0, b1, b2),
	}, nil
}
