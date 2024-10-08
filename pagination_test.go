package pagination

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var fetchableData = []*Item[time.Time, int]{
	{Key: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC), Value: 1},
	{Key: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), Value: 2},
	{Key: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Value: 3},
}

func fetch(ctx context.Context, limit int, lastTime *time.Time) ([]*Item[time.Time, int], bool, error) {
	limit += 1
	r := []*Item[time.Time, int]{}
	hasNext := false
	for _, v := range fetchableData {
		if v.Key.Before(*lastTime) {
			if len(r) == limit {
				hasNext = true
				break
			}
			r = append(r, v)
		}
	}
	return r, hasNext, nil
}

func process(ctx context.Context, data *Item[time.Time, int]) (*Item[time.Time, int], bool) {
	if data.Value == 3 {
		return nil, false
	}
	return data, true
}

func pagFunc() PaginateFunc[time.Time, int] {
	return FetchAndProcess(fetch, ProcessConcurrently(process))
}

func produceNext(ctx context.Context, result *PaginationResult[time.Time, int]) string {
	return result.Items[len(result.Items)-1].Key.String()
}

func TestPagination(t *testing.T) {
	now := time.Now()
	t.Run("try getting 3", func(t *testing.T) {
		p := New(pagFunc(), now)
		pagResult, err := p.Paginate(context.TODO(), 3)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(pagResult.Items))
		assert.Equal(t, false, pagResult.HasNext)
	})
	t.Run("try getting 3 by 1", func(t *testing.T) {
		p := New(pagFunc(), now)
		pagResult, err := p.Paginate(context.TODO(), 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(pagResult.Items))
		assert.Equal(t, true, pagResult.HasNext)
		pagResult, err = p.Paginate(context.TODO(), 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(pagResult.Items))
		assert.Equal(t, true, pagResult.HasNext)
	})
	t.Run("try getting 2", func(t *testing.T) {
		p := New(pagFunc(), now)
		pagResult, err := p.Paginate(context.TODO(), 2)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(pagResult.Items))
		assert.Equal(t, true, pagResult.HasNext)
	})
	t.Run("try timeout", func(t *testing.T) {
		p := New(pagFunc(), now, WithTimeout(0))
		pagResult, err := p.Paginate(context.TODO(), 2)
		assert.Error(t, err, ErrResourceExhausted)
		assert.Equal(t, 0, len(pagResult.Items))
	})
	t.Run("try next func", func(t *testing.T) {
		p := New(pagFunc(), now, WithNextFunc(produceNext))
		pagResult, err := p.Paginate(context.TODO(), 2)
		assert.NoError(t, err, nil)
		assert.Equal(t, 2, len(pagResult.Items))
		assert.Equal(t, true, pagResult.HasNext)
		assert.Equal(t, "2020-02-01 00:00:00 +0000 UTC", *pagResult.Next)
	})
}
