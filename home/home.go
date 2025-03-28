package forum

import (
	"github.com/labstack/echo/v4"

        structs "forum/structs"
        utils "forum/utils"
)

func GetHomePage(c echo.Context) error {
        response := structs.HomeResponse{}

        topics, err := utils.GetTopics(c)
        if err.IsError() {
                err.HandleError(c)
                response.Status.Error = err.Message
                return c.Render(err.Status, "home", response)
        }
        response.Topics = topics
        response.User = c.Get("user").(structs.User)

        return c.Render(200, "home", response)
}
