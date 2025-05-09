package forum

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/russross/blackfriday/v2"

	structs "forum/structs"
	utils "forum/utils"
)

func GetTopicPage(c echo.Context) error {
        var response structs.TopicResponse

        topic, error := utils.GetTopic(c)
        error.HandleError(c)

        response.Topic = topic
        user, ok := c.Get("user").(structs.User)
        if ok {
                response.User = user
        }

        messages, error := utils.GetMessages(c)
        error.HandleError(c)
        response.Messages = messages

	conf, err := utils.GetConfig()
	if err != nil {
		c.Logger().Error("Error retrieving config", err)
		conf.MessagesPerPage = 20
	}

	currentPage, error := utils.GetCurrentPage(c)
        error.HandleError(c)
	response.Page.CurrentPage = currentPage + 1
	response.Page.TotalPage = (len(messages) / conf.MessagesPerPage) + 1

        return c.Render(200, "topic", response)
}

func PostMessage(c echo.Context) error {
        var response structs.TopicResponse
        var message structs.Message

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
        db := c.Get("db").(*sql.DB)

        if message.Content == "" {
                c.Logger().Error("Message empty")
                response.Status.Error = "Message must be filled"
                return c.Render(422, "topic-form", response)
        }

        _, err := db.Exec(`
        INSERT INTO message (UUID, TopicUUID, Content, CreatedBy, CreationTime) 
        VALUES (?, ?, ?, ?, ?)
        `, message.UUID, message.TopicUUID, message.Content, user.UUID, time.Now())
        if err != nil {
                c.Logger().Error("Error inserting into message: ", err)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topic-form", response)
        }

        message, error := utils.GetMessage(c, message.UUID)
        error.HandleError(c)

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
        DELETE FROM message
        WHERE uuid = ? 
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
        SELECT Content 
        FROM message 
        WHERE UUID = ?
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
        DELETE FROM topic
        WHERE uuid = ?
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

	currentPage, err := utils.GetCurrentPage(c)
        err.HandleError(c)
	response.Page.CurrentPage = currentPage + 1

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

	_, error := db.Exec(`
        INSERT INTO topic (UUID, Name, Description, CreatedBy, CreationTime) 
        VALUES (?, ?, ?, ?, ?)
        `, topic.UUID, topic.Name, topic.Description, user.UUID, time.Now())
        if error != nil {
                c.Logger().Error("Error inserting response Topic: ", error)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "home-form", response)
        }

        row := db.QueryRow(`
        SELECT CreatedByUsername, CreatedByUUID, NmbMessages 
        FROM topicInfo
        WHERE UUID = ?
        `, topic.UUID)
        error = row.Scan(&topic.CreatedByUsername, &topic.CreatedByUUID, &topic.NmbMessages)
        if error != nil {
                c.Logger().Error("Error fetching new topic: ", error)
                response.Status.Error = "Something went wrong. Please try again later."
                return c.Render(500, "home-form", response)
        }

        response.Topics = append(response.Topics, topic)

        c.Response().Header().Set("HX-Redirect", "/page/1")
        return c.NoContent(200)
}
