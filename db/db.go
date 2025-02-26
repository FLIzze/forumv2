package forum

import (
	"database/sql"
	"fmt"
)

func ConnectDb() *sql.DB {
        db, err := sql.Open("mysql", "admin:1231@/forum")
        if err != nil {
                fmt.Printf("Error connecting to database: %s", err)
        }

        return db
}

func CreateTable(db *sql.DB) {
        _, err := db.Exec(`
                CREATE TABLE IF NOT EXISTS topic (
                        UUID varchar(37) NOT NULL, 
                        Name varchar(25) NOT NULL, 
                        Description text NOT NULL
                )
        `)
        if err != nil {
                panic(err)
        }
}
