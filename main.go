package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
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
	squareSize = pxSize / boardSize
)

func load(filePath string) (*image.NRGBA, error) {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		return nil, fmt.Errorf("Cannot read file: %v", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode file: %v", err)
	}
	return img.(*image.NRGBA), nil
}

const (
	piecesImg = "ChessPiecesArray.png"
	pieceSize = 60
	// NOTE: this assumes that square size is bigger than piece size
	piecePadding = (squareSize - pieceSize) / 2
)

func main() {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{imgWidth, imgHeight}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	pieces, err := load(piecesImg)
	if err != nil {
		log.Fatalf("couldn't load image %s: %v", piecesImg, err)
	}

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
	pieceRect, err := pieceRect("n", "b")
	if err != nil {
		log.Fatal(err)
	}

	sqRect, err := notationRect("a1")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(sqRect)
	draw.Draw(img, sqRect, pieces.SubImage(pieceRect), pieceRect.Min, draw.Over)

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

const (
	pieces = "qkrnbp"
	colors = "bw"
)

func pieceRect(piece, color string) (image.Rectangle, error) {
	piece = strings.ToLower(piece)
	if piece == "" {
		// treat empty string as pawn
		piece = "p"
	}

	if len(piece) > 1 || !strings.Contains(pieces, piece) {
		return image.Rectangle{}, fmt.Errorf("invalid piece: %s", strconv.Quote(piece))
	}

	pieceIdx := strings.Index(pieces, piece)

	color = strings.ToLower(color)
	if len(color) > 1 || !strings.Contains(colors, color) {
		return image.Rectangle{}, fmt.Errorf("invalid color: %s", strconv.Quote(color))
	}

	colorIdx := strings.Index(colors, color)

	rect := image.Rect(
		pieceIdx*pieceSize,
		colorIdx*pieceSize,
		(pieceIdx+1)*pieceSize,
		(colorIdx+1)*pieceSize,
	)
	return rect, nil
}

var (
	Files = "abcdefgh"
	Ranks = "12345678"
)

// notationRect returns a rectangle with corresponding location on the board
// from square notation ("a1", "d3").
func notationRect(notation string) (image.Rectangle, error) {
	if len(notation) != 2 {
		return image.Rectangle{}, fmt.Errorf(
			"Wrong square notation format %s, expected format \"b7\"",
			notation,
		)
	}

	notation = strings.ToLower(notation)

	fileStr := string(notation[0])
	rankStr := string(notation[1])

	if !strings.Contains(Ranks, rankStr) || !strings.Contains(Files, fileStr) {
		return image.Rectangle{}, fmt.Errorf(
			"Wrong square notation %s, rank/file not found",
			notation,
		)
	}

	file := strings.Index(Files, fileStr)

	// a1 shuold be bottom of board but because we are drawing from top left
	// then we have to make sure that a1 corresponding to index 0, 7 (instead of 0, 0)
	rank, _ := strconv.Atoi(rankStr)
	rank = boardSize - rank

	rect := image.Rect(
		(file*squareSize)+piecePadding,
		(rank*squareSize)+piecePadding,
		(file+1)*squareSize,
		(rank+1)*squareSize,
	)
	return rect, nil
}
