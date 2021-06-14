package main

const (
	qPrefix = `
select
  history_items.id, title, url
from history_items
  inner join history_visits on history_visits.history_item = history_items.id
where
`

	qPostfix = `
	group by url
order by visit_time desc
limit 40
`

	dbFilePath = "./History.db"
)
