package main

import (
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"

	"github.com/apex/go-apex"
)

type Message struct {
	Body MessageBody `json:"body"`
}

type MessageBody struct {
	Pixels [][]Pixel `json:"pixels"`
	Path   string    `json:"path"`
}

type Pixel struct {
	R int
	G int
	B int
	A int
}

func getFileAndGetPixels(path string) ([][]Pixel, error) {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	response, err := http.Get(path)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	pixels, err := getPixels(response.Body)

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
	return Pixel{
		int(r / 257),
		int(g / 257),
		int(b / 257),
		int(a / 257),
	}
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m Message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		pixels, err := getFileAndGetPixels(m.Body.Path)

		if err != nil {
			return nil, err
		}

		m.Body.Pixels = pixels

		return m, nil
	})
}
