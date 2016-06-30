package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	DocumentRoot string
	TemplateRoot string
	StaticRoot   string
}

func assert(err error) {
	if err != nil {
		log.Fatalf("Fatal error encountered: %s\r\n", err.Error())
	}
}

func main() {
	cfgpath := flag.String("config", "conf", "Path to conf directory")
	bind := flag.String("bind", "127.0.0.1:6060", "Address:Port to bind to")
	flag.Parse()

	// Read and parse config file

	file, err := os.Open(filepath.Join(*cfgpath, "config.json"))
	assert(err)

	var cfg Config
	assert(json.NewDecoder(file).Decode(&cfg))

	file.Close()

	// Make http handlers and listen

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.StaticRoot))))
	http.Handle("/", mkHandler(cfg))

	log.Printf("Listening to %s\r\n", *bind)
	assert(http.ListenAndServe(*bind, nil))
}
