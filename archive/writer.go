package archive

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

type Writer interface {
	StartProcessing()
	StopProcessing()
	Save(*Entry)
	HasPendingWrites() bool
}

type fileWriter struct {
	rootDir string
	queue   chan *Entry
	quit    chan struct{}
}

func NewWriter(dir string) Writer {
	return &fileWriter{
		rootDir: dir,
		queue:   make(chan *Entry, 100),
		quit:    make(chan struct{}),
	}
}

func (w *fileWriter) StartProcessing() {
	go func() {
		for {
			select {
			case entry := <-w.queue:
				writeEntry(w.rootDir, entry)

			case <-w.quit:
				return
			}
		}
	}()
}

func (w *fileWriter) StopProcessing() {
	go func() {
		w.quit <- struct{}{}
	}()
}

func (w *fileWriter) Save(entry *Entry) {
	w.queue <- entry
}

func (w *fileWriter) HasPendingWrites() bool {
	return len(w.queue) > 0
}

func getArchive(rootDir, vID, device string) string {
	dir := path.Join(rootDir, time.Now().Format("2006-01-02"))

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return path.Join(dir, fmt.Sprintf("%s|%s", vID, device))
}

func writeEntry(rootDir string, entry *Entry) {
	file := getArchive(rootDir, entry.VisitorID, entry.DeviceName)

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := io.Writer(f)

	b, err := json.Marshal(entry)
	if err != nil {
		log.Fatal(err)
	}

	b = append(b, ',', '\n')

	w.Write(b)
}
