package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/davidamey/omnitureproxy/logs"
)

var myT *testing.T

var testRequest string = "mid=visitorid&pageName=pagename&c.&a.&DeviceName=devicename&.a&.c"
var expectedLE *logs.LogEntry = &logs.LogEntry{
	Time:      time.Now(),
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

// Mock socket
type mockSocket struct{}

func (s *mockSocket) BroadcastTo(room, message string, args ...interface{}) {
	le := args[0].(*logs.LogEntry)

	if diff := Diff(expectedLE, le); diff != "" {
		myT.Errorf("invalid log entry. differs at %s", diff)
	}
}

// Mock logger
type mockLogger struct{}

func (l *mockLogger) StartProcessing()       {}
func (l *mockLogger) StopProcessing()        {}
func (l *mockLogger) HasPendingWrites() bool { return false }
func (l *mockLogger) Save(le *logs.LogEntry) {
	if diff := Diff(expectedLE, le); diff != "" {
		myT.Errorf("invalid log entry. differs at %s", diff)
	}
}

// Mock Proxy
type mockProxy struct{}

func (p *mockProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// test stuff
	w.Write([]byte("mockresponse"))
}

// Tests
func TestProxy(t *testing.T) {
	// Setup
	myT = t
	r := httptest.NewRequest("GET", "http://omnitureurl.com", strings.NewReader(testRequest))
	w := httptest.NewRecorder()
	p := NewProxy(&mockSocket{}, &mockLogger{}, &mockProxy{})

	// Test
	p.Handle(w, r)

	// Assert
	if w.Body.String() != "mockresponse" {
		t.Errorf("Invalid response:\nexpected: %q\ngot: %q", "mockresponse", w.Body.String())
	}
}

// Helpers
func Diff(expected, actual *logs.LogEntry) string {
	leA, leB := *expected, *actual

	if leA.Time.Sub(leB.Time).Seconds() > 1 {
		return fmt.Sprintf("Time (%q, %q)", leA.Time, leB.Time)
	}

	if leA.VisitorID != leB.VisitorID {
		return fmt.Sprintf("VisitorID (%q, %q)", leA.VisitorID, leB.VisitorID)
	}

	if leA.PageName != leB.PageName {
		return fmt.Sprintf("PageName (%q, %q)", leA.PageName, leB.PageName)
	}

	if len(leA.AdditionalData) != len(leB.AdditionalData) {
		return fmt.Sprintf("AdditionalData length (%q, %q)", leA.AdditionalData, leB.AdditionalData)
	}

	for k := range leA.AdditionalData {
		if leA.AdditionalData[k] != leB.AdditionalData[k] {
			return "AdditionalData key: " + k
		}
	}

	if len(leA.ContextData) != len(leB.ContextData) {
		return fmt.Sprintf("ContextData length (%q, %q)", leA.ContextData, leB.ContextData)
	}

	for k := range leA.ContextData {
		if leA.ContextData[k] != leB.ContextData[k] {
			return "ContextData key: " + k
		}
	}

	return ""
}
