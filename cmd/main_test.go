package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDbPanicsWhenQueryParamsPresent(t *testing.T) {
	assert.Panics(t, func() {
		initDb(EnvConfig{
			MySqlConnectionString: "root:password@location?parseTime=true",
		})
	})
}