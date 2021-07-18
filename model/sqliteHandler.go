package model

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteHandler struct {
	db *sql.DB
}

func newSqliteHandler(dbpath string) DBHandler {
	database, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			sessionId STRING,
			name      TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		);`)
	statement.Exec()
	return &sqliteHandler{db: database}
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}

func (s *sqliteHandler) GetTodos(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=?", sessionId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}
	return todos
}

func (s *sqliteHandler) AddTodo(name string, sessionId string) *Todo {
	statement, err := s.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	if err != nil {
		panic(err)
	}
	result, err := statement.Exec(sessionId, name, false)
	if err != nil {
		panic(err)
	}
	id, _ := result.LastInsertId()
	todo := &Todo{
		ID:        int(id),
		Name:      name,
		Completed: false,
		CreatedAt: time.Now(),
	}
	return todo
}

func (s *sqliteHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
	if err != nil {
		log.Print(err)
		return false
	}
	rst, err := stmt.Exec(id)
	if err != nil {
		log.Print(err)
		return false
	}
	cnt, _ := rst.RowsAffected()
	if cnt != 1 {
		log.Println("wrong id deleted cnt : ", cnt)
	}
	return cnt > 0
}

func (s *sqliteHandler) UpdateTodo(id int) interface{} {
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE id=?", id)
	if err != nil {
		log.Print(err)
		return Success{false}
	}
	todo := Todo{}
	if rows.Next() {
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
	}
	todo.Completed = !todo.Completed
	rows.Close()
	stmt, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
	if err != nil {
		log.Print(err)
		return Success{false}
	}
	rst, err := stmt.Exec(todo.Completed, id)
	if err != nil {
		log.Print(err)
		return Success{false}
	}
	cnt, _ := rst.RowsAffected()
	if cnt != 1 {
		log.Println("wrong id exist cnt : ", cnt)
		return Success{false}
	}
	return todo
}
