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
	ListenPort string
	TargetURL  string
	ArchiveDir string
	AssetsDir  string

	hub socket.Hub
}

func (op *OmnitureProxy) Notify(entry *archive.Entry) {
	// fmt.Printf("notify: %s\n", entry.PageName)
	op.hub.SendToAll(entry)
}

func (op *OmnitureProxy) Start() {
	log.Println("Starting")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(op.AssetsDir)))
	mux.Handle("/api/", api.NewApi())

	// socket
	op.hub = socket.NewHub()
	mux.HandleFunc("/ws", op.hub.ServeWS)

	// archive
	archiver := archive.NewWriter(op.ArchiveDir)
	archiver.StartProcessing()
	defer archiver.StopProcessing()

	// proxy
	p := proxy.NewProxy(archiver, op, op.TargetURL)
	mux.HandleFunc("/b/ss/", p.Handle)

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(mux)

	// start the server
	s := &http.Server{
		Addr:           detectAddress(),
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
	// log.Fatal(s.ListenAndServeTLS("cantor.cer", "cantor.key"))
}

func detectAddress(addr ...string) string {
	if len(addr) > 0 {
		return addr[0]
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":3000"
}

func main() {
	const (
		defaultPort       = ":3001"
		portUsage         = "server port e.g. ':3000' or ':8080'"
		defaultTarget     = "http://localhost:5000"
		targetUsage       = "redirect url e.g. 'http://localhost:3000'"
		defaultArchiveDir = "_archive"
		archiveDirUsage   = "folder to log to e.g. 'archive'"
		defaultAssetsDir  = "_assets"
		assetsDirUsage    = "folder to server static site from, e.g. 'assets'"
	)

	port := flag.String("port", defaultPort, portUsage)
	url := flag.String("url", defaultTarget, targetUsage)
	archiveDir := flag.String("archive", defaultArchiveDir, archiveDirUsage)
	assetsDir := flag.String("assets", defaultAssetsDir, assetsDirUsage)

	flag.Parse()

	fmt.Println("server will run on:", *port)
	fmt.Println("redirecting to:", *url)
	fmt.Println("archiving to:", *archiveDir)
	fmt.Println("serving site from:", *assetsDir)

	op := &OmnitureProxy{
		ListenPort: *port,
		TargetURL:  *url,
		ArchiveDir: *archiveDir,
		AssetsDir:  *assetsDir,
	}

	op.Start()

	for {
		// Stay alive forever...is this good practice?
	}
}
