package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vafrcor/go-http-image-manipulation/helpers"
)

// func validateUpload() {

// }

func basicResize(c echo.Context) error {
	now := time.Now()
	ts := now.UnixNano()

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Source
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	defer src.Close()

	// Destination
	cwd, _ := os.Getwd()
	baseUploadPath := filepath.Join(cwd, "storages", "uploads")
	uploadPath := filepath.Join(baseUploadPath, fmt.Sprintf("%d", ts))
	outputPath := filepath.Join(cwd, "storages", "public")
	_ = os.Mkdir(uploadPath, os.ModePerm) // 0755
	tempFilepath := filepath.Join(uploadPath, file.Filename)
	dst, err := os.Create(tempFilepath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// validate
	tempFile, err := os.Open(tempFilepath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	defer tempFile.Close()
	_, imageType, err := image.Decode(tempFile)
	if err != nil {
		// return err
		return c.JSON(http.StatusBadRequest, err)
	}

	if imageType != "png" {
		return c.JSON(http.StatusBadRequest, "Only accept image using PNG format")
	}

	// do main logic
	ir := helpers.ImageResizer{}
	output, err := ir.Resize(cwd, uploadPath, outputPath, file.Filename, 640, 480, 100, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("Output: %s://%s/static%s", helpers.GetEchoRequestScheme(c), c.Request().Host, strings.Replace(output, outputPath, "", 100)))
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
	e.POST("/basic-resizer", basicResize)

	// Run the application
	e.Logger.Fatal(e.Start(":9000"))
}
