package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// TODO:
// - Create "GetBatch" method
// - Don't check for a ttl parameter, or set a specific ttl key. Just rely on redis's own ttl mechanism, which will return false if the key is stale. Investigate this.

const ttlDays int64 = 30

type Redis struct {
	client *redis.Client
}

type CacheResponse struct {
	Colors []string
	Counts []int
	Hit    bool
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

func (r Redis) GetBatch(keys []string) (map[string]CacheResponse, error) {
	res := make(map[string]CacheResponse)

	redSlice, err := r.client.MGet(context.TODO(), keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to mget from redis: %w", err)
	}

	for i, key := range keys {
		if redSlice[i] == nil {
			res[key] = CacheResponse{
				Hit: false,
			}
			continue
		}

		vals, ok := redSlice[i].(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert redis response to string")
		}

		colors, counts := r.parseRedisOut(vals)

		res[key] = CacheResponse{
			Colors: colors,
			Counts: counts,
			Hit:    true,
		}
	}

	return res, nil
}

// Gets a value from redis given a key. If the key is stale, cacheHit is returned as false
func (r Redis) Get(key string) CacheResponse {
	// Get the value
	vals, err := r.client.Get(context.TODO(), key).Result()
	if err != nil || len(vals) == 0 {
		return CacheResponse{Hit: false}
	}
	colors, counts := r.parseRedisOut(vals)
	return CacheResponse{Colors: colors, Counts: counts, Hit: true}
}

func (r Redis) Set(key string, colors []string, counts []int) {
	ctx := context.Background()

	val := r.parseRedisIn(colors, counts)
	resInt := r.client.Set(ctx, key, val, time.Duration(ttlDays*24)*time.Hour)
	if resInt.Err() != nil || resInt.Val() == "" {
		// TODO: Handle error here!!
		panic("unhandled redis setting error!")
	}
}

func (r Redis) parseRedisOut(vals string) ([]string, []int) {
	colors := []string{}
	counts := []int{}

	slc := strings.Split(vals, ",")

	if len(slc) != 3 {
		panic("weird length of redis value")
	}

	for _, c := range slc {
		count, err := strconv.Atoi(c[7:])
		if err != nil {
			fmt.Println("Error converting count to int")
			count = 2500 // Literally out of thin air
		}

		colors = append(colors, c[:7])
		counts = append(counts, count)

	}
	return colors, counts
}

func (r Redis) parseRedisIn(colors []string, counts []int) string {
	slc := []string{"", "", ""}
	for i := 0; i < 3; i++ {
		if i >= len(colors) {
			slc[i] = "XXXXXXX0000"
			continue
		}
		if len(colors[i]) != 7 { // used during development
			panic("weird color length")
		}
		if counts[i] > 9999 {
			panic("count is too high")
		}

		slc[i] = colors[i] + fmt.Sprintf("%04d", counts[i])
	}
	return strings.Join(slc, ",")
}
