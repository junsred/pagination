package pagination

import (
	"context"
	"errors"
	"time"
)

type PaginateFunc[K comparable, V any] func(context.Context, int, *K) ([]*Item[K, V], bool, error)
type NextFunc[K comparable, V any] func(context.Context, *PaginationResult[K, V]) string

type pagOptions struct {
	maxProcessTime time.Duration
	nextFunc       any
}

type PaginationResult[K comparable, V any] struct {
	Items   []*Item[K, V]
	HasNext bool
	Next    *string
}

type Pagination[K comparable, V any] struct {
	options      pagOptions
	paginateFunc PaginateFunc[K, V]
	lastKey      K
}

type Item[K comparable, V any] struct {
	Key   K
	Value V
}

const (
	// max time to give a pagination to complete
	maxProcessTime = 5 * time.Second
)

var (
	ErrResourceExhausted = errors.New("pagination: resource exhausted")
)

type PagOption interface {
	apply(*pagOptions)
}

type funcPagOption struct {
	f func(*pagOptions)
}

func (fdo *funcPagOption) apply(do *pagOptions) {
	fdo.f(do)
}

func newFuncPagOption(f func(*pagOptions)) *funcPagOption {
	return &funcPagOption{
		f: f,
	}
}

func WithTimeout(s time.Duration) PagOption {
	return newFuncPagOption(func(o *pagOptions) {
		o.maxProcessTime = s
	})
}

func WithNextFunc[K comparable, V any](f NextFunc[K, V]) PagOption {
	return newFuncPagOption(func(o *pagOptions) {
		o.nextFunc = f
	})
}

func New[K comparable, V any](pagFunc PaginateFunc[K, V], lastKey K, opts ...PagOption) *Pagination[K, V] {
	p := &Pagination[K, V]{
		options: pagOptions{
			maxProcessTime: maxProcessTime,
		},
		paginateFunc: pagFunc,
		lastKey:      lastKey,
	}
	for _, opt := range opts {
		opt.apply(&p.options)
	}
	return p
}

// Paginate returns a pagination result data list and a boolean indicating if there is a next page
//
//	type Item = pagination.Item[int, string]
//	var items = []*Item{{Key: 1, Value: "Item 1"}, {Key: 2, Value: "Item 2"}}
//	func pagFunc(ctx context.Context, needed int, _ *int) ([]*Item, bool, error) {
//		if needed > len(items) {
//			needed = len(items)
//		}
//		returnItems := items[:needed]
//		items = items[needed:]
//		return returnItems, len(items) > 0, nil
//	}
//	p := pagination.New(pagFunc, 0)
//	pagResult, err := p.Paginate(context.TODO(), 2)
func (p *Pagination[K, V]) Paginate(ctx context.Context, length int) (*PaginationResult[K, V], error) {
	ctx, cancel := context.WithTimeout(ctx, p.options.maxProcessTime)
	defer cancel()
	pagResult := &PaginationResult[K, V]{
		Items:   make([]*Item[K, V], 0, length),
		HasNext: true,
	}
	for pagResult.HasNext && len(pagResult.Items) < length {
		if ctx.Err() != nil {
			return pagResult, ErrResourceExhausted
		}
		var processedItems []*Item[K, V]
		var err error
		processedItems, pagResult.HasNext, err = p.paginateFunc(ctx, length-len(pagResult.Items), &p.lastKey)
		if err != nil {
			return pagResult, err
		}
		pagResult.Items = append(pagResult.Items, processedItems...)
	}
	if len(pagResult.Items) > length {
		pagResult.HasNext = true
		pagResult.Items = pagResult.Items[:length]
	}
	if len(pagResult.Items) > 0 {
		p.lastKey = pagResult.Items[len(pagResult.Items)-1].Key
	}
	if pagResult.HasNext && p.options.nextFunc != nil {
		next := p.options.nextFunc.(NextFunc[K, V])(ctx, pagResult)
		pagResult.Next = &next
	}
	return pagResult, nil
}
