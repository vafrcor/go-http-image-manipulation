package helpers

import (
	"errors"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gocv.io/x/gocv"
)

type ImageManipulationOptions struct {
	BasePath       string  `json:"base_path"`
	InputPath      string  `json:"input_path"`
	OutputPath     string  `json:"output_path"`
	FileName       string  `json:"file_name"`
	InputFilePath  string  `json:"input_file_path"`
	OutputFilePath string  `json:"output_file_path"`
	Width          float64 `json:"width"`
	Height         float64 `json:"height"`
	Quality        int     `json:"quality"`
	Debug          bool    `json:"debug"`
}

func (iro *ImageManipulationOptions) init(basePath string, inputPath string, outputPath string, filename string, width float64, height float64, quality int, targetImageFormat string, debug bool) (bool, error) {
	cwd, _ := os.Getwd()

	if width == 0 {
		iro.Width = 100
	} else {
		iro.Width = width
	}
	if height == 0 {
		iro.Height = 100
	} else {
		iro.Height = height
	}
	if quality == 0 {
		iro.Quality = 80
	} else {
		iro.Quality = quality
	}
	if basePath == "" {
		iro.BasePath = filepath.Join(cwd, "files")
	} else {
		iro.BasePath = basePath
	}
	if inputPath == "" {
		iro.InputPath = filepath.Join(iro.BasePath, "input")
	} else {
		iro.InputPath = inputPath
	}
	if outputPath == "" {
		iro.OutputPath = filepath.Join(iro.BasePath, "output")
	} else {
		iro.OutputPath = outputPath
	}
	if len(filename) == 0 {
		return false, errors.New("invalid file name")
	} else {
		iro.FileName = filename
	}
	now := time.Now()
	ts := now.UnixNano()
	mt := strings.Split(filename, ".")
	imgFormat := mt[len(mt)-1]
	iro.InputFilePath = filepath.Join(iro.InputPath, iro.FileName)
	iro.OutputFilePath = filepath.Join(iro.OutputPath, strings.Replace(iro.FileName, fmt.Sprintf(".%s", imgFormat), fmt.Sprintf("-%v-%d.%s", ts, iro.Quality, targetImageFormat), 1))
	iro.Debug = debug
	return true, nil
}

type ImageManipulation struct {
	options ImageManipulationOptions
}

func (ir *ImageManipulation) Resize(basePath string, inputPath string, outputPath string, filename string, width float64, height float64, quality int, debug bool) (string, error) {
	// set options value
	mt := strings.Split(filename, ".")
	imgFormat := mt[len(mt)-1]
	_, err := ir.options.init(basePath, inputPath, outputPath, filename, width, height, quality, imgFormat, debug)
	if err != nil {
		return "", err
	}
	// main logic
	src := gocv.IMRead(ir.options.InputFilePath, gocv.IMReadColor)
	if src.Empty() {
		return "", errors.New("failed to read input file")
	}
	transform := gocv.NewMat()
	fx := ir.options.Width / float64(src.Cols())
	fy := ir.options.Height / float64(src.Rows())
	gocv.Resize(src, &transform, image.Point{}, fx, fy, gocv.InterpolationCubic)

	if ok := gocv.IMWrite(ir.options.OutputFilePath, transform); !ok {
		return "", errors.New("failed to write output file")
	}
	return ir.options.OutputFilePath, nil
}

func (ir *ImageManipulation) Compress(basePath string, inputPath string, outputPath string, filename string, quality int, debug bool) (string, error) {
	// set options value
	_, err := ir.options.init(basePath, inputPath, outputPath, filename, -1, -1, quality, "jpeg", debug)
	if err != nil {
		return "", err
	}
	// main logic
	src := gocv.IMRead(ir.options.InputFilePath, gocv.IMReadColor)
	if src.Empty() {
		return "", errors.New("failed to read input file")
	}
	if ok := gocv.IMWriteWithParams(ir.options.OutputFilePath, src, []int{gocv.IMWriteJpegQuality, ir.options.Quality}); !ok {
		return "", errors.New("failed to write output file")
	}
	return ir.options.OutputFilePath, nil
}

func (ir *ImageManipulation) PngToJpeg(basePath string, inputPath string, outputPath string, filename string, debug bool) (string, error) {
	// set options value
	_, err := ir.options.init(basePath, inputPath, outputPath, filename, -1, -1, 100, "jpeg", debug)
	if err != nil {
		return "", err
	}
	src := gocv.IMRead(ir.options.InputFilePath, gocv.IMReadColor)
	if src.Empty() {
		return "", errors.New("failed to read input file")
	}
	// main logic
	if ok := gocv.IMWriteWithParams(ir.options.OutputFilePath, src, []int{gocv.IMWriteJpegQuality, ir.options.Quality}); !ok {
		return "", errors.New("failed to write output file")
	}
	return ir.options.OutputFilePath, nil
}
