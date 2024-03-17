package helpers

import (
	"net/url"

	"github.com/labstack/echo/v4"
)

func GetBaseUrl(inputUrl string) string {
	baseUrl, err := url.Parse(inputUrl)
	if err != nil {
		return ""
	}
	baseUrl.Path = ""
	baseUrl.Fragment = ""
	return baseUrl.String()
}

func GetEchoRequestScheme(c echo.Context) string {
	httpscheme := "https"
	if c.Request().TLS == nil {
		httpscheme = "http"
	}
	return httpscheme
}
