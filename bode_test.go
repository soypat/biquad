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
