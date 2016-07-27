package domain

func (c *Root) GetAccount(id string) (*Account, int, bool) {
	for i, acc := range c.Accounts {
		if id == acc.Id {
			return &c.Accounts[i], i, true
		}
	}
	return nil, 0, false
}

func (c *Root) IsAccountInUse(id string) bool {

	if target, _, found := c.GetAccount(id); found {
		if len(target.Transactions) > 0 {
			return true
		}

		for _, acc := range c.Accounts {
			for _, tx := range acc.Transactions {
				if tx.RemoteAccountId == id {
					return true
				}
			}
		}
	}

	return false
}

func (c *Root) RemoveAccount(id string) bool {
	if _, i, found := c.GetAccount(id); found {
		c.Accounts = append(c.Accounts[:i], c.Accounts[i + 1:]...)
		return true
	}
	return false
}
