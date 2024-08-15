package sqlite

import (
	"database/sql"

	. "github.com/slh335/shoppinglistserver"
)

type ListService struct {
	DB *sql.DB
}

func (m *ListService) Get(id int) (list List, err error) {
	stmt := `
		SELECT lists.id, lists.name, users.id, users.username
		FROM lists
		INNER JOIN users ON lists.creator_id=users.id
		WHERE lists.id=?`
	row := m.DB.QueryRow(stmt, id)

	err = row.Scan(&list.Id, &list.Name, &list.Creator.Id, &list.Creator.Username)
	if err != nil {
		return list, err
	}

	return list, nil
}

func (m *ListService) All(userId int) (lists []List, err error) {
	stmt := `
		SELECT lists.id, lists.name, users.id, users.username
		FROM lists
		INNER JOIN list_members ON lists.id=list_members.list_id
		INNER JOIN users ON lists.creator_id=users.id
		WHERE users.id=?`
	rows, err := m.DB.Query(stmt, userId)
	if err != nil {
		return lists, err
	}
	defer rows.Close()

	for rows.Next() {
		var list List
		err = rows.Scan(&list.Id, &list.Name, &list.Creator.Id, &list.Creator.Username)
		if err != nil {
			return lists, err
		}
		lists = append(lists, list)
	}

	err = rows.Err()
	if err != nil {
		return lists, err
	}
	return lists, nil
}

func (m *ListService) Add(creator User, name string) (list List, err error) {
	stmt := "INSERT INTO lists (name, creator_id) VALUES (?, ?)"
	res, err := m.DB.Exec(stmt, name, creator.Id)
	if err != nil {
		return list, err
	}

	lastInsertId, _ := res.LastInsertId()
	list = List{
		Id:   int(lastInsertId),
		Name: name,
		Creator: User{
			Id:       creator.Id,
			Username: creator.Username,
		},
	}

	err = m.Join(list.Id, creator.Id)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (m *ListService) Delete(listId int) (err error) {
	stmt := "DELETE FROM lists WHERE id=?"
	_, err = m.DB.Exec(stmt, listId)
	if err != nil {
		return err
	}
	return nil
}

func (m *ListService) Join(listId, userId int) (err error) {
	stmt := "INSERT INTO list_members (list_id, user_id) VALUES (?, ?)"
	_, err = m.DB.Exec(stmt, listId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *ListService) Leave(listId, userId int) (err error) {
	stmt := "DELETE FROM list_members WHERE list_id=? AND user_id=?"
	_, err = m.DB.Exec(stmt, listId, userId)
	if err != nil {
		return err
	}
	return nil
}
