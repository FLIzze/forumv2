package main

import (
	"html/template"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	dbi "forum/db"
	er404 "forum/er404"
	home "forum/home"
	mw "forum/middleware"
	topic "forum/topic"
	user "forum/user"
)

type Templates struct {
        templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
        return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
        safeHTML := template.FuncMap{
                "safeHTML": func(s string) template.HTML {
                        return template.HTML(s) 
                },
        }

        return &Templates{
                templates: template.Must(template.New("").Funcs(safeHTML).ParseGlob("views/*.html")),
        }
}

func main() {
        err := godotenv.Load(".env")
        if err != nil {
                log.Fatal("Error loading .env file")        
        }

        e := echo.New()
        e.Static("/css", "css")  
        e.Static("/src", "src")

        e.Use(middleware.Logger())

        e.Renderer = newTemplate()

        db := dbi.HandleDbSetup()

        e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
                return func(c echo.Context) error {
                        c.Set("db", db)
                        return next(c)
                }
        })

        e.Use(mw.Auth)

        e.GET("/", home.GetHomePage)
        e.GET("/topic/:uuid", topic.GetTopic) 
        e.GET("/login", user.GetLogin)
        e.GET("/register", user.GetRegister)
        e.GET("/*", er404.Get404)

        e.POST("/login", user.PostLogin)
        e.POST("/register", user.PostRegister)

        authGroup := e.Group("")
        authGroup.Use(mw.RequireAuth) 

        { 
                authGroup.GET("/me", user.GetMeProfil)
                authGroup.GET("/user/:username", user.GetProfil)
                authGroup.GET("/postTopic", topic.GetPostTopic)

                authGroup.POST("/topic", topic.PostTopic)
                authGroup.POST("/message", topic.PostMessage)
                authGroup.POST("/quote", topic.QuoteMessage)
                authGroup.POST("/logout", user.LogOut)

                authGroup.DELETE("/message", topic.DeleteMessage)
                authGroup.DELETE("/topic", topic.DeleteTopic)
        }

        PORT := os.Getenv("PORT")
        if PORT == "" {
                PORT = "8080"  
        }
        e.Logger.Fatal(e.Start(":"+ PORT))
}
