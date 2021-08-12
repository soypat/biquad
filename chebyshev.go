package biquad

import (
	"math"
)

// Chebyshev filter. May represent Lowpass and other
// Chebyshev filter implementations.
type Chebyshev struct {
	blt
}

// From wikipedia:
// Chebyshev type I Transfer function:
//  H(s) = (1/(2*e)) / (  s^2 - s*(sp_1+sp_2) + sp_1*sp_2)
// where sp_m = j*cosh(asinh(1/e)/n)*cos(theta_m) - sinh(asinh(1/e)/n)*sin(theta_m)
// where theta_m = pi*(2*m - 1) / (2*n)
// We apply bilinear transformation as seen in butterworth file (and take into account non unity frequncy cutoff wc)
//  H(z) = (wc^2/ (2*e) * (1+2z^-1+z^-2) / (1-2z^{-1}+z^{-2} - wc*(sp1+sp2)*(1-z^{-2}) +wc^2*sp1*sp2*(1+2*z^{-1}+z^{-2}))
// where wc is the pre-warped (analog) cutoff frequency.
//  H(z)_denominator = 1-wc(sp1+sp2)+wc^2*sp1*sp2 + z^{-1}*(-2+2*wc^2*sp1*sp2) + z^{-2}*(1+wc*(sp1+sp2)+wc^2*sp1*sp2)
// We may now construct the BLT after applying frequency warp correction.

// NewChebyshevLPType1 creates a low pass Chebyshev Type I filter from
//  Fs: sampling frequency
//  fh: -3dB attenuation frequency. This frequency is lower than the cutoff in Chebyshev lowpass filters.
// Not guaranteed to have peak unity gain.
func newChebyshevLPType1(Fs, fh, e float64) (*Chebyshev, error) {
	const n = 2. // Filter order.
	switch {
	case fh >= Fs:
		return nil, ErrBadWorkingFreq
	case fh <= 0 || Fs <= 0:
		return nil, ErrBadFreq
	}
	// The -3dB stopband start wh is related to wc by:
	// wh = wc * cosh(acosh(1/e)/n)
	// thus wc = wh / cosh(acosh(1/e)/n)
	wc := 2 * math.Pi * fh / math.Cosh(math.Acosh(1/e)/n)
	td := math.Tan(wc / (2 * Fs)) // wc_a = 2/Ts * tan(wc_d * Ts / 2)
	sp1 := sparametric1(2, 1, true, e)
	sp2 := sparametric1(2, 2, true, e)
	tdc := complex(td, 0)
	var ( // TODO fix whatever dont work here
		b0 = td * td / (2 * e)
		b1 = 2 * b0
		b2 = b0
		a0 = real(1 - tdc*sp1 + sp2 + tdc*tdc*sp1*sp2)
		a1 = real(-2 + 2*tdc*tdc*sp1*sp2)
		a2 = real(1 + tdc*(sp1+sp2) + tdc*tdc*sp1*sp2)
	)
	return &Chebyshev{
		blt: newBLT(a0, a1, a2, b0, b1, b2), // Do I need a new complex valued BLT type?
	}, nil
}

// calculate parametric pole of chebyshev type 1 transfer function.
func sparametric1(n, m int, neg bool, e float64) complex128 {
	N := float64(n)
	M := float64(m)
	tm := math.Pi * (2*M - 1) / (2 * N)
	imagpart := math.Cosh(math.Asinh(1/e)/N) * math.Cos(tm)
	realpart := math.Sinh(math.Asinh(1/e)/N) * math.Sin(tm)
	if neg {
		realpart *= -1
	}
	return complex(realpart, imagpart)
}
