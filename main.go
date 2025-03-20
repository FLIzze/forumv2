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
	// "github.com/microcosm-cc/bluemonday"

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
        // p := bluemonday.UGCPolicy()

        safeHTML := template.FuncMap{
                "safeHTML": func(s string) template.HTML {
                        // sanitized := p.Sanitize(s) 
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
        e.Use(middleware.Logger())

        e.Static("/css", "css")

        e.Renderer = newTemplate()

        db := dbi.HandleDbSetup()

        e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
                return func(c echo.Context) error {
                        c.Set("db", db)
                        return next(c)
                }
        })
        e.Use(mw.AuthMiddleware)

        e.GET("/", home.GetHomePage)
        e.GET("/topic/:uuid", topic.GetTopic) 
        e.GET("/login", user.GetLogin)
        e.GET("/register", user.GetRegister)
        e.GET("/me", user.GetMeProfil)
        e.GET("/user/:username", user.GetProfil)
        e.GET("/*", er404.Get404)

        e.POST("/login", user.PostLogin)
        e.POST("/register", user.PostRegister)

        e.POST("/topic", home.PostTopic)
        e.POST("/message", topic.PostMessage)
        e.POST("/logout", user.LogOut)

        e.DELETE("/message", topic.DeleteMessage)
        e.DELETE("/topic", home.DeleteTopic)

        PORT := os.Getenv("PORT")
        if PORT == "" {
                PORT = "8080"  
        }
        e.Logger.Fatal(e.Start(":"+ PORT))
}
