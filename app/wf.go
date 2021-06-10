package main

import (
	"database/sql"
	"errors"
	"os"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
)

type HistoryItem struct {
	ID    int64
	Title *string
	URL   string
}

func wfRunner(wf *aw.Workflow) func() {
	return func() {
		defer wf.SendFeedback()

		items, err := flow(strings.Join(wf.Args(), " "))
		if err != nil {
			var itemErr *Error
			if errors.As(err, &itemErr) {
				wf.NewWarningItem(itemErr.title, itemErr.Error())
			} else {
				wf.NewWarningItem("Unknown error has happened", err.Error())
			}

			return
		}

		if len(items) == 0 {
			wf.NewItem("Nothing was found")

			return
		}

		for _, item := range items {
			title := item.URL
			addSubtitle := false

			if item.Title != nil && len(*item.Title) > 0 {
				title = *item.Title
				addSubtitle = true
			}

			wfItem := wf.NewItem(title).
				Arg(item.URL).
				UID(strconv.FormatInt(item.ID, 16)).
				Valid(true)

			if addSubtitle {
				wfItem.Subtitle(item.URL)
			}
		}
	}
}

func flow(search string) ([]HistoryItem, error) {
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		return nil, Error{title: "Could not detect Safari history file"}
	}

	db, err := sql.Open("sqlite3_custom", dbFilePath)
	if err != nil {
		return nil, Error{title: "Could not open the DB", message: err}
	}

	search = prepSearch(search)

	cursor, err := db.Query(query, search, search)
	if err != nil {
		return nil, Error{title: "Could not query Safari History", message: err}
	}

	defer func() { _ = cursor.Close() }()

	his := make([]HistoryItem, 0, 40)

	for cursor.Next() {
		hi := HistoryItem{Title: new(string)}

		if err = cursor.Scan(&hi.ID, &hi.Title, &hi.URL); err != nil {
			return nil, Error{title: "Failed to scan a history item", message: err}
		}

		his = append(his, hi)
	}

	if cursor.Err() != nil {
		return nil, Error{title: "Cursor error", message: cursor.Err()}
	}

	return his, nil
}

func prepSearch(query string) string {
	if len(query) == 0 {
		return "%"
	}

	return "%" + strings.ToLower(query) + "%"
}

type Error struct {
	title   string
	message error
}

func (e Error) Error() string {
	return e.title + ": " + e.message.Error()
}
