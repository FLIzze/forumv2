package forum

import (
        "database/sql"
)

func ConnectDb() (*sql.DB, error) {
        db, err := sql.Open("mysql", "admin:1231@/forum")
        return db, err
}

func CreateTable(db *sql.DB) error {
        _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS topic (
                UUID varchar(37) NOT NULL, 
                Name varchar(25) NOT NULL, 
                Description text NOT NULL,
                CreatedBy varchar(37) NOT NULL
        )
        `)
        if err != nil {
                return err
        }

        _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS message (
                UUID varchar(37) NOT NULL,
                TopicUUID varchar(37) NOT NULL,
                Content text NOT NULL,
                CreatedBy varchar(37) NOT NULL
        )
        `)
        if err != nil {
                return err
        }

        _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS user (
                UUID varchar(37) NOT NULL,
                Username varchar(17) NOT NULL,
                Email varchar(254) NOT NULL,
                Password varchar(60) NOT NULL
        )
        `)
        if err != nil {
                return err 
        }

        _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS session (
                SessionUUID varchar(37),
                UserUUID varchar(37) NOT NULL,
                Connected int DEFAULT 0
        )
        `)

        return err
}

func CreateView(db *sql.DB) error {
        _, err := db.Exec(`
        CREATE VIEW IF NOT EXISTS userSession AS SELECT 
                u.UUID AS UserUUID,
                u.Username,
                u.Email,
                s.SessionUUID,
                s.Connected
        FROM user u
        JOIN session s ON u.UUID = s.UserUUID;
        `)
        if err != nil {
                return err 
        }

        _, err = db.Exec(`
        CREATE VIEW IF NOT EXISTS topicInfo AS SELECT      
                t.UUID,
                t.Name,
                t.Description,
                u.Username as CreatedBy,
                COUNT(m.UUID) AS NmbMessages
        FROM topic t 
        JOIN user u ON t.CreatedBy = u.UUID 
        LEFT JOIN message m ON t.UUID = m.TopicUUID 
        GROUP BY t.UUID, t.Name, t.Description, u.Username;
        `)
        if err != nil {
                return err 
        }

        _, err = db.Exec(`
        CREATE VIEW IF NOT EXISTS messageInfo AS SELECT 
                m.TopicUUID,
                m.Content, 
                u.Username as CreatedBy 
        FROM message m 
        JOIN user u ON m.CreatedBy = u.UUID;
        `)

        return err
}
