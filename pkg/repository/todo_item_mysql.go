package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vladislavtinishov/todo-app"
	"strings"
)

type TodoItemMysql struct {
	db *sqlx.DB
}

func NewTodoItemMysql(db *sqlx.DB) *TodoItemMysql {
	return &TodoItemMysql{db: db}
}

func (r *TodoItemMysql) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createItemQuery := fmt.Sprintf("insert into %s (title, description) values (?, ?)", todoItemsTable)
	row, err := tx.Exec(createItemQuery, item.Title, item.Description)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	itemId, err := row.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("insert into %s (list_id, item_id) values(?, ?)", listsItemsTable)
	row, err = tx.Exec(createListItemsQuery, listId, itemId)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return int(itemId), nil
}

func (r *TodoItemMysql) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem

	query := fmt.Sprintf(`
		select ti.id, ti.title, ti.description, ti.done from %s ti inner join %s li on li.item_id = ti.id 
		inner join %s ul on ul.list_id = li.list_id
		where li.list_id = ? and ul.user_id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemMysql) GetById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem

	query := fmt.Sprintf(`
		select ti.id, ti.title, ti.description, ti.done from %s ti inner join %s li on li.item_id = ti.id 
		inner join %s ul on ul.list_id = li.list_id
		where ti.id = ? and ul.user_id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItemMysql) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`delete ti from %s ti inner join %s li on li.item_id = ti.id 
									inner join %s ul on ul.list_id = li.list_id
									where ti.id = ? and ul.user_id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable)

	_, err := r.db.Exec(query, itemId, userId)

	return err
}

func (r *TodoItemMysql) Update(userId, itemId int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	if input.Title != nil {
		setValues = append(setValues, "title=?")
		args = append(args, *input.Title)
	}

	if input.Description != nil {
		setValues = append(setValues, "description=?")
		args = append(args, *input.Description)
	}

	if input.Description != nil {
		setValues = append(setValues, "description=?")
		args = append(args, *input.Description)
	}

	if input.Done != nil {
		setValues = append(setValues, "done=?")
		args = append(args, *input.Done)
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("update %s ti inner join %s li on li.item_id = ti.id inner join %s ul on ul.list_id = li.list_id set %s where ti.id = ? and ul.user_id = ?", todoItemsTable, listsItemsTable, usersListsTable, setQuery)

	args = append(args, userId, itemId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debug("args: %s", args)
	fmt.Println(query)
	_, err := r.db.Exec(query, args...)

	return err
}

func (r *TodoItemMysql) SetDoneStatus(userId, itemId, status int) error {
	query := fmt.Sprintf("update %s ti inner join %s li on li.item_id = ti.id inner join %s ul on ul.list_id = li.list_id set done = ? where ti.id = ? and ul.user_id = ?", todoItemsTable, listsItemsTable, usersListsTable)

	_, err := r.db.Exec(query, status, itemId, userId)

	return err
}
