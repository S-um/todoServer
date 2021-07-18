package model

import (
	"time"
)

type Success struct {
	Success bool `json:"success"`
}

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	SessionID string    `json:"session_id"`
}

type DBHandler interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(name string, sessionId string) *Todo
	RemoveTodo(id int) bool
	UpdateTodo(id int) interface{}
	Close()
}

func NewDBHandler(dbpath string) DBHandler {
	return newSqliteHandler(dbpath)
}
