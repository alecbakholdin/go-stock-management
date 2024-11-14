package task

import (
	"stock-management/internal/web"
	"time"

	"github.com/labstack/echo/v4"
)

type Executor[T1, T2 any] interface {
	Fetch() ([]T1, error)
	Map(t1 T1) (T2)
	Save([]T2) error
}

type Task interface {
}

type TaskExecutor[T1, T2 any] struct {
	e          *echo.Echo
	title      string
	status     string
	Progress   int
	inProgress bool
	ex         Executor[T1, T2]
}

func New[T1, T2 any](e *echo.Echo, title string, ex Executor[T1, T2]) *TaskExecutor[T1, T2] {
	return &TaskExecutor[T1, T2]{
		e:     e,
		title: title,
		ex:    ex,
	}
}

func (t *TaskExecutor[T1, T2]) Execute() {
	defer t.ResetDelay()
	t.inProgress = true
	t.status = "Fetching data"
	data, err := t.ex.Fetch()
	if err != nil {
		t.Error("Error fetching for " + t.title)
		t.status = "Error fetching from source"
		return
	}

	t.status = "Mapping rows"
	output := make([]T2, len(data))
	for i, row := range data {
		output[i] = t.ex.Map(row)
	}

	t.status = "Saving rows to DB"
	if err := t.ex.Save(output); err != nil {
		t.Error("Error saving to DB for " + t.title + ": " + err.Error())
		t.status = "Error saving to DB"
		return
	}
	t.status = "Finished"
	t.inProgress = false
}

func (t *TaskExecutor[T1, T2]) Info(i ...interface{}) {
	t.e.Logger.Info(i...)
}

func (t *TaskExecutor[T1, T2]) Error(i ...interface{}) {
	t.e.Logger.Error(i...)
}

func (t *TaskExecutor[T1, T2]) Title() string {
	return t.title
}

func (t *TaskExecutor[T1, T2]) InProgress() bool {
	return t.inProgress
}

func (t *TaskExecutor[T1, T2]) Status() string {
	return t.status
}

func (t *TaskExecutor[T1, T2]) ResetDelay() {
	go func() {
		time.Sleep(time.Second * 2)
		t.Reset()
	}()
}

func (t *TaskExecutor[T1, T2]) Reset() {
	t.inProgress = false
	t.Progress = 0
	t.status = ""
	t.inProgress = false
}

func (t *TaskExecutor[T1, T2]) Render(c echo.Context) error {
	return web.RenderOk(c, TaskRow(t))
}
