package must

import (
	"fmt"
	"time"
)

func MustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("error loading location %s: %s", name, err.Error()))
	}
	return loc
}
