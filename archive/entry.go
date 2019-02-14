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
	prefix := make(prefix, 0, 5)
	isCtx := false

	for _, p := range parts {
		switch {
		case p == "c.":
			isCtx = true
		case p == ".c":
			isCtx = false
		case strings.HasSuffix(p, "."):
			prefix.push(p)
		case strings.HasPrefix(p, "."):
			prefix.pop()
		default:
			// If we're here then it's a simple k=v pair
			k, v := splitKVP(p)

			// We extract some variables that can be used to index the log entry but we leave them where they were too.
			// mid is the visitor ID
			switch {
			case k == "mid":
				result.VisitorID = v
			case k == "pageName" && !isCtx:
				result.PageName = v
			case k == "DeviceName":
				result.DeviceName = v
			}

			fullKey := prefix.applyTo(k)
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
type prefix []string

func (p *prefix) push(s string) {
	*p = append(*p, s)
}

func (p *prefix) pop() {
	*p = (*p)[:len((*p))-1]
}

func (p *prefix) applyTo(k string) string {
	return strings.Join(*p, "") + k
}

// splitKVP returns two parts split on first `=`
func splitKVP(s string) (key, val string) {
	parts := strings.SplitN(s, "=", 2)
	key = parts[0]
	if len(parts) == 2 {
		val = parts[1]
	}
	return
}
