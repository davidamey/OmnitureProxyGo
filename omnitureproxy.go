package main

import (
	"flag"
	"fmt"

	"github.com/davidamey/omnitureproxy/lib"
)

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
