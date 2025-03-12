package main

import (
	"html/template"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	er404 "forum/er404"
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
                os.Exit(1)
        }
        defer db.Close()

        // err = dbi.CreateTable(db)
        // if err != nil {
        //         fmt.Printf("Error creating table: %s", err)
        // }

        e := echo.New()
        e.Use(middleware.Logger())

        e.Renderer = newTemplate()

        e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
                return func(c echo.Context) error {
                        c.Set("db", db)
                        return next(c)
                }
        })

        e.GET("/", home.HomePage)
        e.POST("/postTopic", home.PostTopic)
        // e.DELETE("/topic/:uuid", home.DeleteTopic)

        e.GET("/topic/:uuid", topic.GetTopic) 
        e.POST("/postMessage", topic.PostMessage)

        e.GET("/*", er404.Get404)

        e.Logger.Fatal(e.Start(":42069"))
}
