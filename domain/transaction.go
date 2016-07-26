package domain

type GetTransactionsResult struct {
	BaseAmount   int
	Transactions []Transaction
}

func (acc *Account) GetTransactions(f DateFilter, p PageFilter) GetTransactionsResult {
	result := make([]Transaction, 0)
	balance := acc.StartingBalance
	offsetCount := p.Offset
	limitCount := p.Limit

	for _, tx := range acc.Transactions {
		txAmount := sum(tx.Details)
		if f.FromDate != nil && tx.Date < *f.FromDate {
			balance += txAmount
			continue
		}
		if f.ToDate != nil && tx.Date < *f.ToDate {
			if offsetCount--; offsetCount >= 0 {
				balance += txAmount
				continue
			}

			if limitCount--; limitCount < 0 {
				break
			}
			result = append(result, tx)
		} else {
			break
		}
	}

	return GetTransactionsResult{BaseAmount:balance, Transactions:result}
}

func (acc *Account) GetTransaction(id string) (Transaction, int, bool) {
	for i, tx := range acc.Transactions {
		if tx.Id == id {
			return tx, i, true
		}
	}

	return Transaction{}, 0, false
}

func (acc *Account) RemoveTransaction(id string) bool {
	if _, i, found := acc.GetTransaction(id); found {
		acc.Transactions = append(acc.Transactions[:i], acc.Transactions[i + 1:]...)
		return true
	}
	return false
}

func sum(in []TransactionSpecification) int {
	result := 0
	for _, txS := range in {
		result += txS.Amount
	}
	return result
}