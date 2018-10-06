package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/dutchcoders/goftp"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

var (
	server    = kingpin.Flag("server", "The FTP host to connect to").Short('s').Required().String()
	username  = kingpin.Flag("username", "The FTP username").Short('u').Required().String()
	password  = kingpin.Flag("password", "The FTP password").Short('p').Required().String()
	dir       = kingpin.Flag("directory", "The directory to watch").Short('d').Required().ExistingDir()
	verbose   = kingpin.Flag("verbose", "Enable verbose output").Short('v').Bool()
	overwrite = kingpin.Flag("overwrite", "Overwrite existing files").Short('f').Bool()
)

var absPath string // the absolute path to *dir

func uploadDir(ftp *goftp.FTP) {
	if *overwrite {
		log.Debug("Removing existing files from server!")
		ftp.Dele("/")
	}
	err := ftp.Upload(absPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Synchronized directory (%s) with server %s!", *dir, *server)
}

func main() {
	kingpin.Parse()
	kingpin.CommandLine.Name = "go-autoftp"
	kingpin.CommandLine.Author("Daniel Malik <mail@fronbasal.de>")

	log.SetLevel(log.InfoLevel)
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	if !filepath.IsAbs(*dir) {
		var err error
		absPath, err = filepath.Abs(*dir)
		if err != nil {
			log.Fatal("An error occurred while attempting to resolve the directory: ", err)
		}
	} else {
		absPath = *dir
	}
	log.Debugf("Using path %s! ", absPath)

	var ftp *goftp.FTP
	var err error
	if ftp, err = goftp.Connect(*server); err != nil {
		log.Fatal(err)
	}
	if err := ftp.Login(*username, *password); err != nil {
		log.Fatal(err)
	}
	defer ftp.Close()
	log.Info("Connected to server ", *server)
	uploadDir(ftp)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	log.Debug("Watching directory!")
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Debugf("A change occurred - syncing directory %s to %s (event: %s)", *dir, *server, event.String())
				uploadDir(ftp)
			case err := <-watcher.Errors:
				log.Warn("An error occurred while attempting to watch the given directory: ", err)
			}
		}
	}()

	err = watcher.Add(*dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
