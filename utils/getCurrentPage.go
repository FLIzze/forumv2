package forum

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetCurrentPage(c echo.Context) (int, error) {
        strPage := c.Param("nmb")
        intPage, err := strconv.Atoi(strPage)
        if err != nil || intPage < 1 {
                return -1, err
        }

        return intPage - 1, nil
}
