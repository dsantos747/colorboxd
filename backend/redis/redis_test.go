package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func setupMockRedis() *miniredis.Miniredis {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return mr
}

func TestSomething(t *testing.T) {
	mr := setupMockRedis()
	defer mr.Close()

	// Write some tests
}
