package main

import (
	"image"
	"image/png"
	"io"
	"os"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/apex/go-apex"
)

type Message struct {
	Pixels [][]Pixel `json:"pixels"`
	Path string `json:"path"`
}

type Pixel struct {
	R int
	G int
	B int
	A int
}

func openFileAndGetPixels(path string) ([][]Pixel, error) {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	pixels, err := getPixels(file)

	if err != nil {
		return nil, err
	}

	return pixels, nil
}

// Many thanks to this SO post: https://goo.gl/CDW8Qq
func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel

	for y := 0; y < height; y++ {
		var row []Pixel

		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}

		pixels = append(pixels, row)
	}

	return pixels, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func PixelsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi")
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m Message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		pixels, err := openFileAndGetPixels(m.Path)

		if err != nil {
			return nil, err
		}

		m.Pixels = pixels

		return m, nil
	})
}
