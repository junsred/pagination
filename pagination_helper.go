package pagination

import (
	"context"
	"sync"
)

type ProcessFunc[K comparable, V any] func(context.Context, []*Item[K, V], int) ([]*Item[K, V], error)
type ProcessConcurrentFunc[K comparable, V any] func(context.Context, *Item[K, V]) (*Item[K, V], bool)

// FetchAndProcess divides the fetching and processing of data
//
//	func(ctx context.Context, limit int) ([]T, bool, error)
func FetchAndProcess[K comparable, V any](
	fetchFunc PaginateFunc[K, V],
	processFunc ProcessFunc[K, V]) PaginateFunc[K, V] {
	return func(ctx context.Context, needed int, lastKey *K) ([]*Item[K, V], bool, error) {
		fetchedData, hasNext, err := fetchFunc(ctx, needed, lastKey)
		if err != nil {
			return nil, hasNext, err
		}
		if len(fetchedData) == 0 {
			return nil, hasNext, nil
		}
		*lastKey = fetchedData[len(fetchedData)-1].Key
		processedData, err := processFunc(ctx, fetchedData, needed)
		if !hasNext && len(fetchedData) > needed {
			if len(processedData) >= needed {
				hasNext = true
			}
		}
		return processedData, hasNext, err
	}
}

// ProcessConcurrently calls the process function concurrently for every item
// context will be cancelled when needed items are processed
//
//	func(ctx context.Context, item T) (T, bool)
func ProcessConcurrently[K comparable, V any](process ProcessConcurrentFunc[K, V]) ProcessFunc[K, V] {
	return func(ctx context.Context, fetchedData []*Item[K, V], needed int) ([]*Item[K, V], error) {
		clean := func(items []*Item[K, V]) []*Item[K, V] {
			var clean []*Item[K, V]
			for _, item := range items {
				if item != nil {
					clean = append(clean, item)
				}
			}
			return clean
		}
		dataLen := len(fetchedData)
		processedData := make([]*Item[K, V], dataLen)
		var wg sync.WaitGroup
		type Pair struct {
			i    int
			data *Item[K, V]
		}
		tasks := make(chan Pair, dataLen)
		worker := func() {
			defer wg.Done()
			for task := range tasks {
				if ctx.Err() != nil {
					return
				}
				if pData, ok := process(ctx, task.data); ok {
					processedData[task.i] = pData
					return
				}
			}
		}
		for range needed {
			wg.Add(1)
			go worker()
		}
		for i, data := range fetchedData {
			tasks <- Pair{
				i:    i,
				data: data,
			}
		}
		close(tasks)
		wg.Wait()
		return clean(processedData), nil
	}
}
