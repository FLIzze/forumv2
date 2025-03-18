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
                Name varchar(25) NOT NULL, 
                Description text NOT NULL,
                CreatedBy varchar(37) NULL,
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
                Content text NOT NULL,
                CreationTime DATETIME NOT NULL,
                TopicUUID varchar(37) NULL,
                CreatedBy varchar(37) NULL,
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
                u.Username AS CreatedByUsername,
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
        LEFT JOIN 
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
                m.UUID, 
                m.TopicUUID,
                m.Content, 
                u.Username as CreatedByUsername,
                u.UUID as CreatedByUUID
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

func CreateIndex(db *sql.DB) error {
        _, err := db.Exec(`
        CREATE INDEX 
                idx_user_username 
        ON 
                user (Username);
        `)
        return err
}
