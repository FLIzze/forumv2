package forum

import (
	"github.com/labstack/echo/v4"
        "github.com/google/uuid"
        dbi "forum/db"
)

type TopicResponse struct {
        Error string
        Subject Subject
        Messages []Message
}

type Subject struct {
        UUID string
        Name string
        Description string
}

type Message struct {
        UUID string
        TopicUUID string
        Content string
}

func GetTopic(c echo.Context) error {
        response := TopicResponse{}

        URI := c.Param("uuid")

        db, err := dbi.ConnectDb()
        if err != nil {
                c.Logger().Error("Error connecting to db: ", err)
                response.Error = "Could not connect to database."
                return c.Render(500, "topic", response)
        }
        defer db.Close()

        row := db.QueryRow(`
                SELECT UUID, Name, Description 
                FROM topic
                WHERE UUID = ?
        `, URI)

        err = row.Scan(&response.Subject.UUID, &response.Subject.Name, &response.Subject.Description)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                response.Error = "Could not retrieve topic."
                return c.Render(500, "topic", response)
        }

        return c.Render(200, "topic", response)
}

func PostMessage(c echo.Context) error {
        response := TopicResponse{}
        message := Message{}

        db, err := dbi.ConnectDb()
        if err != nil {
                c.Logger().Error("Error connecting to db: ", err)
                response.Error = "Could not connect to database."
                return c.Render(500, "topic-form", response)
        }
        defer db.Close()

        message.UUID = uuid.New().String()
        message.TopicUUID = c.Param("topicUUID")
        message.Content = c.FormValue("message")

        if message.Content == "" {
                c.Logger().Error("Message empty")
                response.Error = "Message must be filled"
                return c.Render(422, "topic-form", response)
        }

        _, err = db.Exec(`
        INSERT INTO message (UUID, TopicUUID, Content) 
        VALUES (?, ?, ?)
        `, message.UUID, message.TopicUUID, message.Content)
        if err != nil {
                c.Logger().Error("Error retrieving response.Topic: ", err)
                response.Error = "Error retrieving response.Topic."
                return c.Render(500, "topic-form", response)
        }

        return c.Render(200, "topic", nil)
}
