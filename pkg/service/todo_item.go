package service

import (
	"github.com/vladislavtinishov/todo-app"
	"github.com/vladislavtinishov/todo-app/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	todoList repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, todoList repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, todoList: todoList}
}

func (s *TodoItemService) Create(userId, listId int, input todo.TodoItem) (int, error) {
	_, err := s.todoList.GetById(userId, listId)

	if err != nil {
		return 0, err
	}
	return s.repo.Create(listId, input)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}

func (s *TodoItemService) SetDoneStatus(userId, itemId, status int) error {
	return s.repo.SetDoneStatus(userId, itemId, status)
}
