package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

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

		colors, counts, err := r.parseRedisOut(vals)
		if err != nil {
			fmt.Println("error parsing output from redis: %w", err)
			res[key] = CacheResponse{
				Hit: false,
			}
			continue
		}

		res[key] = CacheResponse{
			Colors: colors,
			Counts: counts,
			Hit:    true,
		}
	}

	return res, nil
}

// Gets a value from redis given a key. If the key is stale, cacheHit is returned as false
func (r Redis) Get(key string) (CacheResponse, error) {
	// Get the value
	vals, err := r.client.Get(context.TODO(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return CacheResponse{Hit: false}, nil
		}
		return CacheResponse{}, fmt.Errorf("error getting from redis: %w", err) // If logs show nil err, then val == ""
	}
	colors, counts, err := r.parseRedisOut(vals)
	if err != nil {
		return CacheResponse{Hit: false}, fmt.Errorf("error parsing output from redis: %w", err)
	}

	return CacheResponse{Colors: colors, Counts: counts, Hit: true}, nil
}

func (r Redis) Set(key string, colors []string, counts []int) error {
	if !strings.Contains(key, "_") {
		return fmt.Errorf("invalid redis key format")
	}

	ctx := context.Background()

	val, err := r.parseRedisIn(colors, counts)
	if err != nil {
		return fmt.Errorf("error parsing redis input: %w", err)
	}

	resInt := r.client.Set(ctx, key, val, time.Duration(ttlDays*24)*time.Hour)
	if resInt.Err() != nil || resInt.Val() == "" {
		return fmt.Errorf("error setting to redis: %w", resInt.Err()) // If logs show nil err, then val == ""
	}
	return nil
}

func (r Redis) parseRedisOut(vals string) ([]string, []int, error) {
	colors := []string{}
	counts := []int{}

	slc := strings.Split(vals, ",")

	if len(slc) != 3 {
		return nil, nil, fmt.Errorf("unexpected length of value fetched from redis; length %d", len(slc))
	}

	for _, c := range slc {
		count, err := strconv.Atoi(c[7:])
		if err != nil || count < 0 {
			return nil, nil, fmt.Errorf("invalid count post-conversion: %w", err)
		}

		if c[:7] != "XXXXXXX" {
			colors = append(colors, c[:7])
		}
		if count != 0 {
			counts = append(counts, count)
		}

	}
	return colors, counts, nil
}

func (r Redis) parseRedisIn(colors []string, counts []int) (string, error) {
	slc := []string{"", "", ""}
	for i := 0; i < 3; i++ {
		if i >= min(len(colors), len(counts)) {
			slc[i] = "XXXXXXX0000"
			continue
		}
		if len(colors[i]) != 7 { // used during development
			return "", fmt.Errorf("weird color string length: %d", len(colors[i]))
		}
		if counts[i] > 9999 || counts[i] < 0 {
			return "", fmt.Errorf("color count is out of range: %d", counts[i])
		}

		slc[i] = colors[i] + fmt.Sprintf("%04d", counts[i])
	}
	return strings.Join(slc, ","), nil
}
