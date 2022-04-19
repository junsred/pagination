package main

import (
	"context"
	"log"

	"github.com/junsred/pagination"
)

var items = []interface{}{
	1, 2, 3, 4,
}

func pagFunc(lastItem *int) pagination.PaginateFunc {
	return pagination.FetchHelper(
		func(ctx context.Context, fetchLimit int) ([]interface{}, error) {
			l := *lastItem
			for in, item := range items { // this simulates a database query
				itm, _ := item.(int)
				if itm > l {
					limit := in + fetchLimit
					if limit > len(items) {
						limit = len(items)
					}
					r := items[in:limit]
					if len(r) != 0 {
						*lastItem = r[len(r)-1].(int)
					}
					return r, nil
				}
			}
			return nil, nil
		},
		func(ctx context.Context, allData []interface{}, needed int) ([]interface{}, error) {
			process := func(data interface{}) interface{} {
				item, _ := data.(int)
				if item == 1 {
					return nil
				}
				return data
			}
			a := pagination.ProcessConcurrently(allData, needed, process)
			return a, nil
		},
		0,
	)
}

func main() {
	ctx := context.TODO()
	lastItem := new(int)
	p := pagination.New(pagFunc(lastItem), 2)
	pagResult, err := p.Paginate(ctx)
	log.Println(pagResult, err)
	if len(pagResult.Data) > 0 && pagResult.HasNext {
		*lastItem = pagResult.Data[len(pagResult.Data)-1].(int)
		pagResult, err = p.Paginate(ctx)
		log.Println(pagResult, err)
	}
}
