package main

import (
	"html/template"
	"io"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	home "forum/home"
	topic "forum/topic"

        "fmt"
        dbi "forum/db"
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
        db, err := dbi.ConnectDb()
        if err != nil {
                fmt.Printf("Error connecting to database: %s", err)
        }

        err = dbi.CreateTable(db)

        e := echo.New()
        e.Use(middleware.Logger())

        e.Renderer = newTemplate()

        e.GET("/topic/:uuid", topic.GetTopic) 

        e.GET("/", home.GetHomePage)
        e.POST("/postTopic", home.PostTopic)
        // e.DELETE("/topic/:uuid", home.DeleteTopic)

        e.POST("/postMessage/:topicUUID", topic.PostMessage)

        e.Logger.Fatal(e.Start(":42069"))
}
