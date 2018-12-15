package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
)

var (
	bind     string
	peer     string
	username string
)

func init() {
	currentUser, err := user.Current()
	checkError(err, "init")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<host>:<port>]", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.StringVar(&bind, "b", ":1337", "bind to port and interface")
	flag.StringVar(&username, "u", currentUser.Username, "set username")

	flag.Parse()
}

func checkError(err error, funcName string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error:%s-----in func:%s", err.Error(), funcName)
		os.Exit(1)
	}
}

func main() {
	var peer string
	if len(flag.Args()) == 1 {
		peer = flag.Arg(0)
	}

	NewPeer(username, bind, peer).Run()
}
