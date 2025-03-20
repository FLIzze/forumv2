package forum

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"

	structs "forum/structs"
	utils "forum/utils"
)

func GetTopic(c echo.Context) error {
        response := structs.TopicResponse{}

        UUID := c.Param("uuid")
        user, ok := c.Get("user").(structs.User)
        if ok {
                response.User = user
        }

        db := c.Get("db").(*sql.DB)

        row := db.QueryRow(`
        SELECT 
                UUID, Name, Description, CreatedByUsername
        FROM 
                topicInfo
        WHERE 
                UUID = ?
        `, UUID)

        err := row.Scan(&response.Subject.UUID, &response.Subject.Name, &response.Subject.Description, 
                                                                        &response.Subject.CreatedByUsername)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                response.Status.Error = "Could not retrieve topic."
                return c.Render(500, "topic", response)
        }

        rows, err := db.Query(`
        SELECT 
                UUID, Content, CreatedByUsername, CreatedByUUID, CreationTime
        FROM 
                messageInfo
        WHERE 
                TopicUUID = ?
        `, UUID)
        if err != nil {
                c.Logger().Error("Error retrieving topic message: ", err)
                response.Status.Error = "Could not retrieve topic message."
                return c.Render(500, "topic", response)
        }
        defer rows.Close()


        for rows.Next() {
                message := structs.Message{}
                err := rows.Scan(&message.UUID, &message.Content, &message.CreatedByUsername, &message.CreatedByUUID, 
                                                                                                &message.CreationTime)
                if err != nil {
                        c.Logger().Error("Error retrieving topic message from column: ", err)
                        response.Status.Error = "Something went wrong. Please try again later."
                        return c.Render(500, "topic", response)
                }

                message.FormattedCreationTime = utils.FormatDate(message.CreationTime)
                response.Messages = append(response.Messages, message)
        }

        return c.Render(200, "topic", response)
}

func PostMessage(c echo.Context) error {
        response := structs.TopicResponse{}
        message := structs.Message{}

        message.UUID = utils.Uuid()
        message.TopicUUID = c.FormValue("uuid")
        message.Content = utils.RenderMarkdown(c.FormValue("message"))

        if message.Content == "" {
                c.Logger().Error("Message empty")
                response.Status.Error = "Message must be filled"
                return c.Render(422, "topic-status", response)
        }

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(structs.User)
        if ok {
                response.User = user
        }

        _, err := db.Exec(`
        INSERT INTO message (UUID, TopicUUID, Content, CreatedBy, CreationTime) 
        VALUES (?, ?, ?, ?, ?)
        `, message.UUID, message.TopicUUID, message.Content, user.UUID, time.Now())
        if err != nil {
                c.Logger().Error("Error inserting into message", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topic-status", response)
        }

        row := db.QueryRow(`
        SELECT 
                CreatedByUsername, CreatedByUUID
        FROM 
                messageInfo
        WHERE
                uuid = ?
        `, message.UUID)
        err = row.Scan(&message.CreatedByUsername, &message.CreatedByUUID)
        if err != nil {
                c.Logger().Error("Error retrieving from messageInfo", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topic-status", response)
        }

        response.Messages = append(response.Messages, message)
        response.Status.Success = "Message sucessfully posted."

        c.Render(200, "topic-status", response)
        return c.Render(200, "oob-message", response)
}

func DeleteMessage(c echo.Context) error {
        response := structs.TopicResponse{}

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(structs.User)
        if !ok {
                response.Status.Error = "You must be logged in to delete a topic."
                return c.Render(401, "topic-status", response)
        }

        createdBy := c.FormValue("createdBy")
        if createdBy != user.UUID {
                response.Status.Error = "You must own the message to delete it."
                return c.Render(401, "topic-status", response)
        }

        uuid := c.FormValue("uuid")

        _, err := db.Exec(`
        DELETE FROM 
                message
        WHERE 
               uuid = ? 
        `, uuid)
        if err != nil {
                c.Logger().Error("Error deleting from message", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topic-status", response)
        }

        return c.NoContent(200)
}
