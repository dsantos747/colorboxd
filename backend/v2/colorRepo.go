package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type PosterColors struct {
	posterURL   string
	colors      []colorful.Color
	counts      []int
	lastQueried time.Time
}

type ColorRepo struct {
	l            *slog.Logger
	posterColors *sync.Map
}

// NewColorRepo returns a repo with a store
func NewColorRepo(l *slog.Logger) (*ColorRepo, error) {
	cr := &ColorRepo{
		l:            l,
		posterColors: &sync.Map{},
	}

	return cr, nil
}

func (cr *ColorRepo) Get(ctx context.Context, posterURLs []string) map[string]PosterColors {
	pc := make(map[string]PosterColors)

	for _, url := range posterURLs {
		p, _ := cr.posterColors.Load(url)
		poster, err := cr.assertPosterColorType(p, url)
		if err != nil {
			continue
		}
		poster.lastQueried = time.Now()
		cr.posterColors.Store(url, poster)
		pc[poster.posterURL] = poster
	}

	return pc
}

func (cr *ColorRepo) Set(ctx context.Context, posterURL string, posterColors []Color) error {
	poster := PosterColors{
		posterURL:   posterURL,
		colors:      make([]colorful.Color, len(posterColors)),
		counts:      make([]int, len(posterColors)),
		lastQueried: time.Now(),
	}

	//Convert Color to posterColors
	for i, col := range posterColors {
		poster.colors[i] = col.rgb
		poster.counts[i] = col.count
	}

	cr.posterColors.Store(posterURL, poster)

	return nil
}

func (cr *ColorRepo) assertPosterColorType(p any, url string) (PosterColors, error) {
	switch poster := p.(type) {
	case *PosterColors:
		if poster == nil {
			cr.posterColors.Delete(url)
			break
		}
		return *poster, nil
	case PosterColors:
		return poster, nil
	}
	return PosterColors{}, fmt.Errorf("invalid color, deleted from store%s", "")
}

// getStoreLength returns the current length of the store
func (cr *ColorRepo) GetStoreLength() int {
	count := 0
	cr.posterColors.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}
