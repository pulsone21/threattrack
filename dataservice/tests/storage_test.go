package tests

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/pulsone21/threattrack/dataService/storage"
	"github.com/stretchr/testify/assert"
)

func TestLoadRawSQL(t *testing.T) {
	loadSqlCases := []struct {
		name        string
		subpath     string
		expectedSql string
		expectedErr bool
	}{
		{"Successful SQL Load", "get.sql", "SELECT * FROM TEST", false},
		{"Error SQL Load", "failing_test.sql", "", true},
	}
	for _, tC := range loadSqlCases {
		t.Run(tC.name, func(t *testing.T) {
			sql, err := storage.LoadRawSQL(tC.subpath)
			assert.Equal(t, tC.expectedSql, sql, fmt.Sprintf("TestCase %s Expected: '%s' Actual: '%s'", tC.name, tC.expectedSql, sql))
			assert.Equal(t, tC.expectedErr, err != nil, fmt.Sprintf("TestCase %s Expected: '%t' Actual: '%t'", tC.name, tC.expectedErr, err != nil))
		})
	}
}

func TestFinalizeMySql(t *testing.T) {
	finalizeSqlCases := []struct {
		name        string
		rawSql      string
		entity      string
		params      *storage.QueryParameter
		expectedSql string
	}{
		{"With WHERE Statement", "SELECT * \nFROM TEST\n %s LIMIT ? OFFSET ?;", "incidents", &storage.QueryParameter{Limit: 20, Offset: 0, Query: map[string]string{"id": "123"}}, "SELECT * \nFROM TEST\n WHERE incidents.id=\"123\" LIMIT ? OFFSET ?;"},
		{"Without WHERE Statement", "SELECT * \nFROM TEST\n %s LIMIT ? OFFSET ?;", "incidents", &storage.QueryParameter{Limit: 20, Offset: 0, Query: map[string]string{}}, "SELECT * \nFROM TEST\n  LIMIT ? OFFSET ?;"},
	}

	for _, tC := range finalizeSqlCases {
		t.Run(tC.name, func(t *testing.T) {
			actual := storage.FinalizeSQL(tC.rawSql, tC.entity, *tC.params)
			assert.Equal(t, tC.expectedSql, actual, fmt.Sprintf("TestCase %s Expected: '%s', Actual: '%s'", tC.name, tC.expectedSql, actual))
		})
	}
}

func makeUrlValues(vals map[string]string) url.Values {
	urlV := make(url.Values)
	for k, v := range vals {
		urlV.Add(k, v)
	}
	return urlV
}

func TestExtractUrlQueries(t *testing.T) {
	testTable := []struct {
		name           string
		urlValues      url.Values
		withParams     bool
		expectedLimit  int
		expectedOffset int
		expectedMap    map[string]string
	}{
		{"Given Limit", makeUrlValues(map[string]string{"limit": "20", "offset": "0"}), true, 20, 0, map[string]string{}},
		{"Empty UrlQuery", makeUrlValues(map[string]string{}), true, 200, 0, map[string]string{}},
		{"Given UrlParams, Params allowed", makeUrlValues(map[string]string{"status": "open"}), true, 200, 0, map[string]string{"status": "open"}},
		{"Given UrlParams but no Params allowed", makeUrlValues(map[string]string{"status": "open"}), false, 200, 0, map[string]string{}},
	}
	for _, tC := range testTable {
		t.Run(tC.name, func(t *testing.T) {
			limit, offset, params := storage.ExtractUrlQueries(tC.urlValues, tC.withParams)
			assert.Equal(t, tC.expectedLimit, limit, fmt.Sprintf("TestCase %s", tC.name))
			assert.Equal(t, tC.expectedOffset, offset, fmt.Sprintf("TestCase %s", tC.name))
			assert.Equal(t, tC.expectedMap, params, fmt.Sprintf("TestCase %s", tC.name))
		})
	}
}
