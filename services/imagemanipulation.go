package services

import (
	"errors"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gocv.io/x/gocv"
)

type ImageManipulationOptions struct {
	BasePath        string  `json:"base_path"`
	InputPath       string  `json:"input_path"`
	OutputPath      string  `json:"output_path"`
	FileName        string  `json:"file_name"`
	InputFilePath   string  `json:"input_file_path"`
	OutputFilePath  string  `json:"output_file_path"`
	Width           float64 `json:"width"`
	Height          float64 `json:"height"`
	Quality         int     `json:"quality"`
	KeepAspectRatio bool    `json:"keep_aspect_ratio"`
	Debug           bool    `json:"debug"`
}

func (imo *ImageManipulationOptions) init(basePath string, inputPath string, outputPath string, filename string, width float64, height float64, quality int, targetImageFormat string, keepAspecRatio bool, debug bool) (bool, error) {
	cwd, _ := os.Getwd()

	if width == 0 {
		imo.Width = 100
	} else {
		imo.Width = width
	}
	if height == 0 {
		imo.Height = 100
	} else {
		imo.Height = height
	}
	if quality == 0 {
		imo.Quality = 80
	} else {
		imo.Quality = quality
	}
	if basePath == "" {
		imo.BasePath = filepath.Join(cwd, "files")
	} else {
		imo.BasePath = basePath
	}
	if inputPath == "" {
		imo.InputPath = filepath.Join(imo.BasePath, "input")
	} else {
		imo.InputPath = inputPath
	}
	if outputPath == "" {
		imo.OutputPath = filepath.Join(imo.BasePath, "output")
	} else {
		imo.OutputPath = outputPath
	}
	if len(filename) == 0 {
		return false, errors.New("invalid file name")
	} else {
		imo.FileName = filename
	}
	now := time.Now()
	ts := now.UnixNano()
	mt := strings.Split(filename, ".")
	imgFormat := mt[len(mt)-1]
	imo.InputFilePath = filepath.Join(imo.InputPath, imo.FileName)
	imo.OutputFilePath = filepath.Join(imo.OutputPath, strings.Replace(imo.FileName, fmt.Sprintf(".%s", imgFormat), fmt.Sprintf("-%v-%d.%s", ts, imo.Quality, targetImageFormat), 1))
	imo.KeepAspectRatio = keepAspecRatio
	imo.Debug = debug
	return true, nil
}

type ImageManipulation struct {
	options ImageManipulationOptions
}

func (im *ImageManipulation) CalculateAspectRatioFit(srcWidth int, srcHeight int, targetWidth int, targetHeight int) map[string]float64 {
	// source: https://opensourcehacker.com/2011/12/01/calculate-aspect-ratio-conserving-resize-for-images-in-javascript/
	r := map[string]float64{}
	ratio := math.Min(float64(targetWidth)/float64(srcWidth), float64(targetHeight)/float64(srcHeight))
	r["width"] = float64(srcWidth) * ratio
	r["height"] = float64(srcHeight) * ratio
	return r
}

func (im *ImageManipulation) PngToJpeg(basePath string, inputPath string, outputPath string, filename string, debug bool) (string, error) {
	// set options value
	_, err := im.options.init(basePath, inputPath, outputPath, filename, -1, -1, 100, "jpeg", true, debug)
	if err != nil {
		return "", err
	}
	src := gocv.IMRead(im.options.InputFilePath, gocv.IMReadColor)
	if src.Empty() {
		return "", errors.New("failed to read input file")
	}
	// main logic
	if ok := gocv.IMWriteWithParams(im.options.OutputFilePath, src, []int{gocv.IMWriteJpegQuality, im.options.Quality}); !ok {
		return "", errors.New("failed to write output file")
	}
	return im.options.OutputFilePath, nil
}

func (im *ImageManipulation) Resize(basePath string, inputPath string, outputPath string, filename string, width float64, height float64, quality int, keepAspecRatio bool, debug bool) (string, error) {
	// set options value
	mt := strings.Split(filename, ".")
	imgFormat := mt[len(mt)-1]
	_, err := im.options.init(basePath, inputPath, outputPath, filename, width, height, quality, imgFormat, keepAspecRatio, debug)
	if err != nil {
		return "", err
	}
	// main logic
	src := gocv.IMRead(im.options.InputFilePath, gocv.IMReadColor)
	if src.Empty() {
		return "", errors.New("failed to read input file")
	}
	transform := gocv.NewMat()

	var fx, fy float64
	if im.options.KeepAspectRatio {
		fitSize := im.CalculateAspectRatioFit(src.Cols(), src.Rows(), int(im.options.Width), int(im.options.Height))
		// fmt.Printf("KEEP ASPECT RATIO: %#v \n", fitSize)
		fx = float64(int(fitSize["width"])) / float64(src.Cols())
		fy = float64(int(fitSize["height"])) / float64(src.Rows())
	} else {
		fx = im.options.Width / float64(src.Cols())
		fy = im.options.Height / float64(src.Rows())
	}

	gocv.Resize(src, &transform, image.Point{}, fx, fy, gocv.InterpolationCubic)

	if ok := gocv.IMWrite(im.options.OutputFilePath, transform); !ok {
		return "", errors.New("failed to write output file")
	}
	return im.options.OutputFilePath, nil
}

func (im *ImageManipulation) Compress(basePath string, inputPath string, outputPath string, filename string, quality int, debug bool) (string, error) {
	// set options value
	_, err := im.options.init(basePath, inputPath, outputPath, filename, -1, -1, quality, "jpeg", true, debug)
	if err != nil {
		return "", err
	}
	// main logic
	src := gocv.IMRead(im.options.InputFilePath, gocv.IMReadColor)
	if src.Empty() {
		return "", errors.New("failed to read input file")
	}
	if ok := gocv.IMWriteWithParams(im.options.OutputFilePath, src, []int{gocv.IMWriteJpegQuality, im.options.Quality}); !ok {
		return "", errors.New("failed to write output file")
	}
	return im.options.OutputFilePath, nil
}
