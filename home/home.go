package forum

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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
                25
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

func PostTopic(c echo.Context) error {
        response := structs.HomeResponse{}
        topic := structs.Topic{}

        topic.UUID = uuid.New().String()
        topic.Name = c.FormValue("name")
        topic.Description = c.FormValue("message")
        user, ok := c.Get("user").(structs.User)
        if !ok {
                response.Status.Error = "You must be logged in to post a topic."
                return c.Render(401, "home-form", response)
        } else {
                response.User = user
        }

        if topic.Name == "" || topic.Description == "" {
                response.Status.Error = "Name and description must be filled."
                return c.Render(422, "home-form", response)
        }

        db := c.Get("db").(*sql.DB)

        _, err := db.Exec(`
        INSERT INTO topic (UUID, Name, Description, CreatedBy, CreationTime) 
        VALUES (?, ?, ?, ?, ?)
        `, topic.UUID, topic.Name, topic.Description, user.UUID, time.Now())
        if err != nil {
                c.Logger().Error("Error inserting response Topic: ", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "home-form", response)
        }

        row := db.QueryRow(`
        SELECT 
                CreatedByUsername, CreatedByUUID, NmbMessages 
        FROM 
                topicInfo
        WHERE 
                UUID = ?
        `, topic.UUID)
        err = row.Scan(&topic.CreatedByUsername, &topic.CreatedByUUID, &topic.NmbMessages)
        if err != nil {
                c.Logger().Error("Error fetching new topic: ", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "home-form", response)
        }

        response.Topics = append(response.Topics, topic)

        c.Render(200, "oob-topic", response)
        return c.Render(200, "home-form", response)
}
