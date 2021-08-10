package biquad

import (
	"math"
	"strconv"
)

// Different methods of calculating alpha.
// see http://shepazu.github.io/Audio-EQ-Cookbook/audio-eq-cookbook.html
type alphaCalc struct{}

// the bandwidth in octaves (between -3 dB frequencies
// for BPF and notch or between midpoint (dBgain/2) gain frequencies for peaking EQ)
func (a alphaCalc) bw(w0, BW float64) (alpha float64) {
	sin := math.Sin(w0)
	sharg := math.Ln2 / 2 * BW * w0 / sin
	if sharg <= -1 || sharg >= 1 {
		panic("bad arguments to BW alpha calculation. got ln2/2*bw*w0/sn == " + strconv.FormatFloat(sharg, 'e', 6, 64))
	}
	return sin * math.Sinh(sharg)
}

// the EE kind of definition, except for peakingEQ in which A*Q is the classic EE Q.
// That adjustment in definition was made so that a boost of N dB followed by a cut of
// N dB for identical Q and f0/Fs results in a precisely flat unity gain filter or "wire".
func (a alphaCalc) q(w0, Q float64) (alpha float64) {
	return math.Sin(w0) / (2 * Q)
}

// a "shelf slope" parameter (for shelving EQ only). When S = 1, the shelf slope is as steep
// as it can be and remain monotonically increasing or decreasing gain with frequency. The shelf slope,
// in dB/octave, remains proportional to S for all other values for a fixed f0/Fs and dBgain
func (a alphaCalc) s(w0, A, S float64) (alpha float64) {
	return math.Sin(w0) / 2 * math.Sqrt((A+1/A)*(1/S-1)+2)
}

// A = sqrt(10^(DBgain/20))
// 1/Q = 2 * sinh(ln2/2 * BW * w0 / sin(w0))
