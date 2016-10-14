package api

import "net/http"

func csvHandler(w http.ResponseWriter, r *http.Request) {
	// match := rgxDate.FindStringSubmatch(r.URL.Path)

	// if match == nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("Bad request."))
	// 	return
	// }

	// entries := GetLogEntriesForYMD(match[1], match[2], match[3])
	// lines := make([]string, len(*entries)+1)
	// lines[0] = "Device Name,Date,Page,Load Time"
	// for i, e := range *entries {
	// 	lines[i+1] = fmt.Sprintf("%s,%s,%s,%s", e.ContextData["a.DeviceName"], e.Time, e.PageName, e.ContextData["dba.page.pageinfo.page_load_time"])
	// }

	// output, err := url.QueryUnescape(strings.Join(lines, "\n"))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "text/plain")
	// fmt.Fprint(w, output)
}
