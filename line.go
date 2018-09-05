package main

import (
	"image/color"
	"image/draw"
)

func drawLine(img draw.Image, x1, y1, x2, y2 int, col color.Color) {
	if y1 == y2 {
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		for x := x1; x <= x2; x++ {
			img.Set(x, y1, col)
		}
		return
	}

	if x1 == x2 {
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for y := y1; y <= y2; y++ {
			img.Set(x1, y, col)
		}
		return
	}

	// TODO: Bresenham
}
