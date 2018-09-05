package main

import (
	"math"
	"math/cmplx"
)

func dft(input []float64) []complex128 {
	output := make([]complex128, len(input))

	arg := -2.0 * math.Pi / float64(len(input))
	for k := 0; k < len(input); k++ {
		r, i := 0.0, 0.0
		for n := 0; n < len(input); n++ {
			r += input[n] * math.Cos(arg*float64(n)*float64(k))
			i += input[n] * math.Sin(arg*float64(n)*float64(k))
		}
		output[k] = complex(r, i)
	}
	return output
}

func hfft(samples []float64, freqs []complex128, n, step int) {
	if n == 1 {
		freqs[0] = complex(samples[0], 0)
		return
	}

	half := n / 2

	hfft(samples, freqs, half, 2*step)
	hfft(samples[step:], freqs[half:], half, 2*step)

	for k := 0; k < half; k++ {
		a := -2 * math.Pi * float64(k) / float64(n)
		e := cmplx.Rect(1, a) * freqs[k+half]

		freqs[k], freqs[k+half] = freqs[k]+e, freqs[k]-e
	}
}

func fft(samples []float64) []complex128 {
	n := len(samples)
	freqs := make([]complex128, n)
	hfft(samples, freqs, n, 1)
	return freqs
}
