package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// 32-bit bitmask to determine the color of squares (all 1's represent black)
const ColorBitboard uint32 = 0xAA55AA55

var Purple = color.RGBA{0x71, 0x03, 0x8A, 0xFF}

func IdxColor(idx int) color.Color {
	// Convert index to 32-bit representation since the board pattern is the same
	if idx >= 32 {
		idx -= 32
	}
	// Compute if square is black or white based on index
	// and its intersection to the color bitmask
	if (ColorBitboard>>idx)&1 != 0 {
		return color.White
	}
	return Purple
}

const (
	pxSize     = 640
	imgWidth   = pxSize
	imgHeight  = pxSize
	boardSize  = 8
	squareSize = pxSize / 8
)

func main() {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{imgWidth, imgHeight}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	// cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < imgWidth; x++ {
		xIdx := x / squareSize

		for y := 0; y < imgHeight; y++ {
			yIdx := y / squareSize
			sqIdx := (xIdx * boardSize) + yIdx

			img.Set(x, y, IdxColor(sqIdx))
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}
