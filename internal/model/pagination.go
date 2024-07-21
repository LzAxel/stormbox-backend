package model

import (
	"chat-backend/internal/apperror"
)

const (
	MaxLimit     = 100
	DefaultLimit = 10
)

type Pagination struct {
	Offset uint64 `query:"offset"`
	Limit  uint64 `query:"limit"`
}

type FullPagination struct {
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
	Count  uint64 `json:"count"`
	Total  uint64 `json:"total"`
}

func NewPagination(offset, limit uint64) (Pagination, error) {
	if limit > MaxLimit {
		return Pagination{}, apperror.GetErrMaxPaginationLimit(MaxLimit)
	}

	return Pagination{
		Offset: offset,
		Limit:  limit,
	}, nil
}

func (p *Pagination) ToFullPagination(total uint64) FullPagination {
	return FullPagination{
		Offset: p.Offset,
		Limit:  p.Limit,
		Count:  p.Limit,
		Total:  total,
	}
}
