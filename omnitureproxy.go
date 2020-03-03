package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"omnitureproxy/api"
	"omnitureproxy/archive"
	"omnitureproxy/proxy"
	"omnitureproxy/socket"
	"github.com/urfave/negroni"
)

type OmnitureProxy struct {
	ArchiveDir string
	TargetURL  string

	hub socket.Hub
}

func (op *OmnitureProxy) Notify(entry *archive.Entry) {
	// fmt.Printf("notify: %s\n", entry.PageName)
	op.hub.SendToAll(entry)
}

func (op *OmnitureProxy) Start(port string) {
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
	s := &http.Server{
		Addr:           port,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("listening on", port)
	log.Fatal(s.ListenAndServe())
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":3000"
}

func mustGetEnv(key, msg string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	log.Fatal(msg)
	return ""
}

func main() {
	port := ":" + mustGetEnv("PORT", "No http port set (PORT)")
	archiveDir := mustGetEnv("ARCHIVE_DIR", "No archive dir set (ARCHIVE_DIR)")
	// redirectURL := mustGetEnv("REDIRECT_URL", "No redirect URL set (REDIRECT_URL) ")
	redirectURL := os.Getenv("REDIRECT_URL")

	log.Println("archiving to:", archiveDir)
	if redirectURL == "" {
		log.Println("not redirecting")
	} else {
		log.Println("redirecting to:", redirectURL)
	}

	op := &OmnitureProxy{
		ArchiveDir: archiveDir,
		TargetURL:  redirectURL,
	}

	op.Start(port)

	for {
		// Stayin' alive
	}
}
