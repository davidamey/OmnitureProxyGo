package logs

import (
	"io/ioutil"
	"path"
	"regexp"
)

type Fetcher interface {
	GetLogDates() []string
	GetVisitorsForDate(string) []string
	GetLog(string, string) []byte
}

type fileFetcher struct {
	rootDir string
}

var validDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func NewFetcher(dir string) Fetcher {
	return &fileFetcher{rootDir: dir}
}

func (f *fileFetcher) GetLogDates() []string {
	files, _ := ioutil.ReadDir(f.rootDir)

	var dates []string
	for _, f := range files {
		if f.IsDir() && validDate.MatchString(f.Name()) {
			dates = append(dates, f.Name())
		}
	}
	return dates
}

func (f *fileFetcher) GetVisitorsForDate(date string) []string {
	// todo: check date dir exists?

	files, _ := ioutil.ReadDir(path.Join(f.rootDir, date))

	var visitors []string
	for _, f := range files {
		visitors = append(visitors, f.Name())
	}
	return visitors
}

func (f *fileFetcher) GetLog(date, visitor string) []byte {
	logFile := path.Join(f.rootDir, date, visitor)
	raw, _ := ioutil.ReadFile(logFile)

	// To make this valid JSON, we need to remove the last `,` and wrap in `[]`.
	return append(append([]byte{'['}, raw[:len(raw)-1]...), ']')
}
