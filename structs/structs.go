package forum

import (
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
        SessionUUID string
        UUID string
        Username string
        Email string
        NmbMessagesPosted int
        NmbTopicsCreated int
        LastMessage *string
        LastTopic *string
        CreationTime *time.Time
        FormattedCreationTime string
}

type HomeResponse struct {
        Status Status
        Topics []Topic
        User User
        Page Page
}

type Page struct {
        TotalPage int
        CurrentPage int
}

type Topic struct {
        UUID string
        Name string
        Description string
        CreatedByUsername string
        CreatedByUUID string
        NmbMessages int
        LastMessage *time.Time
        FormattedLastMessage string
        CreationTime *time.Time
        FormattedCreationTime string
}

type TopicResponse struct {
        Status Status
        Topic Topic
        Messages []Message
        User User
        QuotedMessage string
        Page Page
}

type Message struct {
        UUID string
        TopicUUID string
        Content string
        CreatedByUsername string
        CreatedByUUID string
        CreationTime *time.Time
        FormattedCreationTime string
}

type ProfilResponse struct {
        Status Status
        User User
        UserProfil User
}

type Status struct {
        Error string
        Success string
}

type Config struct {
        Host string
        Port string
        TopicsPerPage int
        MessagesPerPage int
}

type Error struct {
        Err error
        Status int
        Message string
}

func NewError(err error, status int, message string) Error {
        return Error {
                Err: err,
                Status: status,
                Message: message,
        }
}

func (e Error) IsError() bool {
        return e.Err != nil
}

func (e Error) HandleError(c echo.Context) error {
        switch e.Status {
        case 404:
                return c.Render(e.Status, "404", nil)
        case 401:
                //later on notification center
                e.Message = "You need to login to perform this action."
                return c.Render(e.Status, "404", nil)
        case 422:
                //later on notification center
                e.Message = "Invalid Input."
                return c.Render(e.Status, "404", nil)
        case 500:
                //later on notification center
                e.Message = "Something went wrong. Please wait and try again."
                return c.Render(e.Status, "404", nil)
        default:
                return nil
        }
}
