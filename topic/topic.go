package forum

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"

	structs "forum/structs"
	utils "forum/utils"
)

func GetTopicPage(c echo.Context) error {
        var response structs.TopicResponse

        topic, err := utils.GetTopic(c)
        if err.IsError() {
                err.HandleError(c)
                response.Status.Error = err.Message
                return c.Render(err.Status, "topic", response)
        }
        response.Topic = topic
        response.User = c.Get("user").(structs.User)

        messages, err := utils.GetMessages(c)
        if err.IsError() {
                err.HandleError(c)
                response.Status.Error = err.Message
                return c.Render(err.Status, "topic", response)
        }
        response.Messages = messages

        return c.Render(200, "topic", response)
}

func PostMessage(c echo.Context) error {
        response := structs.TopicResponse{}
        message := structs.Message{}

        message.UUID = utils.Uuid()
        message.TopicUUID = c.FormValue("uuid")
        messageContent := c.FormValue("message")
        message.Content = string(blackfriday.Run([]byte(messageContent)))
        date := time.Now()
        message.FormattedCreationTime = utils.FormatDate(&date)

        user, ok := c.Get("user").(structs.User)
        if ok {
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
        err = row.Scan(&message.CreatedByUsername, &message.CreatedByUUID, &response.Topic.UUID)
        if err != nil {
                c.Logger().Error("Error retrieving from messageInfo: ", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topic-form", response)
        }

        response.Messages = append(response.Messages, message)

        c.Render(200, "oob-message", response)
        return c.Render(200, "topic-form", response)
}

func DeleteMessage(c echo.Context) error {
        response := structs.TopicResponse{}

        db := c.Get("db").(*sql.DB)
        user := c.Get("user").(structs.User)

        uuid := c.FormValue("uuid")
        createdByUUID := c.FormValue("createdBy")

        if createdByUUID != user.UUID {
                return c.NoContent(403)
        }

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

func QuoteMessage(c echo.Context) error {
        response := structs.TopicResponse{}

        db := c.Get("db").(*sql.DB)
        messageUUID := c.FormValue("uuid")

        var quotedContent string
        row := db.QueryRow(`
        SELECT 
                Content 
        FROM 
                message 
        WHERE 
                UUID = ?
        `, messageUUID)

        err := row.Scan(&quotedContent)
        if err != nil {
                c.Logger().Error("Error retrieving message content:", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topic-form", response)
        }

        quotedContent = "> " + quotedContent 
        return c.JSON(200, map[string]string{"quotedContent": quotedContent})
}

func DeleteTopic(c echo.Context) error {
        response := structs.HomeResponse{}

        db := c.Get("db").(*sql.DB)
        user := c.Get("user").(structs.User)

        topicUUID := c.FormValue("uuid")
        createdByUUID := c.FormValue("createdBy")

        if createdByUUID != user.UUID {
                return c.NoContent(403)
        }

        _, err := db.Exec(`
        DELETE FROM
                topic
        WHERE
                uuid = ?
        `, topicUUID)
        if err != nil {
                c.Logger().Error("Error deleting from topic", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "home-form", response)
        }

        c.Response().Header().Set("HX-Redirect", "/page/1")
        return c.NoContent(200)
}

func GetPostTopic(c echo.Context) error {
        response := structs.HomeResponse{}

        user, ok := c.Get("user").(structs.User)
        if ok {
                response.User = user
        }

        return c.Render(200, "postTopic", response)
}

func PostTopic(c echo.Context) error {
        response := structs.HomeResponse{}
        topic := structs.Topic{}

        topic.UUID = utils.Uuid()
        topic.Name = c.FormValue("name")
        topic.Description = c.FormValue("message")

        user, ok := c.Get("user").(structs.User)
        if ok {
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

        c.Response().Header().Set("HX-Redirect", "/page/1")
        return c.NoContent(200)
}
