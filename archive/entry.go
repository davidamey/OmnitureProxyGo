package archive

import (
	"strings"
	"time"
)

type Entry struct {
	Time           time.Time         `json:"time"`
	VisitorID      string            `json:"visitorID"`
	DeviceName     string            `json:"deviceName"`
	PageName       string            `json:"pageName"`
	AdditionalData map[string]string `json:"additionalData"`
	ContextData    map[string]string `json:"contextData"`
}

func EntryFromBytes(b []byte) *Entry {
	return EntryFromString(string(b))
}

func EntryFromString(raw string) *Entry {
	result := Entry{
		DeviceName:     "unknown",
		Time:           time.Now(),
		AdditionalData: make(map[string]string),
		ContextData:    make(map[string]string),
	}
	parts := strings.Split(raw, "&")
	prefix := prefix{parts: make([]string, 0, 5)}
	isCtx := false

	for _, p := range parts {
		if p == "c." {
			isCtx = true
		} else if p == ".c" {
			isCtx = false
		} else if strings.HasSuffix(p, ".") {
			prefix.push(p)
		} else if strings.HasPrefix(p, ".") {
			prefix.pop()
		} else {
			// If we're here then it's a simple k=v pair
			k, v := splitKVP(p)

			// We extract some variables that can be used to index the log entry but we leave them where they were too.
			// mid is the visitor ID
			if k == "mid" {
				result.VisitorID = v
			} else if k == "pageName" && !isCtx {
				result.PageName = v
			} else if k == "DeviceName" {
				result.DeviceName = v
			}

			fullKey := strings.Join(prefix.parts, "") + k
			if isCtx {
				result.ContextData[fullKey] = v
			} else {
				result.AdditionalData[fullKey] = v
			}
		}
	}

	return &result
}

// prefix is simply a wrapper around a string array that is eventually 'joined'
type prefix struct {
	parts []string
}

func (p *prefix) push(s string) {
	p.parts = append(p.parts, s)
}

func (p *prefix) pop() {
	p.parts = p.parts[:len(p.parts)-1]
}

// splitKVP returns two parts split on first `=`
func splitKVP(s string) (string, string) {
	parts := strings.SplitN(s, "=", 2)

	if len(parts) == 2 {
		return parts[0], parts[1]
	} else {
		return parts[0], ""
	}
}
