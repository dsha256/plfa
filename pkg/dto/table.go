package dto

import "encoding/json"

type Last20Results struct {
	Time               string `json:"time"`
	Result             int    `json:"result"`
	Color              string `json:"color"`
	GameId             string `json:"gameId"`
	PowerUpList        []any  `json:"powerUpList"`
	PowerUpMultipliers []any  `json:"powerUpMultipliers"`
}

type TableLimits struct {
	Ranges     []float64 `json:"ranges"`
	MinBet     float64   `json:"minBet"`
	MaxBet     float64   `json:"maxBet"`
	MaxPlayers int       `json:"maxPlayers"`
}

type Dealer struct {
	Name string `json:"name"`
}

type PragmaticTable struct {
	TotalSeatedPlayers        int             `json:"totalSeatedPlayers"`
	Last20Results             []Last20Results `json:"last20Results"`
	TableId                   string          `json:"tableId"`
	TableName                 string          `json:"tableName"`
	NewTable                  bool            `json:"newTable"`
	LanguageSpecificTableInfo string          `json:"languageSpecificTableInfo"`
	TableImage                string          `json:"tableImage"`
	TableLimits               TableLimits     `json:"tableLimits"`
	Dealer                    Dealer          `json:"dealer"`
	TableOpen                 bool            `json:"tableOpen"`
	TableType                 string          `json:"tableType"`
	TableSubtype              string          `json:"tableSubtype"`
	Currency                  string          `json:"currency"`
}

// Bytes2PT unmarshal bytes array into PragmaticTable.
func Bytes2PT(data []byte) (PragmaticTable, error) {
	var pt PragmaticTable
	err := json.Unmarshal(data, &pt)
	if err != nil {
		return PragmaticTable{}, err
	}
	return pt, nil
}

type PragmaticTableWithID struct {
	// tID = 100; cID = 200 => TableAndCurrencyID = "100:200"
	TableAndCurrencyID string         `json:"tableAndCurrencyID,omitempty"`
	PragmaticTable     PragmaticTable `json:"pragmaticTable"`
}

// CombTbAndCurrIDs combines tableID and currencyID in the following format: "$(tableID):$(currencyID)"
// Example: tableID = 100, currencyID = USD ==> "100:USD"
func CombTbAndCurrIDs(tableID, currencyID string) string {
	return tableID + ":" + currencyID
}
