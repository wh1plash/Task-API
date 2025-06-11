package api

import (
	"fmt"
	"net/http"
	"task/store"
	"task/types"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TaskHandler struct {
	TaskStore store.TaskStore
}

func NewTaskHandler(taskStore store.TaskStore) *TaskHandler {
	return &TaskHandler{
		TaskStore: taskStore,
	}
}

func (h TaskHandler) HandleGetTasks(c *gin.Context) {
	res, err := h.TaskStore.GetTasks(c.Request.Context())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no tasks found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": "All tasks",
		"task":     res,
	})
}

func (h TaskHandler) HandlePutTask(c *gin.Context) {
	id := c.Param("id")

	var params *types.TaskParams
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read JSON body",
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		errs := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range errs {
			errors[e.Field()] = fmt.Sprintf("failed on '%s' tag", e.Tag())
		}
		Err := NewValidationError(errors)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": Err,
		})
		return
	}

	res, err := h.TaskStore.UpdateTask(c.Request.Context(), id, params)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("task with id {%s} not found", c.Param("id")),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "task updated",
		"task":     res,
	})

}

func (h TaskHandler) HandleGetTask(c *gin.Context) {
	res, err := h.TaskStore.GetTask(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("task with id {%s} not found", c.Param("id")),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "successful",
		"task":     res,
	})

}

func (h TaskHandler) HandleDeleteTask(c *gin.Context) {
	id := c.Param("id")

	res, err := h.TaskStore.DeleteTask(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("task with id {%s} not found", c.Param("id")),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "task deleted",
		"task":     res,
	})
}

func (h TaskHandler) HandlePostTask(c *gin.Context) {
	var params types.TaskParams
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read JSON body",
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		errs := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range errs {
			errors[e.Field()] = fmt.Sprintf("failed on '%s' tag", e.Tag())
		}
		Err := NewValidationError(errors)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": Err,
		})
		return
	}

	task, err := types.NewTaskFromParams(params)
	if err != nil {
		c.Error(err)
		return
	}

	insTask, err := h.TaskStore.InsertTask(c.Request.Context(), task)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "created successfully",
		"task":     insTask,
	})

}

func NewValidationError(errors map[string]string) ValidationError {
	return ValidationError{
		Status: http.StatusUnprocessableEntity,
		Errors: errors,
	}
}

type ValidationError struct {
	Status int               `json:"status"`
	Errors map[string]string `json:"errors"`
}

func (e ValidationError) Error() string {
	return "validation failed"
}
