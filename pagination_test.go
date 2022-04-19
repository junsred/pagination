package pagination

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var fetchableData = []interface{}{
	&TimeItem{Data: 1, Time: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)},
	&TimeItem{Data: 2, Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
}

func F(lastTime *time.Time) ([]interface{}, error) {
	r := []interface{}{}
	for _, v := range fetchableData {
		item, _ := v.(*TimeItem)
		if item.Time.Before(*lastTime) {
			r = append(r, v)
			*lastTime = v.(*TimeItem).Time
		}
	}
	return r, nil
}

func pagFunc(lastTime *time.Time) func(ctx context.Context, needed int) ([]interface{}, error) {
	return func(ctx context.Context, needed int) ([]interface{}, error) {
		fetchedData, err := F(lastTime)
		if err != nil {
			return nil, err
		}
		if len(fetchedData) == 0 {
			return nil, nil
		}
		process := func(data interface{}) interface{} {
			TimeItem, ok := data.(*TimeItem)
			if !ok {
				return nil
			}
			return TimeItem
		}
		return ProcessConcurrently(fetchedData, needed, process), nil
	}
}

func TestPagination(t *testing.T) {
	t.Run("try getting all", func(t *testing.T) {
		now := time.Now()
		p := New(pagFunc(&now), 2)
		pagResult, err := p.Paginate(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, 2, len(pagResult.Data))
		assert.Equal(t, false, pagResult.HasNext)
	})
	t.Run("try getting half", func(t *testing.T) {
		now := time.Now()
		lastTime := &now
		p := New(pagFunc(lastTime), 1)
		pagResult, err := p.Paginate(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(pagResult.Data))
		assert.Equal(t, true, pagResult.HasNext)
		if len(pagResult.Data) == 0 {
			return
		}
		*lastTime = pagResult.Data[0].(*TimeItem).Time
		pagResult, err = p.Paginate(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(pagResult.Data))
		assert.Equal(t, false, pagResult.HasNext)
	})
}
