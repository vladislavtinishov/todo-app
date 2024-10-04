package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vladislavtinishov/todo-app"
	"strings"
)

type TodoListMysql struct {
	db *sqlx.DB
}

func NewTodoListMysql(db *sqlx.DB) *TodoListMysql {
	return &TodoListMysql{db: db}
}

func (r *TodoListMysql) Create(userId int, todoList todo.TodoList) (int, error) {
	tx, err := r.db.Begin()

	if err != nil {
		return 0, err
	}
	createListQuery := fmt.Sprintf("insert into %s (title, description) values (?, ?)", todoListsTable)
	result, err := r.db.Exec(createListQuery, todoList.Title, todoList.Description)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	listId, err := result.LastInsertId()

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("insert into %s (user_id, list_id) values (?, ?)", usersListsTable)
	_, err = r.db.Exec(createUsersListQuery, userId, listId)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(listId), tx.Commit()
}

func (r *TodoListMysql) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	getAllQuery := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s ul on tl.id = ul.list_id where ul.user_id = ?", todoListsTable, usersListsTable)
	err := r.db.Select(&lists, getAllQuery, userId)

	return lists, err
}

func (r *TodoListMysql) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s ul on tl.id = ul.list_id where ul.user_id = ? and tl.id = ?", todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	if err != nil {
		return list, err
	}

	return list, nil
}

func (r *TodoListMysql) Delete(userId, listId int) error {
	query := fmt.Sprintf("delete tl from %s tl inner join %s ul on tl.id = ul.list_id where ul.user_id = ? and ul.list_id = ?", todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListMysql) Update(userId, listId int, input todo.UpdateListInput) error {
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

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("update %s tl inner join %s ul on tl.id = ul.list_id set %s where ul.user_id = ? and tl.id = ?", todoListsTable, usersListsTable, setQuery)

	args = append(args, userId, listId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debug("args: %s", args)
	fmt.Println(query)
	_, err := r.db.Exec(query, args...)

	return err
}
