package api

import (
	"net/http"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	if s := r.FormValue("now"); s != "" {
		t, err := time.Parse(dateFormat, s)
		if err != nil {
			http.Error(w, "invalid now: "+err.Error(), http.StatusBadRequest)
			return
		}
		now = t
	}

	result, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(result))
}
