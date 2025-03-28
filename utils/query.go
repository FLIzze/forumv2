package forum

import (
	"database/sql"

        "github.com/labstack/echo/v4"
	"github.com/russross/blackfriday/v2"

        structs "forum/structs"
)

func GetTopics(c echo.Context) ([]structs.Topic, structs.Error) {
        var topic structs.Topic
        var topics []structs.Topic

        conf, err := GetConfig()
        if err != nil {
                c.Logger().Info("Error retrieving config: ", err)
                conf.TopicsPerPage = 30
        }

        page, err := GetCurrentPage(c)
        if err != nil {
                c.Logger().Info("Error retrieving current page: ", err)
                page = 0
        }

        db := c.Get("db").(*sql.DB)

        OFFSET := conf.TopicsPerPage * page
        LIMIT := conf.TopicsPerPage

        rows, err := db.Query(`
        SELECT UUID, Name, Description, CreatedByUsername, CreatedByUUID, NmbMessages, LastMessage, CreationTime
        FROM topicInfo
        LIMIT ?
        OFFSET ?
        `, LIMIT, OFFSET)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                return nil, structs.NewError(err, 500, "Something went wrong. Please try again Later.")
        }
        defer rows.Close()

        for rows.Next() {
                err := rows.Scan(&topic.UUID, &topic.Name, &topic.Description, &topic.CreatedByUsername, 
                        &topic.CreatedByUUID, &topic.NmbMessages, &topic.LastMessage, &topic.CreationTime)
                if err != nil {
                        c.Logger().Error("Error scanning row", err)
                        return nil, structs.NewError(err, 422, "Something went wrong. Please try again Later.")
                }

                topic.FormattedCreationTime = FormatDate(topic.CreationTime)
                topic.FormattedLastMessage = FormatDate(topic.LastMessage)

                topics = append(topics, topic)
        }

        return topics, structs.NewError(nil, 200, "")
}

func GetTopic(c echo.Context) (structs.Topic, structs.Error) {
        var topic structs.Topic
        var plainContent string

        db := c.Get("db").(*sql.DB)
        uuid := c.Param("uuid")

        row := db.QueryRow(`
        SELECT UUID, Name, Description, CreatedByUsername, CreatedByUUID, LastMessage, CreationTime
        FROM topicInfo
        WHERE UUID = ?
        `, uuid)

        err := row.Scan(&topic.UUID, &topic.Name, &plainContent, &topic.CreatedByUsername, &topic.CreatedByUUID, 
                                                                        &topic.LastMessage, &topic.CreationTime)
        if err != nil {
                c.Logger().Error("Error scanning row", err)
                return topic, structs.NewError(err, 404, "")
        }

        htmlContent := string(blackfriday.Run([]byte(plainContent)))
        topic.Description = htmlContent
        topic.FormattedCreationTime = FormatDate(topic.CreationTime)
        topic.FormattedLastMessage = FormatDate(topic.LastMessage)

        return topic, structs.NewError(nil, 200, "")
}

func GetMessages(c echo.Context) ([]structs.Message, structs.Error) {
        var messages []structs.Message
        var message structs.Message
        var plainContent string

        uuid := c.Param("uuid")

        page, err := GetCurrentPage(c)
        if err != nil {
                c.Logger().Info("Error retrieving current page: ", err)
                page = 0
        }

        conf, err := GetConfig()
        if err != nil {
                c.Logger().Info("Error retrieving config: ", err)
                conf.MessagesPerPage = 30
        }

        OFFSET := conf.MessagesPerPage * page
        LIMIT := conf.MessagesPerPage

        db := c.Get("db").(*sql.DB)

        rows, err := db.Query(`
        SELECT UUID, Content, CreatedByUsername, CreatedByUUID, CreationTime
        FROM messageInfo
        WHERE TopicUUID = ?
        LIMIT ?
        OFFSET ?
        `, uuid, LIMIT, OFFSET)
        if err != nil {
                c.Logger().Error("Error retrieving topic message: ", err)
                return messages, structs.NewError(err, 500, "Something went wrong. Please try again later.")
        }
        defer rows.Close()

        for rows.Next() {
                err := rows.Scan(&message.UUID, &plainContent, &message.CreatedByUsername, &message.CreatedByUUID, 
                &message.CreationTime)
                if err != nil {
                        c.Logger().Error("Error retrieving topic message from column: ", err)
                        return messages, structs.NewError(err, 500, "Something went wrong. Please try again later.")
                }

                htmlContent := string(blackfriday.Run([]byte(plainContent)))
                message.Content = htmlContent
                message.FormattedCreationTime = FormatDate(message.CreationTime)

                messages = append(messages, message)
        }

        return messages, structs.NewError(nil, 200, "")
}
