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
        CreationTime time.Time
        LastMessage *time.Time
}
