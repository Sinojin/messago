package messago

type Message struct {
	Sender Sender
	Message string
}

type Sender struct {
	SenderName string
}

func MessageInit(SenderName , message string) (*Message) {
	//Creating message
	return &Message{Sender:Sender{SenderName:SenderName},Message:message}

}

