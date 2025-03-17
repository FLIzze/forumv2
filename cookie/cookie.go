package forum

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func PostCookie(c echo.Context, value string) {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = value
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}

func GetCookie(c echo.Context) (*http.Cookie, error) {
        cookie, err := c.Cookie("session")
        return cookie, err
}

func RemoveCookie(c echo.Context) error {
    cookie := new(http.Cookie)
    cookie.Name = "session"
    cookie.Value = "" 
    cookie.Expires = time.Unix(0, 0) 
    cookie.Path = "/" 
    c.SetCookie(cookie)
    return nil
}
