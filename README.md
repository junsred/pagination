# Pagination

[![Go Report Card](https://goreportcard.com/badge/github.com/junsred/pagination)](https://goreportcard.com/report/github.com/junsred/pagination)

Simple, fast and safe pagination tool for go. Now with generics!

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

type Item = pagination.Item[int, string]

var items = []*Item{
	{Key: 1, Value: "Item 1"},
	{Key: 2, Value: "Item 2"},
	{Key: 3, Value: "Item 3"},
}

func fetch(ctx context.Context, limit int, lastItem *int) ([]*Item, bool, error) {
	for i, item := range items {
		if item.Key > *lastItem {
			if limit > len(items)-i {
				limit = len(items) - i
			}
			returnItems := items[i : i+limit]
			return returnItems, len(items[i+limit:]) > 0, nil
		}
	}
	return nil, false, nil
}

func main() {
	ctx := context.TODO()
	var lastItem int
	p := pagination.New(fetch, lastItem)
	pagResult, _ := p.Paginate(ctx, 2)
	for _, item := range pagResult.Items {
		log.Println("first page item", item)
	}
	pagResult, _ = p.Paginate(ctx, 2)
	for _, item := range pagResult.Items {
		log.Println("second page item", item)
	}
}

```
