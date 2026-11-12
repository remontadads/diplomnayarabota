package api

import (
	"net/http"
	"strconv"

	"diplomnayarabota/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.FormValue("search")

	limit := 50
	if s := r.FormValue("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			limit = n
		}
	}

	tasks, err := db.Tasks(search, limit)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения задач: " + err.Error()})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
