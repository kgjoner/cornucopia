package htypes

type Pagination struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	HasNext bool `json:"hasNext"`
}

type PaginationCreationFields struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func NewPagination(c *PaginationCreationFields) *Pagination {
	if c.Limit == 0 {
		c.Limit = 20
	}

	return &Pagination{
		Page:  c.Page,
		Limit: c.Limit,
	}
}

func (p Pagination) Offset() int {
	return p.Page * p.Limit
}

type PaginatedData[T any] struct {
	Data    []T  `json:"data"`
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	HasNext bool `json:"hasNext"`
}

func NewPaginatedData[T any](p Pagination, data []T) *PaginatedData[T] {
	return &PaginatedData[T]{
		Data:    data,
		Page:    p.Page,
		Limit:   p.Limit,
		HasNext: p.HasNext,
	}
}

func TransformPaginatedData[K, T any](data *PaginatedData[T], transformer func(data T) K) *PaginatedData[K] {
	outputData := []K{}
	for _, value := range data.Data {
		outputData = append(outputData, transformer(value))
	}

	return &PaginatedData[K]{
		Data: outputData,
		Page: data.Page,
		Limit: data.Limit,
		HasNext: data.HasNext,
	}
}
