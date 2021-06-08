package main

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
	_ "github.com/mattn/go-sqlite3"
)

type HistoryItem struct {
	ID    int64
	Title string
	URL   string
}

//goland:noinspection SqlNoDataSourceInspection
const query = `-- noinspection SqlResolve
select
       history_items.id, title, url
from history_items
    inner join history_visits on history_visits.history_item = history_items.id
where title like ? or url like ?
	group by url
order by visit_time desc
limit 40 collate nocase
`

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

	wf := aw.New()

	wf.Run(func() {
		defer wf.SendFeedback()

		shf := "./History.db"
		if _, err = os.Stat(shf); os.IsNotExist(err) {
			wf.NewWarningItem("Could not detect Safari history file", err.Error())

			return
		}

		db, err := sql.Open("sqlite3", shf)
		if err != nil {
			wf.NewWarningItem("Could not open the DB", err.Error())

			return
		}

		arg := strings.Join(wf.Args(), " ")
		if len(arg) == 0 {
			arg = "%"
		} else {
			arg = "%" + arg + "%"
		}

		cursor, err := db.Query(query, arg, arg)
		if err != nil {
			wf.NewWarningItem("Could not query Safari History", err.Error())

			return
		}

		defer func() { _ = cursor.Close() }()

		for cursor.Next() {
			var hi HistoryItem

			if err = cursor.Scan(&hi.ID, &hi.Title, &hi.URL); err != nil {
				wf.NewWarningItem("Failed to scan an history item", err.Error())

				return
			}

			wf.NewItem(hi.Title).UID(strconv.FormatInt(hi.ID, 16)).Arg(hi.URL).Valid(true)
		}

		if cursor.Err() != nil {
			wf.NewWarningItem("Cursor error", cursor.Err().Error())
		}
	})
}
