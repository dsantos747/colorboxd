package redis

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

// TO TEST
// GetBatch method works
// Input with insufficient colors or counts
// Input with invalid color or count format
//

func TestSomething(t *testing.T) {
	assert := assert.New(t)

	s := miniredis.RunT(t)

	rc := New(fmt.Sprintf("redis://%s", s.Addr()))

	rc.Set("test_key", []string{"#FF0000"}, []int{1})

	res := rc.Get("test_key")

	assert.True(res.Hit)
	assert.Equal([]string{"#FF0000"}, res.Colors)
	assert.Equal([]int{1}, res.Counts)
}

func TestSomething2(t *testing.T) {
	assert := assert.New(t)

	s := miniredis.RunT(t)

	rc := New(fmt.Sprintf("redis://%s", s.Addr()))

	rc.Set("test_key", []string{"#FF0000", "#FF0000", "#FF0000"}, []int{1, 2, 3})

	res := rc.Get("test_key")

	assert.True(res.Hit)
	assert.Equal([]string{"#FF0000", "#FF0000", "#FF0000"}, res.Colors)
	assert.Equal([]int{1, 2, 3}, res.Counts)
}
