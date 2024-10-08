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
	if limit > len(items) {
		limit = len(items)
	}
	for in, item := range items {
		if item.Key > *lastItem {
			if limit > len(items)-in {
				limit = len(items) - in
			}
			returnItems := items[in : in+limit]
			return returnItems, len(items) > in+limit, nil
		}
	}
	return nil, false, nil
}

func main() {
	ctx := context.TODO()
	var lastItem int
	p := pagination.New(fetch, lastItem)
	pagResult, err := p.Paginate(ctx, 2)
	log.Println(pagResult, err)
	pagResult, err = p.Paginate(ctx, 2)
	log.Println(pagResult, err)
}
