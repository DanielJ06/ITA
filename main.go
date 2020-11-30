package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"reflect"

	"github.com/nfnt/resize"
)

var ASCIISTR = "MND8OZ$7I?+=~:,.."

func ReadImage() (image.Image, int) {
	// Define width flag
	width := flag.Int("w", 80, "Use -w <width>")
	// Define image path flag
	path := flag.String("p", "gopher.png", "Use -p <imagePath>")
	// Parse the flags after they're defined
	flag.Parse()

	// Read the named file
	f, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()
	return img, *width
}

func ResizeImage(img image.Image, width int) (image.Image, int, int) {
	size := img.Bounds()
	heigth := (size.Max.Y * width * 10) / (size.Max.X * 16)
	img = resize.Resize(uint(width), uint(heigth), img, resize.Lanczos2)

	return img, width, heigth
}

func Convert(img image.Image, w int, h int) []byte {
	table := []byte(ASCIISTR)
	buff := new(bytes.Buffer)

	for i := 0; i < h; i++ {
		for x := 0; x < w; x++ {
			// Return every pixel applying the GrayModel
			g := color.GrayModel.Convert(img.At(x, i))
			y := reflect.ValueOf(g).FieldByName("Y").Uint()
			pos := int(y * 16 / 255)
			_ = buff.WriteByte(table[pos])
		}
		_ = buff.WriteByte('\n')
	}
	return buff.Bytes()
}

func main() {
	res := Convert(ResizeImage(ReadImage()))
	fmt.Print(string(res))
}
