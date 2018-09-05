package main

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	"os"
)

type PngImage struct {
	img *Image128
}

func (p *PngImage) ColorModel() color.Model {
	return color.RGBA64Model
}

func (p *PngImage) Bounds() image.Rectangle {
	return p.img.Bounds()
}

func (p *PngImage) At(x int, y int) color.Color {
	c := p.img.At(x, y)
	r, g, b, a := c.RGBA()
	return color.RGBA64{
		R: uint16(r >> 16),
		G: uint16(g >> 16),
		B: uint16(b >> 16),
		A: uint16(a >> 16),
	}
}

func savePng(img *Image128, fileName string) error {
	pi := &PngImage{img: img}

	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	bf := bufio.NewWriter(outFile)

	err = png.Encode(bf, pi)
	if err != nil {
		return err
	}

	err = bf.Flush()
	if err != nil {
		return err
	}

	return nil
}
