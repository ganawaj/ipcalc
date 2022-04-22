package main

import (
	"flag"
	"log"
	"os"

	"github.com/ganawaj/ipcalc/internal/api"
)

// var (
// 	version string
// 	os_ver  string
// 	os_arc  string
// 	go_ver  string
// 	git_sha string
// )

func main() {

	var cfg api.ServerConfig

	flag.IntVar(&cfg.Port, "port", 4000, "HTTP listen port")
	flag.StringVar(&cfg.Address, "address", "", "HTTP listening IP")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LUTC|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ltime|log.Llongfile)

	server := &api.Server{
		Config:        cfg,
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
	}

	err := server.ListenAndServe()
	errorLog.Fatal(err)
}