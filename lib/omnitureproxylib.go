package omnitureproxy

import (
	"fmt"
	"net/http"

	"github.com/davidamey/omnitureproxy/lib/logs"
	"github.com/davidamey/omnitureproxy/lib/proxy"
	"github.com/davidamey/omnitureproxy/lib/sockets"
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
