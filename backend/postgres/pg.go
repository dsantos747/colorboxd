package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CacheEntry represents a cache record
type CacheEntry struct {
	ID         string
	Color1     string
	Color2     string
	Color3     string
	Value1     int
	Value2     int
	Value3     int
	LastAccess time.Time
	CreatedAt  time.Time
}

type PGservice struct {
	l    *slog.Logger
	Pool *pgxpool.Pool
}

// Connect to Neon Postgres
func New(l *slog.Logger, url string) (*PGservice, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return &PGservice{l, pool}, nil
}

// Insert a new cache entry
func (s PGservice) InsertCacheEntry(ctx context.Context, entry CacheEntry) error {
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO colors_cache (poster_id, color1, color2, color3, value1, value2, value3, last_accessed, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
		ON CONFLICT (poster_id) DO UPDATE 
		SET color1 = EXCLUDED.color1, color2 = EXCLUDED.color2, color3 = EXCLUDED.color3,
		    value1 = EXCLUDED.value1, value2 = EXCLUDED.value2, value3 = EXCLUDED.value3,
		    last_accessed = now()
	`, entry.ID, entry.Color1, entry.Color2, entry.Color3, entry.Value1, entry.Value2, entry.Value3)
	return err
}

// Retrieve a cache entry and update `last_accessed`
func (s PGservice) GetCacheEntry(ctx context.Context, id string) (*CacheEntry, error) {
	var entry CacheEntry
	err := s.Pool.QueryRow(ctx, `
		UPDATE colors_cache SET last_accessed = now()
		WHERE poster_id = $1
		RETURNING poster_id, color1, color2, color3, value1, value2, value3, last_accessed, created_at
	`, id).Scan(
		&entry.ID, &entry.Color1, &entry.Color2, &entry.Color3,
		&entry.Value1, &entry.Value2, &entry.Value3,
		&entry.LastAccess, &entry.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s PGservice) InsertCacheBatch(ctx context.Context, entries []CacheEntry) error {
	if len(entries) == 0 {
		return nil
	}

	// Construct the bulk insert SQL
	sql := `
		INSERT INTO colors_cache (poster_id, color1, color2, color3, value1, value2, value3, last_accessed, created_at)
		VALUES %s
		ON CONFLICT (poster_id) DO UPDATE 
		SET color1 = EXCLUDED.color1, color2 = EXCLUDED.color2, color3 = EXCLUDED.color3,
		    value1 = EXCLUDED.value1, value2 = EXCLUDED.value2, value3 = EXCLUDED.value3,
		    last_accessed = now();
	`

	// Build placeholders dynamically
	values := make([]string, len(entries))
	args := []interface{}{}
	for i, entry := range entries {
		values[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, now(), now())",
			i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7)
		args = append(args, entry.ID, entry.Color1, entry.Color2, entry.Color3, entry.Value1, entry.Value2, entry.Value3)
	}

	// Final query
	query := fmt.Sprintf(sql, strings.Join(values, ","))
	_, err := s.Pool.Exec(ctx, query, args...)
	return err
}
func (s PGservice) GetCacheBatch(ctx context.Context, ids []string) (map[string]CacheEntry, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// Construct placeholders dynamically
	placeholders := []string{}
	args := []interface{}{}
	for i, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		UPDATE colors_cache 
		SET last_accessed = now()
		WHERE poster_id = ANY(ARRAY[%s]::TEXT[])
		RETURNING poster_id, color1, color2, color3, value1, value2, value3, last_accessed, created_at;
	`, strings.Join(placeholders, ","))

	rows, err := s.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Use a map for fast lookup of cached results
	cacheMap := make(map[string]CacheEntry)

	for rows.Next() {
		var entry CacheEntry
		err := rows.Scan(
			&entry.ID, &entry.Color1, &entry.Color2, &entry.Color3,
			&entry.Value1, &entry.Value2, &entry.Value3,
			&entry.LastAccess, &entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cacheMap[entry.ID] = entry
	}
	return cacheMap, nil
}

// Manually trigger LRU eviction (optional)
func (s PGservice) EvictOldCache(ctx context.Context) {
	rows, err := s.Pool.Exec(ctx, `
		DELETE FROM colors_cache WHERE poster_id IN (
			SELECT poster_id FROM colors_cache ORDER BY last_accessed ASC OFFSET 1000
		);
	`)
	if err != nil {
		s.l.Error("failed to evict rows", "error", err)
		return
	}
	if rows.RowsAffected() > 0 {
		s.l.Info("evicted rows", "rows_affected", rows.RowsAffected())
	}
}
