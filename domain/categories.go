package domain

func (c *Root) GetCategory(id string) (*Category, int, bool) {
	for i, cat := range c.Categories {
		if id == cat.Id {
			return &c.Categories[i], i, true
		}
	}
	return nil, 0, false
}

func (c *Root) RemoveCategory(id string) bool {
	if _, i, found := c.GetCategory(id); found {
		c.Categories = append(c.Categories[:i], c.Categories[i + 1:]...)
		return true
	}
	return false
}

func (c *Root) GetTransactionsByCategory(f DateFilter) (result map[string][]TransactionSplit) {
	result = make(map[string][]TransactionSplit, 0)

	for _, account := range c.Accounts {
		for _, tx := range account.Transactions {
			if f.FromDate != nil && tx.Date < *f.FromDate {
				continue
			}
			if f.ToDate != nil && tx.Date >= *f.ToDate {
				continue
			}

			for _, txDetail := range tx.Splits {
				specs, found := result[txDetail.CategoryId]
				if !found {
					specs = make([]TransactionSplit, 0)
				}
				result[txDetail.CategoryId] = append(specs, txDetail)
			}
		}
	}

	return
}
