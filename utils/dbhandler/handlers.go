package dbhandler

import (
	"database/sql"

	"github.com/kgjoner/cornucopia/helpers/htypes"
)

func HandleListQuery[T any](rows *sql.Rows, pag *htypes.Pagination, dest func(item *T) []any) (*htypes.PaginatedData[T], error) {
	items := []T{}
	for rows.Next() {
		var item T
		err := rows.Scan(dest(&item)...)
		if err == sql.ErrNoRows {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(items) > pag.Limit {
		pag.HasNext = true
		items = items[:pag.Limit]
	}

	return htypes.NewPaginatedData(*pag, items), nil
}

func HandleListQueryWithoutPagination[T any](rows *sql.Rows, dest func(item *T) []any) ([]T, error) {
	items := []T{}
	for rows.Next() {
		var item T
		err := rows.Scan(dest(&item)...)
		if err == sql.ErrNoRows {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
