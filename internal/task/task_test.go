package task

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestResetTask(t *testing.T) {
	task := New(echo.New(), "title", "/testing", &testEx{fetch: []inputType{}})
	task.inProgress.Store(true)
	task.status = "status"
	task.Reset()

	assert.Equal(t, "title", task.title)
	assert.Equal(t, "Idle", task.status)
	assert.Equal(t, false, task.InProgress())
}

func TestExecuteTask(t *testing.T) {
	testExecutor := &testEx{fetch: []inputType{{fieldOne: "one", fieldTwo: "two"}}}
	task := New(echo.New(), "title", "/testing", testExecutor)
	task.Execute()
	assert.ElementsMatch(t, []inputType{{fieldOne: "one", fieldTwo: "two"}}, testExecutor.written)
}

type testEx struct {
	fetch   []inputType
	written []inputType
}

type inputType struct {
	fieldOne string
	fieldTwo string
}

func (t *testEx) Fetch() ([]inputType, error) {
	return t.fetch, nil
}

func (t *testEx) Save(i []inputType) error {
	t.written = i
	return nil
}
