package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vafrcor/go-http-image-manipulation/helpers"
	"github.com/vafrcor/go-http-image-manipulation/models"
)

func ValidateImageFileUpload(c echo.Context) (map[string]string, error) {
	data := map[string]string{
		"cwd":              "",
		"base_upload_path": "",
		"upload_path":      "",
		"output_path":      "",
		"filename":         "",
	}
	now := time.Now()
	ts := now.UnixNano()

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	// Validate Source
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Move File into destination directory
	cwd, _ := os.Getwd()

	baseUploadPath := filepath.Join(cwd, "storages", "uploads")
	uploadPath := filepath.Join(baseUploadPath, fmt.Sprintf("%d", ts))
	outputPath := filepath.Join(cwd, "storages", "public")
	_ = os.Mkdir(uploadPath, os.ModePerm)
	tempFilepath := filepath.Join(uploadPath, file.Filename)
	dst, err := os.Create(tempFilepath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Validate MimeType
	tempFile, err := os.Open(tempFilepath)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	_, imageType, err := image.Decode(tempFile)
	if err != nil {
		return nil, err
	}

	if imageType != "png" {
		fmt.Printf("Invalid mime %s\n", imageType)
		return nil, errors.New("only accept image using PNG format")
	}

	// Return data for next process
	data["filename"] = file.Filename
	data["cwd"] = cwd
	data["base_upload_path"] = baseUploadPath
	data["upload_path"] = uploadPath
	data["output_path"] = outputPath
	return data, nil
}

func ImageConvertToJpeg(c echo.Context) error {
	data, err := ValidateImageFileUpload(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: err.Error(),
			Status:  false,
		})
	}
	// fmt.Printf("DATA: %#v\n", data)
	ir := helpers.ImageResizer{}
	output, err := ir.Resize(data["cwd"], data["upload_path"], data["output_path"], data["filename"], 640, 480, 100, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &models.Response{
		Message: "Ok",
		Status:  true,
		Output:  fmt.Sprintf("Output: %s://%s/static%s", helpers.GetEchoRequestScheme(c), c.Request().Host, strings.Replace(output, data["output_path"], "", 100)),
	})
}

func ImageResize(c echo.Context) error {
	width := c.FormValue("width")
	height := c.FormValue("height")
	if width == "" || height == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "Invalid width or height",
			Status:  false,
		})
	}
	widthFloat, _ := strconv.ParseFloat(width, 64)
	heightFloat, _ := strconv.ParseFloat(height, 64)
	if widthFloat < 0 || heightFloat < 0 {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "Invalid width or height",
			Status:  false,
		})
	}

	data, err := ValidateImageFileUpload(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: err.Error(),
			Status:  false,
		})
	}
	// fmt.Printf("DATA: %#v\n", data)
	ir := helpers.ImageResizer{}
	output, err := ir.Resize(data["cwd"], data["upload_path"], data["output_path"], data["filename"], widthFloat, heightFloat, 100, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &models.Response{
		Message: "Ok",
		Status:  true,
		Output:  fmt.Sprintf("Output: %s://%s/static%s", helpers.GetEchoRequestScheme(c), c.Request().Host, strings.Replace(output, data["output_path"], "", 100)),
	})
}

func ImageCompress(c echo.Context) error {
	quality := c.FormValue("quality")
	queryError := ""
	if quality == "" {
		queryError = "Invalid quality (must between 0 - 100)"
	}
	qualityInt, err := strconv.Atoi(quality)
	if qualityInt < 0 || err != nil {
		queryError = "Invalid quality (must between 0 - 100)"
	}
	if queryError != "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: queryError,
			Status:  false,
		})
	}
	data, err := ValidateImageFileUpload(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: err.Error(),
			Status:  false,
		})
	}
	// fmt.Printf("DATA: %#v\n", data)
	ir := helpers.ImageResizer{}
	output, err := ir.Resize(data["cwd"], data["upload_path"], data["output_path"], data["filename"], 640, 480, qualityInt, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &models.Response{
		Message: "Ok",
		Status:  true,
		Output:  fmt.Sprintf("Output: %s://%s/static%s", helpers.GetEchoRequestScheme(c), c.Request().Host, strings.Replace(output, data["output_path"], "", 100)),
	})
}

func main() {
	// Initial setup
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "storages/public")

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<!DOCTYPE html><body><h1>Welcome to HTTP Image Resizer by vin(PNG to JPG)</h1><br><p>There are 3 available endpoints:<ol><li>Basic Image Resizer </li><li>Image Resizerwith specific dimension</li><li>Image Resizer with compression</li></ol></p></body></html>")
	})
	e.POST("/image-png-to-jpeg", ImageConvertToJpeg)
	e.POST("/image-resize", ImageResize)
	e.POST("/image-compression", ImageCompress)

	// Run the application
	e.Logger.Fatal(e.Start(":9000"))
}
