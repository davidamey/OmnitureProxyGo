package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/davidamey/omnitureproxy/archive"
)

type Proxier interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Socket interface {
	BroadcastTo(string, string, ...interface{})
}

type proxyWrapper struct {
	archiver archive.Writer
	socket   Socket
	proxy    Proxier
}

func (p *proxyWrapper) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	// Reading the body clears the reader so put a new one in its place
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	fmt.Println("broadcasting")
	entry := archive.EntryFromBytes(body)
	p.archiver.Save(entry)
	p.socket.BroadcastTo("clients", "entry", entry)

	if p.proxy == nil {
		w.WriteHeader(200)
		fmt.Println(w, "No proxy")
	} else {
		p.proxy.ServeHTTP(w, r)
	}
}

func NewProxier(target string) Proxier {
	if target == "" {
		return nil
	} else {
		url, _ := url.Parse(target)
		return httputil.NewSingleHostReverseProxy(url)
	}
}

func NewProxy(socket Socket, archiver archive.Writer, proxy Proxier) *proxyWrapper {
	return &proxyWrapper{
		archiver: archiver,
		socket:   socket,
		proxy:    proxy,
	}
}
