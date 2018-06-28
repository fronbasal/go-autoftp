package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/dutchcoders/goftp"
	"log"
	"fmt"
	"github.com/fsnotify/fsnotify"
)

var (
	server   = kingpin.Flag("server", "The FTP host to connect to").Required().String()
	username = kingpin.Flag("username", "The FTP username").Required().String()
	password = kingpin.Flag("password", "The FTP password").Required().String()
	dir      = kingpin.Flag("directory", "The directory to watch").Required().ExistingDir()
)

func uploadDir(ftp *goftp.FTP) {
	err := ftp.Upload("./" + *dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Uploaded dir " + *dir)
}

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
	fmt.Printf("Connected to %s as %s \n", *server, *username)
	uploadDir(ftp)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("Got event: " + event.String())
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("Modified file, uploading directory.")
					uploadDir(ftp)
				}
			case err := <-watcher.Errors:
				log.Println("Error: " + err.Error())
			}
		}
	}()

	err = watcher.Add(*dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
