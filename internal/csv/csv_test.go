package csv

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderMap(t *testing.T) {
	type testStruct struct {
		FirstField string `csv:"First"`
		Second string
		Third string `csv:"Eight"`
	}

	headerMap := GetHeaderMap([]string{"Second", "Third", "Eight", "First"}, reflect.TypeOf(testStruct{}))
	assert.Equal(t, 4, len(headerMap))
	assert.Equal(t, "Second", headerMap[0].Name)
	assert.Nil(t, headerMap[1])
	assert.Equal(t, "Third", headerMap[2].Name)
	assert.Equal(t, "FirstField", headerMap[3].Name)
}

func TestParse(t *testing.T) {
	type testStruct struct {
		AnInt int	 `csv:"Int"`
		AFloat float32 `csv:"Float"`
		AString string `csv:"String"`
	}
	inputStr := `String,Float,Int,Ignore
string,1.01,3,ignore
stringtwo,1.02,4,ignore
`

	outputArr, err := Parse(strings.NewReader(inputStr), &testStruct{})
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, []testStruct{{3,1.01,"string"},{4,1.02,"stringtwo"}}, outputArr)
	}
}