package biquad

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
)

func (b blt) getH() func(z complex128) complex128 {
	b0d, b1d, b2d := complex(b.b0d, 0), complex(b.b1d, 0), complex(b.b2d, 0)
	a1d, a2d := complex(b.a1d, 0), complex(b.a2d, 0)
	return func(z complex128) complex128 {
		return (b0d + b1d/z + b2d/(z*z)) / (1 + a1d/z + a2d/(z*z))
	}
}

// Plot the BodePlot of a discrete time transfer function H and the sampling period.
func plotBode(plotname string, ts float64, H func(z complex128) complex128) {
	p := plot.New()
	xy := noNaN{bodeFunc{H, ts}}
	bode, err := plotter.NewLine(xy)
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Bode Plot"
	p.Y.Label.Text = "real( H( e^{jwT} ) )"
	p.X.Label.Text = "Frequency [rad/s]"
	p.Add(bode)
	p.X.Scale = plot.LogScale{}
	p.Save(30*font.Centimeter, 20*font.Centimeter, plotname)
}

// Plot discrete time systems from H(z)
type bodeFunc struct {
	H  func(complex128) complex128
	ts float64 // sampling period
}

func (b bodeFunc) Len() int { return 400 }

// z -> e^{jwT} == cos(wT) + j*sin(wT)
func (b bodeFunc) XY(i int) (x, y float64) {
	w := float64(i) + 0.1
	x = math.Pow(w, 2)
	y = real(b.H(complex(math.Cos(w*b.ts), math.Sin(w*b.ts))))
	return x, y
}

type noNaN struct {
	plotter.XYer
}

func (n noNaN) XY(i int) (x, y float64) {
	x, y = n.XYer.XY(i)
	if y != y {
		y = 0.1
	}
	return x, y
}
