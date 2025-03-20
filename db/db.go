package forum

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnectDb() (*sql.DB, error) {
        USER := os.Getenv("DB_USER")
        PASSWORD := os.Getenv("DB_PASSWORD")

        if USER == "" || PASSWORD == "" {
                log.Fatal("DB_USER or DB_PASSWORD is not set in the environment variables")
        }

        credentials := fmt.Sprintf("%s:%s@/forum?parseTime=true", USER, PASSWORD)

        db, err := sql.Open("mysql", credentials)
        if err != nil {
                return nil, err
        }

        return db, nil
}

func CreateTable(db *sql.DB) error {
        _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS user (
                UUID varchar(37) NOT NULL PRIMARY KEY,
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
        CREATE TABLE IF NOT EXISTS topic (
                UUID varchar(37) NOT NULL PRIMARY KEY, 
                Name varchar(95) NOT NULL, 
                Description LONGTEXT NOT NULL,
                CreatedBy varchar(37) NULL,
                Pinned int DEFAULT 0,
                FOREIGN KEY (CreatedBy) REFERENCES user(UUID) ON DELETE SET NULL,
                CreationTime DATETIME NOT NULL
        )
        `)
        if err != nil {
                return err
        }

        _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS message (
                UUID varchar(37) NOT NULL PRIMARY KEY,
                Content LONGTEXT NOT NULL,
                CreationTime DATETIME NOT NULL,
                TopicUUID varchar(37) NULL,
                CreatedBy varchar(37) NULL,
                Pinned int DEFAULT 0,
                FOREIGN KEY (TopicUUID) REFERENCES topic(UUID) ON DELETE CASCADE,
                FOREIGN KEY (CreatedBy) REFERENCES user(UUID) ON DELETE SET NULL
        )
        `)
        if err != nil {
                return err
        }

        _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS session (
                SessionUUID varchar(37) PRIMARY KEY,
                UserUUID varchar(37) NULL,
                FOREIGN KEY (UserUUID) REFERENCES user(UUID) ON DELETE SET NULL,
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
                t.CreatedBy AS CreatedByUUID,
                t.Pinned,
                u.Username AS CreatedByUsername,
                COUNT(m.UUID) AS NmbMessages,
                COALESCE(
                        (SELECT 
                                m2.CreationTime
                                FROM 
                                        message m2 
                                WHERE 
                                        m2.TopicUUID = t.UUID 
                                ORDER BY 
                                        m2.CreationTime DESC 
                        LIMIT 1),
                        t.CreationTime
                ) AS LastMessage
        FROM 
                topic t 
        LEFT JOIN 
                user u ON t.CreatedBy = u.UUID 
        LEFT JOIN 
                message m ON t.UUID = m.TopicUUID 
        GROUP BY 
                t.UUID, t.Name, t.Description, u.Username
        ORDER BY 
                LastMessage DESC
        `)
        if err != nil {
                return err 
        }

        _, err = db.Exec(`
        CREATE VIEW IF NOT EXISTS messageInfo AS 
        SELECT 
                m.UUID, 
                m.TopicUUID,
                m.Content, 
                m.CreationTime,
                m.Pinned,
                u.Username as CreatedByUsername,
                u.UUID as CreatedByUUID
        FROM 
                message m 
        JOIN 
                user u ON m.CreatedBy = u.UUID
        ORDER BY
               CreationTime ASC 
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
                        SELECT MAX(m2.Content) 
                        FROM message m2 
                        WHERE m2.CreatedBy = u.UUID
                ) AS LastMessage,
                (
                        SELECT MAX(t2.Name) 
                        FROM topic t2 
                        WHERE t2.CreatedBy = u.UUID
                ) AS LastTopic
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

func HandleDbSetup() *sql.DB {
        db, err := ConnectDb()
        if err != nil {
                log.Fatal("Error connecting to db", err)        
        }

        err = CreateTable(db)
        if err != nil {
                log.Fatal("Error creating tables:", err)
        }

        err = CreateView(db)
        if err != nil {
                log.Fatal("Error creating views:", err)
        }

        return db
}
