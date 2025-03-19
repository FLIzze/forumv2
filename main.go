package main

import (
	"html/template"
	"io"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
        "github.com/Masterminds/sprig/v3"

	er404 "forum/er404"
	home "forum/home"
	topic "forum/topic"
        user "forum/user"
        mw "forum/middleware"
)

type Templates struct {
        templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
        return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
        funcMap := sprig.FuncMap()

	return &Templates{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html")),
	}
}

func main() {
        e := echo.New()
        e.Use(middleware.Logger())

        e.Static("/css", "css")

        e.Renderer = newTemplate()

        e.Use(mw.DBMiddleware)
        e.Use(mw.AuthMiddleware)

        e.GET("/", home.GetHomePage)
        e.GET("/topic/:uuid", topic.GetTopic) 
        e.GET("/login", user.GetLogin)
        e.GET("/register", user.GetRegister)
        e.GET("/*", er404.Get404)

        e.POST("/login", user.PostLogin)
        e.POST("/register", user.PostRegister)

        e.POST("/topic", home.PostTopic)
        e.POST("/message", topic.PostMessage)
        e.POST("/logout", user.LogOut)

        e.DELETE("/message", topic.DeleteMessage)
        e.DELETE("/topic", home.DeleteTopic)

        e.GET("/user/:username", user.GetProfil)

        e.Logger.Fatal(e.Start(":42069"))
}
