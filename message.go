package main

import (
	"encoding/json"
)

// MessageType ...
type MessageType int

const (
	// MessageNormal standard messages exchanged between users
	MessageNormal MessageType = iota
)

// Message ...
type Message struct {
	Kind int    `json:"kind"`
	User string `json:"user"`
	Data string `json:"data"`
}

// DecodeMessage ...
func DecodeMessage(b []byte) *Message {
	var m Message
	err := json.Unmarshal(b, &m)
	checkError(err, "DecodeMessage")
	return &m
}

// Bytes ...
func (m *Message) Bytes() []byte {
	b, err := json.Marshal(m)
	checkError(err, "Message.Bytes()")
	return b
}
