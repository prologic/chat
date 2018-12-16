package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// MessageType ...
type MessageType int

const (
	// MessageNormal standard messages exchanged between users
	MessageNormal MessageType = iota

	// MessageHello MessageType
	MessageHello

	// MessageKey MessageType
	MessageKey
)

// Message ...
type Message struct {
	Kind MessageType `json:"kind"`
	Addr string      `json:"addr"`
	User string      `json:"user"`
	Data string      `json:"data"`
}

// DecodeMessage ...
func DecodeMessage(b []byte) (msg *Message, err error) {
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Errorf("error decoding messages: %s", err)
	}
	return
}

// Bytes ...
func (m *Message) Bytes() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Errorf("error encoding message: %s", err)
		return []byte{}, err
	}
	return b, nil
}
