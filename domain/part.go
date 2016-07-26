package domain

func (c *Root) GetPart(id string) (Party, int, bool) {
	for i, part := range c.Parties {
		if id == part.Id {
			return part, i, true
		}
	}
	return Party{}, 0, false
}
