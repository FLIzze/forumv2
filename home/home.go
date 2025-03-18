package forum

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	user "forum/user"
)

type HomeResponse struct {
        Error string
        Success string
        Topics []Topic
        User user.User
}

type Topic struct {
        UUID string
        Name string
        Description string
        CreatedByUsername string
        CreatedByUUID string
        NmbMessages int
}

func GetHomePage(c echo.Context) error {
        response := HomeResponse{}
        topic := Topic{}

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(user.User)
        if ok {
                response.User = user
        }

        rows, err := db.Query(`
        SELECT 
                UUID, Name, Description, CreatedByUsername, CreatedByUUID, NmbMessages 
        FROM 
                topicInfo
        `)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                response.Error = "Could not retrieve topic."
                return c.Render(500, "home", response)
        }
        defer rows.Close()

        for rows.Next() {
                err := rows.Scan(&topic.UUID, &topic.Name, &topic.Description, &topic.CreatedByUsername, 
                                                                &topic.CreatedByUUID, &topic.NmbMessages)
                if err != nil {
                        c.Logger().Error("Error scanning row", err)
                        response.Error = "Something went wrong. Please try again later."
                        return c.Render(422, "home", response)
                }

                response.Topics = append(response.Topics, topic)
        }

        return c.Render(200, "home", response)
}

func PostTopic(c echo.Context) error {
        response := HomeResponse{}
        topic := Topic{}

        topic.UUID = uuid.New().String()
        topic.Name = c.FormValue("name")
        topic.Description = c.FormValue("description")

        if topic.Name == "" || topic.Description == "" {
                c.Logger().Error("Name and/or description empty")
                response.Error = "Name and/or description must be filled."
                return c.Render(422, "topics-form", response)
        }

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(user.User)
        if !ok {
                response.Error = "You must be logged in to post a topic."
                return c.Render(401, "topics-form", response)
        } else {
                response.User = user
        }

        _, err := db.Exec(`
        INSERT INTO topic (UUID, Name, Description, CreatedBy, CreationTime) 
        VALUES (?, ?, ?, ?, ?)
        `, topic.UUID, topic.Name, topic.Description, user.UUID, time.Now())
        if err != nil {
                c.Logger().Error("Error inserting response Topic: ", err)
                response.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topics-form", response)
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
                response.Error = "Something went wrong. Please try again later."
                return c.Render(500, "topics-form", response)
        }

        response.Topics = append(response.Topics, topic)
        response.Success = "Topic succesfully created."

        c.Render(200, "topics-form", response)
        return c.Render(200, "oob-topic", response)
}

func DeleteTopic(c echo.Context) error {
        response := HomeResponse{}

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(user.User)
        if !ok {
                response.Error = "You must be logged in to delete a topic."
                return c.Render(401, "home", response)
        } else {
                response.User = user
        }

        uuid := c.FormValue("uuid")

        _, err := db.Exec(`
        DELETE FROM
                topic
        WHERE 
                uuid = ?
        `, uuid)
        if err != nil {
                c.Logger().Error("Error deleting from topic", err)
                response.Error = "Something went wrong. Please try again later."
                return c.Render(500, "home", response)
        }

        response.Success = "Topic succesfully deleted."

        return c.Render(200, "home", response)
}
