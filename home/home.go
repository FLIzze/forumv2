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

	currentPage, error := utils.GetCurrentPage(c)
        error.HandleError(c)
	response.Page.CurrentPage = currentPage + 1

	response.Page.TotalPage = len(topics)
        response.Topics = topics

        user, ok := c.Get("user").(structs.User)
        if ok {
                response.User = user
        }

        return c.Render(200, "home", response)
}
