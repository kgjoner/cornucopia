package dbhandler

import (
	"database/sql"

	"github.com/kgjoner/cornucopia/helpers/htypes"
	"github.com/kgjoner/cornucopia/utils/structop"
)

func HandleSingleQuery[K, T any](result T, err error) (*K, error) {
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var resp K
	err = structop.New(result).Copy(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func HandleListQuery[K, T any](result []T, err error) func(*htypes.Pagination) (*htypes.PaginatedData[K], error) {
	return func(pag *htypes.Pagination) (*htypes.PaginatedData[K], error) {
		if err != nil {
			if err.Error() == sql.ErrNoRows.Error() {
				return nil, nil
			} else {
				return nil, err
			}
		}

		if len(result) > pag.Limit {
			pag.HasNext = true
			result = result[:pag.Limit]
		}

		var resp []K
		err := structop.CopySlice(result, &resp)
		if err != nil {
			return nil, err
		}

		return htypes.NewPaginatedData(*pag, resp), nil
	}
}

func HandleListQueryWithoutPagination[K, T any](result []T, err error) ([]K, error) {
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var resp []K
	err = structop.CopySlice(result, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
