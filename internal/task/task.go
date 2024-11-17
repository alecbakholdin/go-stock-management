package task

import (
	"context"
	"fmt"
	"stock-management/internal/models"
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

type TaskHistoryTable interface {
	SaveTaskHistory(context.Context, models.SaveTaskHistoryParams) error
	GetLatestTaskHistory(ctx context.Context, task_name string) (models.TaskHistory, error)
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
	q          TaskHistoryTable
	title      string
	status     string
	inProgress atomic.Bool
	urlPath    string
	ex         Executor[T]
}

func New[T any](q TaskHistoryTable, title, urlPath string, ex Executor[T]) Task {
	return &TaskExecutor[T]{
		q:       q,
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
	history := &models.SaveTaskHistoryParams{
		TaskName:   t.Title(),
		StartTime:  time.Now(),
		TaskStatus: models.TaskHistoryTaskStatusFailed,
	}
	defer t.ResetDelay()
	defer t.saveHistory(history)

	data, err := t.ex.Fetch()
	if err != nil {
		t.status = fmt.Sprintf("error fetching from source: %s", err.Error())
		log.Errorf("During %s - %s", t.Title(), t.status)
		return
	}

	t.status = "Saving rows to DB"
	n, err := t.ex.Save(data)
	if err != nil {
		t.status = fmt.Sprintf("error saving to database: %s", err.Error())
		log.Errorf("During %s - %s", t.Title(), t.status)
		return
	}

	history.TaskStatus = models.TaskHistoryTaskStatusSucceeded
	t.status = fmt.Sprintf("Saved %d rows to database", n)
	log.Infof("%s for task %s", t.status, t.Title())
}

func (t *TaskExecutor[T]) saveHistory(history *models.SaveTaskHistoryParams) {
	history.Details = t.status
	history.EndTime = time.Now()
	if err := t.q.SaveTaskHistory(context.Background(), *history); err != nil {
		log.Errorf("Failed to save history for %s: %s", t.Title(), err.Error())
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
