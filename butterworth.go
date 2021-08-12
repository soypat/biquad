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
//  fc: cutoff frequency of approximately -3dB gain
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

//https://tttapa.github.io/Pages/Mathematics/Systems-and-Control-Theory/Analog-Filters/Butterworth-Filters.html
// Normalized Lowpass butterworth filters defined in frequency domain:
//  H(jw) = 1 / sqrt(1+w^{2*n})
// where n is the filter order.
// In order to determine the transfer function, we'll start from the frequency response squared.
//  |H(jw)|^2 = H(jw) * conj(H(jw)) = H(jw) * H(-jw) = 1 / (1 + w^{2*n})
// We use identity  s = jw -> w = s/j
//  H(jw) * H(-jw) = H(s) * H(-s) =  1 / (1 + (s/j)^{2*n})
// ...
// Normalized Butterworth functions
//  Hlp(s) = 1/Bn(s/wc)   Hhp(s) = s^n / wc^n * Bn(s/wc)
// where Bn is
//  Bn(s) = PROD^{n/2-1}_{k=0}  s^2 - 2*s*cos(2*pi*(2*k+n+1)/(4*n)) + 1
// Hlp(s) is calculated above. Calculating Hhp(s) now:
//  Hhp(s) = s^2 / (s^2 + s*sqrt(2)*wc + wc^2)
//  wc_a = 2/Ts * tan(wc_d * Ts / 2)
// Applying the bilinear transformation:
//  H(z) = ((1-z^-1)/(1+z^-1))^2  / ( ((1-z^-1)/(1+z^-1))^2 + td*sqrt(2)*((1-z^-1)/(1+z^-1)) + td^2)   [2]
// Multiply top and bottom by (1+z^-1)^2 and group
//  H(z) = (1 - 2*z^{-1} + z^{-2}) / ( 1+td*sqrt(2)+td^2 + z^{-1}*(-2+td^2) + z^{-2}*(1-td*sqrt(2)+td^2) )
// The BLT coefficients are then readily available above for a High pass Butterworth filter.

// NewButterworthHP creates a high pass Butterworth filter from
//  Fs: sampling frequency
//  fc: cutoff frequency of approximately -3dB gain
// Not guaranteed to have peak unity gain.
func NewButterworthHP(Fs, fc float64) (*ButterWorth, error) {
	switch {
	case fc >= Fs:
		return nil, ErrBadWorkingFreq
	case fc <= 0 || Fs <= 0:
		return nil, ErrBadFreq
	}
	fc *= 1
	// digital cutoff frequency
	wc := 2 * math.Pi * fc
	td := math.Tan(wc / (2 * Fs))
	var (
		b0 = 1.
		b1 = -2.
		b2 = b0
		a0 = 1 + td*math.Sqrt2 + td*td // 1+td*sqrt(2)+td^2
		a1 = -2 + 2*td*td
		a2 = 1 - td*math.Sqrt2 + td*td
	)
	return &ButterWorth{
		blt: newBLT(a0, a1, a2, b0, b1, b2),
	}, nil
}
