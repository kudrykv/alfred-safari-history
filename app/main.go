package main

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("utf8lower", strings.ToLower, true)
		},
	})
}

func main() {
	var err error

	if os.Getenv("alfred_workflow_data") == "" {
		if err = os.Setenv("alfred_workflow_data", "./tmp/data"); err != nil {
			panic(err)
		}
	}

	if os.Getenv("alfred_workflow_cache") == "" {
		if err = os.Setenv("alfred_workflow_cache", "./tmp/cache"); err != nil {
			panic(err)
		}
	}

	wf := aw.New(aw.MaxResults(0))

	wf.Run(func() {
		defer wf.SendFeedback()

		shf := "./History.db"
		if _, err = os.Stat(shf); os.IsNotExist(err) {
			wf.NewWarningItem("Could not detect Safari history file", err.Error())

			return
		}

		db, err := sql.Open("sqlite3_custom", shf)
		if err != nil {
			wf.NewWarningItem("Could not open the DB", err.Error())

			return
		}

		arg := strings.Join(wf.Args(), " ")
		if len(arg) == 0 {
			arg = "%"
		} else {
			arg = "%" + strings.ToLower(arg) + "%"
		}

		cursor, err := db.Query(query, arg, arg)
		if err != nil {
			wf.NewWarningItem("Could not query Safari History", err.Error())

			return
		}

		defer func() { _ = cursor.Close() }()

		hits := false

		for cursor.Next() {
			hits = true

			hi := HistoryItem{Title: new(string)}

			if err = cursor.Scan(&hi.ID, &hi.Title, &hi.URL); err != nil {
				wf.NewWarningItem("Failed to scan an history item", err.Error())

				return
			}

			var item *aw.Item
			if hi.Title == nil {
				item = wf.NewItem("")
			} else {
				item = wf.NewItem(*hi.Title)
			}

			item.UID(strconv.FormatInt(hi.ID, 16)).Arg(hi.URL).Valid(true)
		}

		if cursor.Err() != nil {
			wf.NewWarningItem("Cursor error", cursor.Err().Error())
		}

		if !hits {
			wf.NewItem("Nothing was found")
		}
	})
}
