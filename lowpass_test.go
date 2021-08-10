package biquad

import (
	"encoding/csv"
	"image/color"
	"os"
	"strconv"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
)

func TestFilter(t *testing.T) {
	const N = 12000
	fp, err := os.Open("testdata/noisy_peak.csv")
	if err != nil {
		t.Fatal(err)
	}
	r := csv.NewReader(fp)
	r.FieldsPerRecord = 3
	r.Read() // skip header
	sl1 := make([]float64, N)
	sl2 := make([]float64, N)
	i := 0
	for i < N {
		rec, err := r.Read()
		if err != nil {
			break
		}
		sl1[i], err = strconv.ParseFloat(rec[1], 32)
		if err != nil {
			t.Fatal(err)
		}
		sl2[i], err = strconv.ParseFloat(rec[2], 32)
		if err != nil {
			t.Fatal(err)
		}
		if sl2[i] != sl2[i] || sl1[i] != sl1[i] { // is nan
			t.Errorf("float parsed NaN on line %d", i+1)
			return
		}
		i++
	}
	const (
		Ts = 40. / float64(N) // Sampling Period
		Fs = 1 / Ts           // Sampling frequency
		f0 = Fs / 100.        // Working frequency
		T0 = 1 / f0           // Working period
	)
	s1 := &dData{ts: Ts, v: sl1}
	s2 := &dData{ts: Ts, v: sl2}

	lp, err := NewLowPass(Fs, f0, 1)
	if err != nil {
		t.Fatal(err)
	}
	s1filtered, err := lp.Filter(s1)
	if err != nil {
		t.Fatal(err)
	}
	lp.Filter(s2)

	p := plot.New()
	p.Legend = plot.NewLegend()
	// original signal.
	l1, err := plotter.NewLine(s1)
	l1.Color = color.RGBA{R: 255, A: 255}
	if err != nil {
		t.Error(err)
	}
	// filtered signal
	l1f, err := plotter.NewLine(s1filtered)
	if err != nil {
		t.Error(err)
	}

	p.Add(l1, l1f)
	p.Legend.Add("Original signal [red]", l1)
	p.Legend.Add("Filtered Signal", l1f)
	cm := font.Centimeter
	p.X.Max = 5
	p.Y.Max = -350
	p.Y.Min = -1200
	err = p.Save(40*cm, 20*cm, "testdata/out.svg")
	if err != nil {
		t.Error(err)
	}
}

// evenly spaced data
type dData struct {
	// sampling period
	ts float64
	v  []float64
}

func (d *dData) XY(i int) (t, y float64) {
	return d.ts * float64(i), d.v[i]
}
func (d *dData) Len() int { return len(d.v) }
