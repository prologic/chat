package main

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// MessageHandler ...
type MessageHandler func(msg *Message) (err error)

// Peer ...
type Peer struct {
	bind     string
	peer     string
	username string

	conn net.PacketConn

	running bool
	outbox  chan *Message
	inbox   chan *Message
	handler MessageHandler
}

// NewPeer ...
func NewPeer(username, bind, peer string) *Peer {
	return &Peer{
		bind:     bind,
		peer:     peer,
		username: username,

		running: true,
		outbox:  make(chan *Message, 20),
		inbox:   make(chan *Message),
	}
}

// Start ...
func (p *Peer) Start() {
	conn, err := net.ListenPacket("udp4", p.bind)
	if err != nil {
		log.Fatalf("error binding to interface %s: %s", p.bind, err)
	}

	p.conn = conn

	go p.loop()
	go p.readpump()
	go p.writepump()
}

// SendMessage ...
func (p *Peer) SendMessage(msg string) error {
	p.outbox <- &Message{User: p.username, Data: msg}
	return nil
}

func (p *Peer) writepump() {
	for p.running {
		msg := <-p.outbox
		raddr, err := net.ResolveUDPAddr("udp4", p.peer)
		if err != nil {
			log.Errorf("error resolving peer address %s: %s", peer, err)
		} else {
			b, err := msg.Bytes()
			if err == nil {
				_, err = p.conn.WriteTo(b, raddr)
				if err != nil {
					log.Errorf("error sending messages to peer %s: %s", p.peer, err)
				}
			}
		}
	}
}

func (p *Peer) readpump() {
	buf := make([]byte, 4096)

	for p.running {
		n, addr, err := p.conn.ReadFrom(buf)
		if err != nil {
			log.Errorf("error reading from peer %s: %s", addr, err)
		} else {
			p.peer = addr.String()
			msg, err := DecodeMessage(buf[:n])
			if err == nil {
				p.inbox <- msg
			}
		}
	}
}

// OnMessage ...
func (p *Peer) OnMessage(handler MessageHandler) {
	p.handler = handler
}

func (p *Peer) loop() {
	for p.running {
		msg := <-p.inbox
		if p.handler != nil {
			err := p.handler(msg)
			if err != nil {
				log.Errorf("error handling message %v: %s", msg, err)
			}

		}
	}
}
