package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"

	"github.com/xigh/go-wavreader"
)

var (
	OFFSET = flag.Uint64("offset", 0, "sey begin of samples")
	LENGTH = flag.Uint64("length", 0, "set number of samples [0 means all]")

	RATIO = flag.Float64("ratio", 0.8, "set ratio")

	WIDTH  = flag.Uint("width", 2048, "set width")
	HEIGHT = flag.Uint("height", 450, "set height")

	HIDEAVG    = flag.Bool("hideavg", false, "hide average")
	HIDERULERS = flag.Bool("hiderulers", false, "hide rulers")

	OUT = flag.String("out", "out.png", "set output filename")

	BINS      = flag.Uint("bins", 512, "set freq bins")
	PREEMP    = flag.Float64("preemp", 0.95, "pre-emphasis")
	RECTANGLE = flag.Bool("rectangle", false, "use rectangle window")
	DFT       = flag.Bool("dft", false, "use dft instead of fft")
	LOG10     = flag.Bool("log10", false, "pretty")
	MAG       = flag.Bool("mag", false, "mag")

	BG0 = flag.String("BG0", "000", "set background color 0")
	BG1 = flag.String("BG1", "333", "set background color 1")

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

	r, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	wr, err := wavreader.New(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s: %dHz, %d channels, %d samples, %v\n",
		name, wr.Rate(), wr.Chans(), wr.Len(), wr.Duration())

	start := *OFFSET
	if start > wr.Len() {
		log.Fatalf("offset bigger than file")
	}

	length := *LENGTH
	if start+length > wr.Len() {
		log.Printf("length too long\n")
		length = wr.Len() - start
	}

	samples := make([]float64, length)
	for i := uint64(0); i < length; i++ {
		s, err := wr.At(0, start+i)
		if err != nil {
			log.Fatal(err)
		}
		samples[i] = float64(s)
	}

	if *PREEMP > 0 {
		for i := len(samples) - 1; i > 0; i-- {
			samples[i] = samples[i] - *PREEMP*samples[i-1]
		}
	}

	W := int(*WIDTH)
	H := int(*HEIGHT)
	B := int(*BINS)

	bounds := image.Rect(-20, -20, W+20, H+40+B)
	img := NewImage128(bounds)

	bg0 := ParseColor(*BG0)
	fmt.Printf("bg0: %.8x\n", bg0)
	draw.Draw(img, img.Bounds(), image.NewUniform(bg0), image.ZP, draw.Src)

	fmt.Println("drawwav:")
	i0 := img.Sub(image.Rect(0, 0, W, H))
	drawwav(i0, samples)

	fmt.Println("drawfft:")
	i1 := img.Sub(image.Rect(0, H+20, W, H+20+B))
	drawfft(i1, samples, wr.Rate(), uint32(B))

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
