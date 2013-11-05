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
//
// Crop the image test-input.jpg 20% from top, 10% from bottom, and write the output to test-output.png
// imago-crop -i test-input.jpg -o test-output.png -outformat png -top 20% -bottom 10%
//
// Crop the image test-input.jpg 25% px from left, ensuring at least 300px width,
// and write the output to test-output.jpg
// imago-crop -i test-input.jpg -o test-output.jpg -left 25% -minwidth 300
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/draw"
	_ "image/png"
	_ "image/gif"
	"image/jpeg"
	"image"
	"io"
	"os"
	"strings"
	"strconv"
)

var _imago_infilename, _imago_outfilename string
var _imago_outformat string
var _imago_left_offset, _imago_top_offset, _imago_right_offset, _imago_bottom_offset string
var _imago_left_offset_px, _imago_top_offset_px, _imago_right_offset_px, _imago_bottom_offset_px int
var _imago_left_offset_perc, _imago_top_offset_perc, _imago_right_offset_perc, _imago_bottom_offset_perc float64
var _imago_jpeg_quality, _imago_min_width int

// Note to people who may read this for learning:
// There are multiple sophisticated command-line option parsers available for Go
// at http://code.google.com/p/go-wiki/wiki/Projects#Command-line_Option_Parsers
// Do take a look at these and decide whether you want to simply use 'flag' package
// or want to instead use the existing available parsing packages.
func init() {

	flag.StringVar(&_imago_infilename, "i", "", "input filename")
	flag.StringVar(&_imago_outfilename, "o", "", "output filename")
	flag.StringVar(&_imago_left_offset, "left", "0", "left offset")
	flag.StringVar(&_imago_top_offset, "top", "0", "top offset")
	flag.StringVar(&_imago_right_offset, "right", "0", "right offset")
	flag.StringVar(&_imago_bottom_offset, "bottom", "0", "bottom offset")
	flag.IntVar(&_imago_min_width, "minwidth", 0, "min width")
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

	// Parse left offset. Process percentage if specified.
	if _imago_left_offset != "" {
		if strings.HasSuffix(_imago_left_offset, "%") {
			offset := strings.TrimSuffix(_imago_left_offset, "%")
			_imago_left_offset_perc, _ = strconv.ParseFloat(offset, 64)
			if _imago_left_offset_perc < 0 || _imago_left_offset_perc > 100 {
				fmt.Println("Left offset percetange must be between 0 and 100.")
				paramsOkay = false
			}
		} else {
			result64, _ := strconv.ParseInt(_imago_left_offset, 10, 0)
			_imago_left_offset_px = int(result64)
		}
	}

	// Parse top offset. Process percentage if specified.
	if _imago_top_offset != "" {
		if strings.HasSuffix(_imago_top_offset, "%") {
			offset := strings.TrimSuffix(_imago_top_offset, "%")
			_imago_top_offset_perc, _ = strconv.ParseFloat(offset, 64)
			if _imago_top_offset_perc < 0 || _imago_top_offset_perc > 100 {
				fmt.Println("Top offset percetange must be between 0 and 100.")
				paramsOkay = false
			}
		} else {
			result64, _ := strconv.ParseInt(_imago_top_offset, 10, 0)
			_imago_top_offset_px = int(result64)
		}
	}

	// Parse bottom offset. Process percentage if specified.
	if _imago_bottom_offset != "" {
		if strings.HasSuffix(_imago_bottom_offset, "%") {
			offset := strings.TrimSuffix(_imago_bottom_offset, "%")
			_imago_bottom_offset_perc, _ = strconv.ParseFloat(offset, 64)
			if _imago_bottom_offset_perc < 0 || _imago_bottom_offset_perc > 100 {
				fmt.Println("Bottom offset percetange must be between 0 and 100.")
				paramsOkay = false
			}
		} else {
			result64, _ := strconv.ParseInt(_imago_bottom_offset, 10, 0)
			_imago_bottom_offset_px = int(result64)
		}
	}

	// Parse right offset. Process percentage if specified.
	if _imago_right_offset != "" {
		if strings.HasSuffix(_imago_right_offset, "%") {
			offset := strings.TrimSuffix(_imago_right_offset, "%")
			_imago_right_offset_perc, _ = strconv.ParseFloat(offset, 64)
			if _imago_right_offset_perc < 0 || _imago_right_offset_perc > 100 {
				fmt.Println("Right offset percetange must be between 0 and 100.")
				paramsOkay = false
			}
		} else {
			result64, _ := strconv.ParseInt(_imago_right_offset, 10, 0)
			_imago_right_offset_px = int(result64)
		}
	}

	if !paramsOkay {
		fmt.Print("Usage: imago-crop -i input_file_path -o output_file_path [-left left_offset] [-top top_offset] [-right right_offset] [-bottom bottom_offset] [-outformat output_format] [-minwidth minimum_width] [-jpegqual jpeg_quality]\n\n")
		fmt.Print("All offsets can be specified in pixels (e.g. 300) or as percentage (e.g. 20%).\n")
		fmt.Print("Default output_format is jpeg (can be 'jpeg' or 'png').\n")
		fmt.Print("Default jpeg_quality is 85.\n")
		os.Exit(1)
	}

}

// The cropping function
func crop(r io.Reader, w io.Writer, left_offset_pixels, top_offset_pixels, right_offset_pixels, bottom_offset_pixels, min_width_pixels int) {

	img, _, err := image.Decode(r)
	if err != nil {
		fmt.Println("Invalid input file path, or Invalid image.")
		fmt.Println(err.Error())
		return
	}

	b := img.Bounds()

	// Check if we are within the min-width, if specified by the user
	if b.Max.X - left_offset_pixels - right_offset_pixels < min_width_pixels {
		diff_pixels := min_width_pixels - (b.Max.X - left_offset_pixels - right_offset_pixels)
		left_offset_pixels = left_offset_pixels - (diff_pixels / 2)
		right_offset_pixels = right_offset_pixels - (diff_pixels / 2)
	}

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

	fi_pre, errfi := os.Open(_imago_infilename)
	if errfi != nil {
		panic(errfi)
	}
	defer func() {
		if errfi := fi_pre.Close(); errfi != nil {
			panic(errfi)
		}
	}()

	r_pre := bufio.NewReader(fi_pre)

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

	img, _, err := image.Decode(r_pre)
	if err != nil {
		fmt.Println("Invalid input file path, or Invalid image.")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Pre-compute the offsets in case percentages were specified for one or more offsets

	b := img.Bounds()
	if _imago_left_offset_px == 0 && _imago_left_offset_perc != 0.0 {
			_imago_left_offset_px = (int) ((float64(b.Max.X) * _imago_left_offset_perc) / 100.0)
	}
	if _imago_top_offset_px == 0 && _imago_top_offset_perc != 0.0 {
			_imago_top_offset_px = (int) ((float64(b.Max.Y) * _imago_top_offset_perc) / 100.0)
	}
	if _imago_bottom_offset_px == 0 && _imago_bottom_offset_perc != 0.0 {
			_imago_bottom_offset_px = (int) ((float64(b.Max.Y) * _imago_bottom_offset_perc) / 100.0)
	}
	if _imago_right_offset_px == 0 && _imago_right_offset_perc != 0.0 {
			_imago_right_offset_px = (int) ((float64(b.Max.X) * _imago_right_offset_perc) / 100.0)
	}

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

	crop(r, w, _imago_left_offset_px, _imago_top_offset_px, _imago_right_offset_px, _imago_bottom_offset_px, _imago_min_width)

	// TODO: Return appropriate exit code in case of failure
}
