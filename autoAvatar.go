package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/png"
	"log"
	"os"
)

import _ "fmt"

func createHash(key string) string { // string -> md5 hash
	data := []byte(key)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func lastBits(hash string) []byte { // hash string -> []byte{last bits of all bytes}
	var bits []byte
	for _, b := range []byte(hash) {
		bits = append(bits, b&1) // побитово умножаю на 00000001
	}
	return bits
}

func mirrorBits(bits []byte) []byte { // prepearing array of bytes to convert in image
	for i := 4; i > 0; i-- {
		bits = append(bits, bits[i*8-8:i*8]...) // копирую по 8 битов с конца и добавляю в конец
	}
	return bits
}

func generateImage(bits []byte, colors [2]int, size, sideBlocks int) *image.RGBA { // create image
	scale := size / sideBlocks
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for x := 0; x < sideBlocks; x++ {
		for y := 0; y < sideBlocks; y++ {
			idx := colors[bits[x*8+y]] // по 8 бит от начала до конца
			col := palette.Plan9[idx]
			startPoint := image.Point{x * scale, y * scale}
			endPoint := image.Point{x*scale + scale, y*scale + scale}
			rectangle := image.Rectangle{startPoint, endPoint}
			draw.Draw(img, rectangle, &image.Uniform{col}, image.Point{}, draw.Src)
		}
	}
	return img
}

func saveImage(img *image.RGBA, filename string) { // download image in file

	f, err := os.Create(fmt.Sprintf("%s.png", filename))
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func selectColors(bits []byte) [2]int { // select color [0 - 256]
	var firstcolor int
	for _, b := range bits {
		firstcolor += int(b)
	}
	return [2]int{firstcolor, 256 - firstcolor}

}

func main() {
	inputString := flag.String("input", "Example", "input data")
	flag.Parse()
	hash := createHash(*inputString)
	bits := lastBits(hash)
	colors := selectColors(bits)
	outputFile := "avatar"
	imageSize := 256
	blocks := 8
	img := generateImage(mirrorBits(bits), colors, imageSize, blocks)
	saveImage(img, outputFile)
}
