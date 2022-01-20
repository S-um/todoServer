package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sir/todos/model"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var rd *render.Render = render.New()

type Success struct {
	Success bool `json:"success"`
}

type AppHandler struct {
	http.Handler
	db model.DBHandler
}

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	bearToken := r.Header.Get("Authorization")
	trArr := strings.Split(bearToken, " ")
	if len(trArr) != 2 {
		log.Println("token error")
		rd.JSON(w, http.StatusUnauthorized, []*model.Todo{})
		return
	}
	log.Println("Token : [" + trArr[0] + "] [" + trArr[1] + "]")
	sessionid := new(model.Todo)
	if err := json.NewDecoder(r.Body).Decode(sessionid); err != nil {
		fmt.Print(err)
	}
	list := a.db.GetTodos(sessionid.SessionID)
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodosHandler(w http.ResponseWriter, r *http.Request) {
	reqtodo := new(model.Todo)
	if err := json.NewDecoder(r.Body).Decode(reqtodo); err != nil {
		fmt.Print(err)
	}
	name := reqtodo.Name
	newTodo := a.db.AddTodo(name, reqtodo.SessionID)
	rd.JSON(w, http.StatusCreated, newTodo)
}

func (a *AppHandler) deleteTodosHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId, exist := vars["id"]
	if !exist {
		w.WriteHeader(http.StatusOK)
		rd.JSON(w, http.StatusOK, Success{false})
		return
	}
	id, _ := strconv.Atoi(strId)
	ok := a.db.RemoveTodo(id)
	rd.JSON(w, http.StatusOK, Success{ok})
}

func (a *AppHandler) updateTodosHandler(w http.ResponseWriter, r *http.Request) {
	todo := new(model.Todo)
	if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
		fmt.Println(err)
		rd.JSON(w, http.StatusOK, Success{false})
		return
	}
	jsonValue := a.db.UpdateTodo(todo.ID)
	rd.JSON(w, http.StatusOK, jsonValue)
}

func (a *AppHandler) Close() {
	a.db.Close()
}

func MakeHandler(dbpath string) *AppHandler {
	r := mux.NewRouter()
	a := &AppHandler{
		Handler: r,
		db:      model.NewDBHandler(dbpath),
	}

	r.HandleFunc("/", a.indexHandler)
	r.HandleFunc("/todos", a.getTodosHandler).Methods("MYGET")
	r.HandleFunc("/todos", a.addTodosHandler).Methods("POST")
	r.HandleFunc("/todos", a.updateTodosHandler).Methods("PUT")
	r.HandleFunc("/todos/{id:[0-9]+}", a.deleteTodosHandler).Methods("DELETE")

	return a
}
