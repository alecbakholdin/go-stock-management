package task

import (
	"stock-management/internal/web"
	"time"

	"github.com/labstack/echo/v4"
)

type Executor[T any] interface {
	Fetch() ([]T, error)
	Save([]T) error
}

type Task interface {
}

type TaskExecutor[T any] struct {
	e          *echo.Echo
	title      string
	status     string
	Progress   int
	inProgress bool
	ex         Executor[T]
}

func New[T any](e *echo.Echo, title string, ex Executor[T]) *TaskExecutor[T] {
	return &TaskExecutor[T]{
		e:     e,
		title: title,
		ex:    ex,
	}
}

func (t *TaskExecutor[T]) Execute() {
	defer t.ResetDelay()
	t.inProgress = true
	t.status = "Fetching data"
	data, err := t.ex.Fetch()
	if err != nil {
		t.Error("Error fetching for " + t.title)
		t.status = "Error fetching from source"
		return
	}

	t.status = "Saving rows to DB"
	if err := t.ex.Save(data); err != nil {
		t.Error("Error saving to DB for " + t.title + ": " + err.Error())
		t.status = "Error saving to DB"
		return
	}
	t.status = "Finished"
	t.inProgress = false
}

func (t *TaskExecutor[T]) Info(i ...interface{}) {
	t.e.Logger.Info(i...)
}

func (t *TaskExecutor[T]) Error(i ...interface{}) {
	t.e.Logger.Error(i...)
}

func (t *TaskExecutor[T]) Title() string {
	return t.title
}

func (t *TaskExecutor[T]) InProgress() bool {
	return t.inProgress
}

func (t *TaskExecutor[T]) Status() string {
	return t.status
}

func (t *TaskExecutor[T]) ResetDelay() {
	go func() {
		time.Sleep(time.Second * 2)
		t.Reset()
	}()
}

func (t *TaskExecutor[T]) Reset() {
	t.inProgress = false
	t.Progress = 0
	t.status = ""
	t.inProgress = false
}

func (t *TaskExecutor[T]) Render(c echo.Context) error {
	return web.RenderOk(c, TaskRow(t))
}
