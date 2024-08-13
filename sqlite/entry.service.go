package sqlite

import (
	"database/sql"

	. "github.com/slh335/shoppinglistserver"
)

type EntryService struct {
	DB *sql.DB
}

func (m *EntryService) All() (entries []Entry, err error) {
	stmt := "SELECT * FROM entries"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry Entry
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

func (m *EntryService) Complete(id int, completed bool) (updated bool, err error) {
	stmt := "UPDATE entries SET completed=? WHERE id=?"
	res, err := m.DB.Exec(stmt, completed, id)
	if err != nil {
		return false, err
	}

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected > 0, nil
}

func (m *EntryService) Insert(text, category string) (id int, err error) {
	stmt := "INSERT INTO entries (text, category) VALUES (?, ?)"
	res, err := m.DB.Exec(stmt, text, category)
	if err != nil {
		return -1, err
	}

	lastInsertId, _ := res.LastInsertId()
	return int(lastInsertId), nil
}

func (m *EntryService) Update(id int, text, category string) (updated bool, err error) {
	stmt := "UPDATE entries SET text=?, category=? WHERE id=?"
	res, err := m.DB.Exec(stmt, text, category, id)
	if err != nil {
		return false, err
	}

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected > 0, nil
}

func (m *EntryService) Delete(id int) (deleted bool, err error) {
	stmt := "DELETE FROM entries WHERE id=?"
	res, err := m.DB.Exec(stmt, id)
	if err != nil {
		return false, err
	}

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected > 0, nil
}
