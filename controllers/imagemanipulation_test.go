package controllers

import (
	"bytes"
	"encoding/json"
	"image"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"testing"

	_ "golang.org/x/image/bmp"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vafrcor/go-http-image-manipulation/models"
)

func TestImageManipulationImageConvertPngToJpeg(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample-test.png")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// writer.WriteField("bu", "HFL")
	// writer.WriteField("wk", "10")
	part, _ := writer.CreateFormFile("file", "sample-test.png")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	png.Encode(body, imageData)
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-png-to-jpeg")

	if assert.NoError(t, ImageConvertPngToJpeg(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, true, len(rec.Body.String()) > 0)
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
	}
}

func TestImageManipulationImageConvertPngToJpegInvalidFileFormat(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample.gif")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "sample.gif")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	gif.Encode(body, imageData, &gif.Options{NumColors: 256})
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-png-to-jpeg")

	if assert.NoError(t, ImageConvertPngToJpeg(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, true, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		if err == nil {
			assert.True(t, data.Status == false)
			assert.Equal(t, "only accept image using specific format (png)", data.Message)
		}
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
	}
}

func TestImageManipulationImageConvertPngToJpegInvalidFileContent(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample-test.png")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "sample-test.png")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	// imageData, _, _ := image.Decode(testFile)

	// png.Encode(body, imageData)
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-png-to-jpeg")

	if assert.NoError(t, ImageConvertPngToJpeg(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, true, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		if err == nil {
			assert.True(t, data.Status == false)
			assert.Equal(t, "image: unknown format", data.Message)
		}
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
	}
}

func TestImageManipulationImageResize(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample-test.png")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("width", "480")
	writer.WriteField("height", "320")
	// writer.WriteField("keep_aspect_ratio", "1")
	part, _ := writer.CreateFormFile("file", "sample-test.png")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	png.Encode(body, imageData)
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-resize")

	if assert.NoError(t, ImageResize(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, true, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		if err == nil {
			assert.True(t, data.Status)
		}
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
		// {"message":"invalid width or height","status":false,"data":null}
	}
}

func TestImageManipulationImageResizeInvalidFileFormat(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample.gif")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("width", "480")
	writer.WriteField("height", "360")
	writer.WriteField("keep_aspect_ratio", "1")
	part, _ := writer.CreateFormFile("file", "sample.gif")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	gif.Encode(body, imageData, &gif.Options{NumColors: 256})
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-resize")

	if assert.NoError(t, ImageResize(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, true, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		if err == nil {
			assert.True(t, data.Status == false)
			assert.Equal(t, "only accept image using specific format (png,jpg,jpeg,bmp)", data.Message)
		}
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
	}
}

func TestImageManipulationImageResizeInvalidDimension1(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample-test.png")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("width", "")
	writer.WriteField("height", "")
	writer.WriteField("keep_aspect_ratio", "1")
	part, _ := writer.CreateFormFile("file", "sample-test.png")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	png.Encode(body, imageData)
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-resize")

	if assert.NoError(t, ImageResize(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, true, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		if err == nil {
			assert.True(t, data.Status == false)
			assert.Equal(t, "invalid width or height", data.Message)
		}
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
	}
}

func TestImageManipulationImageResizeInvalidDimension2(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample-test.png")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("width", "-480")
	writer.WriteField("height", "-320")
	writer.WriteField("keep_aspect_ratio", "1")
	part, _ := writer.CreateFormFile("file", "sample-test.png")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	png.Encode(body, imageData)
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-resize")

	if assert.NoError(t, ImageResize(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, true, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		if err == nil {
			assert.True(t, data.Status == false)
			assert.Equal(t, "invalid width or height", data.Message)
		}
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
	}
}

func TestImageManipulationImageResizeInvalidKeepAspecRatio(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample-test.png")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("width", "480")
	writer.WriteField("height", "320")
	writer.WriteField("keep_aspect_ratio", "4")
	part, _ := writer.CreateFormFile("file", "sample-test.png")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	png.Encode(body, imageData)
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-resize")

	if assert.NoError(t, ImageResize(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, true, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		if err == nil {
			assert.True(t, data.Status == false)
			assert.Equal(t, "invalid keep_aspect_ratio option value (choose either 1 or 0)", data.Message)
		}
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
	}
}

func TestImageManipulationImageCompress(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample-test.png")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("quality", "80")
	part, _ := writer.CreateFormFile("file", "sample-test.png")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	png.Encode(body, imageData)
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-compression")

	if assert.NoError(t, ImageCompress(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, true, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		if err == nil {
			assert.True(t, data.Status)
		}
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
		// {"message":"invalid width or height","status":false,"data":null}
	}
}

func TestImageManipulationImageCompressInvalidQuality(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample-test.png")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("quality", "")
	part, _ := writer.CreateFormFile("file", "sample-test.png")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	png.Encode(body, imageData)
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-compression")

	if assert.NoError(t, ImageCompress(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		if err == nil {
			assert.True(t, data.Status == false)
			assert.Equal(t, "invalid quality (must between 1 - 100)", data.Message)
		}
	}
}

func TestImageManipulationImageCompressInvalidFileFormat(t *testing.T) {
	// Setup
	e := echo.New()
	cwd, _ := os.Getwd()
	rootDir := filepath.Clean(filepath.Join(cwd, ".."))
	testFilePath := filepath.Join(rootDir, "storages", "test", "sample.gif")
	// fmt.Printf("TEST FILE PATH: %s\n", testFilePath)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("quality", "80")
	part, _ := writer.CreateFormFile("file", "sample.gif")
	testFile, _ := os.Open(testFilePath)
	defer testFile.Close()
	imageData, _, _ := image.Decode(testFile)

	gif.Encode(body, imageData, &gif.Options{NumColors: 256})
	part.Write([]byte(body.Bytes()))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/image-compression")

	if assert.NoError(t, ImageCompress(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, len(rec.Body.String()) > 0)
		jsonData := []byte(rec.Body.Bytes())
		var data models.Response
		err := json.Unmarshal(jsonData, &data)
		// fmt.Printf("RESPONSE BODY: %s\n", rec.Body.String())
		if err == nil {
			assert.True(t, data.Status == false)
			assert.Equal(t, "only accept image using specific format (png,jpg,jpeg,bmp)", data.Message)
		}
	}
}
