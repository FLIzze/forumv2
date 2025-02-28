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
                Description text NOT NULL
        )
        `)
        if err != nil {
                return err
        }

        _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS message (
                UUID varchar(37) NOT NULL,
                TopicUUID varchar(37) NOT NULL,
                Content text NOT NULL
        )
        `)
        return err
}
