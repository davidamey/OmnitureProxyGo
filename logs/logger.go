package logs

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"time"
)

type Logger interface {
	StartProcessing()
	StopProcessing()
	Save(*LogEntry)
	HasPendingWrites() bool
}

type fileLogger struct {
	rootDir string
	queue   chan *LogEntry
	quit    chan bool
}

func NewLogger(dir string) Logger {
	return &fileLogger{
		rootDir: dir,
		queue:   make(chan *LogEntry, 100),
		quit:    make(chan bool),
	}
}

func (l *fileLogger) StartProcessing() {
	go func() {
		for {
			select {
			case le := <-l.queue:
				writeLog(l.rootDir, le)

			case <-l.quit:
				return
			}
		}
	}()
}

func (l *fileLogger) StopProcessing() {
	go func() {
		l.quit <- true
	}()
}

func (l *fileLogger) Save(le *LogEntry) {
	l.queue <- le
}

func (l *fileLogger) HasPendingWrites() bool {
	return len(l.queue) > 0
}

func getLogFile(rootDir, vID string) string {
	logDir := path.Join(rootDir, time.Now().Format("2006-01-02"))

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	return path.Join(logDir, vID)
}

func writeLog(rootDir string, le *LogEntry) {
	file := getLogFile(rootDir, le.VisitorID)

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := io.Writer(f)

	b, err := json.Marshal(le)
	if err != nil {
		log.Fatal(err)
	}

	b = append(b, ',')

	w.Write(b)
}
