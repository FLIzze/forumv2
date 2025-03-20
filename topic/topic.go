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
        subject := structs.Subject{}

        UUID := c.Param("uuid")
        user, ok := c.Get("user").(structs.User)
        if ok {
                response.User = user
        }

        db := c.Get("db").(*sql.DB)

        row := db.QueryRow(`
        SELECT 
                UUID, Name, Description, CreatedByUsername, LastMessage
        FROM 
                topicInfo
        WHERE 
                UUID = ?
        `, UUID)

        err := row.Scan(&subject.UUID, &subject.Name, &subject.Description, &subject.CreatedByUsername, &subject.LastMessage)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                response.Status.Error = "Could not retrieve topic."
                return c.Render(500, "topic", response)
        }
        subject.FormattedLastMessage = utils.FormatDate(subject.LastMessage)
        response.Subject = subject

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
        message.Content = c.FormValue("message")

        user, ok := c.Get("user").(structs.User)
        if !ok {
                response.Status.Error = "You must be logged in to post a message."
                return c.Render(401, "topic-form", response)
        } else {
                response.User = user
        }

        if message.Content == "" {
                c.Logger().Error("Message empty")
                response.Status.Error = "Message must be filled"
                return c.Render(422, "topic-form", response)
        }

        db := c.Get("db").(*sql.DB)

        _, err := db.Exec(`
        INSERT INTO message (UUID, TopicUUID, Content, CreatedBy, CreationTime) 
        VALUES (?, ?, ?, ?, ?)
        `, message.UUID, message.TopicUUID, message.Content, user.UUID, time.Now())
        if err != nil {
                c.Logger().Error("Error inserting into message: ", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topic-form", response)
        }

        row := db.QueryRow(`
        SELECT 
                CreatedByUsername, CreatedByUUID, TopicUUID
        FROM 
                messageInfo
        WHERE
                uuid = ?
        `, message.UUID)
        err = row.Scan(&message.CreatedByUsername, &message.CreatedByUUID, &response.Subject.UUID)
        if err != nil {
                c.Logger().Error("Error retrieving from messageInfo: ", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topic-form", response)
        }

        response.Messages = append(response.Messages, message)
        response.Status.Success = "Message sucessfully posted."

        c.Render(200, "oob-message", response)
        return c.Render(200, "topic-form", response)
}

func DeleteMessage(c echo.Context) error {
        response := structs.TopicResponse{}

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(structs.User)
        if !ok {
                response.Status.Error = "You must be logged in to delete a topic."
                return c.Render(401, "topic-form", response)
        }

        createdBy := c.FormValue("createdBy")
        if createdBy != user.UUID {
                response.Status.Error = "You must own the message to delete it."
                return c.Render(401, "topic-form", response)
        }

        uuid := c.FormValue("uuid")

        _, err := db.Exec(`
        DELETE FROM 
                message
        WHERE 
               uuid = ? 
        `, uuid)
        if err != nil {
                c.Logger().Error("Error deleting from message: ", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topic-form", response)
        }

        return c.NoContent(200)
}
