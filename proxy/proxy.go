package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"omnitureproxy/archive"
)

type Notifier interface {
	Notify(*archive.Entry)
}

type Proxier interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type proxy struct {
	archiver archive.Writer
	notifier Notifier
	proxier  Proxier
}

func (p *proxy) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	// Reading the body clears the reader so put a new one in its place
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	entry := archive.EntryFromBytes(body)
	p.archiver.Save(entry)
	p.notifier.Notify(entry)

	if p.proxier == nil {
		w.WriteHeader(200)
		fmt.Println(w, "No proxy")
	} else {
		p.proxier.ServeHTTP(w, r)
	}
}

func newProxier(target string) Proxier {
	if target == "" {
		return nil
	}
	url, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(url)
}

func New(archiver archive.Writer, notifier Notifier, proxyURL string) *proxy {
	return &proxy{archiver, notifier, newProxier(proxyURL)}
}
