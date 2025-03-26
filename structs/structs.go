package forum

import (
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
