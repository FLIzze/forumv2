package forum

import (
	"database/sql"

	"github.com/labstack/echo/v4"

        structs "forum/structs"
        utils "forum/utils"
)

func GetHomePage(c echo.Context) error {
        response := structs.HomeResponse{}
        topic := structs.Topic{}

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(structs.User)
        if ok {
                response.User = user
        }

        rows, err := db.Query(`
        SELECT 
                UUID, Name, Description, CreatedByUsername, CreatedByUUID, NmbMessages, LastMessage, CreationTime
        FROM 
                topicInfo
        LIMIT
                35
        `)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "home", response)
        }
        defer rows.Close()

        for rows.Next() {
                err := rows.Scan(&topic.UUID, &topic.Name, &topic.Description, &topic.CreatedByUsername, 
                        &topic.CreatedByUUID, &topic.NmbMessages, &topic.LastMessage, &topic.CreationTime)
                if err != nil {
                        c.Logger().Error("Error scanning row", err)
                        response.Status.Error = "Something went wrong. Please try again later."
                        return c.Render(422, "home", response)
                }

                topic.FormattedCreationTime = utils.FormatDate(topic.CreationTime)
                topic.FormattedLastMessage = utils.FormatDate(topic.LastMessage)
                response.Topics = append(response.Topics, topic)
        }

        return c.Render(200, "home", response)
}
