package pagination

import (
	"context"
	"sync"
)

func ProcessConcurrently(fetchedData []interface{},
	needed int,
	process func(interface{}) interface{}) []interface{} {
	clean := func(items []interface{}) []interface{} {
		var clean []interface{}
		for _, item := range items {
			if item != nil {
				clean = append(clean, item)
			}
		}
		return clean
	}
	check := func(items []interface{}) ([]interface{}, bool) {
		cItems := clean(items)
		if len(cItems) > needed {
			return cItems, true
		}
		return cItems, false
	}
	dataLen := len(fetchedData)
	allData := make([]interface{}, dataLen, dataLen)
	var mutex sync.RWMutex
	var wg sync.WaitGroup
	for i, invite := range fetchedData {
		wg.Add(1)
		go func(i int, data interface{}) {
			defer wg.Done()
			pData := process(data)
			if pData == nil {
				return
			}
			mutex.Lock()
			allData[i] = pData
			mutex.Unlock()
		}(i, invite)
		if (i+1)%(needed) == 0 {
			wg.Wait()
			if cItems, hasNext := check(allData); hasNext {
				return cItems
			}
		}
	}
	wg.Wait()
	cItems, _ := check(allData)
	return cItems
}

func FetchHelper(
	fetchFunc func(context.Context, int) ([]interface{}, error),
	process func(context.Context, []interface{}, int) ([]interface{}, error),
	fetchMore int) PaginateFunc {
	outOfData := false
	return func(ctx context.Context, needed int) ([]interface{}, error) {
		if outOfData {
			return nil, nil
		}
		fetchLimit := needed + fetchMore
		fetchedData, err := fetchFunc(ctx, fetchLimit)
		if err != nil {
			return nil, err
		}
		if len(fetchedData) == 0 {
			return nil, nil
		}
		if len(fetchedData) < fetchLimit {
			outOfData = true
		}
		return process(ctx, fetchedData, needed)
	}
}
