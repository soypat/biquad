package biquad

import "testing"

func TestBLTBode(t *testing.T) {
	const (
		fs = 100.
		ts = 1 / fs
	)
	lp, err := NewNotch(fs, 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	H := lp.getH()
	plotBode("notch_f0=4.png", ts, H)
}

func TestButterBode(t *testing.T) {
	//https://www.robots.ox.ac.uk/~sjrob/Teaching/SP/l6.pdf
	const (
		fs = 100.
		ts = 1 / fs
	)
	// Design a digital low-pass Butterworth filter with a 3dB cut-off frequency of 2kHz
	// and minimum attenuation of 30dB at 4.25kHz for a sampling rate of 10kHz
	b, err := NewButterworthLP(1e4, 2e2)
	if err != nil {
		t.Fatal(err)
	}
	H := b.getH()
	plotBode("butter_wc=2k.png", ts, H)
}
