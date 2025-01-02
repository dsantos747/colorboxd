package redis

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

// This test does most of the testing of formatting / basic redis interaction
func TestGetSet(t *testing.T) {
	assert := assert.New(t)
	s := miniredis.RunT(t)
	rc := New(fmt.Sprintf("redis://%s", s.Addr()))

	testCases := []struct {
		name       string
		key        string
		colors     []string
		counts     []int
		hit        bool
		colors_out []string
		counts_out []int
		errStrSet  string
		errStrGet  string
	}{
		{
			name:      "Successful Set and Get",
			key:       "testKey_1",
			colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			counts:    []int{3000, 200, 10},
			hit:       true,
			errStrSet: "",
			errStrGet: "",
		},
		{
			name:      "Fail due to bad color format",
			key:       "testKey_2",
			colors:    []string{"badFormat", "terribleFormat", "youCan'tExpectMeToBelieveThisIsAHexColour"},
			counts:    []int{3000, 200, 10},
			hit:       false,
			errStrSet: "weird color string length",
			errStrGet: "",
		},
		{
			name:      "Fail due to bad count format - < 0",
			key:       "testKey_3",
			colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			counts:    []int{-20, 1, 1},
			hit:       false,
			errStrSet: "color count is out of range",
			errStrGet: "",
		},
		{
			name:      "Fail due to bad count format - > 9999",
			key:       "testKey_4",
			colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			counts:    []int{10000, 1, 1},
			hit:       false,
			errStrSet: "color count is out of range",
			errStrGet: "",
		},
		{
			name:       "Successfully set and get, insufficient length",
			key:        "testKey_5",
			colors:     []string{"#FF0000"},
			counts:     []int{3000, 200, 10},
			hit:        true,
			colors_out: []string{"#FF0000"},
			counts_out: []int{3000},
			errStrSet:  "",
			errStrGet:  "",
		},
		{
			name:      "Successfully overwrite a key, and get",
			key:       "testKey_5",
			colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			counts:    []int{3000, 200, 10},
			hit:       true,
			errStrSet: "",
			errStrGet: "",
		},
		{
			name:      "Fail with invalid key format",
			key:       "badKeyName",
			colors:    []string{"#FF0000", "#00FF00", "#0000FF"},
			counts:    []int{3000, 200, 10},
			hit:       false,
			errStrSet: "invalid redis key format",
			errStrGet: "",
		},
	}

	for _, tc := range testCases {
		if tc.colors_out == nil {
			tc.colors_out = tc.colors
		}
		if tc.counts_out == nil {
			tc.counts_out = tc.counts
		}

		err := rc.Set(tc.key, tc.colors, tc.counts)
		if tc.errStrSet != "" {
			assert.ErrorContains(err, tc.errStrSet)
		} else {
			assert.Nil(err)
		}

		res, err := rc.Get(tc.key)
		assert.Equal(tc.hit, res.Hit)
		if tc.errStrGet != "" {
			assert.ErrorContains(err, tc.errStrGet)
		}
		if tc.hit {
			assert.Nil(err)
			assert.Equal(tc.colors_out, res.Colors)
			assert.Equal(tc.counts_out, res.Counts)
		}

	}
}

// This test builds upon the previous knowledge that our basic redis implementation works, and verifies that the batch check works
func TestGetSetBatch(t *testing.T) {
	assert := assert.New(t)
	s := miniredis.RunT(t)
	rc := New(fmt.Sprintf("redis://%s", s.Addr()))

	testCases := []struct {
		name      string
		keys      []string
		colors    [][]string
		counts    [][]int
		keys_out  []string
		res       map[string]CacheResponse
		errStrSet string
		errStrGet string
	}{
		{
			name:   "Successful Set and Get",
			keys:   []string{"testKey1_1", "testKey1_2", "testKey1_3"},
			colors: [][]string{{"#FF0000", "#FF0000", "#FF0000"}, {"#00FF00", "#00FF00", "#00FF00"}, {"#0000FF", "#0000FF", "#0000FF"}},
			counts: [][]int{{1000, 100, 10}, {2000, 200, 20}, {3000, 300, 30}},
			res: map[string]CacheResponse{
				"testKey1_1": {Colors: []string{"#FF0000", "#FF0000", "#FF0000"}, Counts: []int{1000, 100, 10}, Hit: true},
				"testKey1_2": {Colors: []string{"#00FF00", "#00FF00", "#00FF00"}, Counts: []int{2000, 200, 20}, Hit: true},
				"testKey1_3": {Colors: []string{"#0000FF", "#0000FF", "#0000FF"}, Counts: []int{3000, 300, 30}, Hit: true},
			},
			errStrSet: "",
			errStrGet: "",
		},
		{
			name:     "Successful Set and Get - try get some nonexistent keys",
			keys:     []string{"testKey2_1", "testKey2_2", "testKey2_3"},
			colors:   [][]string{{"#FF0000", "#FF0000", "#FF0000"}, {"#00FF00", "#00FF00", "#00FF00"}, {"#0000FF", "#0000FF", "#0000FF"}},
			counts:   [][]int{{1000, 100, 10}, {2000, 200, 20}, {3000, 300, 30}},
			keys_out: []string{"testKey2_1", "badKey1", "badKey2"},
			res: map[string]CacheResponse{
				"testKey2_1": {Colors: []string{"#FF0000", "#FF0000", "#FF0000"}, Counts: []int{1000, 100, 10}, Hit: true},
				"badKey1":    {Hit: false},
				"badKey2":    {Hit: false},
			},
			errStrSet: "",
			errStrGet: "",
		},
	}

	for _, tc := range testCases {
		if tc.keys_out == nil {
			tc.keys_out = tc.keys
		}

		var err error
		for i, k := range tc.keys {
			err = rc.Set(k, tc.colors[i], tc.counts[i])
			if err != nil {
				break
			}
		}
		if tc.errStrSet != "" {
			assert.ErrorContains(err, tc.errStrSet)
		} else {
			assert.Nil(err)
		}

		res, err := rc.GetBatch(tc.keys_out)
		if tc.errStrGet != "" {
			assert.ErrorContains(err, tc.errStrGet)
		} else {
			assert.Nil(err)
			assert.True(reflect.DeepEqual(tc.res, res))
			fmt.Println(tc.res)
			fmt.Println(res)
		}

	}
}
