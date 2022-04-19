Example usage


```go
func pagFunc() PaginateFunc {
	return pagination.FetchHelper(
		func(ctx context.Context, fetchLimit int) ([]interface{}, error) {
			//fetch from db
		},
		func(ctx context.Context, allData []interface{}, needed int) ([]interface{}, error) {
			//filter fetched data
		},
		0,
	)
}

p := pag.New(pagFunc, 20)
pagResult, err := p.Paginate(ctx)
```
