package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"stock-management/internal/models"
	"stock-management/internal/util/must"
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
		status:  "",
	}
}

func (t *TaskExecutor[T]) Execute() {
	if !t.inProgress.CompareAndSwap(false, true) {
		return
	}

	t.status = "fetching data"
	go t.fetchAndSave()
}

var loc = must.MustLoadLocation("America/New_York")
func (t *TaskExecutor[T]) fetchAndSave() {
	history := &models.SaveTaskHistoryParams{
		TaskName:   t.Title(),
		StartTime:  time.Now().In(loc),
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

	t.status = "saving rows to DB"
	n, err := t.ex.Save(data)
	if err != nil {
		t.status = fmt.Sprintf("error saving to database: %s", err.Error())
		log.Errorf("During %s - %s", t.Title(), t.status)
		return
	}

	history.TaskStatus = models.TaskHistoryTaskStatusSucceeded
	t.status = fmt.Sprintf("saved %d rows to database", n)
	log.Infof("%s for task %s", t.status, t.Title())
}

func (t *TaskExecutor[T]) saveHistory(history *models.SaveTaskHistoryParams) {
	history.Details = t.status
	history.EndTime = time.Now().In(loc)
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
	if t.status != "" {
		return t.status
	}
	if history, err := t.q.GetLatestTaskHistory(context.Background(), t.Title()); err != nil && !errors.Is(err, sql.ErrNoRows){
		log.Errorf("During %s - error getting latest task history: %s", t.Title(), err.Error())
	} else if err == nil {
		startStr := history.StartTime.In(loc).Format("Jan 2, 2006 3:04 PM")
		return fmt.Sprintf("last executed %s", startStr)
	}
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
	t.status = ""
	t.inProgress.Store(false)
}

func (t *TaskExecutor[T]) Render(c echo.Context) error {
	return web.RenderOk(c, TaskRow(t))
}
