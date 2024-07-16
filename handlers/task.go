package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/task-manager/app"
	"github.com/task-manager/data"
	"github.com/task-manager/models"
)

// GetTasks godoc
// @Summary Get all tasks
// @Description Get all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Router /tasks [get]
func GetTasks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tasks, err := data.GetTasks(app)
		if err != nil {
			log.Errorf("couldn't get tasks from database: %s",
				err.Error())
			http.Error(w,
				"couldn't get tasks",
				http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(tasks)

		if err != nil {
			log.Errorf("couldn't marshal response: %s",
				err.Error())
			http.Error(w,
				"couldn't marshal response",
				http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

}

// AddTask godoc
// @Summary Create a task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task"
// @Success 200 {object} models.Task
// @Router /task [post]
func AddTask(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("couldn't read request body: %v", err)
			http.Error(w,
				"couldn't read body",
				http.StatusInternalServerError)
			return
		}
		err = r.Body.Close()
		if err != nil {
			log.Errorf("couldn't close body: %v", err)
			http.Error(w,
				"couldn't close body",
				http.StatusInternalServerError)
			return
		}
		var task models.Task
		err = json.Unmarshal(body, &task)
		if err != nil {
			log.Errorf("couldn't unmarshal payload: %v", err)
			http.Error(w,
				"couldn't unmarshal payload",
				http.StatusInternalServerError)
			return

		}
		if task.Title == nil || task.Description == nil {
			http.Error(w,
				"missing parameters",
				http.StatusBadRequest)
			return
		}

		AddedTask, err := data.AddTask(app, task)
		if err != nil {
			log.Errorf("couldn't add task to database: %s",
				err.Error())
			http.Error(w,
				"couldn't add task",
				http.StatusInternalServerError)
			return
		}

		err = app.RedisDB().SetJson(strconv.Itoa(AddedTask.ID), AddedTask)
		if err != nil {
			log.Warnf("couldn't insert task into redis: %v", err)
		}
		log.Info("task was added successfully")
		response, err := json.Marshal(AddedTask)

		if err != nil {
			log.Errorf("couldn't marshal response: %s",
				err.Error())
			http.Error(w,
				"couldn't marshal response",
				http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update task details
// @Tags tasks
// @Accept json
// @Param task body models.Task true "Task"
// @Success 200
// @Failure 404
// @Router /task [patch]
func EditTask(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("couldn't read request body: %v", err)
			http.Error(w,
				"couldn't read body",
				http.StatusInternalServerError)
			return
		}
		err = r.Body.Close()
		if err != nil {
			log.Errorf("couldn't close body: %v", err)
			http.Error(w,
				"couldn't close body",
				http.StatusInternalServerError)
			return
		}
		var task models.Task
		err = json.Unmarshal(body, &task)
		if err != nil {
			log.Errorf("couldn't unmarshal payload: %v", err)
			http.Error(w,
				"couldn't unmarshal payload",
				http.StatusInternalServerError)
			return

		}
		if task.ID == 0 {
			log.Error("task id is missing")
			http.Error(w,
				"task id is missing",
				http.StatusBadRequest)
			return
		}

		err = data.EditTask(app, task)
		if err != nil {
			log.Errorf("couldn't edit users in database: %s",
				err.Error())
			http.Error(w,
				"couldn't edit task details",
				http.StatusInternalServerError)
			return
		}

		//remove old value from redis
		err = app.RedisDB().Del(strconv.Itoa(task.ID))
		if err != nil {
			log.Warnf("couldn't delete task from redis: %v", err)
		}

		log.Info("task was edited successfully")
		w.WriteHeader(200)
	}

}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200
// @Failure 404
// @Router /task/{id} [delete]
func DeleteTask(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		err := data.DeleteTask(app, id)
		if err != nil {
			log.Errorf("couldn't delete task %s",
				err.Error())
			if err.Error() == sql.ErrNoRows.Error() {
				http.Error(w,
					"couldn't delete the task",
					http.StatusNotFound)

			} else {
				http.Error(w,
					"couldn't delete the task",
					http.StatusInternalServerError)

			}
			return

		}
		log.Info("task was deleted successfully")
		err = app.RedisDB().Del(id)
		if err != nil {
			log.Warnf("couldn't delete task from redis: %v", err)
		}
		w.WriteHeader(200)
	}
}

// GetTaskByID godoc
// @Summary Get a task
// @Description Get a task by its ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404
// @Router /task/{id} [get]
func GetTaskByID(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]
		task := models.Task{}
		//get task from redis if exists
		var GotValueFromRedis bool

		result, err := app.RedisDB().Get(id)
		if err == nil {
			log.Warnf("couldn't get data from  redis:%v", err)
		} else if result == "" {
			log.Warn("data doesn't exist in redis")
		} else {
			err := json.Unmarshal([]byte(result), &task)
			if err != nil {
				log.Errorf("couldn't unmarshal value :%v", err)
			}
			GotValueFromRedis = true
		}

		if !GotValueFromRedis {
			task, err = data.GetTaskByID(app, id)
			if err != nil {
				log.Errorf("couldn't get task from database: %s",
					err.Error())
				if err.Error() == sql.ErrNoRows.Error() {
					http.Error(w,
						"task not found",
						http.StatusNotFound)
					return
				}
				http.Error(w,
					"couldn't get task",
					http.StatusInternalServerError)
				return
			}
			//set value in redis for next time
			err = app.RedisDB().SetJson(strconv.Itoa(task.ID), task)
			if err != nil {
				log.Warnf("couldn't insert task into redis: %v", err)
			}
		}

		response, err := json.Marshal(task)

		if err != nil {
			log.Errorf("couldn't marshal response: %s",
				err.Error())
			http.Error(w,
				"couldn't marshal response",
				http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

}
