package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// Peer ...
type Peer struct {
	bind     string
	peer     string
	username string

	conn net.PacketConn

	running bool
	outbox  chan *Message
	inbox   chan *Message
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

// Run ...
func (p *Peer) Run() {
	//raddr, err := net.ResolveUDPAddr("udp4", p.peer)
	//checkError(err, "client")

	conn, err := net.ListenPacket("udp4", p.bind)
	checkError(err, "main")
	p.conn = conn
	defer p.conn.Close()

	if p.peer != "" {
		p.SendMessage(p.peer, "hello")
	}

	go p.printMessage()
	go p.readpump()
	go p.writepump()

	p.getMessage()
}

// SendMessage ...
func (p *Peer) SendMessage(peer, msg string) {
	raddr, err := net.ResolveUDPAddr("udp4", peer)
	checkError(err, "client.SendMessage")

	m := Message{User: p.username, Data: msg}
	_, err = p.conn.WriteTo(m.Bytes(), raddr)
	checkError(err, "client.SendMessage")
}

func (p *Peer) writepump() {
	for p.running {
		msg := <-p.outbox
		raddr, err := net.ResolveUDPAddr("udp4", p.peer)
		checkError(err, "client")
		_, err = p.conn.WriteTo(msg.Bytes(), raddr)
		checkError(err, "sendMessage")
	}
}

func (p *Peer) readpump() {
	buf := make([]byte, 4096)

	for p.running {
		n, addr, err := p.conn.ReadFrom(buf)
		p.peer = addr.String()
		checkError(err, "readpump")
		msg := DecodeMessage(buf[:n])
		p.inbox <- msg
	}
}

func (p *Peer) getMessage() {
	reader := bufio.NewReader(os.Stdin)
	for p.running {
		fmt.Print(">>> ")
		msg, err := reader.ReadString('\n')
		checkError(err, "getMessage")
		p.outbox <- &Message{User: p.username, Data: msg}
	}
}

func (p *Peer) printMessage() {
	for p.running {
		msg := <-p.inbox
		fmt.Printf("<%s> %s\n", msg.User, msg.Data)
	}
}
