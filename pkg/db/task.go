package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat)
	          VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(search string, limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = 50
	}

	var (
		rows *sql.Rows
		err  error
	)

	switch {
	case search == "":
		rows, err = db.Query(
			`SELECT id, date, title, comment, repeat
			 FROM scheduler
			 ORDER BY date ASC
			 LIMIT ?`, limit)

	case isSearchDate(search):
		rows, err = db.Query(
			`SELECT id, date, title, comment, repeat
			 FROM scheduler
			 WHERE date = ?
			 ORDER BY date ASC
			 LIMIT ?`, toDBDate(search), limit)

	default:
		like := "%" + search + "%"
		rows, err = db.Query(
			`SELECT id, date, title, comment, repeat
			 FROM scheduler
			 WHERE title LIKE ? OR comment LIKE ?
			 ORDER BY date ASC
			 LIMIT ?`, like, like, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0, limit)
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func isSearchDate(s string) bool {
	_, err := time.Parse("02.01.2006", s)
	return err == nil
}

func toDBDate(s string) string {
	t, _ := time.Parse("02.01.2006", s)
	return t.Format("20060102")
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача с id %d не найдена", task.ID)
	}

	return nil
}

func GetTask(id string) (*Task, error) {
	var t Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func DeleteTask(id string) error {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача с id %s не найдена", id)
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача с id %s не найдена", id)
	}
	return nil
}
