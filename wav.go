package main

import (
	"encoding/binary"
	"os"

	"github.com/cryptix/wav"
)

type Wav struct {
	f *os.File
	r *wav.Reader
}

func (w *Wav) Close() {
	w.f.Close()
}

func (w *Wav) GetSampleCount() uint32 {
	return w.r.GetSampleCount()
}

func (w *Wav) GetSampleRate() uint32 {
	return w.r.GetFile().SampleRate
}

func (w *Wav) GetBits() uint16 {
	return w.r.GetFile().SignificantBits
}

func (w *Wav) GetChannels() uint16 {
	return w.r.GetFile().Channels
}

func (w *Wav) GetFormat() uint16 {
	return w.r.GetFile().AudioFormat
}

func (w *Wav) GetSamplesAt(offset uint32, count uint32) []float64 {
	if w.GetFormat() != 1 {
		panic("todo")
	}

	if w.GetChannels() != 1 {
		panic("todo")
	}

	_, err := w.f.Seek(int64(w.r.FirstSampleOffset()+offset*2), os.SEEK_SET)
	if err != nil {
		panic(err)
	}

	raw := make([]int16, count)
	err = binary.Read(w.f, binary.LittleEndian, raw)
	if err != nil {
		panic(err)
	}

	samples := make([]float64, count)
	for i, s := range raw {
		samples[i] = mapRange(float64(s), -32768, 32767, -.5, .5)
	}

	return samples
}

func OpenWav(name string) (*Wav, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	r, err := wav.NewReader(f, info.Size())
	if err != nil {
		return nil, err
	}

	return &Wav{
		f: f,
		r: r,
	}, nil
}
