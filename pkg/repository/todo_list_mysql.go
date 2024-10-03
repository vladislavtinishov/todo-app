package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vladislavtinishov/todo-app"
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
