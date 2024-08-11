package sqlite

import (
	"database/sql"

	"github.com/slh335/einkaufsliste-server/models"
)

type EntryModel struct {
	DB *sql.DB
}

func (m *EntryModel) All() (entries []models.Entry, err error) {
	rows, err := m.DB.Query("SELECT * FROM entries")
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.Entry
		err = rows.Scan(&entry.Id, &entry.Text, &entry.Category, &entry.Completed)
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}

	err = rows.Err()
	if err != nil {
		return entries, err
	}
	return entries, nil
}
