package forum

import (
	"encoding/json"
	"log"
	"os"

	structs "forum/structs"
)

func GetConfig() (structs.Config, error) {
        config := structs.Config{}
        file, err := os.Open("./conf.json")
        defer file.Close()
        if err != nil {
		log.Println("Error reading conf.json")
                return config, err
        }

        decoder := json.NewDecoder(file)
        err = decoder.Decode(&config)
        if err != nil {
		log.Println("Error decoding conf.json")
                return config, err
        }

        return config, nil
}
