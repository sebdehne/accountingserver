package domain

func (c *Root) GetAccount(id string) (Account, int, bool) {
	for i, acc := range c.Accounts {
		if id == acc.Id {
			return acc, i, true
		}
	}
	return Account{}, 0, false
}

