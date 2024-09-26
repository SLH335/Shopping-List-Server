package sqlite

import (
	"database/sql"
	"time"

	. "github.com/slh335/shoppinglistserver"
)

type EntryService struct {
	DB *sql.DB
}

func (m *EntryService) Get(id int) (entry Entry, err error) {
	stmt := "SELECT * FROM entries WHERE id=? ORDER BY created_at"
	row := m.DB.QueryRow(stmt, id)

	var createdAtStr string
	err = row.Scan(&entry.Id, &entry.ListId, &entry.Text, &entry.Category, &entry.OrderIndex, &entry.Completed, &createdAtStr)
	if err != nil {
		return entry, err
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return entry, err
	}
	entry.CreatedAt = createdAt

	return entry, nil
}

func (m *EntryService) All(listId int) (entries []Entry, err error) {
	stmt := "SELECT * FROM entries WHERE list_id=? ORDER BY category, order_index"
	rows, err := m.DB.Query(stmt, listId)
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry Entry
		var createdAtStr string
		err = rows.Scan(&entry.Id, &entry.ListId, &entry.Text, &entry.Category, &entry.OrderIndex, &entry.Completed, &createdAtStr)
		if err != nil {
			return entries, err
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return entries, err
		}
		entry.CreatedAt = createdAt

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

func (m *EntryService) Move(listId int, category string, oldIndex, newIndex int) (updated bool, err error) {
	if oldIndex == newIndex {
		return true, nil
	}

	var res sql.Result
	if oldIndex < newIndex {
		stmt := `UPDATE entries
			SET order_index = CASE
				WHEN order_index=? THEN ?
				WHEN order_index>? AND order_index<? THEN order_index-1
			END
			WHERE list_id=? AND category=? AND order_index>=? AND order_index<?`
		res, err = m.DB.Exec(stmt, oldIndex, newIndex-1, oldIndex, newIndex, listId, category, oldIndex, newIndex)
	} else {
		stmt := `UPDATE entries
			SET order_index = CASE
				WHEN order_index=? THEN ?
				WHEN order_index>=? AND order_index<? THEN order_index+1
			END
			WHERE list_id=? AND category=? AND order_index>=? AND order_index<=?`
		res, err = m.DB.Exec(stmt, oldIndex, newIndex, newIndex, oldIndex, listId, category, newIndex, oldIndex)
	}
	if err != nil {
		return false, err
	}

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected > 0, nil
}

func (m *EntryService) Add(listId int, text, category string) (entry Entry, err error) {
	createdAt := time.Now()
	stmt := `INSERT INTO entries (list_id, text, category, order_index, created_at)
		VALUES (?, ?, ?, (SELECT IFNULL(MAX(order_index), -1) + 1 FROM entries WHERE list_id=? AND category=?), ?)`
	res, err := m.DB.Exec(stmt, listId, text, category, listId, category, createdAt.Format(time.RFC3339))
	if err != nil {
		return entry, err
	}
	lastInsertId, _ := res.LastInsertId()

	stmt = "SELECT IFNULL(MAX(order_index), 0) FROM entries WHERE list_id=? AND category=?"
	row := m.DB.QueryRow(stmt, listId, category)
	var orderIndex int
	err = row.Scan(&orderIndex)
	if err != nil {
		return entry, err
	}

	entry = Entry{
		Id:         int(lastInsertId),
		ListId:     listId,
		Text:       text,
		Category:   category,
		OrderIndex: orderIndex,
		Completed:  false,
		CreatedAt:  createdAt,
	}
	return entry, nil
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
