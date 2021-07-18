package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/sir/todos/model"
	"github.com/stretchr/testify/assert"
)

type sessionID struct {
	SessionID string `json:"session_id"`
}

type createData struct {
	Name      string `json:"name"`
	SessionID string `json:"session_id"`
}

func TestTodos(t *testing.T) {
	os.Remove("./test.db")
	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	// test1
	setTest := new(createData)
	testStr1 := "Test todo"
	setTest.Name = testStr1
	setTest.SessionID = "1"
	data, _ := json.Marshal(setTest)
	resp, err := http.Post(ts.URL+"/todos", "application/json", strings.NewReader(string(data)))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	todo := new(model.Todo)
	err = json.NewDecoder(resp.Body).Decode(todo)
	id1 := todo.ID
	assert.NoError(err)
	assert.Equal(todo.Name, testStr1)

	// test2
	testStr2 := "Test todo2"
	setTest.Name = testStr2
	setTest.SessionID = "1"
	data, _ = json.Marshal(setTest)
	resp, err = http.Post(ts.URL+"/todos", "application/json", strings.NewReader(string(data)))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(todo)
	id2 := todo.ID
	assert.NoError(err)
	assert.Equal(todo.Name, testStr2)

	// get test
	getTest := new(sessionID)
	getTest.SessionID = "1"
	data, _ = json.Marshal(getTest)
	req, _ := http.NewRequest("MYGET", ts.URL+"/todos", strings.NewReader(string(data)))
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal(testStr1, t.Name)
		} else if t.ID == id2 {
			assert.Equal(testStr2, t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

	// change status test
	data, _ = json.Marshal(*todo)
	req, _ = http.NewRequest("PUT", ts.URL+"/todos", strings.NewReader(string(data)))
	req.PostForm = url.Values{
		"id": []string{strconv.Itoa(id2)},
	}
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.NotContains(string(body), "success")
	assert.Contains(string(body), "id")
	assert.Contains(string(body), "name")
	assert.Contains(string(body), "completed")
	assert.Contains(string(body), "true")
	assert.Contains(string(body), "created_at")
	assert.Equal(http.StatusOK, resp.StatusCode)

	// delete test
	req, _ = http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id2), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ = json.Marshal(getTest)
	req, _ = http.NewRequest("MYGET", ts.URL+"/todos", strings.NewReader(string(data)))
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)
	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal(testStr1, t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1"))
		}
	}
}
