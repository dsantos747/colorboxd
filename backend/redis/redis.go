package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(url string) Redis {
	var client *redis.Client
	opt, err := redis.ParseURL(url)
	if err == nil {
		client = redis.NewClient(opt)
	}

	return Redis{
		client: client,
	}
}

// Gets a a value from redis given a key. If the key is stale, cacheHit is returned as false
func (r Redis) Get(key string) ([]string, []int, bool) {
	vals, err := r.client.LRange(context.TODO(), key, 0, -1).Result()
	if err != nil {
		return nil, nil, false
	}
	ts, err := r.client.Get(context.TODO(), key+"_t").Result()
	if err != nil {
		return nil, nil, false
	}
	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil, nil, false
	}
	if time.Now().Unix()-tsInt > 30*24*3600 {
		return nil, nil, false
	}
	colors, counts := r.parseRedisOut(vals)
	return colors, counts, true
}

func (r Redis) Set(key string, colors []string, counts []int) {
	vals := r.parseRedisIn(colors, counts)
	r.client.RPush(context.TODO(), key, vals) // TTL?

	// Set the timestamp of the key
	tsKey := key + "_t"
	r.client.Set(context.TODO(), tsKey, time.Now().Unix(), 0)
}

func (r Redis) parseRedisOut(vals []string) ([]string, []int) {
	colors := make([]string, len(vals))
	counts := make([]int, len(vals))
	for i, c := range vals {
		count, err := strconv.Atoi(c[6:])
		if err != nil {
			fmt.Println("Error converting count to int")
			count = 2500 // Literally out of thin air
		}

		colors[i] = c[:6]
		counts[i] = count

	}
	return colors, counts
}

func (r Redis) parseRedisIn(colors []string, counts []int) [3]string {
	out := [3]string{}
	for i := 0; i < 3; i++ {
		if i >= len(colors) {
			out[i] = ""
		}
		out[i] = colors[i] + strconv.Itoa(counts[i])
	}
	return out
}