package task

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"stock-management/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResetTask(t *testing.T) {
	task := New(&testHistoryTable{}, "title", "/testing", &testEx{fetch: []inputType{}}).(*TaskExecutor[inputType])
	task.inProgress.Store(true)
	task.status = "status"
	task.Reset()

	assert.Equal(t, "title", task.title)
	assert.Equal(t, "Idle", task.status)
	assert.Equal(t, false, task.InProgress())
}

func TestExecuteTask(t *testing.T) {
	t.Parallel()
	testHistoryTable := &testHistoryTable{}
	testExecutor := &testEx{fetch: []inputType{{fieldOne: "one", fieldTwo: "two"}}}
	task := New(testHistoryTable, "title", "/testing", testExecutor)
	task.Execute()
	for task.InProgress() {
		time.Sleep(time.Millisecond)
	}
	assert.ElementsMatch(t, []inputType{{fieldOne: "one", fieldTwo: "two"}}, testExecutor.written)

	if !assert.Equal(t, 1, len(testHistoryTable.savedHistories)) {
		t.FailNow()
	}
	history := testHistoryTable.savedHistories[0]
	assert.Greater(t, history.EndTime, history.StartTime)
	assert.Equal(t, "title", history.TaskName)
	assert.Equal(t, models.TaskHistoryTaskStatusSucceeded, history.TaskStatus)
	assert.Equal(t, "Saved 1 rows to database", history.Details)
}

func TestExecuteTaskFailsFetch(t *testing.T) {
	t.Parallel()
	testHistoryTable := &testHistoryTable{}
	task := New(testHistoryTable, "title", "/testing", &testEx{fetchErr: errors.New("err")})
	task.Execute()
	for task.InProgress() {
		time.Sleep(time.Millisecond)
	}
	if !assert.Equal(t, 1, len(testHistoryTable.savedHistories)) {
		t.FailNow()
	}
	history := testHistoryTable.savedHistories[0]
	assert.Greater(t, history.EndTime, history.StartTime)
	assert.Equal(t, "title", history.TaskName)
	assert.Equal(t, models.TaskHistoryTaskStatusFailed, history.TaskStatus)
	assert.Regexp(t, regexp.MustCompile("^error fetching from source"), history.Details)
}

func TestExecuteTaskFailsSave(t *testing.T) {
	t.Parallel()
	testHistoryTable := &testHistoryTable{}
	task := New(testHistoryTable, "title", "/testing", &testEx{saveErr: errors.New("err")})
	task.Execute()
	for task.InProgress() {
		time.Sleep(time.Millisecond)
	}
	if !assert.Equal(t, 1, len(testHistoryTable.savedHistories)) {
		t.FailNow()
	}
	history := testHistoryTable.savedHistories[0]
	assert.Greater(t, history.EndTime, history.StartTime)
	assert.Equal(t, "title", history.TaskName)
	assert.Equal(t, models.TaskHistoryTaskStatusFailed, history.TaskStatus)
	assert.Regexp(t, regexp.MustCompile("^error saving to database"), history.Details)
}

type testHistoryTable struct {
	savedHistories []models.SaveTaskHistoryParams
	latestHistory  *models.TaskHistory
}

func (t *testHistoryTable) SaveTaskHistory(_ context.Context, task models.SaveTaskHistoryParams) error {
	t.savedHistories = append(t.savedHistories, task)
	return nil
}

func (t *testHistoryTable) GetLatestTaskHistory(ctx context.Context, task_name string) (models.TaskHistory, error) {
	if t.latestHistory == nil {
		return models.TaskHistory{}, sql.ErrNoRows
	}
	return *t.latestHistory, nil
}

type testEx struct {
	fetchErr error
	saveErr error
	fetch   []inputType
	written []inputType
}

type inputType struct {
	fieldOne string
	fieldTwo string
}

func (t *testEx) Fetch() ([]inputType, error) {
	return t.fetch, t.fetchErr
}

func (t *testEx) Save(i []inputType) (int, error) {
	t.written = i
	return len(i), t.saveErr
}
