package main

import (
        "html/template"
        "io"
        "github.com/labstack/echo/v4"
        "github.com/labstack/echo/v4/middleware"
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

func newTopic(name, description string) Topic {
        return Topic{
                Name: name,
                Description: description,
        }
}

type Topics = []Topic

type Data struct {
        Topics Topics
}

type Topic struct {
        Name string
        Description string
}

func newData() Data {
        return Data{
                Topics: []Topic{
                        newTopic("Basketball", "I like balls"),
                        newTopic("Climbing", "Go go"),
                },
        }
}

func main() {
        e := echo.New()
        e.Use(middleware.Logger())

        e.Renderer = newTemplate()
        data := newData()

        e.GET("/", func(c echo.Context) error {
                return c.Render(200, "index", data)
        })

        e.POST("/add-topic", func(c echo.Context) error {
                name := c.FormValue("name")
                description := c.FormValue("description")

                topic := newTopic(name, description)
                data.Topics = append(data.Topics, topic)

                return c.Render(200, "oob-topic", topic)
        })

        e.Logger.Fatal(e.Start(":42069"))
}
