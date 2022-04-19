# Pagination

[![Go Report Card](https://goreportcard.com/badge/github.com/junsred/pagination)](https://goreportcard.com/report/github.com/junsred/pagination)

Simple, fast and safe pagination tool for go.

### Documentation

- [Concurrent example](https://github.com/junsred/pagination/tree/main/_examples/concurrent)
- [Fetch example](https://github.com/junsred/pagination/tree/main/_examples/fetch)
- [Simple example](https://github.com/junsred/pagination/tree/main/_examples/simple)

### Example usage

```go
package main

import (
	"context"
	"log"

	"github.com/junsred/pagination"
)

var items = []interface{}{
	1, 2, 3,
}

func pagFunc() pagination.PaginateFunc {
	lastItem := 0
	return pagination.FetchHelper(
		func(ctx context.Context, fetchLimit int) ([]interface{}, error) {
			for in, item := range items {
				i, _ := item.(int)
				if i > lastItem {
					return items[in:], nil
				}
			}
			return nil, nil
		},
		func(ctx context.Context, allData []interface{}, needed int) ([]interface{}, error) {
			lastItemIndex := len(allData)
			if needed == lastItemIndex {
				lastItemIndex -= 1
			}
			lastItem = allData[lastItemIndex-1].(int)
			return allData, nil
		},
		0,
	)
}

func main() {
	ctx := context.TODO()
	p := pagination.New(pagFunc(), 2)
	pagResult, err := p.Paginate(ctx)
	log.Println(pagResult, err)
	pagResult, err = p.Paginate(ctx)
	log.Println(pagResult, err)
}
```
