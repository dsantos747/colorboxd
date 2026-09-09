package postgres

import (
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func newMockPostgres(t *testing.T) (Postgres, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	t.Cleanup(mock.Close)
	return Postgres{pool: mock}, mock
}

func TestGetSet(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name       string
		key        string
		colors     []string
		counts     []int
		hit        bool
		colors_out []string
		counts_out []int
		errStrSet  string
	}{
		{
			name:   "Successful Set and Get",
			key:    "testKey_1",
			colors: []string{"#FF0000", "#00FF00", "#0000FF"},
			counts: []int{3000, 200, 10},
			hit:    true,
		},
		{
			name:      "Fail due to bad color format",
			key:       "testKey_2",
			colors:    []string{"badFormat", "terribleFormat", "youCan'tExpectMeToBelieveThisIsAHexColour"},
			counts:    []int{3000, 200, 10},
			hit:       false,
			errStrSet: "weird color string length",
		},
		{
			name:      "Fail due to bad count format - < 0",
			key:       "testKey_3",
			colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			counts:    []int{-20, 1, 1},
			hit:       false,
			errStrSet: "color count is out of range",
		},
		{
			name:      "Fail due to bad count format - > 9999",
			key:       "testKey_4",
			colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			counts:    []int{10000, 1, 1},
			hit:       false,
			errStrSet: "color count is out of range",
		},
		{
			name:       "Successfully set and get, insufficient length",
			key:        "testKey_5",
			colors:     []string{"#FF0000"},
			counts:     []int{3000, 200, 10},
			hit:        true,
			colors_out: []string{"#FF0000"},
			counts_out: []int{3000},
		},
		{
			name:      "Fail with invalid key format",
			key:       "badKeyName",
			colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			counts:    []int{3000, 200, 10},
			hit:       false,
			errStrSet: "invalid postgres cache key format",
		},
	}

	p, mock := newMockPostgres(t)

	for _, tc := range testCases {
		if tc.colors_out == nil {
			tc.colors_out = tc.colors
		}
		if tc.counts_out == nil {
			tc.counts_out = tc.counts
		}

		if tc.errStrSet == "" {
			mock.ExpectBatch().ExpectExec("INSERT INTO poster_color_cache").
				WithArgs(tc.key, tc.colors_out, tc.counts_out).
				WillReturnResult(pgxmock.NewResult("INSERT", 1))
		}

		err := p.Set(tc.key, tc.colors, tc.counts)
		if tc.errStrSet != "" {
			assert.ErrorContains(err, tc.errStrSet)
			continue
		}
		assert.Nil(err)

		if tc.hit {
			rows := pgxmock.NewRows([]string{"poster_id", "colors", "counts"}).
				AddRow(tc.key, tc.colors_out, tc.counts_out)
			mock.ExpectQuery("UPDATE poster_color_cache").WithArgs([]string{tc.key}, ttlDays).WillReturnRows(rows)
		} else {
			mock.ExpectQuery("UPDATE poster_color_cache").WithArgs([]string{tc.key}, ttlDays).
				WillReturnRows(pgxmock.NewRows([]string{"poster_id", "colors", "counts"}))
		}

		res, err := p.Get(tc.key)
		assert.Nil(err)
		assert.Equal(tc.hit, res.Hit)
		if tc.hit {
			assert.Equal(tc.colors_out, res.Colors)
			assert.Equal(tc.counts_out, res.Counts)
		}
	}

	assert.NoError(mock.ExpectationsWereMet())
}

func TestGetSetBatch(t *testing.T) {
	assert := assert.New(t)
	p, mock := newMockPostgres(t)

	keys := []string{"testKey1_1", "testKey1_2", "testKey1_3"}
	colors := [][]string{{"#FF0000", "#FF0000", "#FF0000"}, {"#00FF00", "#00FF00", "#00FF00"}, {"#0000FF", "#0000FF", "#0000FF"}}
	counts := [][]int{{1000, 100, 10}, {2000, 200, 20}, {3000, 300, 30}}

	batchExp := mock.ExpectBatch()
	for i, k := range keys {
		batchExp.ExpectExec("INSERT INTO poster_color_cache").
			WithArgs(k, colors[i], counts[i]).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
	}

	err := p.SetBatch(keys, colors, counts)
	assert.NoError(err)

	getKeys := []string{"testKey1_1", "badKey1", "badKey2"}
	rows := pgxmock.NewRows([]string{"poster_id", "colors", "counts"}).
		AddRow(keys[0], colors[0], counts[0])
	mock.ExpectQuery("UPDATE poster_color_cache").WithArgs(getKeys, ttlDays).WillReturnRows(rows)

	res, err := p.GetBatch(getKeys)
	assert.NoError(err)
	assert.Equal(map[string]CacheResponse{
		"testKey1_1": {Colors: colors[0], Counts: counts[0], Hit: true},
		"badKey1":    {Hit: false},
		"badKey2":    {Hit: false},
	}, res)

	assert.NoError(mock.ExpectationsWereMet())
}

func TestSetBatchLengthMismatch(t *testing.T) {
	p, _ := newMockPostgres(t)
	err := p.SetBatch([]string{"a_1"}, [][]string{{"#FF0000"}, {"#00FF00"}}, [][]int{{1}})
	assert.ErrorContains(t, err, "length of keys, colors, and counts do not match")
}
