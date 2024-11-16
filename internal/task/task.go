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
	Execute() 
	GetHandler(echo.Context) error
	PostHandler(echo.Context) error
	Title() string
	Status() string
	InProgress() bool
	UrlPath() string
}

type TaskExecutor[T any] struct {
	title      string
	status     string
	inProgress atomic.Bool
	urlPath    string
	ex         Executor[T]
}

func New[T any](title, urlPath string, ex Executor[T]) Task {
	return &TaskExecutor[T]{
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
		log.Errorf("Error fetching for %s: %s",t.Title(), err)
		t.status = "Error fetching from source"
		return
	}

	t.status = "Saving rows to DB"
	if n, err := t.ex.Save(data); err != nil {
		log.Errorf("Error saving to DB for %s: %s", t.Title(), err.Error())
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
	t.status = "Idle"
	t.inProgress.Store(false)
}

func (t *TaskExecutor[T]) Render(c echo.Context) error {
	return web.RenderOk(c, TaskRow(t))
}
