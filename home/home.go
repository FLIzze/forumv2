package forum

import (
	"github.com/labstack/echo/v4"

        structs "forum/structs"
        utils "forum/utils"
)

func GetHomePage(c echo.Context) error {
        response := structs.HomeResponse{}

        topics, err := utils.GetTopics(c)
        err.HandleError(c)

        response.Topics = topics
        user, ok := c.Get("user").(structs.User)
        if ok {
                response.User = user
        }


        return c.Render(200, "home", response)
}
