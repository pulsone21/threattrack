package tests

import (
	"fmt"
	"testing"

	"github.com/pulsone21/threattrack/lib/entities"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	name string
}

func createTestStruct(name string) *testStruct {
	return &testStruct{
		name: name,
	}
}

func createTestStructs(names []string) *[]testStruct {
	var t []testStruct
	for _, v := range names {
		t = append(t, testStruct{
			name: v,
		})
	}
	return &t
}

func TestNewApiResponse(t *testing.T) {
	ts := createTestStructs([]string{"Hello", "World"})
	exp1 := entities.ApiResponse{
		StatusCode: 200,
		RequestUrl: "/test",
		Data:       *ts,
	}

	t1 := createTestStruct("Hello World")
	t1s := []testStruct{*t1}
	exp2 := entities.ApiResponse{
		StatusCode: 200,
		RequestUrl: "/test",
		Data:       t1s,
	}
	testCases := []struct {
		name       string
		statusCode int
		uri        string
		data       interface{}
		Expected   string
	}{
		{"ApiResponse given a Slice", 200, "/test", *ts, fmt.Sprintf("&%v", exp1)},
		{"ApiResponse given a single interface", 200, "/test", *t1, fmt.Sprintf("&%v", exp2)},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			res := entities.NewApiResponse(tC.statusCode, tC.uri, tC.data)
			actual := fmt.Sprint(res)
			assert.Equal(t, tC.Expected, actual, fmt.Sprintf("TestCase %s Expected: '%s', Actual: '%s'", tC.name, tC.Expected, actual))
		})
	}
}
