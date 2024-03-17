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

type ImageResizerOptions struct {
	BasePath       string  `json:"base_path"`
	InputPath      string  `json:"input_path"`
	OutputPath     string  `json:"output_path"`
	FileName       string  `json:"file_name"`
	InputFilePath  string  `json:"input_file_path"`
	OutputFilePath string  `json:"output_file_path"`
	Width          float64 `json:"width"`
	Height         float64 `json:"height"`
	Quality        int64   `json:"quality"`
	Debug          bool    `json:"debug"`
}

func (iro *ImageResizerOptions) init(basePath string, inputPath string, outputPath string, filename string, width float64, height float64, quality int64, debug bool) (bool, error) {
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
	iro.InputFilePath = filepath.Join(iro.InputPath, iro.FileName)
	iro.OutputFilePath = filepath.Join(iro.OutputPath, strings.Replace(iro.FileName, ".png", fmt.Sprintf("-%v-%d.jpeg", ts, iro.Quality), 1))
	iro.Debug = debug
	return true, nil
}

type ImageResizer struct {
	options ImageResizerOptions
}

func (ir *ImageResizer) Resize(basePath string, inputPath string, outputPath string, filename string, width float64, height float64, quality int64, debug bool) (string, error) {
	// set options value
	_, err := ir.options.init(basePath, inputPath, outputPath, filename, width, height, quality, debug)
	if err != nil {
		return "", err
	}
	src := gocv.IMRead(ir.options.InputFilePath, gocv.IMReadColor)
	if src.Empty() {
		return "", errors.New("failed to read input file")
	}
	transform := gocv.NewMat()
	fx := ir.options.Width / float64(src.Cols())
	fy := ir.options.Height / float64(src.Rows())
	gocv.Resize(src, &transform, image.Point{}, fx, fy, gocv.InterpolationCubic)

	// ok := gocv.IMWrite(ir.options.OutputFilePath, transform)
	if ok := gocv.IMWriteWithParams(ir.options.OutputFilePath, transform, []int{gocv.IMWriteJpegQuality, int(ir.options.Quality)}); !ok {
		return "", errors.New("failed to write output file")
	}
	return ir.options.OutputFilePath, nil
}
