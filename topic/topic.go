package forum

import (
	"database/sql"
	"html/template"
	"time"

	"github.com/labstack/echo/v4"

	user "forum/user"
	utils "forum/utils"
)

type TopicResponse struct {
        Error string
        Subject Subject
        Messages []Message
        User user.User
}

type Subject struct {
        UUID string
        Name string
        Description string
}

type Message struct {
        UUID string
        TopicUUID string
        Content template.HTML
        CreatedBy string
}

func GetTopic(c echo.Context) error {
        response := TopicResponse{}

        UUID := c.Param("uuid")
        user, ok := c.Get("user").(user.User)
        if !ok {
                c.Logger().Debug("User is not logged in")
        } else {
                response.User = user
        }

        db := c.Get("db").(*sql.DB)

        row := db.QueryRow(`
        SELECT UUID, Name, Description 
        FROM topic
        WHERE UUID = ?
        `, UUID)

        err := row.Scan(&response.Subject.UUID, &response.Subject.Name, &response.Subject.Description)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                response.Error = "Could not retrieve topic."
                return c.Render(500, "topic", response)
        }

        rows, err := db.Query(`
        SELECT Content, CreatedBy
        FROM messageInfo
        WHERE TopicUUID = ?
        `, UUID)
        if err != nil {
                c.Logger().Error("Error retrieving topic message: ", err)
                response.Error = "Could not retrieve topic message."
                return c.Render(500, "topic", response)
        }
        defer rows.Close()


        for rows.Next() {
                message := Message{}
                err := rows.Scan(&message.Content, &message.CreatedBy)
                if err != nil {
                        c.Logger().Error("Error retrieving topic message from column: ", err)
                        response.Error = "Internal server error"
                        return c.Render(500, "topic", response)
                }

                response.Messages = append(response.Messages, message)
        }

        return c.Render(200, "topic", response)
}

func PostMessage(c echo.Context) error {
        response := TopicResponse{}
        message := Message{}

        message.UUID = utils.Uuid()
        message.TopicUUID = c.FormValue("uuid")
        message.Content = utils.RenderMarkdown(c.FormValue("message"))

        if message.Content == "" {
                c.Logger().Error("Message empty")
                response.Error = "Message must be filled"
                return c.Render(422, "topic-form", response)
        }

        response.Messages = append(response.Messages, message)

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(user.User)
        if !ok {
                c.Logger().Debug("User is not logged in")
        } else {
                response.User = user
        }

        _, err := db.Exec(`
        INSERT INTO message (UUID, TopicUUID, Content, CreatedBy, CreationTime) 
        VALUES (?, ?, ?, ?, ?)
        `, message.UUID, message.TopicUUID, message.Content, user.UUID, time.Now())
        if err != nil {
                c.Logger().Error("Error retrieving response.Topic: ", err)
                response.Error = "Internal server error"
                return c.Render(500, "topic-form", response)
        }

        c.Render(200, "topic-form", response)
        return c.Render(200, "oob-created-topic", message)
}
