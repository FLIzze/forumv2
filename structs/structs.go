package forum

import (
	"html/template"
	"time"
)

type User struct {
        SessionUUID string
        UUID string
        Username string
        Email string
        NmbMessagesPosted int
        NmbTopicsCreated int
        LastMessage *string
        CreationTime time.Time
}

type HomeResponse struct {
        Status Status
        Topics []Topic
        User User
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
}

type TopicResponse struct {
        Status Status
        Subject Subject
        Messages []Message
        User User
}

type Subject struct {
        UUID string
        Name string
        Description string
        CreatedByUsername string
}

type Message struct {
        UUID string
        TopicUUID string
        Content template.HTML
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
