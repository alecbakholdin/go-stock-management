package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResetTask(t *testing.T) {
	task := New("title", "/testing", &testEx{fetch: []inputType{}}).(*TaskExecutor[inputType])
	task.inProgress.Store(true)
	task.status = "status"
	task.Reset()

	assert.Equal(t, "title", task.title)
	assert.Equal(t, "Idle", task.status)
	assert.Equal(t, false, task.InProgress())
}

func TestExecuteTask(t *testing.T) {
	testExecutor := &testEx{fetch: []inputType{{fieldOne: "one", fieldTwo: "two"}}}
	task := New("title", "/testing", testExecutor)
	task.Execute()
	for task.InProgress() {
		time.Sleep(time.Millisecond)
	}
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

func (t *testEx) Save(i []inputType) (int, error) {
	t.written = i
	return 0, nil
}
