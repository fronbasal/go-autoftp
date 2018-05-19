package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/dutchcoders/goftp"
	"log"
	"fmt"
)

var (
	server   = kingpin.Flag("server", "The FTP host to connect to").Required().String()
	username = kingpin.Flag("username", "The FTP username").Required().String()
	password = kingpin.Flag("password", "The FTP password").Required().String()
	dir      = kingpin.Flag("directory", "The directory to watch").Required().ExistingDir()
)

func main() {
	kingpin.Parse()
	var ftp *goftp.FTP
	var err error
	if ftp, err = goftp.Connect(*server); err != nil {
		log.Fatal(err)
	}
	if err := ftp.Login(*username, *password); err != nil {
		log.Fatal(err)
	}
	defer ftp.Close()
	fmt.Printf("Connected to %s as %s", *server, *username)
}
