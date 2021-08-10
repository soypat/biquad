package biquad

import "errors"

var (
	ErrShortXY        = errors.New("XYer length must be greater than 2 to apply BLT filter")
	ErrBadWorkingFreq = errors.New("working frequency can not be higher than sampling frequency")
	ErrNegBandwidth   = errors.New("bandwidth must be greater than zero")
	ErrBadFreq        = errors.New("zero or negative frequency")
	ErrBadGain        = errors.New("negative or zero gain")
)
