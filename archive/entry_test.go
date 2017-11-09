package archive

import (
	"fmt"
	"testing"
	"time"
)

func TestEntryFromString(t *testing.T) {
	raw := "mid=visitorid&pageName=pagename&c.&a.&DeviceName=devicename&.a&.c"
	want := &Entry{
		Time:       time.Now(),
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

	if diff := Diff(want, EntryFromString(raw)); diff != "" {
		t.Errorf("invalid log entry. differs at %s", diff)
	}
}

func Diff(expected, actual *Entry) string {
	leA, leB := *expected, *actual

	if leA.Time.Sub(leB.Time).Seconds() > 1 {
		return fmt.Sprintf("Time (%q, %q)", leA.Time, leB.Time)
	}

	if leA.VisitorID != leB.VisitorID {
		return fmt.Sprintf("VisitorID (%q, %q)", leA.VisitorID, leB.VisitorID)
	}

	if leA.DeviceName != leB.DeviceName {
		return fmt.Sprintf("DeviceName (%q, %q)", leA.DeviceName, leB.DeviceName)
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
