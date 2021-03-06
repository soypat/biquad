package biquad_test

import (
	"math"

	"github.com/soypat/biquad"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
)

func ExampleLowPass() {
	const (
		pi = math.Pi
		// working frequency. The frequency that matters (we wish to keep it)
		f0     = 2.
		fnoise = 4.
		fs     = 100.   // Sampling frequency.
		N      = 100    // amount of sample points
		ts     = 1 / fs // sampling period
	)
	// We generate waveform data composed of a "working" frequency and a "noise" frequency
	data := make([]float64, N)
	for i := 0; i < N; i++ {
		t := float64(i) * ts
		data[i] = math.Sin(2*pi*f0*t) + math.Sin(2*pi*fnoise*t)
	}
	signal := biquad.MakeSignal(fs, data)
	lp, err := biquad.NewLowPass(fs, f0, 1)
	if err != nil {
		panic(err)
	}

	filtered, err := lp.Filter(signal)
	if err != nil {
		panic(err)
	}
	p := plot.New()
	ls, _ := plotter.NewLine(signal)
	lf, _ := plotter.NewLine(filtered)
	lf.Dashes = []font.Length{0.1 * font.Centimeter, 0.1 * font.Centimeter}
	p.Add(ls, lf)
	p.Legend.Add("original signal", ls)
	p.Legend.Add("filtered signal", lf)
	err = p.Save(30*font.Centimeter, 15*font.Centimeter, "out.png")
	if err != nil {
		panic(err)
	}
}
