package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/nfnt/resize"
)

func scale(img image.Image, w int) (image.Image, int, int) {
	sz := img.Bounds()
	h := (sz.Max.Y * w * 10) / (sz.Max.X * 16)
	img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	return img, w, h
}
func convertColorToRGB(col color.Color) (uint8, uint8, uint8) {
	rgbaColor := color.RGBAModel.Convert(col)
	_r, _g, _b, _ := rgbaColor.RGBA()
	// rgb values are uint8s, I cannot comprehend why the stdlib would return
	// int32s :facepalm:
	r := uint8(_r & 0xFF)
	g := uint8(_g & 0xFF)
	b := uint8(_b & 0xFF)
	return r, g, b
}

// 48;2;r;g;bm - set background colour to rgb
func rgbBackgroundSequence(r, g, b uint8) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// 38;2;r;g;bm - set text colour to rgb
func rgbTextSequence(r, g, b uint8) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func resetColorSequence() string {
	return "\x1b[0m"
}

func ConvertImageToANSI(img image.Image, skip int) string {
	var upperHalfBlock = "▀"
	// We'll just reuse this to increment the loop counters
	skip++
	ansi := resetColorSequence()
	yMax := img.Bounds().Max.Y
	xMax := img.Bounds().Max.X

	sequences := make([]string, yMax)

	for y := img.Bounds().Min.Y; y < yMax; y += 2 * skip {
		sequence := ""
		for x := img.Bounds().Min.X; x < xMax; x += skip {
			upperPix := img.At(x, y)
			lowerPix := img.At(x, y+skip)

			ur, ug, ub := convertColorToRGB(upperPix)
			lr, lg, lb := convertColorToRGB(lowerPix)

			if y+skip >= yMax {
				sequence += resetColorSequence()
			} else {
				sequence += rgbBackgroundSequence(lr, lg, lb)
			}

			sequence += rgbTextSequence(ur, ug, ub)
			sequence += upperHalfBlock

			sequences[y] = sequence
		}
	}

	for y := img.Bounds().Min.Y; y < yMax; y += 2 * skip {
		ansi += sequences[y] + resetColorSequence() + "\n"
	}

	return ansi
}
