package task

import (
	"fmt"
	"stock-management/internal/web"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Executor[T any] interface {
	Fetch() ([]T, error)
	Save([]T) (int, error)
}

type Task interface {
}

type TaskExecutor[T any] struct {
	e          *echo.Echo
	title      string
	status     string
	Progress   int
	inProgress atomic.Bool
	urlPath    string
	ex         Executor[T]
}

func New[T any](e *echo.Echo, title, urlPath string, ex Executor[T]) *TaskExecutor[T] {
	return &TaskExecutor[T]{
		e:       e,
		title:   title,
		urlPath: urlPath,
		ex:      ex,
		status:  "Idle",
	}
}

func (t *TaskExecutor[T]) Execute() {
	if !t.inProgress.CompareAndSwap(false, true) {
		return
	}

	t.status = "Fetching data"
	go t.fetchAndSave()
}
func (t *TaskExecutor[T]) fetchAndSave() {
	defer t.ResetDelay()
	data, err := t.ex.Fetch()
	if err != nil {
		t.Error("Error fetching for "+t.title, ": ", err)
		t.status = "Error fetching from source"
		return
	}

	t.status = "Saving rows to DB"
	if n, err := t.ex.Save(data); err != nil {
		t.Error("Error saving to DB for " + t.title + ": " + err.Error())
		t.status = "Error saving to DB"
		return
	} else {
		t.status = fmt.Sprintf("Saved %d rows to database", n)
		log.Infof("Saved %d rows to database for executor %s", n, t.Title())
	}
}

func (t *TaskExecutor[T]) GetHandler(c echo.Context) error {
	return web.RenderOk(c, TaskRow(t))
}

func (t *TaskExecutor[T]) PostHandler(c echo.Context) error {
	t.Execute()
	return web.RenderOk(c, TaskRow(t))
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
	return t.inProgress.Load()
}

func (t *TaskExecutor[T]) Status() string {
	return t.status
}

func (t *TaskExecutor[T]) UrlPath() string {
	return t.urlPath
}

func (t *TaskExecutor[T]) ResetDelay() {
	go func() {
		time.Sleep(time.Second * 5)
		t.Reset()
	}()
}

func (t *TaskExecutor[T]) Reset() {
	t.Progress = 0
	t.status = "Idle"
	t.inProgress.Store(false)
}

func (t *TaskExecutor[T]) Render(c echo.Context) error {
	return web.RenderOk(c, TaskRow(t))
}
