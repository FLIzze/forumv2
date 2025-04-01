package forum

import (
	"strconv"

	"github.com/labstack/echo/v4"

	structs "forum/structs"
)

func GetCurrentPage(c echo.Context) (int, structs.Error) {
        strPage := c.Param("nmb")
        intPage, err := strconv.Atoi(strPage)
        if err != nil || intPage < 1 {
		c.Logger().Debug("Error retrieving currentPage")
                return 0, structs.NewError(500, err)
        }

	return intPage - 1, structs.NewError(200, err)
}
