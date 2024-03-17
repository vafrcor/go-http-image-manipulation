package helpers

import (
	"github.com/labstack/echo/v4"
)

func GetEchoRequestScheme(c echo.Context) string {
	httpscheme := "https"
	if c.Request().TLS == nil {
		httpscheme = "http"
	}
	return httpscheme
}
