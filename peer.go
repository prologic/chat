package main

import (
	"encoding/base64"
	"net"

	"github.com/monnand/dhkx"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/twofish"
)

// MessageHandler ...
type MessageHandler func(msg *Message) (err error)

// Peer ...
type Peer struct {
	bind     string
	peer     string
	username string

	cipher     Cipher
	group      *dhkx.DHGroup
	privateKey *dhkx.DHKey
	sessionKey []byte

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

	// Get a group. Use the default one would be enough.
	g, err := dhkx.GetGroup(0)
	if err != nil {
		log.Errorf("error getting group: %s", err)
	} else {
		p.group = g

		// Generate a private key from the group.
		// Use the default random number generator.
		priv, err := g.GeneratePrivateKey(nil)
		if err != nil {
			log.Errorf("error generating private key: %s", err)
		} else {
			p.privateKey = priv
		}
	}

	if p.peer != "" {
		p.outbox <- &Message{Kind: MessageHello}

		// Get the public key from the private key.
		pub := p.privateKey.Bytes()

		// Send the public key to Bob.
		msg := base64.StdEncoding.EncodeToString(pub)
		p.outbox <- &Message{Kind: MessageKey, Data: msg}
	}

	go p.loop()
	go p.readpump()
	go p.writepump()
}

// SendMessage ...
func (p *Peer) SendMessage(msg string) error {
	ciphertext, err := encrypt(p.cipher, []byte(msg))
	if err != nil {
		log.Errorf("error encrypting message: %s", err)
		return err
	}

	data := base64.StdEncoding.EncodeToString(ciphertext)
	p.outbox <- &Message{User: p.username, Data: data}
	return nil
}

func (p *Peer) writepump() {
	for p.running {
		msg := <-p.outbox
		msg.User = p.username
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
			msg, err := DecodeMessage(buf[:n])
			if err == nil {
				msg.Addr = addr.String()
				p.inbox <- msg
			}
		}
	}
}

// SetCipher ...
func (p *Peer) SetCipher(cipher Cipher) {
	p.cipher = cipher
}

// SetKey ...
func (p *Peer) SetKey(s []byte) {
	if len(p.sessionKey) == 0 {
		// Get the public key from the private key.
		pub := p.privateKey.Bytes()

		// Send the public key to Bob.
		msg := base64.StdEncoding.EncodeToString(pub)
		p.outbox <- &Message{Kind: MessageKey, Data: msg}

		pubKey := dhkx.NewPublicKey(s)

		// Compute the key
		k, err := p.group.ComputeKey(pubKey, p.privateKey)
		if err != nil {
			log.Errorf("error computing key: %s", err)
		} else {
			// Get the key in the form of []byte
			// We can only use 32bytes of the key (twofish limits)
			key := k.Bytes()[:32]

			p.sessionKey = key

			cipher, err := twofish.NewCipher(key)
			if err != nil {
				log.Fatalf("error initializing cipher: %s", err)
			}
			p.cipher = cipher
		}
	}
}

// SetPeer ...
func (p *Peer) SetPeer(peer string) {
	p.peer = peer
}

// OnMessage ...
func (p *Peer) OnMessage(handler MessageHandler) {
	p.handler = handler
}

func (p *Peer) loop() {
	for p.running {
		msg := <-p.inbox
		if p.handler != nil {
			// Only decrypt normal messages
			if msg.Kind == MessageNormal {
				data, err := base64.StdEncoding.DecodeString(msg.Data)
				if err != nil {
					log.Errorf("error decoding message data: %s", err)
				} else {
					plaintext, err := decrypt(p.cipher, data)
					if err != nil {
						log.Errorf("error decrypting message: %s", err)
					} else {
						msg.Data = string(plaintext)
					}
				}
			}

			err := p.handler(msg)
			if err != nil {
				log.Errorf("error handling message %v: %s", msg, err)
			}
		}
	}
}
