package logs

import (
	"os"
	"path"
	"sort"
	"testing"
)

var fetcherRootDir string = "log_test_dir"

func TestGetLogDates(t *testing.T) {
	// Setup
	os.MkdirAll(path.Join(fetcherRootDir, "2016-01-01"), 0755)
	os.MkdirAll(path.Join(fetcherRootDir, "2015-02-02"), 0755)

	f := NewFetcher(fetcherRootDir)

	// Act
	dates := f.GetLogDates()

	// Assert
	if len(dates) != 2 {
		t.Errorf("incorrect number of dates returned.\n%q", dates)
	}

	sort.Strings(dates)
	want := []string{"2015-02-02", "2016-01-01"}
	for i, _ := range dates {
		if want[i] != dates[i] {
			t.Errorf("date mismatch\nexpected: %q\ngot: %q", want[i], dates[i])
		}
	}
}
