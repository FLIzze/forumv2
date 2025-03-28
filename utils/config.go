package forum

import (
        "encoding/json"
        "os"

        structs "forum/structs"
)

func GetConfig() (structs.Config, error) {
        config := structs.Config{}
        file, err := os.Open("./conf.json")
        defer file.Close()
        if err != nil {
                return config, err
        }

        decoder := json.NewDecoder(file)
        err = decoder.Decode(&config)
        if err != nil {
                return config, err
        }

        return config, nil
}
