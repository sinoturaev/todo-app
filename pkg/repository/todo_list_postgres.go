package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	todo_app "github.com/sinoturaev/todo-app"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo_app.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUserListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo_app.TodoList, error) {
	var lists []todo_app.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo_app.TodoList, error) {
	var list todo_app.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2",
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input todo_app.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	ardId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", ardId))
		args = append(args, *input.Title)
		ardId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("Description=$%d", ardId))
		args = append(args, *input.Description)
		ardId++
	}

	//title=$1
	//description=$1
	//title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, ardId, ardId+1)
	args = append(args, listId, userId)

	logrus.Debug("UpdateQuery: %s", query)
	logrus.Debug("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
