package main

import (
	"encoding/json"
	"log"
)

const (
	connect    = "conn"
	disconnect = "disconn"
	message    = "msg"
)

// Message interface anthing that can goto json
type Message struct {
	Type string
	User string
	Msg  string
}

func (m Message) encode() []byte {
	if data, err := json.Marshal(m); err == nil {
		return data
	}
	log.Fatalln("Failed to parse message", m)
	return []byte("{}")
}

func newMessage(data []byte) Message {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Fatalln("Failed to parse message", err)
		return Message{}
	}
	return msg
}
