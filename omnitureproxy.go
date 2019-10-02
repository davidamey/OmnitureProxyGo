package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davidamey/omnitureproxy/api"
	"github.com/davidamey/omnitureproxy/archive"
	"github.com/davidamey/omnitureproxy/proxy"
	"github.com/davidamey/omnitureproxy/socket"
	"github.com/urfave/negroni"
)

type OmnitureProxy struct {
	TargetURL  string
	ArchiveDir string

	hub socket.Hub
}

func (op *OmnitureProxy) Notify(entry *archive.Entry) {
	// fmt.Printf("notify: %s\n", entry.PageName)
	op.hub.SendToAll(entry)
}

func (op *OmnitureProxy) Start() {
	log.Println("Starting")

	mux := http.NewServeMux()
	mux.Handle("/api/", api.NewApi(op.ArchiveDir))

	// sockets
	op.hub = socket.NewHub()
	go op.hub.Run()
	mux.HandleFunc("/ws", op.hub.ServeWS)

	// archive
	archiver := archive.NewWriter(op.ArchiveDir)
	archiver.StartProcessing()
	defer archiver.StopProcessing()

	// proxy
	p := proxy.New(archiver, op, op.TargetURL)
	mux.HandleFunc("/b/ss/", p.Handle)

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(mux)

	// start the server
	port := getPort()
	s := &http.Server{
		Addr:           port,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("listening on", port)
	log.Fatal(s.ListenAndServe())
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":3000"
}

func main() {
	const (
		defaultTarget     = ""
		targetUsage       = "redirect url e.g. 'https://somewhere:5000' ('' to disable)"
		defaultArchiveDir = "_archive"
		archiveDirUsage   = "folder to log to e.g. '/var/omniture_archive'"
	)

	url := flag.String("url", defaultTarget, targetUsage)
	archiveDir := flag.String("archive", defaultArchiveDir, archiveDirUsage)

	flag.Parse()

	fmt.Println("redirecting to:", *url)
	fmt.Println("archiving to:", *archiveDir)

	op := &OmnitureProxy{
		TargetURL:  *url,
		ArchiveDir: *archiveDir,
	}

	op.Start()

	for {
		// Stayin' alive
	}
}
