## Example usage


```go
package main

import (
	"context"
	"log"

	"github.com/junsred/pagination"
)

var items = []interface{}{
	"1", "2", "3",
}

func pagFunc() pagination.PaginateFunc {
	lastItemIndex := 0
	return pagination.FetchHelper(
		func(ctx context.Context, fetchLimit int) ([]interface{}, error) {
			return items[lastItemIndex:], nil
		},
		func(ctx context.Context, allData []interface{}, needed int) ([]interface{}, error) {
			lastItemIndex += needed - 1
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
