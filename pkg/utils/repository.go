package utils

func CombTbAndCurrIDs(tableID, currencyID string) string {
	return tableID + ":" + currencyID
}
