package forum

import (
        "database/sql"
)

func ConnectDb() (*sql.DB, error) {
        db, err := sql.Open("mysql", "admin:1231@/forum?parseTime=true")
        return db, err
}

func CreateTable(db *sql.DB) error {
        _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS topic (
                UUID varchar(37) NOT NULL UNIQUE, 
                Name varchar(25) NOT NULL, 
                Description text NOT NULL,
                CreatedBy varchar(37) NOT NULL,
                CreationTime DATETIME NOT NULL
        )
        `)
        if err != nil {
                return err
        }

        _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS message (
                UUID varchar(37) NOT NULL UNIQUE,
                TopicUUID varchar(37) NOT NULL,
                Content text NOT NULL,
                CreatedBy varchar(37) NOT NULL,
                CreationTime DATETIME NOT NULL
        )
        `)
        if err != nil {
                return err
        }

        _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS user (
                UUID varchar(37) NOT NULL UNIQUE,
                Username varchar(17) NOT NULL UNIQUE,
                Email varchar(254) NOT NULL UNIQUE,
                Password varchar(60) NOT NULL,
                CreationTime DATETIME NOT NULL
        )
        `)
        if err != nil {
                return err 
        }

        _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS session (
                SessionUUID varchar(37) UNIQUE,
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
        FROM 
                user u
        JOIN 
                session s ON u.UUID = s.UserUUID;
        `)
        if err != nil {
                return err 
        }

        _, err = db.Exec(`
        CREATE VIEW IF NOT EXISTS topicInfo AS 
        SELECT      
                t.UUID,
                t.Name,
                t.Description,
                u.Username AS CreatedBy,
                COUNT(m.UUID) AS NmbMessages,
                (SELECT 
                        m2.CreationTime
                        FROM 
                                message m2 
                        WHERE 
                                m2.TopicUUID = t.UUID 
                        ORDER BY 
                                m2.CreationTime DESC 
                        LIMIT 
                                1
                ) AS LastMessage
        FROM 
                topic t 
        JOIN 
                user u ON t.CreatedBy = u.UUID 
        LEFT JOIN 
                message m ON t.UUID = m.TopicUUID 
        GROUP BY 
                t.UUID, t.Name, t.Description, u.Username;
        `)
        if err != nil {
                return err 
        }

        _, err = db.Exec(`
        CREATE VIEW IF NOT EXISTS messageInfo AS 
        SELECT 
                m.TopicUUID,
                m.Content, 
                u.Username as CreatedBy 
        FROM 
                message m 
        JOIN 
                user u ON m.CreatedBy = u.UUID;
        `)
        if err != nil {
                return err 
        }

        _, err = db.Exec(`
        CREATE VIEW IF NOT EXISTS userInfo AS 
        SELECT
                u.UUID AS UserUUID,
                u.Username,
                u.CreationTime,
                COALESCE(COUNT(DISTINCT m.UUID), 0) AS NmbMessagesPosted,
                COALESCE(COUNT(DISTINCT t.UUID), 0) AS NmbTopicsCreated,
                (
                        SELECT MAX(m2.CreationTime) 
                        FROM message m2 
                        WHERE m2.CreatedBy = u.UUID
                ) AS LastMessage
        FROM 
                user u
        LEFT JOIN 
                message m ON u.UUID = m.CreatedBy
        LEFT JOIN 
                topic t ON u.UUID = t.CreatedBy
        GROUP BY 
                u.UUID, u.Username, u.CreationTime;
        `)

        return err
}
