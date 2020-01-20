package messaging

const TypeInfoAttribute = "TypeInfo"

type Message struct {
	Body              interface{}
	MessageAttributes map[string]interface{}
}

type MessageAttributes interface {
	GetMessageAttributes() map[string]interface{}
}
