package forum

type Topic struct {
        UUID string
        Name string
        Description string
}


type Topics struct {
        Topics []Topic
}

func NewTopic(uuid, name, description string) Topic {
        return Topic{
                UUID: uuid,
                Name: name,
                Description: description,
        }
}
