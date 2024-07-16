package data

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
	"github.com/task-manager/app"
	"github.com/task-manager/models"
)

func GetTasks(app *app.App) (tasks []models.Task, err error) {
	var rows *sql.Rows
	tasks = []models.Task{}
	rows, err = app.PostgresDB().Conn.Query(`SELECT id,
	 title,
	 description,
	 ep(create_time),
	 ep(deadline),
	 ep(update_time) from task`)

	if err != nil {
		log.Errorf("Couldn't query tasks: %v", err)
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID,
			&task.Title,
			&task.Description,
			&task.CreateTime,
			&task.UpdateTime,
			&task.Deadline,
		)
		if err != nil {
			log.Errorf("couldn't scan rows:%v", err)
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func AddTask(app *app.App,
	taskTobeAdded models.Task) (task models.Task,
	err error) {

	err = app.PostgresDB().Conn.QueryRow(`INSERT INTO task ("title","description","deadline") values($1,$2,ts($3)) returning id, ep(create_time), ep(update_time)`,
		taskTobeAdded.Title,
		taskTobeAdded.Description,
		taskTobeAdded.Deadline,
	).Scan(&task.ID, &task.CreateTime, &task.UpdateTime)

	task.Title = taskTobeAdded.Title
	task.Description = taskTobeAdded.Description
	task.Deadline = taskTobeAdded.Deadline

	if err != nil {
		log.Errorf("Couldn't insert task: %v", err)
		return task, err
	}

	return task, nil
}

func DeleteTask(app *app.App,
	id string) error {

	result, err := app.PostgresDB().Conn.Exec(`delete from task where id = $1`, id)

	if err != nil {
		log.Errorf("Couldn't delete task: %v", err)
		return err
	}
	if x, _ := result.RowsAffected(); x == 0 {
		log.Errorf("not found")
		return sql.ErrNoRows
	}

	return nil
}

func GetTaskByID(app *app.App,
	id string) (task models.Task, err error) {

	task = models.Task{}
	err = app.PostgresDB().Conn.QueryRow(`SELECT id,
	 title,
	 description,
	 ep(create_time),
	 ep(deadline),
	 ep(update_time) from task where id = $1`, id).Scan(&task.ID,
		&task.Title,
		&task.Description,
		&task.CreateTime,
		&task.UpdateTime,
		&task.Deadline)

	if err != nil {
		log.Errorf("Couldn't query tasks: %v", err)
		return task, err
	}
	if task.ID == 0 {
		log.Errorf("task doesn't exist")
		return task, sql.ErrNoRows

	}

	return task, nil
}

func EditTask(app *app.App,
	task models.Task) error {

	_, err := app.PostgresDB().Conn.Exec(`update task set title = coalesce($1, title),
	description = coalesce($2, description),
	deadline= coalesce(ts($3), deadline),
	where id = $1`,
		task.ID,
		task.Title,
		task.Description)

	if err != nil {
		log.Errorf("Couldn't patch task: %v", err)
		return err
	}

	return nil
}
