package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	bind     string
	peer     string
	username string
)

func init() {
	var defaultUser string
	currentUser, err := user.Current()
	if err != nil {
		log.Errorf("unable to determine current user: %s", err)
		defaultUser = ""
	} else {
		defaultUser = currentUser.Username
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<host>:<port>]", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.StringVar(&bind, "b", ":1337", "bind to port and interface")
	flag.StringVar(&username, "u", defaultUser, "set username")

	flag.Parse()
}

func main() {
	var peer string
	if len(flag.Args()) == 1 {
		peer = flag.Arg(0)
	}

	p := NewPeer(username, bind, peer)

	p.OnMessage(func(msg *Message) error {
		switch msg.Kind {
		case MessageNormal:
			fmt.Printf("<%s> %s\n", msg.User, msg.Data)
		case MessageHello:
			p.SetPeer(msg.Addr)
		case MessageKey:
			data, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				log.Errorf("error decoding message data: %s", err)
			} else {
				p.SetKey(data)
			}
		default:
			log.Errorf("invalu message type %d", msg.Kind)
		}
		return nil
	})

	p.Start()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf("error reading input: %s", err)
		} else {
			msg = strings.TrimSpace(msg)
			if msg != "" {
				p.SendMessage(msg)
			}
		}
	}
}
