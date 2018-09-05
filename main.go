package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log"
	"time"
)

var (
	OFFSET = flag.Int("offset", 0, "sey begin of samples")
	LENGTH = flag.Int("length", 0, "set number of samples [0 means all]")

	RATIO = flag.Float64("ratio", 0.8, "set ratio")

	WIDTH  = flag.Int("width", 2048, "set width")
	HEIGHT = flag.Int("height", 450, "set height")

	HIDEAVG    = flag.Bool("hideavg", false, "hide average")
	HIDERULERS = flag.Bool("hiderulers", false, "hide rulers")

	BINS = flag.Int("bins", 512, "set freq bins")

	OUT = flag.String("out", "out.png", "set output filename")

	DFT = flag.Bool("dft", false, "use dft instead of fft")

	BG0 = flag.String("BG0", "000000", "set background color 0")
	BG1 = flag.String("BG1", "444444", "set background color 1")
	BG2 = flag.String("BG2", "447744", "set background color 2")

	FG0 = flag.String("FG0", "0972a2", "set forground color 0")
	FG1 = flag.String("FG1", "6b5f7e", "set forground color 1")
	RUL = flag.String("RUL", "a0b0c0", "set rulers color")
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Printf("usage: fft [options] file.wav\n")
		return
	}

	name := flag.Arg(0)

	wav, err := OpenWav(name)
	if err != nil {
		log.Fatalf("OpenWav failed: %v", err)
	}
	defer wav.Close()

	count := wav.GetSampleCount()
	rate := wav.GetSampleRate()

	duration := time.Duration((float64(count) / float64(rate)) * float64(time.Second))
	ft := "unknown"
	switch wav.GetFormat() {
	case 1:
		ft = "signed"
	}

	o := uint32(*OFFSET)

	l := uint32(*LENGTH)
	if l == 0 {
		l = count
	}

	if o > count {
		log.Fatalf("empty duration")
	}

	if o+l > count {
		l = count - o
	}

	fmt.Printf("%s: %dHz %d bits [%s], %d channels, %d/%d samples, %v\n",
		name, rate, wav.GetBits(), ft, wav.GetChannels(), l, count, duration)

	samples := wav.GetSamplesAt(o, l)
	fmt.Printf("samples: %.2d\n", len(samples))

	bounds := image.Rect(-20, -20, *WIDTH+20, *HEIGHT+40+*BINS)
	img := NewImage128(bounds)

	bg0 := ParseColor(*BG0)
	fmt.Printf("bg0: %.8x\n", bg0)
	draw.Draw(img, img.Bounds(), image.NewUniform(bg0), image.ZP, draw.Src)

	fmt.Println("drawwav:")
	i0 := img.Sub(image.Rect(0, 0, *WIDTH, *HEIGHT))
	drawwav(i0, samples)

	fmt.Println("drawfft:")
	i1 := img.Sub(image.Rect(0, *HEIGHT+20, *WIDTH, *HEIGHT+20+*BINS))
	drawfft(i1, samples, rate, uint32(*BINS))

	a0, s0 := img.Stats()
	fmt.Printf("img stats: %d reads, %d writes\n", a0, s0)

	fmt.Printf("saving %q\n", *OUT)

	err = savePng(img, *OUT)
	if err != nil {
		log.Fatalf("savePng failed: %v", err)
	}

	fmt.Printf("saved %q\n", *OUT)

	a1, s1 := img.Stats()
	fmt.Printf("img stats: %d reads, %d writes\n", a1, s1)
}
