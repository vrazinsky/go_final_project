package store

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDb(ctx context.Context) (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbfile := "scheduler.db"
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbfile = envFile
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbfile)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	db, err := sql.Open("sqlite", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	if install {
		_, err = db.ExecContext(ctx, createTableQuery)
		if err != nil {
			return nil, err
		}
		_, err = db.ExecContext(ctx, createIndexQuery)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
