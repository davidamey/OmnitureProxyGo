package archive

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

const rootDir string = "archive_test_dir"

var entry *Entry = &Entry{
	Time:       time.Date(2016, time.June, 28, 8, 13, 0, 0, time.UTC),
	VisitorID:  "visitorid",
	DeviceName: "devicename",
	PageName:   "pagename",
	AdditionalData: map[string]string{
		"mid":      "visitorid",
		"pageName": "pagename",
	},
	ContextData: map[string]string{
		"a.DeviceName": "devicename",
	},
}

func assertArchive(t *testing.T, entry *Entry, want string) {
	archive := path.Join(
		rootDir,
		time.Now().Format("2006-01-02"),
		fmt.Sprintf("%s|%s", entry.VisitorID, entry.DeviceName))

	got, err := ioutil.ReadFile(archive)
	if err != nil {
		t.Errorf("unable to read archive %s", archive)
	}

	if string(got) != want {
		t.Errorf("archive mismatch:\nexpected: %q\ngot: %q", want, string(got))
	}
}

func cleanUp() {
	os.RemoveAll(rootDir)
}

func TestWriterSingleVisitor(t *testing.T) {
	defer cleanUp()

	// Setup
	w := NewWriter(rootDir)
	w.StartProcessing()
	defer w.StopProcessing()

	// Test
	w.Save(entry)

	for w.HasPendingWrites() {
		time.Sleep(time.Second)
	}

	// Assert
	assertArchive(t, entry,
		"{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid\",\"deviceName\":\"devicename\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},\n")
}

func TestWriterSameVisitor(t *testing.T) {
	defer cleanUp()

	// Setup
	w := NewWriter(rootDir)
	w.StartProcessing()
	defer w.StopProcessing()

	// Test
	w.Save(entry)
	w.Save(entry)

	for w.HasPendingWrites() {
		time.Sleep(time.Second)
	}

	// Assert
	assertArchive(t, entry,
		"{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid\",\"deviceName\":\"devicename\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},\n"+
			"{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid\",\"deviceName\":\"devicename\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},\n")
}

func TestWriterMultipleVisitors(t *testing.T) {
	defer cleanUp()

	// Setup
	entry2 := &Entry{
		Time:       time.Date(2016, time.June, 28, 8, 13, 0, 0, time.UTC),
		VisitorID:  "visitorid2",
		DeviceName: "devicename",
		PageName:   "pagename",
		AdditionalData: map[string]string{
			"mid":      "visitorid2",
			"pageName": "pagename",
		},
		ContextData: map[string]string{
			"a.DeviceName": "devicename",
		},
	}

	// Setup
	w := NewWriter(rootDir)
	w.StartProcessing()
	defer w.StopProcessing()

	// Test
	w.Save(entry)
	w.Save(entry2)

	for w.HasPendingWrites() {
		time.Sleep(time.Second)
	}

	// Assert
	assertArchive(t, entry, "{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid\",\"deviceName\":\"devicename\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},\n")
	assertArchive(t, entry2, "{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid2\",\"deviceName\":\"devicename\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid2\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},\n")
}

func BenchmarkWriter(b *testing.B) {
	defer cleanUp()

	// Setup
	w := NewWriter(rootDir)
	w.StartProcessing()
	defer w.StopProcessing()

	// Test
	for i := 0; i < b.N; i++ {
		entry.VisitorID = fmt.Sprintf("visitorId%d", i)
		w.Save(entry)
	}

	for w.HasPendingWrites() {
		time.Sleep(time.Second)
	}
}
