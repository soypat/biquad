package biquad

import (
	"math"
	"testing"
)

func TestBLTBode(t *testing.T) {
	const (
		fs = 100.
		ts = 1 / fs
		f0 = 4
		w0 = math.Pi * 2 * f0
	)
	lp, err := NewNotch(fs, f0, 1)
	if err != nil {
		t.Fatal(err)
	}
	H := lp.getH()
	plotBode("notch_f0=4.png", ts, f0*60, H)
}

func TestButterBode(t *testing.T) {
	//https://www.robots.ox.ac.uk/~sjrob/Teaching/SP/l6.pdf
	const (
		fs = 1e4
		ts = 1 / fs
		f0 = 2e3
		w0 = math.Pi * 2 * f0
	)
	// Design a digital low-pass Butterworth filter with a 3dB cut-off frequency of 2kHz
	// and minimum attenuation of 30dB at 4.25kHz for a sampling rate of 10kHz
	b, err := NewButterworthLP(fs, f0)
	if err != nil {
		t.Fatal(err)
	}
	H := b.getH()
	plotBode("butter_wc=2k.png", ts, f0*60, H)
}

func TestChebyBode(t *testing.T) {
	const (
		fs = 1e4
		ts = 1 / fs
		f0 = 2e2
	)
	b, err := newChebyshevLPType1(fs, f0, 0.1)
	if err != nil {
		t.Fatal(err)
	}
	H := b.getH()
	plotBode("cheby1_wc=2k.png", ts, f0*60, H)
}
