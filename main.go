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
	pxSize    = 800
	boardSize = 8
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
)

func drawBoard(img *image.RGBA) {
	size := img.Bounds().Size()
	width := size.X
	height := size.Y
	squareSize := width / boardSize

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		xIdx := x / squareSize

		for y := 0; y < height; y++ {
			yIdx := y / squareSize
			sqIdx := (xIdx * boardSize) + yIdx

			img.Set(x, y, IdxColor(sqIdx))
		}
	}
}

func drawPiece(color, notation string, img *image.RGBA, squareSize int, pieces *image.NRGBA) {
	if len(notation) < 2 || len(notation) > 3 {
		log.Fatalf("wrong notation: %s", strconv.Quote(notation))
	}

	piece := "" // pawn
	if len(notation) == 3 {
		piece = string(notation[0])
		notation = notation[1:]
	}

	pieceRect, err := pieceRect(piece, color)
	if err != nil {
		log.Fatal(err)
	}

	sqRect, err := notationRect(notation, squareSize)
	if err != nil {
		log.Fatal(err)
	}

	draw.Draw(img, sqRect, pieces.SubImage(pieceRect), pieceRect.Min, draw.Over)
}

func main() {
	position := []struct {
		color string
		loc   []string
	}{
		{"w", []string{"kb6", "a7"}},
		{"b", []string{"ka8", "h2"}},
	}

	img := image.NewRGBA(image.Rect(0, 0, pxSize, pxSize))
	drawBoard(img)

	squareSize := pxSize / boardSize

	pieces, err := load(piecesImg)
	if err != nil {
		log.Fatalf("couldn't load image %s: %v", piecesImg, err)
	}

	for _, side := range position {
		for _, p := range side.loc {
			drawPiece(side.color, p, img, squareSize, pieces)
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image1.png")
	png.Encode(f, img)
}

func pieceRect(piece, color string) (image.Rectangle, error) {
	const (
		pieces = "qkrnbp"
		colors = "bw"
	)

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

// notationRect returns a rectangle with corresponding location on the board
// from square notation ("a1", "d3").
func notationRect(notation string, squareSize int) (image.Rectangle, error) {
	const (
		files = "abcdefgh"
		ranks = "12345678"
	)

	if len(notation) != 2 {
		return image.Rectangle{}, fmt.Errorf(
			"Wrong square notation format %s, expected format \"b7\"",
			notation,
		)
	}

	notation = strings.ToLower(notation)

	fileStr := string(notation[0])
	rankStr := string(notation[1])

	if !strings.Contains(ranks, rankStr) || !strings.Contains(files, fileStr) {
		return image.Rectangle{}, fmt.Errorf(
			"Wrong square notation %s, rank/file not found",
			notation,
		)
	}

	file := strings.Index(files, fileStr)

	// a1 shuold be bottom of board but because we are drawing from top left
	// then we have to make sure that a1 corresponding to index 0, 7 (instead of 0, 0)
	rank, _ := strconv.Atoi(rankStr)
	rank = boardSize - rank

	// NOTE: this assumes that square size is bigger than piece size
	piecePadding := (squareSize - pieceSize) / 2

	rect := image.Rect(
		(file*squareSize)+piecePadding,
		(rank*squareSize)+piecePadding,
		(file+1)*squareSize,
		(rank+1)*squareSize,
	)
	return rect, nil
}
