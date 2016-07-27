package domain

func (c *Root) GetParty(id string) (*Party, int, bool) {
	for i, part := range c.Parties {
		if id == part.Id {
			return &c.Parties[i], i, true
		}
	}
	return nil, 0, false
}

func (c *Root) IsPartyInUse(id string) bool {

	if _, _, found := c.GetParty(id); found {
		for _, acc := range c.Accounts {
			for _, tx := range acc.Transactions {
				if tx.RemotePartyId == id {
					return true
				}
			}
		}
	}

	return false
}

func (c *Root) RemoveParty(id string) bool {
	if _, i, found := c.GetParty(id); found {
		c.Parties = append(c.Parties[:i], c.Parties[i + 1:]...)
		return true
	}
	return false
}
