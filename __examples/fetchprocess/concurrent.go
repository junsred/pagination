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
	{Key: 4, Value: "Item 4"},
	{Key: 5, Value: "Item 5"},
}

func fetch(ctx context.Context, fetchLimit int, lastItem *int) ([]*Item, bool, error) {
	l := *lastItem
	hasNext := false
	for in, item := range items { // this simulates a database query
		if item.Key > l {
			limit := in + fetchLimit
			if limit > len(items) {
				limit = len(items)
			} else if limit < len(items) {
				hasNext = true
			}
			r := items[in:limit]
			return r, hasNext, nil
		}
	}
	return nil, hasNext, nil
}

func process(ctx context.Context, item *Item) (*Item, bool) {
	if item.Value == "Item 1" {
		return nil, false
	}
	return item, true
}

func main() {
	ctx := context.TODO()
	var lastItem int
	p := pagination.New(pagination.FetchAndProcess(
		fetch,
		pagination.ProcessConcurrently(process),
	), lastItem)
	pagResult, err := p.Paginate(ctx, 2)
	log.Println(pagResult, err)
	if len(pagResult.Items) > 0 && pagResult.HasNext {
		pagResult, err = p.Paginate(ctx, 2)
		log.Println(pagResult, err)
	}
}
