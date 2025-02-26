package main

import (
        "html/template"
        "io"

        _ "github.com/go-sql-driver/mysql"
        "github.com/labstack/echo/v4"
        "github.com/labstack/echo/v4/middleware"

        topics "forum/topics"
)

type Templates struct {
        templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
        return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
        return &Templates{
                templates: template.Must(template.ParseGlob("views/*.html")),
        }
}

func main() {
        e := echo.New()
        e.Use(middleware.Logger())

        e.Renderer = newTemplate()

        e.GET("/*", topics.GetTopic) 
        e.GET("/", topics.GetTopics)
        e.POST("/postTopic", topics.PostTopic)

        e.Logger.Fatal(e.Start(":42069"))
}
