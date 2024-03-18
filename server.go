package main

import (
	"html/template"
	"io"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vafrcor/go-http-image-manipulation/controllers"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Initial setup
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "storages/public")
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// Routes Definition
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})
	e.POST("/image-png-to-jpeg", controllers.ImageConvertPngToJpeg)
	e.POST("/image-resize", controllers.ImageResize)
	e.POST("/image-compression", controllers.ImageCompress)

	// Run the application
	e.Logger.Fatal(e.Start(":9000"))
}
