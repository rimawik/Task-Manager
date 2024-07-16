package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/task-manager/app"
	_ "github.com/task-manager/docs"
	"github.com/task-manager/handlers"
)

func NewRouter(app *app.App) http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/v1/tasks", handlers.GetTasks(app)).Methods("GET")
	r.HandleFunc("/v1/task/{id}", handlers.GetTaskByID(app)).Methods("GET")
	r.HandleFunc("/v1/task", handlers.AddTask(app)).Methods("POST")
	r.HandleFunc("/v1/task/{id}", handlers.DeleteTask(app)).Methods("DELETE")
	r.HandleFunc("/v1/task", handlers.EditTask(app)).Methods("PATCH")
	// Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return r
}
