package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ttlDays = 30

// dbPool is the subset of *pgxpool.Pool's interface we depend on, so tests can substitute pgxmock.
type dbPool interface {
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type Postgres struct {
	pool dbPool
}

type CacheResponse struct {
	Colors []string
	Counts []int
	Hit    bool
}

const schema = `
CREATE TABLE IF NOT EXISTS poster_color_cache (
	poster_id TEXT PRIMARY KEY,
	colors TEXT[] NOT NULL,
	counts INTEGER[] NOT NULL,
	last_accessed TIMESTAMPTZ NOT NULL DEFAULT now()
);
`

// New connects to postgres (via a pooled connection, suited to a long-running server process)
// and ensures the cache table exists.
func New(url string) (Postgres, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return Postgres{}, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return Postgres{}, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if _, err := pool.Exec(ctx, schema); err != nil {
		return Postgres{}, fmt.Errorf("failed to ensure poster_color_cache table: %w", err)
	}

	return Postgres{pool: pool}, nil
}

// Get fetches a cache entry by key. A stale (past TTL) or missing entry returns Hit: false.
func (p Postgres) Get(key string) (CacheResponse, error) {
	res, err := p.GetBatch([]string{key})
	if err != nil {
		return CacheResponse{}, err
	}
	return res[key], nil
}

// Set writes a cache entry. colors and counts are truncated to their shared length.
func (p Postgres) Set(key string, colors []string, counts []int) error {
	if !strings.Contains(key, "_") {
		return fmt.Errorf("invalid postgres cache key format")
	}
	return p.SetBatch([]string{key}, [][]string{colors}, [][]int{counts})
}

// GetBatch fetches multiple cache entries by key. Keys that are missing or stale (past TTL)
// come back with Hit: false. Successfully retrieved entries have their last_accessed refreshed.
func (p Postgres) GetBatch(keys []string) (map[string]CacheResponse, error) {
	res := make(map[string]CacheResponse, len(keys))
	for _, k := range keys {
		res[k] = CacheResponse{Hit: false}
	}
	if len(keys) == 0 {
		return res, nil
	}

	ctx := context.Background()
	rows, err := p.pool.Query(ctx, `
		UPDATE poster_color_cache SET last_accessed = now()
		WHERE poster_id = ANY($1) AND last_accessed > now() - make_interval(days => $2)
		RETURNING poster_id, colors, counts
	`, keys, ttlDays)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch from postgres: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var colors []string
		var counts []int
		if err := rows.Scan(&id, &colors, &counts); err != nil {
			return nil, fmt.Errorf("failed to scan postgres row: %w", err)
		}
		res[id] = CacheResponse{Colors: colors, Counts: counts, Hit: true}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed reading postgres rows: %w", err)
	}

	return res, nil
}

// SetBatch writes multiple cache entries. Each colors/counts pair is truncated to its shared length.
func (p Postgres) SetBatch(keys []string, colors [][]string, counts [][]int) error {
	if len(keys) != len(colors) || len(keys) != len(counts) {
		return fmt.Errorf("length of keys, colors, and counts do not match")
	}
	if len(keys) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for i, key := range keys {
		n := min(len(colors[i]), len(counts[i]))
		for _, c := range colors[i][:n] {
			if len(c) != 7 {
				return fmt.Errorf("weird color string length: %d", len(c))
			}
		}
		for _, c := range counts[i][:n] {
			if c < 0 || c > 9999 {
				return fmt.Errorf("color count is out of range: %d", c)
			}
		}
		batch.Queue(`
			INSERT INTO poster_color_cache (poster_id, colors, counts, last_accessed)
			VALUES ($1, $2, $3, now())
			ON CONFLICT (poster_id) DO UPDATE
			SET colors = EXCLUDED.colors, counts = EXCLUDED.counts, last_accessed = now()
		`, key, colors[i][:n], counts[i][:n])
	}

	br := p.pool.SendBatch(context.Background(), batch)
	defer br.Close()
	for range keys {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("error batch-setting to postgres: %w", err)
		}
	}
	return nil
}
