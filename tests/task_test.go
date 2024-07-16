package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/task-manager/app"
	"github.com/task-manager/cache"
	"github.com/task-manager/config"
	"github.com/task-manager/db"
	"github.com/task-manager/handlers"
	"github.com/task-manager/models"
)

func TestGetTasks(t *testing.T) {

	r := mux.NewRouter()

	cfg, err := config.LoadTestConfig()

	if err != nil {
		logrus.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}
	redis := &cache.Rdb{}

	testApp := app.BuildApp(cfg, postgresDB, redis)

	r.HandleFunc("/tasks", handlers.GetTasks(testApp)).Methods("GET")

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var tasks []models.Task
	err = json.NewDecoder(rr.Body).Decode(&tasks)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(tasks), 2)
}

func TestAddTask(t *testing.T) {
	// Create a new router
	r := mux.NewRouter()

	cfg, err := config.LoadTestConfig()

	if err != nil {
		logrus.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}
	redis, err := cache.ConnectToRedis(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}

	testApp := app.BuildApp(cfg, postgresDB, redis)

	r.HandleFunc("/task", handlers.AddTask(testApp)).Methods("POST")

	taskToBeAdded := map[string]interface{}{
		"title":       "Task 1",
		"description": "Description for Task 1",
	}
	taskJSON, _ := json.Marshal(taskToBeAdded)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/task", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var task models.Task
	err = json.NewDecoder(rr.Body).Decode(&task)
	if err != nil {
		t.Fatal(err)
	}
	// Check the response body

	assert.Equal(t, *task.Title,
		taskToBeAdded["title"])
}

func TestDeleteTask(t *testing.T) {

	//prepare db and configs
	r := mux.NewRouter()

	cfg, err := config.LoadTestConfig()

	if err != nil {
		logrus.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}
	redis, err := cache.ConnectToRedis(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't connect to redis: %v", err)
	}

	testApp := app.BuildApp(cfg, postgresDB, redis)

	//test case 1: task not found
	r.HandleFunc("/task/{id}", handlers.DeleteTask(testApp)).Methods("DELETE")

	req, err := http.NewRequest("DELETE", "/task/1000", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	//test case 2: task successfully deleted
	r.HandleFunc("/task/{id}", handlers.DeleteTask(testApp)).Methods("DELETE")

	req, err = http.NewRequest("DELETE", "/task/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

}

func TestGetTask(t *testing.T) {

	//prepare db and configs
	r := mux.NewRouter()

	cfg, err := config.LoadTestConfig()

	if err != nil {
		logrus.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}
	redis, err := cache.ConnectToRedis(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't connect to redis: %v", err)
	}

	testApp := app.BuildApp(cfg, postgresDB, redis)

	//test case 1: task not found
	r.HandleFunc("/task/{id}", handlers.GetTaskByID(testApp)).Methods("GET")

	req, err := http.NewRequest("GET", "/task/1000", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	//test case 2: get task successfully
	r.HandleFunc("/task/{id}", handlers.GetTaskByID(testApp)).Methods("GET")

	req, err = http.NewRequest("GET", "/task/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var task models.Task
	err = json.NewDecoder(rr.Body).Decode(&task)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, task.ID, 2)

}
