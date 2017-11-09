package archive

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

type Reader interface {
	GetDates() []string
	GetVisitorsForDate(string) []*Visitor
	GetArchive(string, string) []byte
}

type Visitor struct {
	VID    string `json:"vid"`
	Device string `json:"device"`
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

func (r *fileReader) GetVisitorsForDate(date string) []*Visitor {
	// todo: check date dir exists?
	files, _ := ioutil.ReadDir(path.Join(r.rootDir, date))

	visitors := []*Visitor{}
	for _, f := range files {
		vid, device := splitLogName(f.Name())
		visitors = append(visitors, &Visitor{vid, device})
	}
	return visitors
}

func (r *fileReader) GetArchive(date, visitor string) []byte {
	archiveDir := path.Join(r.rootDir, date)
	files, _ := ioutil.ReadDir(archiveDir)

	var archive string
	for _, f := range files {
		vid, _ := splitLogName(f.Name())
		if vid == visitor {
			archive = path.Join(archiveDir, f.Name())
		}
	}

	if archive == "" {
		return nil
	}

	raw, _ := ioutil.ReadFile(archive)
	// To make this valid JSON, we need to remove the last `,\n` and wrap in `[]`.
	return append(append([]byte{'['}, raw[:len(raw)-2]...), ']')
}

func splitLogName(s string) (string, string) {
	parts := strings.SplitN(s, "|", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	} else {
		return parts[0], "unknown"
	}
}
