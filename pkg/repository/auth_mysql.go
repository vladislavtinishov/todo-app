package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vladislavtinishov/todo-app"
)

type AuthMysql struct {
	db *sqlx.DB
}

func NewAuthMysql(db *sqlx.DB) *AuthMysql {
	return &AuthMysql{db: db}
}

func (r *AuthMysql) CreateUser(user todo.User) (int, error) {
	query := fmt.Sprintf("insert into %s (name, username, password_hash) values (?, ?, ?)", usersTable)
	result, err := r.db.Exec(query, user.Name, user.Username, user.Password)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *AuthMysql) GetUser(username, password string) (todo.User, error) {
	var user todo.User

	query := fmt.Sprintf("select id from %s where username=? and password_hash=?", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
