package proxy

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"omnitureproxy/archive"
)

var myT *testing.T

var testRequest string = "mid=visitorid&pageName=pagename&c.&a.&DeviceName=devicename&.a&.c"
var expectedLE *archive.Entry = &archive.Entry{
	Time:      time.Now(),
	VisitorID: "1234567890",
	PageName:  "pagename",
	AdditionalData: map[string]string{
		"mid":      "1234567890",
		"pageName": "pagename",
	},
	ContextData: map[string]string{
		"a.DeviceName": "devicename",
	},
}

type mockNotifier struct{}

func (s *mockNotifier) Notify(entry *archive.Entry) {
	if diff := Diff(expectedLE, entry); diff != "" {
		myT.Errorf("invalid log entry. differs at %s", diff)
	}
}

type mockArchiver struct{}

func (l *mockArchiver) StartProcessing()       {}
func (l *mockArchiver) StopProcessing()        {}
func (l *mockArchiver) HasPendingWrites() bool { return false }
func (l *mockArchiver) Save(le *archive.Entry) {
	if diff := Diff(expectedLE, le); diff != "" {
		myT.Errorf("invalid log entry. differs at %s", diff)
	}
}

// Tests
func TestProxy(t *testing.T) {
	// Setup
	myT = t
	r := httptest.NewRequest("GET", "http://omniture-url.com", strings.NewReader(testRequest))
	w := httptest.NewRecorder()
	p := New(&mockArchiver{}, &mockNotifier{}, "")

	// Test
	p.Handle(w, r)

	// Assert
	if w.Body.String() != "" {
		t.Errorf("Invalid response:\nexpected: %q\ngot: %q", "mockresponse", w.Body.String())
	}
}

// Helpers
func Diff(expected, actual *archive.Entry) string {
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
