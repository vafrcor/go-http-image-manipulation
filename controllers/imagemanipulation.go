package controllers

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vafrcor/go-http-image-manipulation/models"
	helpers "github.com/vafrcor/go-http-image-manipulation/services"
	_ "golang.org/x/image/bmp"
)

func ValidateImageFileUpload(c echo.Context, allowedFormat []string) (map[string]string, error) {
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

	if !slices.Contains(allowedFormat, imageType) {
		// fmt.Printf("invalid mime %s\n", imageType)
		msg := fmt.Sprintf("only accept image using specific format (%s)", strings.Join(allowedFormat, ","))
		return nil, errors.New(msg)
	}

	// Return data for next process
	data["filename"] = file.Filename
	data["cwd"] = cwd
	data["base_upload_path"] = baseUploadPath
	data["upload_path"] = uploadPath
	data["output_path"] = outputPath
	return data, nil
}

func ImageConvertPngToJpeg(c echo.Context) error {
	data, err := ValidateImageFileUpload(c, []string{"png"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: err.Error(),
			Status:  false,
		})
	}
	// fmt.Printf("DATA: %#v\n", data)
	im := helpers.ImageManipulation{}
	output, err := im.PngToJpeg(data["cwd"], data["upload_path"], data["output_path"], data["filename"], false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &models.Response{
		Message: "Ok",
		Status:  true,
		Data:    fmt.Sprintf("%s://%s/static%s", helpers.GetEchoRequestScheme(c), c.Request().Host, strings.Replace(output, data["output_path"], "", 100)),
	})
}

func ImageResize(c echo.Context) error {
	keepAspecRatio := c.FormValue("keep_aspec_ratio")
	possibleAR := []string{"0", "1"}
	if keepAspecRatio == "" {
		keepAspecRatio = "1"
	}
	if !slices.Contains(possibleAR, keepAspecRatio) {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "invalid keep_aspec_ratio option value (choose either 1 or 0)",
			Status:  false,
		})
	}
	width := c.FormValue("width")
	height := c.FormValue("height")
	if width == "" || height == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "invalid width or height",
			Status:  false,
		})
	}
	if width == "" || height == "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "invalid width or height",
			Status:  false,
		})
	}
	widthFloat, _ := strconv.ParseFloat(width, 64)
	heightFloat, _ := strconv.ParseFloat(height, 64)
	if widthFloat < 0 || heightFloat < 0 {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "invalid width or height",
			Status:  false,
		})
	}

	data, err := ValidateImageFileUpload(c, []string{"png", "jpg", "jpeg", "bmp"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: err.Error(),
			Status:  false,
		})
	}
	// fmt.Printf("DATA: %#v\n", data)
	im := helpers.ImageManipulation{}
	keepAspectRatioBool, _ := strconv.ParseBool(keepAspecRatio)
	output, err := im.Resize(data["cwd"], data["upload_path"], data["output_path"], data["filename"], widthFloat, heightFloat, 100, keepAspectRatioBool, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &models.Response{
		Message: "Ok",
		Status:  true,
		Data:    fmt.Sprintf("%s://%s/static%s", helpers.GetEchoRequestScheme(c), c.Request().Host, strings.Replace(output, data["output_path"], "", 100)),
	})
}

func ImageCompress(c echo.Context) error {
	quality := c.FormValue("quality")
	queryError := ""
	if quality == "" {
		queryError = "invalid quality (must between 1 - 100)"
	}
	qualityInt, err := strconv.Atoi(quality)
	if qualityInt < 0 || err != nil {
		queryError = "invalid quality (must between 1 - 100)"
	}
	if queryError != "" {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: queryError,
			Status:  false,
		})
	}
	data, err := ValidateImageFileUpload(c, []string{"png", "jpg", "jpeg", "bmp"})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: err.Error(),
			Status:  false,
		})
	}
	// fmt.Printf("DATA: %#v\n", data)
	im := helpers.ImageManipulation{}
	output, err := im.Compress(data["cwd"], data["upload_path"], data["output_path"], data["filename"], qualityInt, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &models.Response{
		Message: "Ok",
		Status:  true,
		Data:    fmt.Sprintf("%s://%s/static%s", helpers.GetEchoRequestScheme(c), c.Request().Host, strings.Replace(output, data["output_path"], "", 100)),
	})
}
