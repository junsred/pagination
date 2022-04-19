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
	return func(ctx context.Context, needed int) ([]interface{}, error) {
		r := []interface{}{}
		for _, v := range items {
			if v.(int) <= lastItem {
				continue
			}
			r = append(r, v)
			if len(r) != needed {
				lastItem = v.(int)
			}
		}
		if len(r) == 0 {
			return nil, nil // stop the loop
		}
		return r, nil
	}
}

func main() {
	ctx := context.TODO()
	p := pagination.New(pagFunc(), 2)
	pagResult, err := p.Paginate(ctx)
	log.Println(pagResult, err)
	pagResult, err = p.Paginate(ctx)
	log.Println(pagResult, err)
}
