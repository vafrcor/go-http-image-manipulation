package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetEchoRequestScheme(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	req1 := httptest.NewRequest(http.MethodGet, "http://localhost:9000", nil)
	rec1 := httptest.NewRecorder()
	c1 := e.NewContext(req1, rec1)
	c1.SetPath("/test")
	assert.Equal("http", GetEchoRequestScheme(c1), "they should be equal")

	req2 := httptest.NewRequest(http.MethodGet, "https://localhost:9000", nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	c2.SetPath("/test")
	assert.Equal("https", GetEchoRequestScheme(c2), "they should be equal")
}
