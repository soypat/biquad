package biquad

import (
	"math"
)

// Butterworth filter. May represent Lowpass and other
// Butterworth filter implementations.
type ButterWorth struct {
	blt
}

// https://www.robots.ox.ac.uk/~sjrob/Teaching/SP/l6.pdf
// Butterworth calculated according to bilinear transformation s = 2/Ts * (1-z^{-1})/(1+z^{-1})
// where H(s) (analog) is
//  H(s) = wc^2 / (s^2 + s*sqrt(2)*wc + wc^2)    [1]
// Ts is the sampling period.
// wc is the analog cutoff frequency in radians/s. To obtain it
// from the digital cutoff frequency we must pre-warp the digital
// frequency:
//  wc_a = 2/Ts * tan(wc_d * Ts / 2)
// Applying the bilinear transformation to [1]:
//  H(z) = td^2 / ( ((1-z^-1)/(1+z^-1))^2 + td*sqrt(2)*((1-z^-1)/(1+z^-1)) + td^2)   [2]
// where td = tan(wc_d * Ts / 2)
// Note: (2/Ts)^2 is canceled out.
// We now multiply denominator and numerator by (1+z^{-1})^2
//  H(z) = td*(1 + 2*z^{-1} + z^{-2}) / ( (1-z^{-1})^2 + td*sqrt(2)*(1-z^{-1})*(1+z^{-1}) + td^2*(1+z^{-1})^2
//  H(z) denominator = 1 - 2z^{-1} + z^{-2} + td*sqrt(2)*(1 - z^{-2}) + td^2 * (1+2z^{-1} + z^{-2})
//   grouping...     = 1 + td*sqrt(2) + td^2 + z^{-1}*( -2 + 2*td^2 )  + z^{-2} * (1 - td*sqrt(2) + td^2)
// The coefficients for BLT application are now readily available from the work above.

// NewButterworthLP creates a low pass Butterworth filter from
//  Fs: sampling frequency
//  fc: cutoff frequency
// Not guaranteed to have peak unity gain.
func NewButterworthLP(Fs, fc float64) (*ButterWorth, error) {
	switch {
	case fc >= Fs:
		return nil, ErrBadWorkingFreq
	case fc <= 0 || Fs <= 0:
		return nil, ErrBadFreq
	}
	// digital cutoff frequency
	wc := 2 * math.Pi * fc
	td := math.Tan(wc / (2 * Fs))
	var (
		b0 = td
		b1 = 2 * td
		b2 = b0
		a0 = 1 + td*math.Sqrt2 + td*td
		a1 = -2 + 2*td*td
		a2 = 1 - td*math.Sqrt2 + td*td
	)
	return &ButterWorth{
		blt: newBLT(a0, a1, a2, b0, b1, b2),
	}, nil
}
