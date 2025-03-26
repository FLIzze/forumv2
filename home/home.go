package forum

import (
	"database/sql"
        "strconv"

	"github.com/labstack/echo/v4"

        structs "forum/structs"
        utils "forum/utils"
)

func GetHomePage(c echo.Context) error {
        response := structs.HomeResponse{}
        topic := structs.Topic{}
        page := structs.Page{}

        strPage := c.Param("nmb")
        intPage, err := strconv.Atoi(strPage)
        if err != nil || intPage < 1 {
                c.Logger().Error("Error during Atoi: ", err)
                response.Status.Error = "Invalid Page number"
                return c.Render(400, "home", response)
        }

        page.CurrentPage = intPage
        intPage -= 1

        MAX_TOPICS_DISPLAYED := 30
        OFFSET := MAX_TOPICS_DISPLAYED * intPage
        LIMIT := MAX_TOPICS_DISPLAYED

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(structs.User)
        if ok {
                response.User = user
        }

        var totalTopics int
        err = db.QueryRow(`SELECT COUNT(*) FROM topicInfo`).Scan(&totalTopics)
        if err != nil {
                c.Logger().Error("Error counting topics: ", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "home", response)
        }

        page.TotalPage = (totalTopics + MAX_TOPICS_DISPLAYED - 1) / MAX_TOPICS_DISPLAYED

        rows, err := db.Query(`
        SELECT 
                UUID, Name, Description, CreatedByUsername, CreatedByUUID, NmbMessages, LastMessage, CreationTime
        FROM 
                topicInfo
        LIMIT
                ?
        OFFSET
                ?
        `, LIMIT, OFFSET)
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
                response.Page = page
        }

        return c.Render(200, "home", response)
}
