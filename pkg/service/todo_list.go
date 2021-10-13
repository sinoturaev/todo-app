package service

import (
	todo_app "github.com/sinoturaev/todo-app"
	"github.com/sinoturaev/todo-app/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list todo_app.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) 	GetAll(userId int) ([]todo_app.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) 	GetById(userId, listId int) (todo_app.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService)  Delete(userId, listID int) error {
	return s.repo.Delete(userId, listID)
}

func (s *TodoListService) Update(userId, listId int, input todo_app.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}