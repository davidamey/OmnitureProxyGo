package logs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

var rootDir string = "log_test_dir"
var le *LogEntry = &LogEntry{
	Time:      time.Date(2016, time.June, 28, 8, 13, 0, 0, time.UTC),
	VisitorID: "visitorid",
	PageName:  "pagename",
	AdditionalData: map[string]string{
		"mid":      "visitorid",
		"pageName": "pagename",
	},
	ContextData: map[string]string{
		"a.DeviceName": "devicename",
	},
}

func assertLog(t *testing.T, vID, want string) {
	logFile := path.Join(rootDir, time.Now().Format("2006-01-02"), vID)

	got, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Errorf("unable to read log file %s", logFile)
	}

	if string(got) != want {
		t.Errorf("log mismatch:\nexpected: %q\ngot: %q", want, string(got))
	}
}

func cleanUp() {
	os.RemoveAll(rootDir)
}

func TestLoggerSingleVisitor(t *testing.T) {
	defer cleanUp()

	// Setup
	l := NewLogger(rootDir)
	l.StartProcessing()
	defer l.StopProcessing()

	// Test
	l.Save(le)

	for l.HasPendingWrites() {
		time.Sleep(time.Second)
	}

	// Assert
	assertLog(t, le.VisitorID, "{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},")
}

func TestLoggerSameVisitor(t *testing.T) {
	defer cleanUp()

	// Setup
	l := NewLogger(rootDir)
	l.StartProcessing()
	defer l.StopProcessing()

	// Test
	l.Save(le)
	l.Save(le)

	for l.HasPendingWrites() {
		time.Sleep(time.Second)
	}

	// Assert
	assertLog(t, le.VisitorID,
		"{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},"+
			"{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},")
}

func TestLoggerMultipleVisitors(t *testing.T) {
	defer cleanUp()

	// Setup
	le2 := &LogEntry{
		Time:      time.Date(2016, time.June, 28, 8, 13, 0, 0, time.UTC),
		VisitorID: "visitorid2",
		PageName:  "pagename",
		AdditionalData: map[string]string{
			"mid":      "visitorid2",
			"pageName": "pagename",
		},
		ContextData: map[string]string{
			"a.DeviceName": "devicename",
		},
	}

	// Setup
	l := NewLogger(rootDir)
	l.StartProcessing()
	defer l.StopProcessing()

	// Test
	l.Save(le)
	l.Save(le2)

	for l.HasPendingWrites() {
		time.Sleep(time.Second)
	}

	// Assert
	assertLog(t, le.VisitorID, "{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},")
	assertLog(t, le2.VisitorID, "{\"time\":\"2016-06-28T08:13:00Z\",\"visitorID\":\"visitorid2\",\"pageName\":\"pagename\",\"additionalData\":{\"mid\":\"visitorid2\",\"pageName\":\"pagename\"},\"contextData\":{\"a.DeviceName\":\"devicename\"}},")
}

func BenchmarkLogger(b *testing.B) {
	defer cleanUp()

	// Setup
	l := NewLogger(rootDir)
	l.StartProcessing()
	defer l.StopProcessing()

	// Test
	for i := 0; i < b.N; i++ {
		le.VisitorID = fmt.Sprintf("visitorId%d", i)
		l.Save(le)
	}

	for l.HasPendingWrites() {
		time.Sleep(time.Second)
	}
}
