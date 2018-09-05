package main

import (
	"image/draw"
	"math/cmplx"
)

func drawfft(img draw.Image, samples []float64, rate, bins uint32) {
	bn := img.Bounds()

	/*
		bg2 := ParseColor(*BG2)
		fmt.Printf("bg2: %.8x\n", bg2)
		draw.Draw(img, bn, image.NewUniform(bg2), image.ZP, draw.Src)
	*/

	gr := New()
	gr.Append(ParseColor("000000"))
	gr.Append(ParseColor("380F6D"))
	gr.Append(ParseColor("B63679"))
	gr.Append(ParseColor("FD9A69"))
	gr.Append(ParseColor("FCF6B8"))

	for x := 1; x < bn.Dx(); x++ {
		n0 := int64(mapRange(float64(x-1), 0, float64(bn.Dx()), 0, float64(len(samples))))
		n1 := int64(mapRange(float64(x-0), 0, float64(bn.Dx()), 0, float64(len(samples))))
		c := n0 + (n1-n0)/2

		sub := make([]float64, bins*2)
		for i := 0; i < len(sub); i += 1 {
			s := 0.0
			n := int(c) - int(bins) + int(i)
			if n >= 0 && n < len(samples) {
				s = samples[n]
			}
			sub[i] = s
		}

		var freqs []complex128
		if *DFT {
			freqs = dft(sub)
		} else {
			freqs = fft(sub)
		}
		for y := 0; y < int(bins); y++ {
			c := freqs[y/2]
			r := cmplx.Abs(c)
			//r = 30 * math.Log10(math.Abs(r)/10e-6)
			//r /= float64(bins)
			img.Set(x, int(bins)-y, gr.ColorAt(r))
		}
	}
}
