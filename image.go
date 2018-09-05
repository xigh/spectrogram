package main

import (
	"image"
	"image/color"
)

type Image128 struct {
	set    int
	at     int
	pix    []Color128
	bounds image.Rectangle
}

func (img *Image128) ColorModel() color.Model {
	return color.RGBA64Model
}

func (img *Image128) Bounds() image.Rectangle {
	return img.bounds
}

func (img *Image128) offset(x, y int) int {
	p := image.Point{x, y}
	if !p.In(img.bounds) {
		return -1
	}
	stride := img.bounds.Dx()
	my := img.bounds.Min.Y
	mx := img.bounds.Min.X
	ny := y - my
	nx := x - mx
	return ny*stride + nx
}

func (img *Image128) At(x int, y int) color.Color {
	o := img.offset(x, y)
	if o < 0 {
		return Color128{}
	}
	img.at++
	return img.pix[o]
}

func (img *Image128) Set(x int, y int, c color.Color) {
	o := img.offset(x, y)
	if o < 0 {
		return
	}
	img.set++
	r, g, b, a := c.RGBA()
	img.pix[o] = Color128{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func (img *Image128) Stats() (int, int) {
	return img.at, img.set
}

func NewImage128(bounds image.Rectangle) *Image128 {
	dx := bounds.Dx()
	dy := bounds.Dy()
	sz := dx * dy
	return &Image128{
		pix:    make([]Color128, sz),
		bounds: bounds,
	}
}

type SubImage128 struct {
	img    *Image128
	bounds image.Rectangle
}

func (sub *SubImage128) ColorModel() color.Model {
	return sub.ColorModel()
}

func (sub *SubImage128) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(sub.bounds.Dx(), sub.bounds.Dy()),
	}
}

func (sub *SubImage128) At(x int, y int) color.Color {
	mx := sub.bounds.Min.X
	my := sub.bounds.Min.Y
	return sub.img.At(x+mx, y+my)
}

func (sub *SubImage128) Set(x int, y int, c color.Color) {
	mx := sub.bounds.Min.X
	my := sub.bounds.Min.Y
	sub.img.Set(x+mx, y+my, c)
}

func (img *Image128) Sub(bounds image.Rectangle) *SubImage128 {
	return &SubImage128{
		img:    img,
		bounds: bounds,
	}
}
