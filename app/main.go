package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

//nolint:gochecknoinits
func init() {
	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			err := conn.RegisterFunc("utf8lower", strings.ToLower, true)
			if err != nil {
				return fmt.Errorf("register utf8lower func: %w", err)
			}

			return nil
		},
	})

	if os.Getenv("alfred_workflow_data") == "" {
		if err := os.Setenv("alfred_workflow_data", "./tmp/data"); err != nil {
			panic(err)
		}
	}

	if os.Getenv("alfred_workflow_cache") == "" {
		if err := os.Setenv("alfred_workflow_cache", "./tmp/cache"); err != nil {
			panic(err)
		}
	}
}

func main() {
	wf := aw.New()

	wf.Run(wfRunner(wf))
}
