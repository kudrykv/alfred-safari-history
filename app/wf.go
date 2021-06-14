package main

import (
	"database/sql"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
)

type HistoryItem struct {
	ID    int64
	Title *string
	URL   string
}

var rmMultSpacesRegexp = regexp.MustCompile(`\s+`)

func wfRunner(wf *aw.Workflow) func() {
	return func() {
		defer wf.SendFeedback()

		items, err := flow(createTerms(wf.Args()))
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

func flow(terms []string) ([]HistoryItem, error) {
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		return nil, Error{title: "Could not detect Safari history file"}
	}

	db, err := sql.Open("sqlite3_custom", dbFilePath)
	if err != nil {
		return nil, Error{title: "Could not open the DB", message: err}
	}

	titles := make([]string, 0, len(terms))
	urls := make([]string, 0, len(terms))

	for i := range terms {
		titles = append(titles, "utf8lower(ifnull(title, '')) like ?"+strconv.Itoa(i+1))
		urls = append(urls, "utf8lower(ifnull(url, '')) like ?"+strconv.Itoa(i+1))
	}

	q := qPrefix + " (" + strings.Join(titles, " and ") + ") or (" + strings.Join(urls, " and ") + ") " + qPostfix

	cursor, err := db.Query(q, prepTerms(terms)...)
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

func createTerms(slice []string) []string {
	join := strings.Join(slice, " ")
	lower := strings.ToLower(join)
	squashSpaces := rmMultSpacesRegexp.ReplaceAllString(lower, " ")

	return strings.Split(squashSpaces, " ")
}

func prepTerms(slice []string) []interface{} {
	if len(slice) == 0 {
		return nil
	}

	out := make([]interface{}, 0, len(slice))

	for i := range slice {
		out = append(out, prepTerm(slice[i]))
	}

	return out
}

func prepTerm(term string) string {
	if len(term) == 0 {
		return "%"
	}

	return "%" + term + "%"
}

type Error struct {
	title   string
	message error
}

func (e Error) Error() string {
	return e.title + ": " + e.message.Error()
}
