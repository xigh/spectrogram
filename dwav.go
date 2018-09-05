package main

import (
	"fmt"
	"image"
	"image/draw"
	"math"
)

func drawwav(img draw.Image, samples []float64) {
	bn := img.Bounds()

	bg1 := ParseColor(*BG1)
	fmt.Printf("bg1: %.8x\n", bg1)
	draw.Draw(img, bn, image.NewUniform(bg1), image.ZP, draw.Src)

	fg0 := ParseColor(*FG0)
	fmt.Printf("fg0: %.8x\n", fg0)
	fg1 := ParseColor(*FG1)
	fmt.Printf("fg1: %.8x\n", fg1)
	fg2 := ParseColor(*RUL)
	fmt.Printf("fg2: %.8x\n", fg2)

	gsum, gmin, gmax := 0.0, 0.0, 0.0
	for i := 0; i < len(samples); i++ {
		gsmp := samples[i]
		gsum += gsmp
		gmax = math.Max(gsmp, gmax)
		gmin = math.Min(gsmp, gmin)
	}
	gavg := gsum / float64(len(samples))
	gabs := math.Max(math.Abs(gmin), math.Abs(gmax))

	fmt.Printf("samples: %d [max: %.2f, min: %.2f, avg: %.2f]\n", len(samples), gmax, gmin, gavg)

	// -------------------------------------------

	middle := bn.Dy() / 2

	if !*HIDERULERS {
		drawLine(img, 0, middle, bn.Dx(), middle, fg2)
		drawLine(img, 0, 0, 0, bn.Dy(), fg2)
	}

	for i := 1; i < bn.Dx(); i += 1 {
		n0 := int64(mapRange(float64(i-1), 0, float64(bn.Dx()), 0, float64(len(samples))))
		n1 := int64(mapRange(float64(i-0), 0, float64(bn.Dx()), 0, float64(len(samples))))

		// -------------------------------------------

		sum, min, max := 0.0, 0.0, 0.0
		for i := n0; i < n1; i++ {
			smp := samples[i]
			sum += math.Abs(smp)
			max = math.Max(smp, max)
			min = math.Min(smp, min)
		}
		avg := sum / float64(n1-n0)

		// -------------------------------------------

		s0 := int(mapRange(*RATIO*min, -gabs, gabs, -float64(middle), float64(middle)))
		s1 := int(mapRange(*RATIO*max, -gabs, gabs, -float64(middle), float64(middle)))
		if s1 != 0 || s0 != 0 {
			drawLine(img, i, middle-s0, i, middle-s1, fg1)
		}

		if !*HIDEAVG {
			s2 := int(mapRange(*RATIO*avg, gmin, gmax, -float64(middle), float64(middle)))
			if s2 != 0 {
				drawLine(img, i, middle-s2, i, middle+s2, fg0)
			}
		}
	}
}
