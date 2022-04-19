package pagination

import (
	"context"
	"errors"
	"time"
)

type PaginateFunc func(context.Context, int) ([]interface{}, error)

type pagOptions struct {
	perPage        int
	paginateFunc   PaginateFunc
	maxProcessTime time.Duration
}

type PaginationResult struct {
	Data    []interface{}
	HasNext bool
}

type Pagination struct {
	options pagOptions
}

type TimeItem struct {
	Data interface{}
	Time time.Time
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

func New(pagFunc PaginateFunc, perPage int, opts ...PagOption) *Pagination {
	p := &Pagination{
		options: pagOptions{
			paginateFunc:   pagFunc,
			perPage:        perPage,
			maxProcessTime: maxProcessTime,
		},
	}
	for _, opt := range opts {
		opt.apply(&p.options)
	}
	return p
}

// Paginate returns a pagination result data list and a boolean indicating if there is a next page
//  func pagFunc() PaginateFunc {
//	 return pagination.FetchHelper(
//		 func(ctx context.Context, fetchLimit int) ([]interface{}, error) {
//			 //fetch from db
//		 },
//		 func(ctx context.Context, allData []interface{}, needed int) ([]interface{}, error) {
//			 //filter fetched data
//		 },
//		 0,
//	 )
//  }
//  p := pag.New(pagFunc, 20)
//  pagResult, err := p.Paginate(ctx)
func (p *Pagination) Paginate(ctx context.Context) (*PaginationResult, error) {
	ctx, cancel := context.WithTimeout(ctx, p.options.maxProcessTime)
	defer cancel()
	pagResult := &PaginationResult{
		Data: make([]interface{}, 0, p.options.perPage),
	}
	for { // fetch until last item is reached
		if ctx.Err() != nil {
			return pagResult, ErrResourceExhausted
		}
		processedData, err := p.options.paginateFunc(ctx, p.options.perPage-len(pagResult.Data)+1)
		if err != nil {
			return pagResult, err
		}
		if processedData == nil { // finish if no data is received
			break
		}
		if len(processedData)+len(pagResult.Data) > p.options.perPage {
			pagResult.Data = append(pagResult.Data, processedData[:p.options.perPage-len(pagResult.Data)]...)
			pagResult.HasNext = true
			break
		}
		pagResult.Data = append(pagResult.Data, processedData...)
	}
	return pagResult, nil
}
