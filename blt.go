package biquad

// Bilinear transform for filter use.
// Comes from the following biquad transfer function (see http://shepazu.github.io/Audio-EQ-Cookbook/audio-eq-cookbook.html):
//  H(z) = (b_0 + b_1*z^{-1} + b_2*z^{-2}) / (a_0 + a_1*z^{-1} + a_2*z^{-2})
type blt struct {
	// 5 coefficients normalized respect a0.
	b0d, b1d, b2d, a1d, a2d float64
	// Circular buffers for state storage. x is measured signal. y is filter result.
	x, y [3]float64
	// points to `n` index in ring buffer.
	ptr uint
}

//  H(z) = (b_0 + b_1*z^{-1} + b_2*z^{-2}) / (a_0 + a_1*z^{-1} + a_2*z^{-2})
func newBLT(a0, a1, a2, b0, b1, b2 float64) blt {
	if a0 == 0 {
		panic("a0 can not be 0")
	}
	return blt{
		a1d: a1 / a0,
		a2d: a2 / a0,
		b0d: b0 / a0,
		b1d: b1 / a0,
		b2d: b2 / a0,
		ptr: 3,
	}
}

// simplest implementation of BLT filter using biquad transfer function
func (b *blt) advance(x float64) {
	var (
		n   = b.ptr % 3
		nm1 = (b.ptr - 1) % 3
		nm2 = (b.ptr - 2) % 3
	)
	b.x[n] = x // Save sample
	b.y[n] = b.b0d*x + b.b1d*b.x[nm1] + b.b2d*b.x[nm2] -
		b.a1d*b.y[nm1] - b.a2d*b.y[nm2] // Save filtered value.
	// adding one to b.ptr shifts values.
	b.ptr++
}

func (b *blt) ynext() float64 {
	return b.y[b.ptr%3]
}

func (b *blt) init(xy Signal) {
	x, y := xy.XY(0)
	b.x[0] = x
	b.x[1] = x
	b.x[2] = x

	b.y[0] = y
	b.y[1] = y
	b.y[2] = y
}

// Filter applies a bilinear transformation filter to a digital
// signal and returns the filtered result. The length of the data must be greater than 2.
func (b *blt) Filter(signal Signal) (Signal, error) {
	var x float64
	N := signal.Len()
	if N < 3 {
		return nil, ErrShortXY
	}
	fval := make([]float64, N)
	b.init(signal)
	for i := 0; i < N; i++ {
		_, x = signal.XY(i)
		b.advance(x)
		fval[i] = b.ynext()
	}
	return filtered{
		Signal: signal,
		fval:   fval,
	}, nil
}

// DiscreteProcess takes in the next signal data point
// and processes it. DiscreteProcess expects data points
// to be evenly spaced out in time.
func (b *blt) DiscreteProcess(x float64) {
	b.advance(x)
}

// YNext returns the last result of the filter given by
// DiscreteProcess.
func (b *blt) YNext() (y float64) {
	return b.ynext()
}
