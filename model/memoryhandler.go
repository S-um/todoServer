package model

import (
	"crypto/rand"
	"math/big"
	"time"
)

type memoryHandler struct {
	todoMap map[int]*Todo
}

func (m *memoryHandler) getRandID() int {
	numStr, err := rand.Int(rand.Reader, big.NewInt(10000))
	num := int(numStr.Int64())
	_, exist := m.todoMap[num]
	if err != nil || exist {
		return -1
	}
	return num
}

func (m *memoryHandler) GetTodos(sessionId string) []*Todo {
	list := []*Todo{}
	for _, value := range m.todoMap {
		list = append(list, value)
	}
	return list
}

func (m *memoryHandler) AddTodo(name string, sessionId string) *Todo {
	id := m.getRandID()
	newTodo := &Todo{id, name, false, time.Now(), ""}
	m.todoMap[id] = newTodo
	return newTodo
}

func (m *memoryHandler) RemoveTodo(id int) bool {
	_, exist := m.todoMap[id]
	if !exist {
		return false
	}
	delete(m.todoMap, id)
	return true
}

func (m *memoryHandler) UpdateTodo(id int) interface{} {
	updateTodo, exist := m.todoMap[id]
	if !exist {
		return Success{false}
	}
	updateTodo.Completed = !updateTodo.Completed
	return updateTodo
}

func (m *memoryHandler) Close() {
}

func newMemoryHandler() DBHandler {
	m := new(memoryHandler)
	m.todoMap = make(map[int]*Todo)

	return m
}
