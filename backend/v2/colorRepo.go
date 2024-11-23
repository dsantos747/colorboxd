package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type PosterColor struct {
	posterID    string
	colors      []colorful.Color
	count       []int
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

func (cr *ColorRepo) Get(ctx context.Context, posterIDs []string) map[string]PosterColor {
	pc := make(map[string]PosterColor)

	for _, id := range posterIDs {
		p, _ := cr.posterColors.Load(id)
		poster, err := cr.assertPosterColorType(p, id)
		if err != nil {
			continue
		}
		poster.lastQueried = time.Now()
		cr.posterColors.Store(id, poster)
		pc[poster.posterID] = poster
	}

	return pc
}

func (cr *ColorRepo) assertPosterColorType(p any, id string) (PosterColor, error) {
	switch poster := p.(type) {
	case *PosterColor:
		if poster == nil {
			cr.posterColors.Delete(id)
			break
		}
		return *poster, nil
	case PosterColor:
		return poster, nil
	}
	return PosterColor{}, fmt.Errorf("invalid color, deleted from store%s", "")
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
