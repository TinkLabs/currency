package pagination

// PaginatedResult represents a common paginated response object for a collection request
type PaginatedResult struct {
	// unexported field to tell whether current response has next page or not
	// it will be used when total field is not available
	hasNext bool `json:"-"`

	Limit   int         `json:"limit"`
	Skip    int         `json:"skip,omitempty"`
	Total   int         `json:"total"`
	Count   int         `json:"count"`
	OrderBy string      `json:"order_by,omitempty"`
	Order   string      `json:"order,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Refs    interface{} `json:"refs,omitempty"`
	Next    *Next       `json:"next,omitempty"`
}

type Next struct {
	Limit int `json:"limit,omitempty"`
	Skip  int `json:"skip,omitempty"`

	// represents position offset for next page
	StartAfter string `json:"start_after,omitempty"`
}

func (b *PaginatedResult) SetLimit(limit int) *PaginatedResult {
	b.Limit = limit
	return b
}

func (b *PaginatedResult) SetTotal(total int) *PaginatedResult {
	b.Total = total
	return b
}

func (b *PaginatedResult) SetCount(count int) *PaginatedResult {
	b.Count = count
	return b
}

func (b *PaginatedResult) SetOrderBy(orderBy string) *PaginatedResult {
	b.OrderBy = orderBy
	return b
}

func (b *PaginatedResult) SetOrder(order string) *PaginatedResult {
	b.Order = order
	return b
}

func (b *PaginatedResult) SetSkip(skip int) *PaginatedResult {
	b.Skip = skip
	return b
}

func (b *PaginatedResult) SetData(data interface{}) *PaginatedResult {
	b.Data = data
	return b
}

func (b *PaginatedResult) SetReferences(refs interface{}) *PaginatedResult {
	b.Refs = refs
	return b
}

func (b *PaginatedResult) SetHasNext(hasNext bool) *PaginatedResult {
	b.hasNext = hasNext
	return b
}

func (b *PaginatedResult) BuildScrollable() *PaginatedResult {
	nextCursor := b.Skip + b.Limit
	if nextCursor < b.Total {
		b.Next = &Next{
			Limit: b.Limit,
			Skip:  nextCursor,
		}
	}
	return b
}

// BuildScrollablePage creates pagination result with skipping number of pages
func (b *PaginatedResult) BuildScrollablePage() *PaginatedResult {
	nextCursor := b.Skip + 1
	if (b.Skip+1)*b.Count < b.Total {
		b.Next = &Next{
			Limit: b.Limit,
			Skip:  nextCursor,
		}
	}
	return b
}

// BuildScrollableWithStartAfter creates pagination result with start_after offset
func (b *PaginatedResult) BuildScrollableWithStartAfter(startAfter string) *PaginatedResult {
	if b.hasNext {
		b.Next = &Next{
			Limit:      b.Limit,
			StartAfter: startAfter,
		}
	}
	return b
}

func New() *PaginatedResult {
	return &PaginatedResult{}
}
