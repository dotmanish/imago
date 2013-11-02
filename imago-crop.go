// Copyright 2013 Manish Malik (manishmalik.name)
// All rights reserved.
// Use of this source code is governed by a BSD (3-Clause) License
// that can be found in the LICENSE file.

// Usage examples:
//
// Crop the image test-input.jpg 100 px from top and 200 px from left and write the output to test-output.jpg
// imago-crop -i test-input.jpg -o test-output.jpg -top 100 -left 200
//
// Crop the image test-input.jpg 50 px from right and write the output to test-output.png
// imago-crop -i test-input.jpg -o test-output.png -outformat png -right 50

package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"os"
	"strings"
)

var _imago_infilename, _imago_outfilename string
var _imago_outformat string
var _imago_left_offset, _imago_top_offset, _imago_right_offset, _imago_bottom_offset int
var _imago_jpeg_quality int

// Note to people who may read this for learning:
// There are multiple sophisticated command-line option parsers available for Go
// at http://code.google.com/p/go-wiki/wiki/Projects#Command-line_Option_Parsers
// Do take a look at these and decide whether you want to simply use 'flag' package
// or want to instead use the existing available parsing packages.
func init() {

	flag.StringVar(&_imago_infilename, "i", "", "input filename")
	flag.StringVar(&_imago_outfilename, "o", "", "output filename")
	flag.IntVar(&_imago_left_offset, "left", 0, "left offset")
	flag.IntVar(&_imago_top_offset, "top", 0, "top offset")
	flag.IntVar(&_imago_right_offset, "right", 0, "right offset")
	flag.IntVar(&_imago_bottom_offset, "bottom", 0, "bottom offset")
	flag.IntVar(&_imago_jpeg_quality, "jpegqual", 85, "jpeg quality")
	flag.StringVar(&_imago_outformat, "outformat", "jpeg", "output format")

}

// Function to initialize and check command line parameters
func initParams() {

	var paramsOkay bool
	paramsOkay = true

	flag.Parse()

	if _imago_outformat == "jpg" {
		_imago_outformat = "jpeg"
	} else if _imago_outformat != "jpeg" && _imago_outformat != "png" {
		fmt.Print("Invalid '-outformat'. Please specify only '-outformat jpeg' or '-outformat png'\n\n")
		paramsOkay = false
	}

	if _imago_infilename == "" || _imago_outfilename == "" {
		fmt.Print("Please specify both '-i (input file path)' and '-o (output file path)'\n\n")
		paramsOkay = false
	}

	if !paramsOkay {
		fmt.Print("Usage: imago-crop -i input_file_path -o output_file_path [-left left_offset] [-top top_offset] [-right right_offset] [-bottom bottom_offset] [-outformat output_format] [-jpegqual jpeg_quality]\n\n")
		fmt.Print("All offsets are specified in pixels.\n")
		fmt.Print("Default output_format is jpeg (can be 'jpeg' or 'png').\n")
		fmt.Print("Default jpeg_quality is 85.\n")
		os.Exit(1)
	}

}

func crop(r io.Reader, w io.Writer, left_offset_pixels, top_offset_pixels, right_offset_pixels, bottom_offset_pixels int) {

	img, _, err := image.Decode(r)
	if err != nil {
		fmt.Print("Invalid input file path, or Invalid image.")
		return
	}

	b := img.Bounds()
	newimg := image.NewRGBA(image.Rectangle{image.Point{b.Min.X, b.Min.Y}, image.Point{b.Max.X - left_offset_pixels - right_offset_pixels, b.Max.Y - top_offset_pixels - bottom_offset_pixels}})

	draw.Draw(newimg, newimg.Bounds(), img, image.Point{left_offset_pixels, top_offset_pixels}, draw.Src)

	outerr := jpeg.Encode(w, newimg, &jpeg.Options{_imago_jpeg_quality})
	if outerr != nil {
		fmt.Print("Invalid output file path, or some error occurred writing to it.")
		return
	}

}

func main() {

	initParams()

	_imago_outformat = strings.ToLower(_imago_outformat)

	fi, errfi := os.Open(_imago_infilename)
	if errfi != nil {
		panic(errfi)
	}
	defer func() {
		if errfi := fi.Close(); errfi != nil {
			panic(errfi)
		}
	}()

	r := bufio.NewReader(fi)

	fo, errfo := os.Create(_imago_outfilename)
	if errfo != nil {
		panic(errfo)
	}
	defer func() {
		if errfo := fo.Close(); errfo != nil {
			panic(errfo)
		}
	}()

	w := bufio.NewWriter(fo)

	crop(r, w, _imago_left_offset, _imago_top_offset, _imago_right_offset, _imago_bottom_offset)

	// TODO: Return appropriate exit code in case of failure
}
