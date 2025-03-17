package main

import (
	"html/template"
	"io"
	"os"
        "fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	er404 "forum/er404"
	home "forum/home"
	topic "forum/topic"
	dbi "forum/db"
        user "forum/user"
        cookie "forum/cookie"
)

type Templates struct {
        templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
        return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	funcMap := template.FuncMap{
		"mod": func(a, b int) int { return a % b }, 
	}

	return &Templates{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html")),
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

        // err = dbi.CreateView(db)
        // if err != nil {
        //         fmt.Printf("Error creating view: %s", err)
        // }

        e := echo.New()
        e.Use(middleware.Logger())

        e.Renderer = newTemplate()

        var currentUser user.User

        e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
                return func(c echo.Context) error {
                        c.Set("db", db)
                        
                        cookie, err := cookie.GetCookie(c)
                        if err != nil {
                                c.Logger().Debug("User is not logged")
                                c.Set("user", nil)
                                return next(c)
                        }

                        currentUser.SessionUUID = cookie.Value

                        row := db.QueryRow(`
                        SELECT UserUUID, Username, Email
                        FROM userSession
                        WHERE SessionUUID = ?
                        `, currentUser.SessionUUID)

                        err = row.Scan(&currentUser.UUID, &currentUser.Username, &currentUser.Email)
                        if err != nil {
                                c.Logger().Error("Error retrieving user from session", err)
                        }

                        c.Set("user", currentUser)

                        return next(c)
                }
        })

        e.GET("/", home.GetHomePage)
        e.POST("/postTopic", home.PostTopic)

        e.GET("/topic/:uuid", topic.GetTopic) 
        e.POST("/postMessage", topic.PostMessage)

        e.GET("/*", er404.Get404)

        e.GET("/login", user.GetLogin)
        e.POST("/login", user.PostLogin)
        e.GET("/register", user.GetRegister)
        e.POST("/register", user.PostRegister)
        e.POST("/logout", user.LogOut)
        e.GET("/user/:username", user.Profil)

        e.Logger.Fatal(e.Start(":42069"))
}
