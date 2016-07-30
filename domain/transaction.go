package domain

func (acc *Account) GetTransaction(id string) (*Transaction, int, bool) {
	for i, tx := range acc.Transactions {
		if tx.Id == id {
			return &acc.Transactions[i], i, true
		}
	}
	return nil, 0, false
}

type GetTransactionsResult struct {
	BaseAmount          int
	Transactions        []Transaction
	SkippedTransactions bool
}

func (acc *Account) AddTransaction(in Transaction) {
	pos := 0

	for _, tx := range acc.Transactions {
		if in.Date < tx.Date {
			break
		}
		pos++
	}

	// insert the transaction at position "pos"
	acc.Transactions = append(acc.Transactions[:pos], append([]Transaction{in}, acc.Transactions[pos:]...)...)
}

func (acc *Account) GetTransactions(f DateFilter, p PageFilter) GetTransactionsResult {
	balance := acc.StartingBalance
	skippedTransactions := false

	dateFiltered := make([]Transaction, 0)
	for _, tx := range acc.Transactions {
		txAmount := sum(tx.Splits)

		if f.ToDate != nil && tx.Date >= *f.ToDate {
			break
		}

		if f.FromDate != nil && tx.Date < *f.FromDate {
			balance += txAmount
			skippedTransactions = true
			continue
		}

		dateFiltered = append(dateFiltered, tx)
	}

	result := make([]Transaction, 0)

	// [0 1 2 3 4 5 6 7 8]
	//              <--->   offset:0,limit:3
	//          <--->       offset:2,limit:3
	//  <----------->       offset:2,limit:>=6
	if p.Offset < len(dateFiltered) {

		// cut the tail which we do not need
		dateFiltered = dateFiltered[:len(dateFiltered) - p.Offset]

		// calculate how many to skip
		skip := 0
		pageSize := p.Limit
		if len(dateFiltered) > pageSize {
			skip = len(dateFiltered) - pageSize
		}

		for _, tx := range dateFiltered {
			txAmount := sum(tx.Splits)

			if skip--; skip >= 0 {
				balance += txAmount
				skippedTransactions = true
				continue
			}

			result = append(result, tx)
		}
	}

	return GetTransactionsResult{BaseAmount:balance, Transactions:result, SkippedTransactions:skippedTransactions}
}

func (acc *Account) RemoveTransaction(id string) bool {
	if _, i, found := acc.GetTransaction(id); found {
		acc.Transactions = append(acc.Transactions[:i], acc.Transactions[i + 1:]...)
		return true
	}
	return false
}

func sum(in []TransactionSplit) int {
	result := 0
	for _, txS := range in {
		result += txS.Amount
	}
	return result
}