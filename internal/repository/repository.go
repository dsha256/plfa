package repository

import (
	"github.com/dsha256/plfa/pkg/dto"
)

type Repository interface {
	AddTable(table dto.PragmaticTable)
	ListTables() ([]dto.PragmaticTableWithID, error)
}

type AggregateRepository interface {
	Repository
}

type Aggregator struct {
	*MapLiveFeedRepository
}

func NewAggregator() *Aggregator {
	return &Aggregator{MapLiveFeedRepository: NewMapLiveFeedRepository()}
}
