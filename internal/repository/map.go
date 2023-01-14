package repository

import (
	"github.com/dsha256/plfa/pkg/dto"
	"github.com/dsha256/plfa/pkg/utils"
	"sync"
)

type MapLiveFeedRepository struct {
	mutex  sync.RWMutex
	tables map[string]dto.PragmaticTable
}

func NewMapLiveFeedRepository() *MapLiveFeedRepository {
	return &MapLiveFeedRepository{
		tables: make(map[string]dto.PragmaticTable),
	}
}

func (m *MapLiveFeedRepository) AddTable(table dto.PragmaticTable) {
	id := utils.CombTbAndCurrIDs(table.TableId, table.Currency)

	m.mutex.Lock()
	m.tables[id] = table
	m.mutex.Unlock()
}

func (m *MapLiveFeedRepository) ListTables() ([]dto.PragmaticTableWithID, error) {
	var pragmaticTables []dto.PragmaticTableWithID

	m.mutex.RLock()
	for k, v := range m.tables {
		pragmaticTables = append(pragmaticTables, dto.PragmaticTableWithID{TableAndCurrencyID: k, PragmaticTable: v})
	}
	m.mutex.RUnlock()

	return pragmaticTables, nil
}
