package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8)      NOT NULL DEFAULT "",
    title   VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT         NOT NULL DEFAULT "",
    repeat  VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("проверка файла БД %q: %w", dbFile, err)
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("открытие БД %q: %w", dbFile, err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return fmt.Errorf("ping БД: %w", err)
	}

	db = conn

	if install {
		if err := execSchema(); err != nil {
			return fmt.Errorf("применение схемы: %w", err)
		}
	}

	return nil
}

func execSchema() error {
	for _, stmt := range strings.Split(schema, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("exec %q: %w", stmt, err)
		}
	}
	return nil
}

func GetDB() *sql.DB {
	return db
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
