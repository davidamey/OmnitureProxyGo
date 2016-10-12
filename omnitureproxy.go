package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/davidamey/omnitureproxy/api"
	"github.com/davidamey/omnitureproxy/lib"
	"github.com/davidamey/omnitureproxy/proxy"
	"github.com/davidamey/omnitureproxy/sockets"
)

type OmnitureProxy struct {
	ListenPort string
	TargetURL  string
	LogDir     string
	AssetsDir  string
}

func (op *OmnitureProxy) Start() {
	fmt.Println("Starting")

	// static site
	fs := http.FileServer(http.Dir(op.AssetsDir))
	http.Handle("/", fs)

	// api
	http.Handle("/api/", api.NewApi())

	// socket
	socket := sockets.NewSocket()
	http.Handle("/socket.io/", socket)

	// logger
	logger := logs.NewLogger(op.LogDir)
	logger.StartProcessing()

	// proxy
	p := proxy.NewProxy(socket, logger, proxy.NewProxier(op.TargetURL))
	http.HandleFunc("/b/ss/", p.Handle)

	// start the server
	http.ListenAndServe(op.ListenPort, nil)
}

func main() {
	const (
		defaultPort      = ":3000"
		portUsage        = "server port e.g. ':3000' or ':8080'"
		defaultTarget    = "https://nationwide.sc.omtrdc.net"
		targetUsage      = "redirect url e.g. 'http://localhost:3000'"
		defaultLogDir    = "logs"
		logDirUsage      = "folder to log to e.g. 'logs'"
		defaultAssetsDir = "assets"
		assetsDirUsage   = "folder to server static site from, e.g. 'assets'"
	)

	port := flag.String("port", defaultPort, portUsage)
	url := flag.String("url", defaultTarget, targetUsage)
	logDir := flag.String("logs", defaultLogDir, logDirUsage)
	assetsDir := flag.String("assets", defaultAssetsDir, assetsDirUsage)

	flag.Parse()

	fmt.Println("server will run on:", *port)
	fmt.Println("redirecting to:", *url)
	fmt.Println("logging to:", *logDir)
	fmt.Println("serving site from:", *assetsDir)

	op := &omnitureproxy.OmnitureProxy{
		ListenPort: *port,
		TargetURL:  *url,
		LogDir:     *logDir,
		AssetsDir:  *assetsDir,
	}

	op.Start()

	for {
		// Stay alive forever...is this good practice?
	}
}
