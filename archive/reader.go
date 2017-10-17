package archive

import (
	"io/ioutil"
	"path"
	"regexp"
)

type Reader interface {
	GetDates() []string
	GetVisitorsForDate(string) []string
	GetArchive(string, string) []byte
}

type fileReader struct {
	rootDir string
}

var validDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func NewReader(dir string) Reader {
	return &fileReader{rootDir: dir}
}

func (r *fileReader) GetDates() []string {
	files, _ := ioutil.ReadDir(r.rootDir)

	dates := []string{}
	for _, f := range files {
		if f.IsDir() && validDate.MatchString(f.Name()) {
			dates = append(dates, f.Name())
		}
	}
	return dates
}

func (r *fileReader) GetVisitorsForDate(date string) []string {
	// todo: check date dir exists?

	files, _ := ioutil.ReadDir(path.Join(r.rootDir, date))

	visitors := []string{}
	for _, f := range files {
		visitors = append(visitors, f.Name())
	}
	return visitors
}

func (r *fileReader) GetArchive(date, visitor string) []byte {
	logFile := path.Join(r.rootDir, date, visitor)
	raw, _ := ioutil.ReadFile(logFile)

	// To make this valid JSON, we need to remove the last `,\n` and wrap in `[]`.
	return append(append([]byte{'['}, raw[:len(raw)-2]...), ']')
}
