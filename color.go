package main

import (
	"log"
	"strconv"
)

type Color128 struct {
	R, G, B, A uint32
}

func (c Color128) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, c.A
}

func NewColor128(r, g, b, a uint32) Color128 {
	return Color128{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func ParseColor(text string) Color128 {
	if text == "transparent" {
		return Color128{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
	}
	var err error
	r, g, b, a := uint64(0), uint64(0), uint64(0), uint64(0xffffffff)
	switch len(text) {
	case 4:
		a, err = strconv.ParseUint(text[3:4], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		a |= a << 4
		a |= a << 8
		a |= a << 16
		fallthrough

	case 3:
		r, err = strconv.ParseUint(text[0:1], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		r |= r << 4
		r |= r << 8
		r |= r << 16

		g, err = strconv.ParseUint(text[1:2], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		g |= g << 4
		g |= g << 8
		g |= g << 16

		b, err = strconv.ParseUint(text[2:3], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		b |= b << 4
		b |= b << 8
		b |= b << 16

	case 8:
		a, err = strconv.ParseUint(text[6:8], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		a |= a << 8
		a |= a << 16
		fallthrough

	case 6:
		r, err = strconv.ParseUint(text[0:2], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		r |= r << 8
		r |= r << 16

		g, err = strconv.ParseUint(text[2:4], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		g |= g << 8
		g |= g << 16

		b, err = strconv.ParseUint(text[4:6], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		b |= b << 8
		b |= b << 16

	default:
		log.Fatalf("invalid color %q", text)
	}
	return Color128{
		R: uint32(r),
		G: uint32(g),
		B: uint32(b),
		A: uint32(a),
	}
}
