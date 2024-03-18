package services

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageManipulationOptions(t *testing.T) {
	assert := assert.New(t)
	im := ImageManipulation{}

	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	baseUploadPath := filepath.Join(rootDir, "storages", "test")
	outputPath := filepath.Join(rootDir, "storages", "public")

	data := map[string]string{
		"cwd":              rootDir,
		"base_upload_path": baseUploadPath,
		"upload_path":      baseUploadPath,
		"output_path":      outputPath,
		"filename":         "sample-test.png",
	}

	options, err := im.options.init(data["cwd"], data["base_upload_path"], data["output_path"], data["filename"], 1024, 768, 100, "png", true, false)
	assert.NotEqual(nil, options, "Options should not be nil")
	assert.Equal(nil, err, "Error should be nil")

	options2, err2 := im.options.init("", "", "", "sample-test.png", 0, 0, 0, "png", false, false)
	assert.NotEqual(nil, options2, "Options 2 should not be nil")
	assert.Equal(nil, err2, "Error 2 should be nil")

	options3, err3 := im.options.init("", "", "", "", 0, 0, 0, "png", false, false)
	assert.Equal(false, options3, "Options 3 should be false")
	assert.Equal("invalid file name", err3.Error(), "Error 3 should not be nil and has an error about invalid file name")
}

func TestImageManipulationCalculateAspectRatioFit(t *testing.T) {
	assert := assert.New(t)
	im := ImageManipulation{}
	calculate := im.CalculateAspectRatioFit(2000, 1500, 1000, 700)
	assert.Equal(float64(933.3333333333334), calculate["width"], "Width should 933.3333333333334")
	assert.Equal(float64(700), calculate["height"], "Width should 700")
}

func TestImageManipulationPngToJpeg(t *testing.T) {
	assert := assert.New(t)
	im := ImageManipulation{}
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	baseUploadPath := filepath.Join(rootDir, "storages", "test")
	outputPath := filepath.Join(rootDir, "storages", "public")

	data := map[string]string{
		"cwd":              rootDir,
		"base_upload_path": baseUploadPath,
		"upload_path":      baseUploadPath,
		"output_path":      outputPath,
		"filename":         "sample-test.png",
	}
	// fmt.Printf("DATA: %#v\n", data)
	process, err := im.PngToJpeg(data["cwd"], data["upload_path"], data["output_path"], data["filename"], false)
	assert.Equal(nil, err, "Error should be nil")
	assert.NotEqual(nil, process, "Output should not be nil (file path)")
	fexist, _ := os.Stat(process)
	assert.True(fexist != nil, "File should exist")
	// remove output file
	e := os.Remove(process)
	if e != nil {
		panic(e)
	}

	process2, err2 := im.PngToJpeg(data["cwd"], data["upload_path"], data["output_path"], "", false)
	assert.Equal("", process2, "Process2 should have false value")
	assert.Equal("invalid file name", err2.Error(), "Error 2 should contain message")
}

func TestImageManipulationResizeKeepAspectRatio(t *testing.T) {
	assert := assert.New(t)
	im := ImageManipulation{}
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	baseUploadPath := filepath.Join(rootDir, "storages", "test")
	outputPath := filepath.Join(rootDir, "storages", "public")

	data := map[string]string{
		"cwd":              rootDir,
		"base_upload_path": baseUploadPath,
		"upload_path":      baseUploadPath,
		"output_path":      outputPath,
		"filename":         "sample-test.png",
	}
	fmt.Printf("DATA: %#v\n", data)
	process, err := im.Resize(data["cwd"], data["upload_path"], data["output_path"], data["filename"], 480, 460, 100, true, false)
	assert.Equal(nil, err, "Error should be nil")
	assert.NotEqual(nil, process, "Output should not be nil (file path)")
	fexist, _ := os.Stat(process)
	assert.True(fexist != nil, "File should exist")
	// remove output file
	e := os.Remove(process)
	if e != nil {
		panic(e)
	}

	process2, err2 := im.Resize(data["cwd"], data["upload_path"], data["output_path"], "", 480, 460, 100, true, false)
	assert.Equal("", process2, "Process2 should have false value")
	assert.Equal("invalid file name", err2.Error(), "Error 2 should contain message")
}

func TestImageManipulationResizeKeepArbitraryRatio(t *testing.T) {
	assert := assert.New(t)
	im := ImageManipulation{}
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	baseUploadPath := filepath.Join(rootDir, "storages", "test")
	outputPath := filepath.Join(rootDir, "storages", "public")

	data := map[string]string{
		"cwd":              rootDir,
		"base_upload_path": baseUploadPath,
		"upload_path":      baseUploadPath,
		"output_path":      outputPath,
		"filename":         "sample-test.png",
	}
	fmt.Printf("DATA: %#v\n", data)
	process, err := im.Resize(data["cwd"], data["upload_path"], data["output_path"], data["filename"], 480, 460, 100, false, false)
	assert.Equal(nil, err, "Error should be nil")
	assert.NotEqual(nil, process, "Output should not be nil (file path)")
	fexist, _ := os.Stat(process)
	assert.True(fexist != nil, "File should exist")
	// remove output file
	e := os.Remove(process)
	if e != nil {
		panic(e)
	}
}

func TestImageManipulationCompress(t *testing.T) {
	assert := assert.New(t)
	im := ImageManipulation{}
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	baseUploadPath := filepath.Join(rootDir, "storages", "test")
	outputPath := filepath.Join(rootDir, "storages", "public")

	data := map[string]string{
		"cwd":              rootDir,
		"base_upload_path": baseUploadPath,
		"upload_path":      baseUploadPath,
		"output_path":      outputPath,
		"filename":         "sample-test.png",
	}
	fmt.Printf("DATA: %#v\n", data)
	process, err := im.Compress(data["cwd"], data["upload_path"], data["output_path"], data["filename"], 80, false)
	assert.Equal(nil, err, "Error should be nil")
	assert.NotEqual(nil, process, "Output should not be nil (file path)")
	fexist, _ := os.Stat(process)
	assert.True(fexist != nil, "File should exist")
	// remove output file
	e := os.Remove(process)
	if e != nil {
		panic(e)
	}

	process2, err2 := im.Compress(data["cwd"], data["upload_path"], data["output_path"], "", 80, false)
	assert.Equal("", process2, "Process2 should have false value")
	assert.Equal("invalid file name", err2.Error(), "Error 2 should contain message")
}
