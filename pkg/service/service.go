package service

import (
	todo_app "github.com/sinoturaev/todo-app"
	"github.com/sinoturaev/todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo_app.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list todo_app.TodoList) (int, error)
	GetAll(userId int) ([]todo_app.TodoList, error)
	GetById(userId, listID int) (todo_app.TodoList, error)
	Delete(userId, listID int) error
	Update(userId, listId int, input todo_app.UpdateListInput) error


}

type TodoItem interface {
	Create(userId, listId int, item todo_app.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo_app.TodoItem, error)
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList: NewTodoListService(repos.TodoList),
	}
}
