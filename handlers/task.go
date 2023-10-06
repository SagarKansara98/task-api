package handlers

import (
	"net/http"
	"strconv"
	"task-api/buisness/task"

	"github.com/labstack/echo/v4"
)

var resource = "Task"

// taskHandlers handle request for task
type taskHandlers struct {
	task task.Task
}

func (t taskHandlers) create(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	newTask := task.Info{}
	err := c.Bind(&newTask)
	if err != nil {
		return err
	}

	err = c.Validate(newTask)
	if err != nil {
		return badRequest(c, err.Error())
	}
	newTask.UserID = userID

	task, err := t.task.Create(c.Request().Context(), newTask)
	if err != nil {
		return Error(c, err, resource)
	}

	return c.JSON(http.StatusOK, task)
}

func (t taskHandlers) query(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	task, err := t.task.Query(c.Request().Context(), userID)
	if err != nil {
		return Error(c, err, resource)
	}

	return c.JSON(http.StatusOK, task)
}

func (t taskHandlers) queryByID(c echo.Context) error {
	taskIDStr := c.Param("task_id")
	if len(taskIDStr) == 0 {
		return badRequest(c, "task id is not provided")
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return badRequest(c, "invalid task id")
	}

	userID := getUserID(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	task, err := t.task.QueryByID(c.Request().Context(), taskID, userID)
	if err != nil {
		return Error(c, err, resource)
	}

	return c.JSON(http.StatusOK, task)
}

func (t taskHandlers) update(c echo.Context) error {
	updateTask := task.Info{}
	err := c.Bind(&updateTask)
	if err != nil {
		return err
	}

	err = c.Validate(updateTask)
	if err != nil {
		return badRequest(c, err.Error())
	}

	updateTask.UserID = getUserID(c)
	if updateTask.UserID == 0 {
		return echo.ErrUnauthorized
	}

	err = t.task.Update(c.Request().Context(), updateTask)
	if err != nil {
		return Error(c, err, resource)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task updated successfully"})
}

func (t taskHandlers) delete(c echo.Context) error {
	taskIDStr := c.Param("task_id")
	if len(taskIDStr) == 0 {
		return badRequest(c, "please provide task id")
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return badRequest(c, "invalid task id")
	}

	userID := getUserID(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	err = t.task.Delete(c.Request().Context(), taskID, userID)
	if err != nil {
		return Error(c, err, resource)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task deleted successfully"})
}

func (t taskHandlers) updateStatus(c echo.Context) error {
	ids := []int{}
	err := c.Bind(&ids)
	if err != nil {
		return err
	}

	userID := getUserID(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	msgChan := make(chan string)
	go t.task.UpdateStatus(c.Request().Context(), ids, msgChan, userID)
	messages := make([]string, 0)
	for s := range msgChan {
		messages = append(messages, s)
	}

	if len(messages) != 0 {
		return c.JSON(http.StatusInternalServerError, messages)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "all task marked as done"})
}
