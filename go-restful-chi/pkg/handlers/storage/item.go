package itemStorage

import (
	"database/sql"
	"example.com/chi-restful/pkg/handlers/model"
	"log"
)

type ItemStorage struct {
	conn *sql.DB
}

func New(conn *sql.DB) *ItemStorage {
	return &ItemStorage{
		conn: conn,
	}
}

func (e ItemStorage) GetAll() ([]model.Item, error) {
	rows, err := e.conn.Query("select id, name from items")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var i model.Item
	items := []model.Item{}
	for rows.Next() {
		rows.Scan(&i.Id, &i.Name)
		items = append(items, i)
	}

	return items, nil
}

func (e ItemStorage) CreateItem(item model.Item) error {

	s := `
	INSERT INTO items (id, name)
	VALUES ($1, $2)
	`

	_, err := e.conn.Exec(s, item.Id, item.Name)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (e ItemStorage) GetById(id string) (model.Item, error) {

	var item model.Item
	err := e.conn.
		QueryRow("SELECT id, name FROM items WHERE id = $1", id).
		Scan(&item.Id, &item.Name)

	if err != nil {
		return item, err
	}
	return item, nil
}

func (e ItemStorage) UpdateById(id string, item model.Item) error {
	_, err := e.conn.Exec("UPDATE items SET name = $1 WHERE id = $2", item.Name, id)
	return err
}

func (e ItemStorage) DeleteById(id string) error {
	_, err := e.conn.Exec("DELETE FROM items WHERE id = $1", id)
	return err
}
