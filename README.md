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

func pagFunc(ctx context.Context, needed int, _ *int) ([]*Item, bool, error) {
	if needed > len(items) {
		needed = len(items)
	}
	returnItems := items[:needed]
	items = items[needed:]
	return returnItems, len(items) > 0, nil
}

func main() {
	p := pagination.New(pagFunc, 0)
	pagResult, err := p.Paginate(context.TODO(), 2)
	log.Println(pagResult.Items, err)
	pagResult, err = p.Paginate(context.TODO(), 2)
	log.Println(pagResult.Items, err)
}

```
